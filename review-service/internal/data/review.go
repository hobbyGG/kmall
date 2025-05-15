package data

import (
	"context"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"review-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
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
