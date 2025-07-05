package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	// 1. 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 2. 创建表
	createTable := `
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(64) NOT NULL,
            age INT
        );
    `
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("建表失败:", err)
	}

	// 3. 插入数据
	insertSQL := "INSERT INTO users (name, age) VALUES (?, ?)"
	res, err := db.Exec(insertSQL, "Alice", 20)
	if err != nil {
		log.Fatal("插入失败:", err)
	}
	id, _ := res.LastInsertId()
	fmt.Println("插入ID:", id)

	// 4. 查询数据
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	defer rows.Close()

	fmt.Println("所有用户：")
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID:%d, Name:%s, Age:%d\n", u.ID, u.Name, u.Age)
	}

	// 5. 更新数据
	_, err = db.Exec("UPDATE users SET age=? WHERE name=?", 21, "Alice")
	if err != nil {
		log.Fatal("更新失败:", err)
	}
	fmt.Println("更新成功")

}
