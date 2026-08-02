// Computer entity and validation rules for create/update payloads.
package computer

import "github.com/phides-code/updamon-backend/internal/domain"

// PathPrefix is the first URL path segment registered on the gateway (no leading slash).
const PathPrefix = "computers"

// TableName is the DynamoDB table physical name (must match template.yml).
const TableName = "UpdamonComputers"

type Computer struct {
	ID        string `json:"id" dynamodbav:"id"`
	Hostname  string `json:"hostname" dynamodbav:"hostname"`
	IP        string `json:"ip" dynamodbav:"ip"`
	Rating    int    `json:"rating" dynamodbav:"rating"`
	CreatedOn uint64 `json:"createdOn" dynamodbav:"createdOn"`
}

type CreateInput struct {
	Hostname string
	IP       string
	Rating   int
}

type UpdateInput struct {
	ID       string
	Hostname string
	IP       string
	Rating   int
}

func validateHostname(hostname string) error {
	return domain.ValidateRequiredString(hostname, domain.DefaultMinStringLength, domain.DefaultMaxStringLength)
}

func validateIP(ip string) error {
	return domain.ValidateIPv4(ip)
}

func validateRating(rating int) error {
	return domain.ValidateRequiredInt(rating, domain.DefaultMinInt, domain.DefaultMaxInt)
}

func ValidateCreateInput(input CreateInput) error {
	if err := validateHostname(input.Hostname); err != nil {
		return err
	}
	if err := validateIP(input.IP); err != nil {
		return err
	}
	return validateRating(input.Rating)
}

func ValidateUpdateInput(input UpdateInput) error {
	if err := domain.ValidateID(input.ID); err != nil {
		return err
	}
	if err := validateHostname(input.Hostname); err != nil {
		return err
	}
	if err := validateIP(input.IP); err != nil {
		return err
	}
	return validateRating(input.Rating)
}
