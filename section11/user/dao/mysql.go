package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)


var db *gorm.DB

func InitMysql(host, port, user, password, dbName string) (err error) {
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbName))
	if err != nil{
		log.Println(err)
		return
	}
	db.SingularTable(true)
	return
}


