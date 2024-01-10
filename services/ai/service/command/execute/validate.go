package execute

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"strings"
	e "warehouseai/ai/errors"
)

// Поддерживаем только формат "Text" в FormData
func validateFormDataPayload(contentType string, rawRequest *bytes.Buffer, originPayload map[string]interface{}) (*bytes.Buffer, *string, *e.HttpErrorResponse) {
	mediaType, params, err := mime.ParseMediaType(contentType)

	if err != nil {
		return nil, nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return nil, nil, e.NewErrorResponse(e.HttpBadRequest, `Invalid Content-Type for this command. No "multipart/" prefix`)
	}

	reader := multipart.NewReader(rawRequest, params["boundary"])
	formData := make(map[string]string)

	// Валидируем форм дату на сходство пейлоаду команды в БД
	for {
		part, err := reader.NextPart()

		// выбрасывает io.EOF когда "части" закончились
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
		}

		fieldDeclaration, found := originPayload[part.FormName()]

		json.Marshal(fieldDeclaration)

		if !found {
			return nil, nil, e.NewErrorResponse(e.HttpBadRequest, "Invalid command payload.")
		}

		keyValue, readErr := io.ReadAll(part)

		if readErr != nil {
			return nil, nil, e.NewErrorResponse(e.HttpInternalError, readErr.Error())
		}

		formData[part.FormName()] = string(keyValue)
	}

	var newBuffer bytes.Buffer
	writer := multipart.NewWriter(&newBuffer)
	defer writer.Close()

	for key, value := range formData {
		writer.WriteField(key, value)
	}
	writer.FormDataContentType()

	rawRequest = &newBuffer
	boundary := writer.FormDataContentType()

	return &newBuffer, &boundary, nil
}
