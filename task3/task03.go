package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

// 使用SQL扩展库进行查询
// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
// 要求 ：
// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func main() {
	db, err := sqlx.Connect("mysql", "root:root@tcp(10.4.8.22:3306)/gorm")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// 查询技术部员工
	techEmployees, err := getEmployeesByDepartment(db, "技术部")
	if err != nil {
		log.Println("查询技术部员工失败:", err)
	} else {
		fmt.Println("技术部员工:")
		for _, emp := range techEmployees {
			fmt.Printf("%+v\n", emp)
		}
	}

	// 查询工资最高员工
	highestPaid, err := getHighestPaidEmployee(db)
	if err != nil {
		log.Println("查询最高工资员工失败:", err)
	} else {
		fmt.Println("\n工资最高员工:")
		fmt.Printf("%+v\n", highestPaid)
	}
}

func getEmployeesByDepartment(db *sqlx.DB, department string) ([]Employee, error) {
	var employees []Employee
	query := `SELECT id, name, department, salary FROM employees WHERE department = ?`
	err := db.Select(&employees, query, department)
	return employees, err
}

func getHighestPaidEmployee(db *sqlx.DB) (Employee, error) {
	var employee Employee
	query := `SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1`
	err := db.Get(&employee, query)
	return employee, err
}
