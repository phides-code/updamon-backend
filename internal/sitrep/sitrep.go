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
	AptLog    string `json:"aptlog" dynamodbav:"aptlog"`
	Last      string `json:"last" dynamodbav:"last"`
	WAP       string `json:"wap" dynamodbav:"wap"`
	Free      string `json:"free" dynamodbav:"free"`
	DF        string `json:"df" dynamodbav:"df"`
	Who       string `json:"who" dynamodbav:"who"`
	Tailscale string `json:"tailscale" dynamodbav:"tailscale"`
	Bluetooth string `json:"bluetooth" dynamodbav:"bluetooth"`
	CreatedOn uint64 `json:"createdOn" dynamodbav:"createdOn"`
}

type CreateInput struct {
	Hostname  string
	AptLog    string
	Last      string
	WAP       string
	Free      string
	DF        string
	Who       string
	Tailscale string
	Bluetooth string
}

func validateShortString(s string) error {
	return domain.ValidateRequiredString(s, domain.DefaultMinStringLength, domain.DefaultMaxStringLength)
}

func validateWAP(wap string) error {
	return domain.ValidateMAC(wap)
}

func validateMediumString(s string) error {
	return domain.ValidateRequiredString(s, domain.DefaultMinStringLength, domain.DefaultMaxMediumStringLength)
}

func validateLongString(s string) error {
	return domain.ValidateRequiredString(s, domain.DefaultMinStringLength, domain.DefaultMaxLongStringLength)
}

func ValidateCreateInput(input CreateInput) error {
	if err := validateShortString(input.Hostname); err != nil {
		return err
	}
	if err := validateLongString(input.AptLog); err != nil {
		return err
	}
	if err := validateLongString(input.Last); err != nil {
		return err
	}
	if err := validateWAP(input.WAP); err != nil {
		return err
	}
	if err := validateMediumString(input.Free); err != nil {
		return err
	}
	if err := validateMediumString(input.DF); err != nil {
		return err
	}
	if err := validateMediumString(input.Who); err != nil {
		return err
	}
	if err := validateShortString(input.Tailscale); err != nil {
		return err
	}
	return validateShortString(input.Bluetooth)
}
