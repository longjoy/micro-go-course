package common

type CommentListVO struct {
	Id          string      `json:"Id"`
	CommentList []CommentVo `json:"list"`
}

type CommentVo struct {
	Id      string  `json:"Id"`
	Desc    string  `json:"Desc"`
	Score   float32 `json:"Score"`
	ReplyId string  `json:"ReplyId"`
}
