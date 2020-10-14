package service

type Service interface {
	// HealthCheck check service health status
	HealthCheck() bool
}

type CommonService struct {

}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s *CommonService) HealthCheck() bool {
	return true
}

func NewCommonService() *CommonService {
	return &CommonService{}
}

