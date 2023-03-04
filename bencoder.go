package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

func Decode(data []byte) (interface{}, error) {

	if string(data[0]) != "d" {
		return nil, fmt.Errorf("corrupted file")
	}

	response, _, err := handleDict(data, 0)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func decodeKeys(data []byte, from int) (interface{}, int, error) {
	switch string(data[from]) {
	case "i":
		return handleInteger(data, from + 1)
	case "l":
		return handleList(data, from)
	case "d":
		return handleDict(data, from)
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return handleString(data, from)
	default:
		return nil, 0, fmt.Errorf("unknown character")
	}
}

func processKey(data []byte, from int) (interface{}, int, error) {
	val, len, err := handleString(data, from)

	if err != nil {
		return nil, 0, fmt.Errorf("error while parsing %s", err.Error())
	}

	return val, len, nil
}

func processValue(data []byte, from int) (interface{}, int, error) {
	val, len, err := decodeKeys(data, from)

	if err != nil {
		return nil, 0, fmt.Errorf("error while parsing %s", err.Error())
	}

	return val, len, nil
}

func handleDict(data []byte, from int) (interface{}, int, error) {

	var (
		idx = 0
		response = make(map[string]interface{})
	)

	for idx = from +1; idx < len(data); {

		if string(data[idx]) == "e" {
			return response, idx + 1, nil
		}

		key, len, err := processKey(data, idx)

		idx = len

		if err != nil {
			return nil, 0, fmt.Errorf("could not parse data: %s", err.Error())
		}

		value, len, err := processValue(data, idx)

		if err != nil {
			return nil, 0, fmt.Errorf("could not parse data: %s", err.Error())
		}

		idx = len

		keyStr := key.(string)

		if keyStr == "creation date" {
			val := time.Unix(int64(value.(int)), 0).Format(time.UnixDate)
			value = val
		}

		response[keyStr] = value
	}

	return response, 0, nil
}

func handleList(data []byte, from int) (interface{}, int, error) {

	var (
		response = make([]interface{}, 0)
		idx = 0
	)

	for idx = from + 1; string(data[idx]) != "e"; {
		val, len, err := processValue(data, idx)

		if err != nil {
			return nil, 0, fmt.Errorf("error parsing list %s", err.Error())
		}

		idx = len

		response = append(response, val)
	}

	return response, idx + 1, nil
}

func handleInteger(data []byte, from int) (interface{}, int, error) {
	var numStr string
	var idx int

	for idx = from; string(data[idx]) != "e"; idx++ {
		numStr += string(data[idx])
	}

	idx += 1
	num, err := strconv.Atoi(numStr)

	if err != nil {
		fmt.Println("could not convert the number")
	}

	return num, idx, nil
}

func handleString(data []byte, from int) (interface{}, int, error) {
	var numStr string
	var idx int

	for idx = from; string(data[idx]) != ":"; idx++ {
		numStr += string(data[idx])
	}

	idx += 1
	num, err := strconv.Atoi(numStr)

	if err != nil {
		fmt.Println("could not convert the number")
	}

	if num > 1000 {
		return hex.EncodeToString(data[idx : idx+num]), idx+num, nil 
	}

	return string(data[idx : idx+num]), idx + num, nil
}
