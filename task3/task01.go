package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

//假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
//要求 ：
//编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
//编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
//编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
//编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。

type Student struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"size:50"`
	Age   int
	Grade string `gorm:"size:10"`
}

var db *gorm.DB

func initDB() {
	dsn := "root:root@tcp(10.4.8.22:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	db.AutoMigrate(&Student{})
}

func createStudent(name string, age int, grade string) {
	student := Student{Name: name, Age: age, Grade: grade}
	result := db.Create(&student)
	if result.Error != nil {
		log.Println("创建失败:", result.Error)
		return
	}
	fmt.Printf("创建成功，ID: %d\n", student.ID)
}

func getStudent() {
	var student []Student
	//result := db.First(&student, id)
	result := db.Where("Age > ?", 18).Find(&student)
	if result.Error != nil {
		log.Println("查询失败:", result.Error)
		return
	}
	for _, s := range student {
		fmt.Printf("学生信息: ID:%d 姓名:%s 年龄:%d 年级:%s\n",
			s.ID, s.Name, s.Age, s.Grade)

	}
}

func updateStudent() {
	result := db.Model(&Student{}).Where("Name = ?", "张三").
		Updates(Student{Grade: "四年级"})
	if result.Error != nil {
		log.Println("更新失败:", result.Error)
		return
	}
	fmt.Printf("更新了%d条记录\n", result.RowsAffected)
}

func deleteStudent() {
	result := db.Where("Age < ? ", 15).Delete(&Student{})
	if result.Error != nil {
		log.Println("删除失败:", result.Error)
		return
	}
	fmt.Printf("删除了%d条记录\n", result.RowsAffected)
}

func listStudents() {
	var students []Student
	result := db.Find(&students)
	if result.Error != nil {
		log.Println("查询失败:", result.Error)
		return
	}
	for _, s := range students {
		fmt.Printf("ID:%d 姓名:%s 年龄:%d 年级:%s\n",
			s.ID, s.Name, s.Age, s.Grade)
	}
}

func main() {
	initDB()

	// 示例操作
	/*	createStudent("张三", 20, "三年级")
		createStudent("李四", 17, "三年级")*/

	getStudent()
	updateStudent()
	listStudents()

	deleteStudent()
}
