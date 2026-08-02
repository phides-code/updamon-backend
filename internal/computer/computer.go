// Computer entity and validation rules for create/update payloads.
package computer

import "github.com/phides-code/updamon-backend/internal/domain"

// PathPrefix is the first URL path segment registered on the gateway (no leading slash).
const PathPrefix = "computers"

// TableName is the DynamoDB table physical name (must match template.yml).
const TableName = "UpdamonComputers"

type Computer struct {
	ID         string `json:"id" dynamodbav:"id"`
	Descriptor string `json:"descriptor" dynamodbav:"descriptor"`
	Rating     int    `json:"rating" dynamodbav:"rating"`
	CreatedOn  uint64 `json:"createdOn" dynamodbav:"createdOn"`
}

type CreateInput struct {
	Descriptor string
	Rating     int
}

type UpdateInput struct {
	ID         string
	Descriptor string
	Rating     int
}

func validateDescriptor(descriptor string) error {
	return domain.ValidateRequiredString(descriptor, domain.DefaultMinStringLength, domain.DefaultMaxStringLength)
}

func validateRating(rating int) error {
	return domain.ValidateRequiredInt(rating, domain.DefaultMinInt, domain.DefaultMaxInt)
}

func ValidateCreateInput(input CreateInput) error {
	if err := validateDescriptor(input.Descriptor); err != nil {
		return err
	}
	return validateRating(input.Rating)
}

func ValidateUpdateInput(input UpdateInput) error {
	if err := domain.ValidateID(input.ID); err != nil {
		return err
	}
	if err := validateDescriptor(input.Descriptor); err != nil {
		return err
	}
	return validateRating(input.Rating)
}
