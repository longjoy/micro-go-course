package user

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type ServiceMiddleware func(service UserService) UserService

type loggingMiddleware struct {
	UserService
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next UserService) UserService {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) CheckPassword(ctx context.Context, username, password string) (ret bool, err error) {

	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "CheckPassword",
			"username", username,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.UserService.CheckPassword(ctx, username, password)
	return ret, err
}
