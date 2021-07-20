package handler

import (
	"encoding/json"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var (
	res []ErrorResponse
)

func GetJSONErrorResponse(rep []*ErrorResponse) (string, error) {

	for k := range rep {
		res = append(res, *rep[k])
	}

	errorResponse, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(errorResponse), nil
}
