package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm-test/randUtils"
	"math/rand"
	"testing"
	"time"
)

type Funtester struct {
	gorm.Model
	Name string
	Age  int
}

var drive *gorm.DB
var err error

func init() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("start")
	drive, err = gorm.Open("mysql", "root:123456@(localhost:3306)/mysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		return
	}
	drive.DB().SetMaxOpenConns(200)
	drive.DB().SetConnMaxLifetime(10 * time.Second)
	drive.DB().SetConnMaxIdleTime(10 * time.Second)
	drive.DB().SetMaxIdleConns(20)
	drive.AutoMigrate(&Funtester{})
	fmt.Println("success")
}

func TestInsert(t *testing.T) {
	value := &Funtester{
		Name: "Funtester" + randUtils.RandAllString(10),
		Age:  17,
	}
	db1 := drive.Create(value)
	fmt.Println(db1.RowsAffected)
	db2 := drive.Select("name", "age").Create(value)
	fmt.Println(db2.RowsAffected)
	//time.Sleep(1000)
	db3 := drive.Omit("age", "name").Create(&Funtester{Name: "fds", Age: 18})
	fmt.Println(db3.RowsAffected)
	//fs := []Funtester{
	//	{Name: "fs" + randUtils.RandAllString(10), Age: 12},
	//	{Name: "fs" + randUtils.RandAllString(10), Age: 13},
	//}
	//drive.Create(&fs)
}

func TestSelect(t *testing.T) {
	fmt.Println("Test Select start")
	var f Funtester
	drive.First(&f, 34)
	last := drive.Last(&f, "age != 1")
	fmt.Printf("查询记录数 %d", last.RowsAffected)
	fmt.Println(f)
	task := drive.Take(&f)
	fmt.Println(task.RowsAffected)
	fmt.Println(f)
}

func TestSelect2(t *testing.T) {
	fmt.Println("Test Select2 start")
	var f Funtester
	var fs []Funtester
	drive.Where("age = ?", 19).First(&f)
	fmt.Printf("查询结果：")
	fmt.Println(f)
	find := drive.Where("name like ?", "Fun%").Find(&fs).Limit(20).Order("id")
	rows, _ := find.Rows()
	defer rows.Close()
	for rows.Next() {
		var ff Funtester
		drive.ScanRows(rows, &ff)
	}
	var f1 Funtester
	drive.Where("name like ?", "fun").Or("age = ?", 12).First(&f1)
	fmt.Printf("first: ")
	fmt.Println(f1)
}

func TestUpdate(t *testing.T) {
	drive.Model(&Funtester{}).Where("id = ?", 241900).Update("name", "updateName333")
}

func TestDelete(t *testing.T) {
	db := drive.Where("id = ?", 241900).Delete(&Funtester{})
	fmt.Println(db.RowsAffected)
}

func TestRollBack(t *testing.T) {
	value := &Funtester{
		Name: "rollback" + randUtils.RandAllString(10),
		Age:  11,
	}
	begin := drive.Begin()
	begin.Create(value)
	fs := &Funtester{
		Name: "fs" + randUtils.RandAllString(10),
		Age:  322231111114,
	}
	err := begin.Create(&fs).Error
	if err != nil {
		fmt.Println("!!!!!!!err:", err)
		begin.Rollback()
	}
	begin.Commit()
	fmt.Println("end")

}

func TestBatchCreate(t *testing.T) {

}
