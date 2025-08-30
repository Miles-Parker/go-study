package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

//：实现类型安全映射
//假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
//要求 ：
//定义一个 Book 结构体，包含与 books 表对应的字段。
//编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func main() {
	db, err := sqlx.Connect("mysql", "root:root@tcp(10.4.8.22:3306)/gorm")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	expensiveBooks, err := getBooksByMinPrice(db, 50.0)
	if err != nil {
		log.Println("查询失败:", err)
		return
	}

	fmt.Println("价格超过50元的书籍:")
	for _, book := range expensiveBooks {
		fmt.Printf("%+v\n", book)
	}
}

func getBooksByMinPrice(db *sqlx.DB, minPrice float64) ([]Book, error) {
	var books []Book
	query := `SELECT id, title, author, price FROM books WHERE price > ? ORDER BY price DESC`
	err := db.Select(&books, query, minPrice)
	return books, err
}
