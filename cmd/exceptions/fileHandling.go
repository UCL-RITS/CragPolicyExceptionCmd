package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

func attach(id uint, filename string) (uint, error) {
	db := getDB()
	exception := &Exception{}
	db.First(&exception, id)
	if exception.ID == 0 {
		fmt.Println("No record of that exception.")
		return 0, errors.New("No record of that exception.")
	}

	basename := filepath.Base(filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	formFile := &FormFile{}
	formFile.FileContents = b
	formFile.FileName = basename
	formFile.ExceptionID = id
	db.Save(&formFile)
	return formFile.ID, nil
}

func getFilesForException(id uint) ([]FormFile, error) {
	db := getDB()
	exception := &Exception{}
	db.First(&exception, id)
	if exception.ID == 0 {
		return nil, errors.New("No record of that exception.")
	}

	formFiles := &[]FormFile{}
	db.Model(&exception).Related(&formFiles)

	return *formFiles, nil
}

func listFilesForException(id uint) {
	files, err := getFilesForException(id)

	if err != nil {
		fmt.Printf("Could not get files for exception %d: %s\n", id, err)
		return
	}

	if len(files) == 0 {
		fmt.Printf("No files for exception %d.\n", id)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Created On", "Filename", "Size"})
	table.SetBorder(false)

	for _, file := range files {
		table.Append([]string{fmt.Sprintf("%d", file.ID),
			stringFromDate(&file.CreatedAt),
			file.FileName,
			fmt.Sprintf("%d", len(file.FileContents)),
		})
	}
	table.Render()
	return
}

// This is for when you want to download a single file and have
//  referred to it directly by ID
func downloadOneFile(fileID uint) {
	db := getDB()
	file := &FormFile{}
	db.First(&file, fileID)
	if file.ID == 0 {
		fmt.Println("No record of that file.")
		return
	}

	targetFilename := file.FileName
	// This is a while loop in any other language
	for fileExists(targetFilename) {
		targetFilename += "_"
	}
	err := writeOutFile(*file, targetFilename)
	if err != nil {
		fmt.Println("Could not write out file to ", targetFilename, ": ", err)
		return
	}
	fmt.Printf("Wrote out file %d to: %s\n", fileID, targetFilename)
	return
}

// This is for when you want all the files for an exception and
//  have referred to the *exception* by ID, not the file
func downloadFilesForException(exceptionID uint) {
	files, err := getFilesForException(exceptionID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		targetFilename := file.FileName
		// This is a while loop in any other language
		for fileExists(targetFilename) {
			targetFilename += "_"
		}
		err := writeOutFile(file, targetFilename)
		if err != nil {
			fmt.Println("Could not write out file to ", targetFilename, ": ", err)
			return
		}
		fmt.Printf("Wrote out file %d to: %s\n", file.ID, targetFilename)
	}
	return
}

func fileExists(filename string) bool {
	// For checking you're not overwriting a file first
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}

func writeOutFile(file FormFile, targetFilename string) error {
	//ioutil.WriteFile(filename string, data []byte, perm os.FileMode) error
	return ioutil.WriteFile(targetFilename, file.FileContents, os.FileMode(0500))
}
