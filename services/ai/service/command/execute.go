package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
	"warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/sirupsen/logrus"
)

func ExecuteFormDataCommand(formData *multipart.Form, command *GetCommandResponse, ai dataservice.AiInterface, logger *logrus.Logger) (*bytes.Buffer, *e.ErrorResponse) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range command.Payload.Payload {
		if value == "file" {
			fileHeader := formData.File[key][0]
			file, err := fileHeader.Open()

			if err != nil {
				logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
				return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
			}

			field, _ := writer.CreateFormFile(key, fileHeader.Filename)
			io.Copy(field, file)

			defer file.Close()
		} else {
			writer.WriteField(key, formData.Value[key][0])
		}
	}

	writer.Close()

	headers := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": fmt.Sprintf("%s %s", string(command.AuthScheme), command.ApiKey),
	}

	response, httpErr := makeHTTPRequest(command.Payload.URL, string(command.Payload.RequestScheme), headers, url.Values{}, body)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute FormData-Command")
		return nil, httpErr
	}

	buffer, httpErr := decodeHTTPResponse(response, command.Payload.OutputType)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute FormData-Command")
		return nil, httpErr
	}

	if err := updateUsageCount(command.AI, ai); err != nil {
		return nil, err
	}

	return buffer, nil
}

func ExecuteJSONCommand(jsonData map[string]interface{}, command *GetCommandResponse, ai dataservice.AiInterface, logger *logrus.Logger) (*bytes.Buffer, *e.ErrorResponse) {
	json, parseErr := json.Marshal(jsonData)

	if parseErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": parseErr.Error()}).Info("Execute JSON-Command")
		return nil, e.NewErrorResponse(e.HttpInternalError, parseErr.Error())
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("%s %s", string(command.AuthScheme), command.ApiKey),
	}

	body := bytes.NewBuffer(json)
	response, httpErr := makeHTTPRequest(command.Payload.URL, string(command.Payload.RequestScheme), headers, url.Values{}, body)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute JSON-Command")
		return nil, httpErr
	}

	buffer, httpErr := decodeHTTPResponse(response, command.Payload.OutputType)

	if httpErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": httpErr.ErrorMessage}).Info("Execute FormData-Command")
		return nil, httpErr
	}

	if err := updateUsageCount(command.AI, ai); err != nil {
		return nil, err
	}

	return buffer, nil
}

func updateUsageCount(existAi *m.AI, ai dataservice.AiInterface) *e.ErrorResponse {
	if err := ai.Update(existAi, map[string]interface{}{"used": existAi.Used + 1}); err != nil {
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}

func makeHTTPRequest(fullUrl string, httpMethod string, headers map[string]string, queryParameters url.Values, body io.Reader) (io.ReadCloser, *e.ErrorResponse) {
	client := http.Client{}

	url, err := url.Parse(fullUrl)
	if err != nil {
		return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	if httpMethod == "GET" {
		q := url.Query()

		for k, v := range queryParameters {
			q.Set(k, strings.Join(v, ","))
		}

		url.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(httpMethod, url.String(), body)
	if err != nil {
		return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	if res == nil {
		return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	return res.Body, nil
}

func decodeHTTPResponse(response io.ReadCloser, outputType m.IOType) (*bytes.Buffer, *e.ErrorResponse) {
	if outputType == m.Image {
		var buffer bytes.Buffer
		img, _, err := image.Decode(response)

		if err != nil {
			return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
		}

		if err := jpeg.Encode(&buffer, img, nil); err != nil {
			return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
		}

		return &buffer, nil
	} else {
		var buffer bytes.Buffer
		json, err := io.ReadAll(response)

		if err != nil {
			return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
		}

		if _, err := buffer.Write(json); err != nil {
			return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
		}

		return &buffer, nil
	}
}
