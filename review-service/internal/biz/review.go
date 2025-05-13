package biz

import (
	"context"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/GenID"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type ReviewRepo interface {
	// 传入的是gen生成的结构体
	// query.query是gen生成的全局查询结构体，该结构体封装了所有表的查询方法
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
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
