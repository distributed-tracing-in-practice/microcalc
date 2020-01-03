package api

import (
	"encoding/json"
	"io"
)

type CalcRequest struct {
	Method   string `json:"method"`
	Operands []int  `json:"operands"`
}

func ParseCalcRequest(body io.Reader) (CalcRequest, error) {
	var parsedRequest CalcRequest

	err := json.NewDecoder(body).Decode(&parsedRequest)
	if err != nil {
		return parsedRequest, err
	}

	return parsedRequest, nil
}
