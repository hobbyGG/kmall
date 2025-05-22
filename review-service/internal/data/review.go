package data

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
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
func (r *ReviewRepo) ListReviewByStoreID(ctx context.Context, storeID int64, page, size int32) ([]*model.ReviewInfo, error) {
	// 从es获取指定商家id的评论数据
	resp, err := r.data.ESCli.Search().
		Index("review_info").
		From(int(page - 1)).
		Size(int(size)).
		Query(&types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{
						Term: map[string]types.TermQuery{
							"store_id": {Value: storeID},
						},
					},
				},
			},
		}).Do(ctx)
	if err != nil {
		r.log.Debugf("es search error: %v", err)
		return nil, err
	}

	// 处理es返回的数据
	infos := make([]*model.ReviewInfo, 0, resp.Hits.Total.Value)
	for _, hit := range resp.Hits.Hits {
		dataJson := hit.Source_
		// 反序列化json数据到结构体
		temp := biz.ReviewInfo{}
		if err := json.Unmarshal(dataJson, &temp); err != nil {
			r.log.Debugf("json unmarshal error: %v", err)
			return nil, err
		}
		infos = append(infos, &model.ReviewInfo{
			ID:             temp.ID,
			CreateBy:       temp.CreateBy,
			CreateAt:       time.Time(temp.CreateAt),
			UpdateAt:       time.Time(temp.UpdateAt),
			Version:        temp.Version,
			DeleteAt:       temp.DeleteAt,
			ReviewID:       temp.ReviewID,
			OrderID:        temp.OrderID,
			StoreID:        temp.StoreID,
			UserID:         temp.UserID,
			Socore:         temp.Socore,
			Content:        temp.Content,
			Status:         temp.Status,
			IsDefault:      temp.IsDefault,
			HasReply:       temp.HasReply,
			ExpressScore:   temp.ExpressScore,
			ServiceScore:   temp.ServiceScore,
			HasMedia:       temp.HasMedia,
			SkuID:          temp.SkuID,
			SpuID:          temp.SpuID,
			Anonymous:      temp.Anonymous,
			Tags:           temp.Tags,
			OpReason:       temp.OpReason,
			OpUser:         temp.OpUser,
			OpRemark:       temp.OpRemark,
			ExtJSON:        temp.ExtJSON,
			CtrlJSON:       temp.CtrlJSON,
			GoodsSnapshoot: temp.GoodsSnapshoot,
		})
	}
	r.log.Debugf("--->total %d es data, get %d, hits len %d\n", resp.Hits.Total.Value, len(infos), len(resp.Hits.Hits))

	return infos, nil
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
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return -1, err
		}
		// 没有错误，可能不存在
	}
	if ret != nil {
		// 如果存在
		if ret.Status != 10 {
			return -1, errors.New("该申诉已经审核过了")
		}
		appeal.AppealID = ret.AppealID
	}

	// 如果为10就update，不存在则创建，其他就返回错误
	r.data.Q.ReviewAppealInfo.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "review_id"}}, // 冲突列,
			DoUpdates: clause.AssignmentColumns([]string{
				"status",
				"content",
				"reason",
			}),
		},
	).Create(appeal)
	return appeal.AppealID, nil
}
func (r *ReviewRepo) OperateAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) (int64, error) {
	// 检查有申诉是否存在
	if _, err := r.data.Q.ReviewAppealInfo.WithContext(ctx).Where(r.data.Q.ReviewAppealInfo.AppealID.Eq(appeal.AppealID)).First(); err != nil {
		return -1, err
	}
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
	return appeal.AppealID, nil
}
