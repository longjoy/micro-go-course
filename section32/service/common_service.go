package service

import "fmt"

type Service interface {

	Index() string

	Sample(username string) string

	Admin(username string)  string

	// HealthCheck check service health status
	HealthCheck() bool
}

type CommonService struct {

}

func (s *CommonService) Index() string {
	return fmt.Sprintf("hello, wecome to index")
}

func (s *CommonService) Sample(username string) string {
	return fmt.Sprintf("hello %s, wecome to sample", username)
}

func (s *CommonService) Admin(username string) string {
	return fmt.Sprintf("hello %s, wecome to admin", username)

}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s *CommonService) HealthCheck() bool {
	return true
}

func NewCommonService() *CommonService {
	return &CommonService{}
}

