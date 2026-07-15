// Computer entity and validation rules for create/update payloads.
package computer

import (
	"net"

	"github.com/phides-code/go-multi-api/internal/domain"
)

const (
	MinHostnameLength = 1
	MaxHostnameLength = 100
)

type Computer struct {
	ID        string `json:"id" dynamodbav:"id"`
	Hostname  string `json:"hostname" dynamodbav:"hostname"`
	IP        string `json:"ip" dynamodbav:"ip"`
	CreatedOn uint64 `json:"createdOn" dynamodbav:"createdOn"`
}

type CreateInput struct {
	Hostname string
	IP       string
}

type UpdateInput struct {
	ID       string
	Hostname string
	IP       string
}

func validateHostname(hostname string) error {
	return domain.ValidateRequiredString(hostname, MinHostnameLength, MaxHostnameLength)
}

func validateIP(ip string) error {
	parsed := net.ParseIP(ip)
	if parsed == nil || parsed.To4() == nil {
		return domain.ErrValidationFailed
	}
	return nil
}

func ValidateCreateInput(input CreateInput) error {
	if err := validateHostname(input.Hostname); err != nil {
		return err
	}
	return validateIP(input.IP)
}

func ValidateUpdateInput(input UpdateInput) error {
	if err := domain.ValidateID(input.ID); err != nil {
		return err
	}
	if err := validateHostname(input.Hostname); err != nil {
		return err
	}
	return validateIP(input.IP)
}
