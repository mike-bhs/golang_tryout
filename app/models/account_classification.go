package models

type AccountClassification struct {
	ComplianceRelationship string `json:"compliance_relationship"`
	RegulatedService       string `json:"regulated_service"`
}

type AccountClassifications []*AccountClassification