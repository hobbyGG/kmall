// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameReviewInfo = "review_info"

// ReviewInfo 评价信息表
type ReviewInfo struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键" json:"id"`                      // 主键
	CreateBy       string     `gorm:"column:create_by;not null;comment:创建方式标识" json:"create_by"`                         // 创建方式标识
	UpdateBy       string     `gorm:"column:update_by;not null;comment:更新方式标识" json:"update_by"`                         // 更新方式标识
	CreateAt       time.Time  `gorm:"column:create_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_at"` // 创建时间
	UpdateAt       time.Time  `gorm:"column:update_at;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_at"` // 更新时间
	DeleteAt       *time.Time `gorm:"column:delete_at;comment:逻辑删除标记" json:"delete_at"`                                  // 逻辑删除标记
	Version        int32      `gorm:"column:version;not null;comment:乐观锁标记" json:"version"`                              // 乐观锁标记
	ReviewID       int64      `gorm:"column:review_id;not null;comment:评论id" json:"review_id"`                           // 评论id
	Content        string     `gorm:"column:content;not null;comment:评价内容" json:"content"`                               // 评价内容
	Socore         int32      `gorm:"column:socore;not null;comment:评分" json:"socore"`                                   // 评分
	ServiceScore   int32      `gorm:"column:service_score;not null;comment:服务评分" json:"service_score"`                   // 服务评分
	ExpressScore   int32      `gorm:"column:express_score;not null;comment:物流评分" json:"express_score"`                   // 物流评分
	HasMedia       int32      `gorm:"column:has_media;not null;comment:是否有图或视频" json:"has_media"`                        // 是否有图或视频
	OrderID        int64      `gorm:"column:order_id;not null;comment:订单id" json:"order_id"`                             // 订单id
	SkuID          int64      `gorm:"column:sku_id;not null;comment:商品id" json:"sku_id"`                                 // 商品id
	SpuID          int64      `gorm:"column:spu_id;not null;comment:货号" json:"spu_id"`                                   // 货号
	StoreID        int64      `gorm:"column:store_id;not null;comment:店铺id" json:"store_id"`                             // 店铺id
	UserID         int64      `gorm:"column:user_id;not null;comment:用户id" json:"user_id"`                               // 用户id
	Anonymous      int32      `gorm:"column:anonymous;not null;comment:是否匿名" json:"anonymous"`                           // 是否匿名
	Tags           string     `gorm:"column:tags;not null;comment:标签json" json:"tags"`                                   // 标签json
	PicInfo        string     `gorm:"column:pic_info;not null;comment:图片信息json" json:"pic_info"`                         // 图片信息json
	VideoInfo      string     `gorm:"column:video_info;not null;comment:视频信息json" json:"video_info"`                     // 视频信息json
	Status         int32      `gorm:"column:status;not null;comment:状态:10-待审核，20-审核通过，30-审核不通过，40-隐藏" json:"status"`     // 状态:10-待审核，20-审核通过，30-审核不通过，40-隐藏
	IsDefault      int32      `gorm:"column:is_default;not null;comment:是否默认评价" json:"is_default"`                       // 是否默认评价
	HasReply       int32      `gorm:"column:has_reply;not null;comment:是否有回复" json:"has_reply"`                          // 是否有回复
	OpReason       string     `gorm:"column:op_reason;not null;comment:审核拒绝原因" json:"op_reason"`                         // 审核拒绝原因
	OpRemark       string     `gorm:"column:op_remark;not null;comment:审核备注" json:"op_remark"`                           // 审核备注
	OpUser         string     `gorm:"column:op_user;not null;comment:审核人" json:"op_user"`                                // 审核人
	GoodsSnapshoot string     `gorm:"column:goods_snapshoot;not null;comment:商品快照" json:"goods_snapshoot"`               // 商品快照
	ExtJSON        string     `gorm:"column:ext_json;not null;comment:扩展信息json" json:"ext_json"`                         // 扩展信息json
	CtrlJSON       string     `gorm:"column:ctrl_json;not null;comment:控制信息json" json:"ctrl_json"`                       // 控制信息json
}

// TableName ReviewInfo's table name
func (*ReviewInfo) TableName() string {
	return TableNameReviewInfo
}
