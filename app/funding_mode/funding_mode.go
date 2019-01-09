package funding_mode

import (
	"fmt"
)

const PROHIBITED = "prohibited"
const COLLECTIONS = "collections"
const RECEIPTS = "receipts"

const FROM_CLIENT = "from_client"
const OBO_CLIENT = "obo_client"
const OBO_CLIENTS_CUSTOMER = "obo_clients_customer"

type AccountClassification struct {
	regulatedService       string
	complianceRelationship string
}

func (a AccountClassification) IsClient() bool {
	return a.complianceRelationship == "client"
}

func (a AccountClassification) IsNonClient() bool {
	return a.complianceRelationship == "non-client"
}

func (a AccountClassification) IsRegulated() bool {
	return a.regulatedService == "regulated"
}

func (a AccountClassification) IsUnregulated() bool {
	return a.regulatedService == "unregulated"
}

type Sender struct {
	classification string
}

func (s Sender) IsAccountHolder() bool {
	return s.classification == "account_holder"
}

func (s Sender) IsNotAccountHolder() bool {
	return s.classification == "not_account_holder"
}

func ClassifyWithoutHouseAccount(sender Sender, account AccountClassification) (string, string) {
	fundingType := calculateFundingType(sender, account)

	if fundingType == PROHIBITED {
		return fundingType, ""
	}

	if fundingType == COLLECTIONS {
		return fundingType, fmt.Sprintf("%s_%s", fundingType, OBO_CLIENT)
	}

	switch fundingType {
	case RECEIPTS:
		if account.IsClient() {
			return fundingType, fmt.Sprintf("%s_%s", fundingType, FROM_CLIENT)
		} else {
			return fundingType, fmt.Sprintf("%s_%s", fundingType, OBO_CLIENT)
		}
	}

	return "", ""
}

func calculateFundingType(sender Sender, account AccountClassification) string {
	switch {
	case account.IsNonClient():
		return PROHIBITED
	case account.IsRegulated() && account.IsClient() && sender.IsNotAccountHolder():
		return PROHIBITED
	case sender.IsNotAccountHolder():
		return COLLECTIONS
	case sender.IsAccountHolder():
		return RECEIPTS
	}

	return ""
}

func ClassifyWithHouseAccount(sender Sender, aClassification, haClassification AccountClassification) (string, string) {
	return "", ""
}

// func calculateFundingType(sender Sender, aClassification, haClassification AccountClassification) string {
// }
