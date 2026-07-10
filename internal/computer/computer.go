// Computer entity and validation rules for create/update payloads.
package computer

import "github.com/phides-code/go-multi-api/internal/domain"

const (
	MinHostnameLength = 1
	MaxHostnameLength = 100
)

type Computer struct {
	ID        string `json:"id" dynamodbav:"id"`
	Hostname  string `json:"hostname" dynamodbav:"hostname"`
	CreatedOn uint64 `json:"createdOn" dynamodbav:"createdOn"`
}

type CreateInput struct {
	Hostname string
}

type UpdateInput struct {
	ID       string
	Hostname string
}

func validateHostname(hostname string) error {
	return domain.ValidateRequiredString(hostname, MinHostnameLength, MaxHostnameLength)
}

func ValidateCreateInput(input CreateInput) error {
	return validateHostname(input.Hostname)
}

func ValidateUpdateInput(input UpdateInput) error {
	if err := domain.ValidateID(input.ID); err != nil {
		return err
	}
	return validateHostname(input.Hostname)
}
