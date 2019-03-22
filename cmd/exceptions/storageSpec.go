package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Note: this is very loose: treats lower and upper-case SI prefixes as equivalent
//       and treats b and B both as bytes
func storageSpecToUint64(text string) (uint64, error) {
	re := regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?) ?([kmgtKMGT]?)(i?)(?:[bB])$`)

	matches := re.FindStringSubmatch(text)

	if matches == nil {
		return 0, errors.New("conversion regex did not match")
	}

	inputNum := matches[1]
	siPrefix := strings.ToLower(matches[2])
	useBinaryPrefixes := (len(matches[3]) == 1)

	siPrefixMultipliers := map[string]uint64{
		"":  1,
		"k": 1e3,
		"m": 1e6,
		"g": 1e9,
		"t": 1e12,
	}
	biPrefixMultipliers := map[string]uint64{
		"":  1,
		"k": 1 << 10,
		"m": 1 << 20,
		"g": 1 << 30,
		"t": 1 << 40,
	}

	var floatResult float64
	var err error

	floatResult, err = strconv.ParseFloat(inputNum, 64)
	if err != nil {
		return 0, nil
	}

	if useBinaryPrefixes == true {
		floatResult *= float64(biPrefixMultipliers[siPrefix])
	} else {
		floatResult *= float64(siPrefixMultipliers[siPrefix])
	}
	return uint64(floatResult), nil
}

func ValidateStorageSpec(val interface{}) error {

	re := regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?) ?([kmgtKMGT]?)(i?)(?:[bB])$`)

	matches := re.FindStringSubmatch(val.(string))

	if matches == nil {
		return errors.New("conversion regex did not match")
	}
	return nil
}

func TidyStorageSpec(spec string) (string, error) {
	re := regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?) ?([kmgtKMGT]?)(i?)(?:[bB])$`)

	matches := re.FindStringSubmatch(spec)

	if matches == nil {
		return "", errors.New("conversion regex did not match")
	}

	newSpec := ""
	newSpec += matches[1]

	switch strings.ToLower(matches[2]) {
	case "k":
		newSpec += "k"
	case "m":
		newSpec += "M"
	case "g":
		newSpec += "G"
	case "t":
		newSpec += "T"
	}
	newSpec += matches[3] + "B"

	return newSpec, nil
}
