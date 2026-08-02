// Sitrep entity and validation rules for create payloads.
package sitrep

import "github.com/phides-code/updamon-backend/internal/domain"

// PathPrefix is the first URL path segment registered on the gateway (no leading slash).
const PathPrefix = "sitreps"

// TableName is the DynamoDB table physical name (must match template.yml).
const TableName = "UpdamonSitreps"

type Sitrep struct {
	ID        string `json:"id" dynamodbav:"id"`
	Hostname  string `json:"hostname" dynamodbav:"hostname"`
	CreatedOn uint64 `json:"createdOn" dynamodbav:"createdOn"`
}

type CreateInput struct {
	Hostname string
}

func validateHostname(hostname string) error {
	return domain.ValidateRequiredString(hostname, domain.DefaultMinStringLength, domain.DefaultMaxStringLength)
}

func ValidateCreateInput(input CreateInput) error {
	return validateHostname(input.Hostname)
}
