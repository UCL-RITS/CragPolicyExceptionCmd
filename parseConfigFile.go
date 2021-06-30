package main

import (
	"bytes"
	"encoding/json" // While I recognise that JSON is not ideal, it was either this or XML without adding another dependency
	"io/ioutil"
	"log"
	"os"
)

// Example Connection Strings:
//  for sqlite3: just a filename
//  for MySQL: "test:blah@tcp(127.0.0.1)/test_exceptions"
// We append our own parameters for MySQL so if you use any here with ? it won't work -- this hasn't been a problem so far
type DBConfig struct {
	// Yes this is barebones, but trying to work out how to handle the connection parameters in a more cunning way was making my head hurt
	DBType             string `json:"db_type"` // to pass to gorm.Open, if you want to use something other than "sqlite3" or "mysql" you'll have to add the drivers over in db.go
	DBConnectionString string `json:"db_connection_string"`
}

func getExampleConfigText() string {
	return `{
		"db_type": "sqlite3",
		"db_connection_string": "./exceptions.db"
}`
}

// genpw servexcep_rw 19072018

func parseDBConfig(filename string) *DBConfig {
	configFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Fatal error: could not open config file: ", err)
	}
	defer configFile.Close()

	var buffer []byte

	buffer, err = ioutil.ReadAll(configFile)

	if err != nil {
		log.Println("Fatal error: could not read config file: " + filename)
		panic(err)
	}

	dbConfig := &DBConfig{}

	err = json.Unmarshal(buffer, dbConfig)
	if err != nil {
		log.Println("Fatal error: could not parse config file: " + filename)
		switch t := err.(type) {
		case *json.SyntaxError:
			log.Printf("Syntax error on line %d\n", 1+bytes.Count(buffer[0:t.Offset], []byte("\n")))
		}
		panic(err)
	}

	return dbConfig
}
