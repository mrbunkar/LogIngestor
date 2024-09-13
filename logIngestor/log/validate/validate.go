package validate

import (
	"fmt"
	"logIngestor/logIngestor/log/logtype"

	"go.mongodb.org/mongo-driver/bson"
)

type Validator struct {
}

func (v *Validator) Validate(log *logtype.Log) error {
	return nil
}

func (v *Validator) ParseTheFilters(filter bson.M) error {

	for f := range filter {
		fmt.Println(f)
	}
	return nil
}
