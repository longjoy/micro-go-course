package dao

import (
	"testing"
)

func TestUserDAOImpl_Save(t *testing.T) {
	userDAO := &UserDAOImpl{}

	err := InitMysql("127.0.0.1", "3306", "root", "xuan", "user")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	user := &UserEntity{
		Username:"aoho",
		Password:"aoho",
		Email:"aoho@mail.com",
	}

	err = userDAO.Save(user)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	t.Logf("new User ID is %d", user.ID)

}

func TestUserDAOImpl_SelectByEmail(t *testing.T) {

	userDAO := &UserDAOImpl{}

	err := InitMysql("127.0.0.1", "3306", "root", "xuan", "user")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}

	user, err := userDAO.SelectByEmail("aoho@mail.com")
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
	t.Logf("result uesrname is %s", user.Username)

}
