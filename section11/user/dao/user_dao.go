package dao

import "time"

type UserEntity struct {

	ID int64
	Username string
	Password string
	Email string
	CreatedAt time.Time
}

func (UserEntity) TableName() string {
	return "user"
}

type UserDAO interface {
	SelectByEmail(email string)(*UserEntity, error)
	Save(user *UserEntity) error
}

type UserDAOImpl struct {

}

func (userDAO *UserDAOImpl) SelectByEmail(email string)(*UserEntity, error) {
	user := &UserEntity{}
	err := db.Where("email = ?", email).First(user).Error
	return user, err
}

func (userDAO *UserDAOImpl) Save(user *UserEntity) error {
	return db.Create(user).Error
}
