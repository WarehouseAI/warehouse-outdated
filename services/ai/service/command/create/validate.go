package create

import (
	"encoding/json"
	"fmt"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"
)

func validateCreateRequest(request *CreateCommandRequest) *e.ErrorResponse {
	for name, value := range request.Payload {
		var fieldParameters m.CommandFieldParams
		jsonData, err := json.Marshal(value)

		if err != nil {
			return e.NewErrorResponse(e.HttpInternalError, err.Error())
		}

		if err := json.Unmarshal(jsonData, &fieldParameters); err != nil {
			return e.NewErrorResponse(e.HttpUnprocessableEntity, fmt.Sprintf(`Invalid decloration for "%s" field.`, name))
		}

		// Валидирование типов
		if valid := isValidFieldClass(string(fieldParameters.Class)); !valid {
			return e.NewErrorResponse(e.HttpUnprocessableEntity,
				fmt.Sprintf(`Invalid class in "%s" field. Class "%s" does not exists.`, name, string(fieldParameters.Class)))
		}

		if valid := isValidDataType(string(fieldParameters.DataType)); !valid {
			return e.NewErrorResponse(e.HttpUnprocessableEntity,
				fmt.Sprintf(`Invalid data_type in "%s" field. Data type "%s" does not exists.`, name, string(fieldParameters.Class)))
		}

		// Если тип данных у полян не текст. тогда тип данных у запроса должен быть FormData
		if fieldParameters.DataType == m.File && request.PayloadType == m.Json {
			return e.NewErrorResponse(e.HttpUnprocessableEntity,
				"Invalid payload_type. If you want to send an image when executing commands - use FormData as payload_type.")
		}

		if fieldParameters.DataType == m.Object && request.PayloadType == m.FormData {
			return e.NewErrorResponse(e.HttpUnprocessableEntity,
				`Invalid payload_type. If you want to use "object" data type - use JSON as payload_type.`)
		}

		// Нельзя использовать классы "permanent" и "optional" для файлов
		if fieldParameters.DataType == m.File && fieldParameters.Class != m.Free {
			return e.NewErrorResponse(e.HttpUnprocessableEntity,
				fmt.Sprintf(`Invalid data_type in %s field. Fields with the data type "Image" must be "free" class.`, name))
		}

		// Для полей свободного класса нельзя добавить возможные значения
		if fieldParameters.Class == m.Free && len(fieldParameters.Values) != 0 {
			return e.NewErrorResponse(e.HttpUnprocessableEntity,
				fmt.Sprintf(`Invalid values for %s field. If you are using a free class for a field, leave the "values" field blank.`, name))
		}

		// Для полей перманентного класса нужно использовать только одно значение
		if fieldParameters.Class == m.Permanent && len(fieldParameters.Values) != 1 {
			return e.NewErrorResponse(e.HttpUnprocessableEntity,
				fmt.Sprintf(`Invalid values for %s field. Only one value is used for permanent fields.`, name))
		}

	}

	return nil
}
