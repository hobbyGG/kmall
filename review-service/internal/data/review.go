package data

import (
	"context"
	"review-service/internal/biz"
	"review-service/internal/data/model"

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
