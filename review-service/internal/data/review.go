package data

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hobbyGG/kmall/review-service/internal/biz"
	"github.com/hobbyGG/kmall/review-service/internal/data/model"
	"github.com/hobbyGG/kmall/review-service/internal/data/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReviewRepo struct {
	data *Data
	log  *log.Helper
}

func NewReviewRepo(data *Data, logger log.Logger) biz.ReviewRepo {
	return &ReviewRepo{data: data, log: log.NewHelper(logger)}
}

func (r *ReviewRepo) SaveReview(ctx context.Context, rinfo *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.Q.ReviewInfo.WithContext(ctx).Create(rinfo)
	return rinfo, err
}
func (r *ReviewRepo) GetReviewByOrderID(ctx context.Context, orderID int64) ([]*model.ReviewInfo, error) {
	return r.data.Q.ReviewInfo.
		WithContext(ctx).
		Where(r.data.Q.ReviewInfo.OrderID.Eq(orderID)).
		Find()
}

func (r *ReviewRepo) GetReviewByReviewID(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.Q.ReviewInfo.
		WithContext(ctx).
		Where(r.data.Q.ReviewInfo.ReviewID.Eq(reviewID)).
		First()
}

func (r *ReviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) error {
	// 回复存储涉及两个表，需要使用事务操作
	r.data.Q.Transaction(func(tx *query.Query) error {
		if err := tx.ReviewReplyInfo.WithContext(ctx).Create(reply); err != nil {
			return err
		}
		if _, err := tx.ReviewInfo.
			WithContext(ctx).
			Where(tx.ReviewInfo.ReviewID.Eq(reply.ReviewID)).
			Update(tx.ReviewInfo.HasReply, 1); err != nil {
			return err
		}
		return nil
	})
	return nil
}

func (r *ReviewRepo) CreateAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) (int64, error) {
	// 检查appeal状态
	ret, err := r.data.Q.ReviewAppealInfo.WithContext(ctx).Where(r.data.Q.ReviewAppealInfo.ReviewID.Eq(appeal.ReviewID)).First()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return -1, err
		}
	}
	if ret != nil {
		appeal.AppealID = ret.AppealID
	}
	if ret.Status != 10 {
		return -1, errors.New("该申诉已经审核过了")
	}

	// 如果为10就update，不存在则创建，其他就返回错误
	r.data.Q.ReviewAppealInfo.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "review_id"}}, // 冲突列,
			DoUpdates: clause.AssignmentColumns([]string{
				"status",
				"content",
				"reason",
				"update_time",
			}),
		},
	).Create(appeal)
	return -1, nil
}
func (r *ReviewRepo) OperateAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) (int64, error) {
	if _, err := r.data.Q.ReviewAppealInfo.WithContext(ctx).Where(r.data.Q.ReviewAppealInfo.AppealID.Eq(appeal.AppealID)).Updates(map[string]any{
		"status":    appeal.Status,
		"op_remark": appeal.OpRemark,
		"op_user":   appeal.OpUser,
	}); err != nil {
		return -1, nil
	}
	if appeal.Status == 20 {
		// 申诉通过则需要隐藏对应的评论
		_, err := r.data.Q.ReviewInfo.WithContext(ctx).Where(r.data.Q.ReviewInfo.ReviewID.Eq(appeal.ReviewID)).Update(r.data.Q.ReviewAppealInfo.Status, 40)
		if err != nil {
			return -1, err
		}
	}
	return 0, nil
}
