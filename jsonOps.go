package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func dumpAllAsJson() {
	var allExceptions []Exception
	db := getDB()
	defer db.Close()
	db.Preload("Comments").Preload("FormFiles").Preload("StatusChanges").Find(&allExceptions)
	jsonBytes, err := json.MarshalIndent(allExceptions, "", " ")

	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
	return
}

func importAllAsJson() {
	var exceptionsImport []Exception
	buffer, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(buffer, &exceptionsImport)

	if err != nil {
		panic(err)
	}

	db := getDB()
	defer db.Close()
	importTransaction := db.Begin()

	for _, e := range exceptionsImport {
		errs := db.Save(&e).GetErrors()
		if len(errs) != 0 {
			log.Print(errs)
			importTransaction.Rollback()
			break
		}
	}
	importTransaction.Commit()

	return
}
