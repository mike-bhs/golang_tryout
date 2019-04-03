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
	fundingType := calculateFundingTypeWithoutHouseAccount(sender, account)

	if fundingType == PROHIBITED {
		return fundingType, ""
	}

	if fundingType == COLLECTIONS {
		return fundingType, fmt.Sprintf("%s_%s", fundingType, OBO_CLIENT)
	}

	if fundingType == RECEIPTS && account.IsClient() {
		return fundingType, fmt.Sprintf("%s_%s", fundingType, FROM_CLIENT)
	}

	if fundingType == RECEIPTS && !account.IsClient() {
		return fundingType, fmt.Sprintf("%s_%s", fundingType, OBO_CLIENT)
	}

	return "", ""
}

func calculateFundingTypeWithoutHouseAccount(sender Sender, account AccountClassification) string {
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

func ClassifyWithHouseAccount(sender Sender, account, houseAccount AccountClassification) (string, string) {
	fundingType := calculateFundingTypeWithHouseAccount(sender, account, houseAccount)

	if fundingType == PROHIBITED {
		return fundingType, ""
	}

	fundingMode := calculateFundingMode(fundingType, sender, account, houseAccount)

	return fundingType, fundingMode
}

func calculateFundingMode(fundingType string, sender Sender, account, houseAccount AccountClassification) string {
	switch fundingType {
	case RECEIPTS:
		if account.IsClient() {
			return fmt.Sprintf("%s_%s", fundingType, FROM_CLIENT)
		} else {
			return fmt.Sprintf("%s_%s", fundingType, OBO_CLIENT)
		}

	case COLLECTIONS:
		if isNestedCollections(account, houseAccount, sender) {
			return fmt.Sprintf("%s_%s", fundingType, OBO_CLIENTS_CUSTOMER)
		} else {
			return fmt.Sprintf("%s_%s", fundingType, OBO_CLIENT)
		}
	}

	return ""
}

func calculateFundingTypeWithHouseAccount(sender Sender, account, houseAccount AccountClassification) string {
	switch {
	case isNoComplianceRelationship(account, houseAccount):
		return PROHIBITED

	case isNestedPaymentWithCollections(account, houseAccount, sender):
		return PROHIBITED

	case isRegulatedAffiliateReceipts(account, houseAccount, sender):
		return PROHIBITED

	case sender.IsNotAccountHolder():
		return COLLECTIONS

	case isCorporateCollections(account, houseAccount, sender):
		return COLLECTIONS

	case sender.IsAccountHolder():
		return RECEIPTS
	}

	return ""
}

func isNoComplianceRelationship(account, houseAccount AccountClassification) bool {
	return account.IsNonClient() && houseAccount.IsNonClient()
}

func isNestedPaymentWithCollections(account, houseAccount AccountClassification, sender Sender) bool {
	return houseAccount.IsRegulated() && account.IsClient() && sender.IsNotAccountHolder()
}

func isRegulatedAffiliateReceipts(account, houseAccount AccountClassification, sender Sender) bool {
	return sender.IsAccountHolder() && account.IsClient() && houseAccount.IsNonClient() && houseAccount.IsRegulated()
}

func isCorporateCollections(account, houseAccount AccountClassification, sender Sender) bool {
	return sender.IsAccountHolder() && account.IsNonClient() && houseAccount.IsClient() && houseAccount.IsUnregulated()
}

func isNestedCollections(account, houseAccount AccountClassification, sender Sender) bool {
	return account.IsNonClient() && houseAccount.IsClient() && houseAccount.IsRegulated() && sender.IsNotAccountHolder()
}
