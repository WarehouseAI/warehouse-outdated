package execute

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"sort"
	"strings"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"
)

type rule func(originField m.AiCommandField, originFieldName string, actualFieldValue interface{}) error

type validator struct {
	rules          []rule
	requiredFields map[string]bool
}

func newValidator(rules []rule, requiredField map[string]bool) validator {
	return validator{
		rules:          rules,
		requiredFields: requiredField,
	}
}

func (v *validator) validateField(originField m.AiCommandField, fieldName string, fieldValue interface{}) []error {
	var errors []error

	for _, rule := range v.rules {
		fmt.Println(fieldName)
		if err := rule(originField, fieldName, fieldValue); err != nil {
			errors = append(errors, err)
		}

		// Удаляем поле из обязательных, после успешной валидации
		if _, found := v.requiredFields[fieldName]; found {
			delete(v.requiredFields, fieldName)
		}
	}

	return errors
}

// Валидация для команд принимающих JSON
func validateJSONPayload(request *ExecuteCommandRequest[map[string]interface{}]) *e.HttpErrorResponse {
	var typedOriginPayload map[string]m.AiCommandField // не инициализируем мапу, так как в нее будет парситься payload из бд
	requiredFields := make(map[string]bool)

	// Конвертим оригинальную типизацию полей для упрощения проверки
	rawOriginPayload, err := json.Marshal(request.Command.Payload)

	if err != nil {
		return e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	if err := json.Unmarshal(rawOriginPayload, &typedOriginPayload); err != nil {
		return e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	// Получаем все обязательные поля
	for key, value := range typedOriginPayload {
		if value.Requirement == m.Require {
			requiredFields[key] = true
		}
	}

	v := newValidator([]rule{
		validateSelectionValue(),
	}, requiredFields)

	for jsonFieldName, jsonFieldValue := range request.Payload {
		fieldDeclaration, found := typedOriginPayload[jsonFieldName]

		if !found {
			return e.NewErrorResponse(e.HttpBadRequest, fmt.Sprintf(`Field "%s" not found in origin command payload.`, jsonFieldName))
		}

		// Валидируем поле по обозначеным правилам
		if errors := v.validateField(fieldDeclaration, jsonFieldName, jsonFieldValue); len(errors) != 0 {
			var messages []string

			for _, err := range errors {
				messages = append(messages, err.Error())
			}

			return e.NewErrorResponseMultiple(e.HttpUnprocessableEntity, messages)
		}
	}

	// Проверяем что обязательные поля прошли проверку
	if len(requiredFields) != 0 {
		unprovidedFields := make([]string, 0, len(requiredFields))

		for key := range requiredFields {
			unprovidedFields = append(unprovidedFields, key)
		}

		sort.Strings(unprovidedFields)

		return e.NewErrorResponse(e.HttpUnprocessableEntity, fmt.Sprintf("Required fields %s not provided", strings.Join(unprovidedFields, "/")))
	}

	return nil
}

func validateFormDataPayload(request *ExecuteCommandRequest[*multipart.Form]) *e.HttpErrorResponse {
	var typedOriginPayload map[string]m.AiCommandField // не инициализируем мапу, так как в нее будет парситься payload из бд
	requiredFields := make(map[string]bool)
	formTextFields := make(map[string][]string)

	// Конвертим оригинальную типизацию полей для упрощения проверки
	rawOriginPayload, err := json.Marshal(request.Command.Payload)

	if err != nil {
		return e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	if err := json.Unmarshal(rawOriginPayload, &typedOriginPayload); err != nil {
		return e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	// Получаем все обязательные поля
	for key, value := range typedOriginPayload {
		if value.Requirement == m.Require {
			requiredFields[key] = true
		}
	}

	// Ищем поля типа File и Text, поля второго типа заносим в отдельную мапу для дальнейшей валидации
	for fieldName, fieldDeclaration := range typedOriginPayload {
		if fieldDeclaration.Data == m.File {
			if _, ok := request.Payload.File[fieldName]; !ok {
				return e.NewErrorResponse(e.HttpBadRequest, fmt.Sprintf(`field "%s" has incorrect. Field must be provided as File type`, fieldName))
			}

			delete(requiredFields, fieldName) // удаляем из обязательных, так как мы проверили его наличие
		}

		if fieldDeclaration.Data != m.File {
			values, ok := request.Payload.Value[fieldName]

			if !ok {
				return e.NewErrorResponse(e.HttpBadRequest, fmt.Sprintf(`field "%s" has incorrect. Field must be provided as Text type`, fieldName))
			}

			formTextFields[fieldName] = values
		}
	}

	v := newValidator([]rule{
		validateSelectionValue(),
	}, requiredFields)

	for formTextFieldName, formTextFieldValue := range formTextFields {
		fieldDeclaration, found := typedOriginPayload[formTextFieldName]

		if !found {
			return e.NewErrorResponse(e.HttpBadRequest, fmt.Sprintf(`Field "%s" not found in origin command payload.`, formTextFieldName))
		}

		// Валидируем поле по обозначеным правилам
		if errors := v.validateField(fieldDeclaration, formTextFieldName, formTextFieldValue); len(errors) != 0 {
			var messages []string

			for _, err := range errors {
				messages = append(messages, err.Error())
			}

			return e.NewErrorResponseMultiple(e.HttpUnprocessableEntity, messages)
		}
	}

	return nil
}

func validateSelectionValue() rule {
	return func(originField m.AiCommandField, fieldName string, fieldValue interface{}) error {
		// Проверяем полученное значение для типа Selection на валидность
		if originField.Type == m.Selection && originField.Data == m.Object {
			var allowedValues []map[string]interface{}
			target := fieldValue.(map[string]interface{})

			for _, value := range originField.Values {
				var allowedValue map[string]interface{}

				jsonBytes, err := json.Marshal(value)

				if err != nil {
					return fmt.Errorf(err.Error())
				}

				if err := json.Unmarshal(jsonBytes, &allowedValue); err != nil {
					fmt.Println(fieldName)
					return fmt.Errorf(err.Error())
				}

				allowedValues = append(allowedValues, allowedValue)
			}

			// Ищем такой же объект в доступных на выбор объектах
			for _, value := range allowedValues {
				for targetKey, targetValue := range target {
					allowedValueParameter, ok := value[targetKey]

					if !ok {
						return fmt.Errorf(`field "%s" has incorrect. There is no "%s" parameter in allowed objects.`, fieldName, targetKey)
					}

					if allowedValueParameter != targetValue {
						return fmt.Errorf(`field "%s" has incorrect. Invalid parameter value "%s" in allowed object.`, fieldName, targetValue)
					}
				}
			}
		}

		if originField.Type == m.Selection && originField.Data == m.String {
			var allowedValues []string
			target := fieldValue.(string)

			for _, value := range originField.Values {
				allowedValue := value.(string)
				allowedValues = append(allowedValues, allowedValue)
			}

			sort.Strings(allowedValues)

			idx := sort.Search(len(allowedValues), func(i int) bool {
				return allowedValues[i] >= target
			})

			if idx == len(allowedValues) {
				return fmt.Errorf(`field "%s" has incorrect. Value "%s" is not exists in allowed values.`, fieldName, target)
			}
		}

		if originField.Type == m.Selection && originField.Data == m.Number {
			var allowedValues []float64
			target := fieldValue.(float64)

			for _, value := range originField.Values {
				allowedValue := value.(float64)
				allowedValues = append(allowedValues, allowedValue)
			}

			sort.Float64s(allowedValues)

			idx := sort.Search(len(allowedValues), func(i int) bool {
				return allowedValues[i] >= target
			})

			if idx == len(allowedValues) {
				return fmt.Errorf(`field "%s" has incorrect. Value "%f" is not exists in allowed values.`, fieldName, target)
			}
		}

		return nil
	}
}
