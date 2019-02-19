package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func destroyTables(db *gorm.DB) {
	errors := db.DropTableIfExists(&Exception{}, &Comment{}, &FormFile{}, &StatusChange{}).GetErrors()

	for _, err := range errors {
		fmt.Printf("%s", err)
	}
}

func createTables(db *gorm.DB) {
	errors := db.CreateTable(&Exception{}, &Comment{}, &FormFile{}, &StatusChange{}).GetErrors()

	for _, err := range errors {
		fmt.Printf("%s", err)
	}
}

func getDB() *gorm.DB {
	dbConfig := parseDBConfig(*configFile)

	if dbConfig.DBType == "mysql" {
		// If you don't pass parseTime=True here for MySQL DBs, time.Times won't work properly
		dbConfig.DBConnectionString += "?charset=utf8&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(dbConfig.DBType, dbConfig.DBConnectionString)
	if err != nil {
		log.Fatalln("Error: could not connect to database.")
	}

	// gormDebugMode is a package-scope variable set in the command-line parsing
	if *gormDebugMode == true {
		return db.Set("gorm:auto_preload", true).Debug()
	}
	return db.Set("gorm.auto_preload", true)
}

func createNoodlingData(db *gorm.DB) {
	nowTime := time.Now()
	aDay, _ := time.ParseDuration("24h")
	aYear, _ := time.ParseDuration("8760h")

	nowPlusADayTime := nowTime.Add(aDay)
	nowPlusAYearTime := nowTime.Add(aYear)

	exception := Exception{Username: "uccaiki", SubmittedDate: &nowTime, StartDate: &nowPlusADayTime, EndDate: &nowPlusAYearTime, Service: "legion", ExceptionType: "quota", ExceptionDetail: "scratch:1TB"}
	db.NewRecord(exception)
	db.Create(&exception)
	exception2 := Exception{Username: "uccaiki", SubmittedDate: &nowTime, Service: "legion", ExceptionType: "quota", ExceptionDetail: "home:500MB"}
	db.NewRecord(exception2)
	db.Create(&exception2)
	exception3 := Exception{Username: "ccspapp", SubmittedDate: &nowTime, Service: "grace", ExceptionType: "queue", ExceptionDetail: "crag7day"}
	db.NewRecord(exception3)
	db.Create(&exception3)
}

func createDB() {
	db := getDB()
	defer db.Close()
	createTables(db)
}

func destroyDB() {
	db := getDB()
	defer db.Close()
	destroyTables(db)
}

func makeNoodles() {
	db := getDB()
	defer db.Close()
	createNoodlingData(db)
}
