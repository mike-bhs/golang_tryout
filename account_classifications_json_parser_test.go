package main

import (
	"github.com/mike-bhs/golang_tryout/app/models"
	. "github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertAccountClassificationJsonToModel_arrayOfAccountClassifications(t *testing.T) {
	jsonString := `{
        "account_classification": [
          { "compliance_relationship": "non-client", "regulated_service": "regulated" },
          { "compliance_relationship": "client", "regulated_service": "unregulated" }
        ]
      }`

	var examples = models.AccountClassifications {
		&models.AccountClassification{ComplianceRelationship: "non-client", RegulatedService: "regulated"},
		&models.AccountClassification{ComplianceRelationship: "client", RegulatedService: "unregulated"},
	}

	res, err := ConvertAccountClassificationJsonToModel(jsonString)

	for i, example := range examples {
		Equal(t, example, res[i])
	}

	Equal(t, nil, err)
}


