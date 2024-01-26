package Requests

import (
	"encoding/json"
	"fmt"
)

type ResponeType int

const (
	ErrorRespone ResponeType = 999
	ValidRespone ResponeType = 200
)

type ResponeInfo struct {
	Type    ResponeType `json:"Type"`
	Respone string      `json:"Data"`
}

func GetResponseInfo(data []byte) (ResponeInfo, error) {
	var response_info ResponeInfo
	err := json.Unmarshal(data, &response_info)
	if err != nil {
		return ResponeInfo{}, fmt.Errorf("error when attempting to encode the response from the server.\nPlease send this info to the developers:\n%s", err)
	}
	return response_info, nil
}
