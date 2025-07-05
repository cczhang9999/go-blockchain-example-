package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Age  int
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 查询
	var users []User
	db.Where("age > ?", 18).Find(&users)
	result := db.Find(&users)
	fmt.Printf("共查到 %d 条数据\n", result.RowsAffected)
	fmt.Printf("users", users)
}
