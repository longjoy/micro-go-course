package service

import (
	"context"
	"log"
)

type CommentListVO struct {
	Id          string
	CommentList []CommentVo
}

type CommentVo struct {
	Id      string
	Desc    string
	Score   float32
	ReplyId string
}

type Service interface {
	GetCommentsList(ctx context.Context, id string) (CommentListVO, error)
	HealthCheck() string
}

func NewGoodsServiceImpl() Service {
	return &CommentsServiceImpl{}
}

type CommentsServiceImpl struct{}

func (service *CommentsServiceImpl) GetCommentsList(ctx context.Context, id string) (CommentListVO, error) {
	comment1 := CommentVo{Id: "1", Desc: "comments", Score: 1.0, ReplyId: "0"}
	comment2 := CommentVo{Id: "2", Desc: "comments", Score: 1.0, ReplyId: "1"}

	list := []CommentVo{comment1, comment2}
	detail := CommentListVO{Id: id, CommentList: list}
	log.Printf(detail.Id)
	return detail, nil
}

func (service *CommentsServiceImpl) HealthCheck() string {
	return "OK"
}
