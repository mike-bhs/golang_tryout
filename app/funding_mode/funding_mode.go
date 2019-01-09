package funding_mode

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

func ClassifyWithoutHouseAccount(sender Sender, aClassification AccountClassification) (string, string) {
	return "", ""
}

func ClassifyWithHouseAccount(sender Sender, aClassification, haClassification AccountClassification) (string, string) {
	return "", ""
}
