package biz

import (
	"strings"
	"time"

	"github.com/google/wire"
	"github.com/hobbyGG/kmall/review-service/internal/data/model"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewReviewUsecase)

type DataTime time.Time

type ReviewInfo struct {
	model.ReviewInfo
	CreateAt     DataTime `json:"create_at"`
	UpdateAt     DataTime `json:"update_at"`
	ID           int64    `json:"id,string"`
	Version      int32    `json:"version,string"`
	ReviewID     int64    `json:"review_id,string"`
	Socore       int32    `json:"socore,string"`
	ServiceScore int32    `json:"service_score,string"`
	ExpressScore int32    `json:"express_score,string"`
	HasMedia     int32    `json:"has_media,string"`
	OrderID      int64    `json:"order_id,string"`
	SkuID        int64    `json:"sku_id,string"`
	SpuID        int64    `json:"spu_id,string"`
	StoreID      int64    `json:"store_id,string"`
	UserID       int64    `json:"user_id,string"`
	Anonymous    int32    `json:"anonymous,string"`
	Status       int32    `json:"status,string"`
	IsDefault    int32    `json:"is_default,string"`
	HasReply     int32    `json:"has_reply,string"`
}

func (dt *DataTime) UnmarshalJSON(data []byte) error {
	timeString := strings.Trim(string(data), `"`)
	t, err := time.Parse(time.DateTime, timeString)
	if err != nil {
		return err
	}
	*dt = DataTime(t)
	return nil
}
