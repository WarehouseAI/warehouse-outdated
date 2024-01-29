package execute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
	d "warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/sirupsen/logrus"
)

type ExecuteCommandRequest[T *multipart.Form | map[string]interface{}] struct {
	AI      *m.AiProduct
	Command *m.AiCommand
	Payload T
}

type ExecuteCommandResponse struct {
	Raw     *bytes.Buffer
	Headers map[string]string
	Status  int
}

type makeRequestResponse struct {
	payload *http.Response
	err     *e.HttpErrorResponse
}

func ExecuteJSONCommand(
	request ExecuteCommandRequest[map[string]interface{}],
	aiRepository d.AiInterface,
	logger *logrus.Logger,
) (*ExecuteCommandResponse, *e.HttpErrorResponse) {
	if err := validateJSONPayload(&request); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.ErrorMessage}).Info("Execute Command")
		return nil, err
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers[request.AI.AuthHeaderName] = request.AI.AuthHeaderContent

	var buffer bytes.Buffer

	if err := json.NewEncoder(&buffer).Encode(request.Payload); err != nil {
		return nil, e.NewErrorResponse(e.HttpInternalError, fmt.Sprintf("Error encoding map to JSON: %s", err))
	}

	executeCtx := context.Background()
	reqResponse, reqErr := makeHTTPRequest(executeCtx, request.Command.URL, request.Command.RequestType, headers, &buffer)

	if reqErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": reqErr.ErrorMessage}).Info("Execute Command")
		return nil, reqErr
	}

	cmdResponse, responseHeaders, responseStatus, decErr := decodeHTTPResponse(reqResponse, request.Command.OutputType)

	if decErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": decErr.ErrorMessage}).Info("Execute Command")
		return nil, decErr
	}

	if err := updateUsageCount(request.AI, aiRepository); err != nil {
		return nil, err
	}

	return &ExecuteCommandResponse{
		Raw:     cmdResponse,
		Headers: *responseHeaders,
		Status:  *responseStatus,
	}, nil
}

func ExecuteFormCommand(
	request ExecuteCommandRequest[*multipart.Form],
	boundary string,
	aiRepository d.AiInterface,
	logger *logrus.Logger,
) (*ExecuteCommandResponse, *e.HttpErrorResponse) {
	if err := validateFormDataPayload(&request); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.ErrorMessage}).Info("Execute Command")
		return nil, err
	}

	headers := make(map[string]string)
	headers["Content-Type"] = fmt.Sprintf("multipart/form-data; boundary=%s", boundary)
	headers[request.AI.AuthHeaderName] = request.AI.AuthHeaderContent

	// Конвертим mutlipart/form-data в буффер
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	defer writer.Close()

	// Iterate over the form fields and add them to the writer
	for fieldName, fieldValues := range request.Payload.Value {
		for _, fieldValue := range fieldValues {
			writer.WriteField(fieldName, fieldValue)
		}
	}

	// Iterate over the form files and add them to the writer
	for fieldName, fileHeaders := range request.Payload.File {
		for _, fileHeader := range fileHeaders {
			fileWriter, err := writer.CreateFormFile(fieldName, fileHeader.Filename)
			if err != nil {
				logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute Command")
				return nil, e.NewErrorResponse(e.HttpInternalError, fmt.Sprintf("Error creating form file: %s", err))
			}

			// Open the file and copy its content to the form
			file, err := fileHeader.Open()
			if err != nil {
				logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute Command")
				return nil, e.NewErrorResponse(e.HttpInternalError, fmt.Sprintf("Error opening file: %s", err))
			}
			defer file.Close()

			_, err = io.Copy(fileWriter, file)
			if err != nil {
				logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute Command")
				return nil, e.NewErrorResponse(e.HttpInternalError, fmt.Sprintf("Error copying file: %s", err))
			}
		}
	}

	executeCtx := context.Background()
	reqResponse, reqErr := makeHTTPRequest(executeCtx, request.Command.URL, request.Command.RequestType, headers, &buffer)

	if reqErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": reqErr.ErrorMessage}).Info("Execute Command")
		return nil, reqErr
	}

	cmdResponse, responseHeaders, responseStatus, decErr := decodeHTTPResponse(reqResponse, request.Command.OutputType)

	if decErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": decErr.ErrorMessage}).Info("Execute Command")
		return nil, decErr
	}

	if err := updateUsageCount(request.AI, aiRepository); err != nil {
		return nil, err
	}

	return &ExecuteCommandResponse{
		Raw:     cmdResponse,
		Headers: *responseHeaders,
		Status:  *responseStatus,
	}, nil
}

func updateUsageCount(existAi *m.AiProduct, ai d.AiInterface) *e.HttpErrorResponse {
	if err := ai.Update(existAi, map[string]interface{}{"used": existAi.Used + 1}); err != nil {
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}

func makeHTTPRequest(executeCtx context.Context, fullUrl string, httpMethod string, headers map[string]string, body *bytes.Buffer) (*http.Response, *e.HttpErrorResponse) {
	httpClient := http.Client{}

	url, err := url.Parse(fullUrl)
	if err != nil {
		return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	req, err := http.NewRequest(httpMethod, url.String(), body)
	if err != nil {
		return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(executeCtx, time.Second*30)
	respch := make(chan makeRequestResponse)
	defer cancel()

	go func() {
		res, err := httpClient.Do(req)

		if err != nil {
			respch <- makeRequestResponse{
				payload: nil,
				err:     e.NewErrorResponse(e.HttpInternalError, err.Error()),
			}
			return
		}

		if res == nil {
			respch <- makeRequestResponse{
				payload: nil,
				err:     e.NewErrorResponse(e.HttpInternalError, "AI return the empty response"),
			}
			return
		}

		respch <- makeRequestResponse{
			payload: res,
			err:     nil,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, e.NewErrorResponse(e.HttpTimeout, "The request has exceeded the waiting time")

		case resp := <-respch:
			return resp.payload, nil
		}
	}
}

// по дефолту возвращаем заголовок - ответ формата JSON
func decodeHTTPResponse(response *http.Response, outputType string) (*bytes.Buffer, *map[string]string, *int, *e.HttpErrorResponse) {
	var buffer bytes.Buffer
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	rawResponse, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, nil, nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	if _, err := buffer.Write(rawResponse); err != nil {
		return nil, nil, nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	if response.StatusCode == 200 {
		if outputType == string(m.Audio) {
			headers["Content-Type"] = "audio/mp3"
			headers["Content-Length"] = strconv.Itoa(len(rawResponse))
		}

		if outputType == string(m.Image) {
			headers["Content-Type"] = "image/png"
		}
	}

	return &buffer, &headers, &response.StatusCode, nil
}
