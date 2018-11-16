package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func createDB(db *gorm.DB) {
	errors1 := db.Debug().AutoMigrate(&Exception{}).GetErrors()
	errors2 := db.Debug().AutoMigrate(&Comment{}).GetErrors()
	errors3 := db.Debug().AutoMigrate(&FormFile{}).GetErrors()

	for _, err := range errors1 {
		fmt.Printf("%s", err)
	}
	for _, err := range errors2 {
		fmt.Printf("%s", err)
	}
	for _, err := range errors3 {
		fmt.Printf("%s", err)
	}
}

func getDB() *gorm.DB {
	// If you don't pass parseTime=True here, time.Times won't work properly
	db, err := gorm.Open("mysql", "test:blah@tcp(127.0.0.1)/test_exceptions?charset=utf8&parseTime=True&loc=Local")

	//db, err := gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		fmt.Println("Error: could not connect to database.")
		panic(err)
	}

	createDB(db)

	return db
}

func createNoodlingData(db *gorm.DB) {
	nowTime := time.Now()
	exception := Exception{Username: "uccaiki", SubmittedDate: &nowTime, Service: "legion", ExceptionType: "quota", ExceptionDetail: "scratch:1TB"}
	db.NewRecord(exception)
	db.Create(&exception)
	exception2 := Exception{Username: "uccaiki", SubmittedDate: &nowTime, Service: "legion", ExceptionType: "quota", ExceptionDetail: "home:500MB"}
	db.NewRecord(exception2)
	db.Create(&exception2)
	exception3 := Exception{Username: "ccspapp", SubmittedDate: &nowTime, Service: "grace", ExceptionType: "queue", ExceptionDetail: "crag7day"}
	db.NewRecord(exception3)
	db.Create(&exception3)
}

func dbsetup() {
	db := getDB()
	createDB(db)
	createNoodlingData(db)
}
