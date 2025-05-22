package biz

import (
	"context"
	"errors"

	v1 "github.com/hobbyGG/kmall/review-service/api/review/v1"
	"github.com/hobbyGG/kmall/review-service/internal/data/model"
	"github.com/hobbyGG/kmall/review-service/pkg/GenID"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ReviewRepo interface {
	// 传入的是gen生成的结构体
	// query.query是gen生成的全局查询结构体，该结构体封装了所有表的查询方法
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
	GetReviewByReviewID(context.Context, int64) (*model.ReviewInfo, error)

	SaveReply(context.Context, *model.ReviewReplyInfo) error
	CreateAppeal(context.Context, *model.ReviewAppealInfo) (int64, error)
	OperateAppeal(context.Context, *model.ReviewAppealInfo) (int64, error)
	ListReviewByStoreID(context.Context, int64, int32, int32) ([]*model.ReviewInfo, error)
}

type ReviewUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

func NewReviewUsecase(repo ReviewRepo, logger log.Logger) *ReviewUsecase {
	return &ReviewUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateReview 为service层提供的业务逻辑方法，传入上下文和评论结构体，将评论存入数据库中
func (uc *ReviewUsecase) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Infof("SaveReview: %v", review)
	// 业务参数校验
	// 不允许重复
	r, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	if len(r) > 0 {
		return nil, v1.ErrorOrderReviewed("订单%v已经评价过了", review.OrderID)
	}

	// 填充id等字段
	review.ReviewID = GenID.Get()

	return uc.repo.SaveReview(ctx, review)
}
func (uc *ReviewUsecase) GetReviewByOrderID(ctx context.Context, orderID int64) (*model.ReviewInfo, error) {
	reviews, err := uc.repo.GetReviewByOrderID(ctx, orderID)
	review := reviews[0]
	return review, err
}
func (uc *ReviewUsecase) GetReviewByReviewID(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return uc.repo.GetReviewByReviewID(ctx, reviewID)
}
func (uc *ReviewUsecase) ListReviewByStoreID(ctx context.Context, storeID int64, page, size int32) ([]*model.ReviewInfo, error) {
	// 简单参数处理
	if storeID <= 0 {
		uc.log.Debugf("[biz] ListReviewByStore failed, Store:%v", storeID)
		return nil, errors.New("StoreID is invalid")
	}
	if page <= 0 {
		uc.log.Debugf("[biz] ListReviewByStoreID failed, page:%v", page)
		return nil, errors.New("page is invalid")
	}
	if size <= 0 {
		uc.log.Debugf("[biz] ListReviewByStoreID failed, size:%v", size)
		return nil, errors.New("size is invalid")
	}

	reviewList, err := uc.repo.ListReviewByStoreID(ctx, storeID, page, size)
	if err != nil {
		uc.log.Debugf("[biz] ListReviewByStoreID failed, err:%v", err)
		return nil, err
	}
	return reviewList, nil
}

func (uc *ReviewUsecase) ReplyReview(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	// 业务参数处理
	// 检查是否存在评论
	review, err := uc.repo.GetReviewByReviewID(ctx, reply.ReviewID)
	if err != nil {
		uc.log.Debugf("GetReviewByOrderID failed, err:%v", err)
		return nil, err
	}
	if review == nil {
		// 评论不存在
		return nil, errors.New("评论不存在")
	}

	// 水平越权处理
	if review.StoreID != reply.StoreID {
		// 评论的商家与回复的商家不同则发生越权
		return nil, errors.New("不能回复其他商家的评论")
	}
	if review.HasReply != 0 {
		// 已经评论过
		return nil, errors.New("不允许重复回复")
	}

	// 将回复存入数据库
	reply.ReplyID = GenID.Get()
	uc.repo.SaveReply(ctx, reply)
	return &model.ReviewReplyInfo{ReplyID: reply.ReplyID, ReviewID: review.ID}, nil
}
func (uc *ReviewUsecase) CreateAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) (int64, error) {
	// 业务参数处理
	// 检查是否存在评论
	// 检查是否越权
	// 检查是否已经申诉过
	review, err := uc.GetReviewByReviewID(ctx, appeal.ReviewID)
	if err != nil {
		// 出错或不存在评论
		uc.log.Debugf("[biz]GetReviewByOrderID failed, err:%v", err)
		return -1, err
	}
	if appeal.StoreID != review.StoreID {
		return -1, errors.New("不能申诉其他商家的评论")
	}
	if review.Status == 40 {
		return -1, errors.New("该评论已经申诉过")
	}

	appeal.AppealID = GenID.Get()
	appeal.Status = 10
	aid, err := uc.repo.CreateAppeal(ctx, appeal)
	if err != nil {
		return -1, err
	}
	return aid, nil
}

func (uc *ReviewUsecase) OperateAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) (int64, error) {
	return uc.repo.OperateAppeal(ctx, appeal)
}
