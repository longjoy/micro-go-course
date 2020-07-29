package service

import (
	"context"
	"github.com/longjoy/micro-go-course/section08/user/dao"
	"github.com/longjoy/micro-go-course/section08/user/redis"
	"testing"
)

func TestUserServiceImpl_Login(t *testing.T) {


	err := dao.InitMysql("127.0.0.1", "3306", "root", "xuan", "user")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("127.0.0.1","6379", "" )
	if err != nil{
		t.Error(err)
		t.FailNow()
	}


	userService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}

	user, err := userService.Login(context.Background(), "aoho@mail.com", "aoho")

	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)

}

func TestUserServiceImpl_Register(t *testing.T) {


	err := dao.InitMysql("127.0.0.1", "3306", "root", "xuan", "user")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("127.0.0.1","6379", "" )
	if err != nil{
		t.Error(err)
		t.FailNow()
	}


	userService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}

	user, err := userService.Register(context.Background(),
		&RegisterUserVO{
			Username:"aoho",
			Password:"aoho",
			Email:"aoho@mail.com",
		})

	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)

}