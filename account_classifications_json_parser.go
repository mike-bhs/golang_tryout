package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mike-bhs/golang_tryout/app/models"
	"github.com/tidwall/gjson"
)

func ConvertAccountClassificationJsonToModel(j string) (models.AccountClassifications, error) {
	var arr models.AccountClassifications

	if !gjson.Valid(j) {
		message := fmt.Sprintf("JSON is invalid. JSON: %s", j)
		err := errors.New(message)

		return arr, err
	}

	accountClassification := gjson.Get(j, "account_classification").Raw
	err := json.Unmarshal([]byte(accountClassification), &arr)

	if err != nil {
		return arr, err
	}

	return arr, nil
}

