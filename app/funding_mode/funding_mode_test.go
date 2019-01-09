package funding_mode

import (
	"testing"
)

var accountData = [][]string{
	// regulated_service account sender funding_type funding_mode

	// receipts
	[]string{"unregulated", "client", "account_holder", "receipts", "receipts_from_client"},
	[]string{"regulated", "client", "account_holder", "receipts", "receipts_from_client"},

	// collections
	[]string{"unregulated", "client", "not_account_holder", "collections", "collections_obo_client"},

	// prohibited
	[]string{"unregulated", "non-client", "account_holder", "prohibited", ""},
	[]string{"regulated", "non-client", "account_holder", "prohibited", ""},
	[]string{"unregulated", "non-client", "not_account_holder", "prohibited", ""},
	[]string{"regulated", "client", "not_account_holder", "prohibited", ""},
	[]string{"regulated", "non-client", "not_account_holder", "prohibited", ""},
}

var houseAccountData = [][]string{
	// regulated_service account house_account sender funding_type funding_mode

	// receipts
	[]string{"unregulated", "client", "client", "account_holder", "receipts", "receipts_from_client"},
	[]string{"unregulated", "client", "non-client", "account_holder", "receipts", "receipts_from_client"},
	[]string{"regulated", "client", "client", "account_holder", "receipts", "receipts_from_client"},
	[]string{"regulated", "non-client", "client", "account_holder", "receipts", "receipts_obo_client"},

	// collections
	[]string{"unregulated", "non-client", "client", "account_holder", "collections", "collections_obo_client"},
	[]string{"unregulated", "client", "client", "not_account_holder", "collections", "collections_obo_client"},
	[]string{"unregulated", "client", "non-client", "not_account_holder", "collections", "collections_obo_client"},
	[]string{"unregulated", "non-client", "client", "not_account_holder", "collections", "collections_obo_client"},
	[]string{"regulated", "non-client", "client", "not_account_holder", "collections", "collections_obo_clients_customer"},

	// prohibited
	[]string{"regulated", "non-client", "non-client", "account_holder", "prohibited", ""},
	[]string{"unregulated", "non-client", "non-client", "account_holder", "prohibited", ""},
	[]string{"regulated", "client", "non-client", "account_holder", "prohibited", ""},
	[]string{"regulated", "non-client", "non-client", "not_account_holder", "prohibited", ""},
	[]string{"unregulated", "non-client", "non-client", "not_account_holder", "prohibited", ""},
	[]string{"regulated", "client", "client", "not_account_holder", "prohibited", ""},
	[]string{"regulated", "client", "non-client", "not_account_holder", "prohibited", ""},
}

func TestClassifyWithoutHouseAccount(t *testing.T) {
	for _, row := range accountData {
		aClassification := AccountClassification{regulatedService: row[0], complianceRelationship: row[1]}
		sender := Sender{classification: row[2]}

		expectedFundingType := row[3]
		expectedFundingMode := row[4]

		actualFundingType, actualFundingMode := ClassifyWithoutHouseAccount(sender, aClassification)

		if expectedFundingType != actualFundingType {
			t.Error("FundingType was calculated wrong for ClassifyWithoutHouseAccount")
		}

		if expectedFundingMode != actualFundingMode {
			t.Error("FundingMode was calculated wrong for ClassifyWithoutHouseAccount")
		}
	}
}

func TestClassifyWithHouseAccount(t *testing.T) {
	for _, row := range houseAccountData {
		aClassification := AccountClassification{regulatedService: row[0], complianceRelationship: row[1]}
		haClassification := AccountClassification{regulatedService: row[0], complianceRelationship: row[2]}
		sender := Sender{classification: row[3]}

		expectedFundingType := row[4]
		expectedFundingMode := row[5]

		actualFundingType, actualFundingMode := ClassifyWithHouseAccount(sender, aClassification, haClassification)

		if expectedFundingType != actualFundingType {
			t.Error("FundingType was calculated wrong for ClassifyWithHouseAccount")
		}

		if expectedFundingMode != actualFundingMode {
			t.Error("FundingMode was calculated wrong for ClassifyWithHouseAccount")
		}
	}
}
