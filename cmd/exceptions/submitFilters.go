package main

import (
	"errors"
	"fmt"
	"strings"
)

// These have global scope so the CLI can tell people
var (
	validServices       = []string{"myriad", "legion", "grace", "aristotle", "thomas", "michael", "kathleen", "none"}
	validExceptionTypes = []string{"quota", "queue", "access", "special", "sharedspace"}

	validServicesString       = strings.Join(validServices, ", ")
	validExceptionTypesString = strings.Join(validExceptionTypes, ", ")
)

func filterSubmittedUsername(name string) (string, error) {
	name = strings.ToLower(name)
	allowedChars := "abcdefghijklmnopqrstuvxyz1234567890"
	valid := true
	errorSlice := []string{}

	if len(name) != 7 {
		valid = false
		errorSlice = append(errorSlice, fmt.Sprintf("Username is incorrect length: %d, not 7", len(name)))
	}

	for _, v := range name {
		if !strings.ContainsRune(allowedChars, v) {
			valid = false
			errorSlice = append(errorSlice, fmt.Sprintf("Invalid character in username: %q", v))
		}
	}

	var returnError error
	if len(errorSlice) != 0 {
		returnError = errors.New(strings.Join(errorSlice, "; "))
	} else {
		returnError = nil
	}

	if !valid {
		// Blank the name var on error to avoid accidental usage of invalid name
		name = ""
	}
	return name, returnError
}

func filterSubmittedService(service string) (string, error) {
	service = strings.ToLower(service)
	valid := false

	for _, v := range validServices {
		if service == v {
			valid = true
			break
		}
	}

	var returnError error
	returnError = nil
	if !valid {
		// Blank the service var on error to avoid accidental usage of invalid service
		service = ""
		errorMsg := fmt.Sprintf("Invalid service, must be: %s", validServicesString)
		returnError = errors.New(errorMsg)
	}
	return service, returnError
}

func filterSubmittedExceptionType(exceptionType string) (string, error) {
	exceptionType = strings.ToLower(exceptionType)
	valid := false

	for _, v := range validExceptionTypes {
		if exceptionType == v {
			valid = true
		}
	}

	var returnError error
	returnError = nil
	if !valid {
		// Blank the service var on error to avoid accidental usage of invalid service
		exceptionType = ""
		errorMsg := fmt.Sprintf("Invalid exception type, must be: %s", validExceptionTypesString)
		returnError = errors.New(errorMsg)
	}
	return exceptionType, returnError
}
