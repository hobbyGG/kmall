package service

import (
	"context"

	pb "github.com/hobbyGG/kmall/review-service/api/review/v1"
	"github.com/hobbyGG/kmall/review-service/internal/biz"
	"github.com/hobbyGG/kmall/review-service/internal/data/model"
)

type ReviewService struct {
	pb.UnimplementedReviewServer

	uc *biz.ReviewUsecase
}

func NewReviewService(uc *biz.ReviewUsecase) *ReviewService {
	return &ReviewService{uc: uc}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	// 参数处理
	// 1.validate基本参数处理
	// 调用biz层逻辑
	var anonymous int32
	if req.Anonymous {
		anonymous = 1
	}
	review := &model.ReviewInfo{
		UserID:       req.UserID,
		OrderID:      req.OrderID,
		StoreID:      req.StoreID,
		Socore:       req.Score,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpressScore,
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    anonymous,
	}
	r, err := s.uc.SaveReview(ctx, review)
	if err != nil {
		return nil, err
	}
	// 返回响应
	return &pb.CreateReviewReply{ReviewID: r.ID}, nil
}
func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.UpdateReviewReply, error) {
	return &pb.UpdateReviewReply{}, nil
}
func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewReply, error) {
	return &pb.DeleteReviewReply{}, nil
}
func (s *ReviewService) GetReviewByRID(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	rid := req.ReviewID
	// 调用biz层服务
	review, err := s.uc.GetReviewByReviewID(ctx, rid)
	if err != nil {
		return nil, err
	}
	return &pb.GetReviewReply{
		ReviewID:  review.ID,
		UserID:    review.UserID,
		OrderID:   review.OrderID,
		StoreID:   review.StoreID,
		Score:     review.Socore,
		Content:   review.Content,
		PicInfo:   review.PicInfo,
		VideoInfo: review.VideoInfo,
	}, nil
}
func (s *ReviewService) ListReview(ctx context.Context, req *pb.ListReviewRequest) (*pb.ListReviewReply, error) {
	return &pb.ListReviewReply{}, nil
}

func (s *ReviewService) CreateAppeal(ctx context.Context, req *pb.CreateAppealRequest) (*pb.CreateAppealReply, error) {
	// 参数处理
	appeal := &model.ReviewAppealInfo{
		ReviewID: req.ReviewID,
		StoreID:  req.StoreID,
		Content:  req.Content,
		Reason:   req.Reason,
	}
	aid, err := s.uc.CreateAppeal(ctx, appeal)
	if err != nil {
		return nil, err
	}
	return &pb.CreateAppealReply{AppealID: aid}, nil
}
func (s *ReviewService) OperateAppeal(ctx context.Context, req *pb.OperateAppealRequest) (*pb.OperateAppealReply, error) {
	// 参数处理
	appeal := &model.ReviewAppealInfo{
		AppealID: req.AppealID,
		ReviewID: req.ReviewID,
		StoreID:  req.StoreID,
		Status:   req.Status,
		OpRemark: req.OpRemark,
		OpUser:   req.OpUser,
	}
	aid, err := s.uc.OperateAppeal(ctx, appeal)
	if err != nil {
		return nil, err
	}
	return &pb.OperateAppealReply{AppealID: aid}, nil
}
func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	// 基本处理
	reply := &model.ReviewReplyInfo{
		ReviewID:  req.ReviewID,
		StoreID:   req.StoreID,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	}
	// 调用biz层服务
	reply, err := s.uc.ReplyReview(ctx, reply)
	if err != nil {
		return nil, err
	}

	// 包装返回参数
	return &pb.ReplyReviewReply{ReplyID: reply.ID}, nil
}
