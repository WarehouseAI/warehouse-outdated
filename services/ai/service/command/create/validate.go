package create

import (
	"encoding/json"
	"fmt"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"
)

type rule func(field m.AiCommandField, fieldName string) error

type validator struct {
	rules []rule
}

func newValidator(rules []rule) validator {
	return validator{
		rules: rules,
	}
}

func (v *validator) validateField(field m.AiCommandField, fieldName string) []error {
	var errors []error

	for _, rule := range v.rules {
		if err := rule(field, fieldName); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func validateRequest(request *CreateCommandRequest) *e.HttpErrorResponse {
	v := newValidator([]rule{
		validateFieldData(),
		validateFieldRequirement(),
		validateFieldType(),
		validatePayloadCompability(request),
		validateSelectionType(),
		validateInputType(),
		validateDefaultRequirement(),
		validateFileDataIsNotDefault(),
	})

	for name, value := range request.Payload {
		var field m.AiCommandField
		fieldJson, err := json.Marshal(value)

		if err != nil {
			return e.NewErrorResponse(e.HttpInternalError, err.Error())
		}

		if err := json.Unmarshal(fieldJson, &field); err != nil {
			return e.NewErrorResponse(e.HttpInternalError, err.Error())
		}

		fmt.Println(name, field)

		if errors := v.validateField(field, name); len(errors) != 0 {
			var messages []string

			for _, err := range errors {
				messages = append(messages, err.Error())
			}

			return e.NewErrorResponseMultiple(e.HttpUnprocessableEntity, messages)
		}
	}

	return nil
}

func validateFieldType() rule {
	return func(field m.AiCommandField, fieldName string) error {
		if field.Type != m.Input && field.Type != m.Selection {
			return fmt.Errorf(`field "%s" has incorrect. Use input/selection in "type" parameter.`, fieldName)
		}

		return nil
	}
}

func validateFieldRequirement() rule {
	return func(field m.AiCommandField, fieldName string) error {
		if field.Requirement != m.Default && field.Requirement != m.Require && field.Requirement != m.Optional {
			return fmt.Errorf(`field "%s" is incorrect. Use default/require/optional in "requirement" parameter.`, fieldName)
		}

		return nil
	}
}

func validateFieldData() rule {
	return func(field m.AiCommandField, fieldName string) error {
		if field.Data != m.Bool && field.Data != m.File && field.Data != m.Number && field.Data != m.Object && field.Data != m.String {
			return fmt.Errorf(`field "%s" is incorrect. At now support string/number/file/bool/object in "data" parameter.`, fieldName)
		}

		return nil
	}
}

func validatePayloadCompability(newCommand *CreateCommandRequest) rule {
	return func(field m.AiCommandField, fieldName string) error {
		if newCommand.PayloadType == m.Json && field.Data == m.File {
			return fmt.Errorf(`field "%s" is incorrect. JSON payload type is not support files, use FormData instead.`, fieldName)
		}

		return nil
	}
}

func validateSelectionType() rule {
	return func(field m.AiCommandField, fieldName string) error {
		if field.Type == m.Selection && len(field.Values) == 0 {
			return fmt.Errorf(`field "%s" is incorrect. Add at least one value to be selected in the "values" parameter.`, fieldName)
		}

		return nil
	}
}

func validateInputType() rule {
	return func(field m.AiCommandField, fieldName string) error {
		if field.Type == m.Input && len(field.Values) != 0 {
			return fmt.Errorf(`field "%s" is incorrect. Do not provide a "values" parameter if the field type is "input".`, fieldName)
		}

		return nil
	}
}

func validateDefaultRequirement() rule {
	return func(field m.AiCommandField, fieldName string) error {
		if field.Requirement == m.Default && len(field.Values) != 0 {
			return fmt.Errorf(`field "%s" is incorrect. Add one standard value that will be added automatically when the request is submitted.`, fieldName)
		}

		return nil
	}
}

func validateFileDataIsNotDefault() rule {
	return func(field m.AiCommandField, fieldName string) error {
		if field.Requirement == m.Default && field.Data == m.File {
			return fmt.Errorf(`field "%s" is incorrect. The file data cannot be a "default" value.`, fieldName)
		}

		return nil
	}
}
