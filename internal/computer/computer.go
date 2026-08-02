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
	OS        string `json:"os" dynamodbav:"os"`
	Kernel    string `json:"kernel" dynamodbav:"kernel"`
	Model     string `json:"model" dynamodbav:"model"`
	RAM       string `json:"ram" dynamodbav:"ram"`
	CPU       string `json:"cpu" dynamodbav:"cpu"`
	Storage   string `json:"storage" dynamodbav:"storage"`
	CreatedOn uint64 `json:"createdOn" dynamodbav:"createdOn"`
}

type CreateInput struct {
	Hostname string
	IP       string
	OS       string
	Kernel   string
	Model    string
	RAM      string
	CPU      string
	Storage  string
}

type UpdateInput struct {
	ID       string
	Hostname string
	IP       string
	OS       string
	Kernel   string
	Model    string
	RAM      string
	CPU      string
	Storage  string
}

func validateRequiredString(s string) error {
	return domain.ValidateRequiredString(s, domain.DefaultMinStringLength, domain.DefaultMaxStringLength)
}

func validateIP(ip string) error {
	return domain.ValidateIPv4(ip)
}

func ValidateCreateInput(input CreateInput) error {
	if err := validateRequiredString(input.Hostname); err != nil {
		return err
	}
	if err := validateIP(input.IP); err != nil {
		return err
	}
	if err := validateRequiredString(input.OS); err != nil {
		return err
	}
	if err := validateRequiredString(input.Kernel); err != nil {
		return err
	}
	if err := validateRequiredString(input.Model); err != nil {
		return err
	}
	if err := validateRequiredString(input.RAM); err != nil {
		return err
	}
	if err := validateRequiredString(input.CPU); err != nil {
		return err
	}
	return validateRequiredString(input.Storage)
}

func ValidateUpdateInput(input UpdateInput) error {
	if err := domain.ValidateID(input.ID); err != nil {
		return err
	}
	if err := validateRequiredString(input.Hostname); err != nil {
		return err
	}
	if err := validateIP(input.IP); err != nil {
		return err
	}
	if err := validateRequiredString(input.OS); err != nil {
		return err
	}
	if err := validateRequiredString(input.Kernel); err != nil {
		return err
	}
	if err := validateRequiredString(input.Model); err != nil {
		return err
	}
	if err := validateRequiredString(input.RAM); err != nil {
		return err
	}
	if err := validateRequiredString(input.CPU); err != nil {
		return err
	}
	return validateRequiredString(input.Storage)
}
