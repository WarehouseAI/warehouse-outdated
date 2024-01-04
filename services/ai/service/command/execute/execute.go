package execute

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	d "warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"
	"warehouseai/ai/service/command/get"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type ExecuteCommandRequest struct {
	AiID        string
	CommandName string
	Raw         []byte
	ContentType string
}

type ExecuteCommandResponse struct {
	Raw     *bytes.Buffer
	Headers map[string]string
	Status  int
}

type makeRequestResponse struct {
	payload *http.Response
	err     *e.ErrorResponse
}

func ExecuteCommand(request ExecuteCommandRequest, aiRepository d.AiInterface, logger *logrus.Logger) (*ExecuteCommandResponse, *e.ErrorResponse) {
	ai, dbErr := get.GetCommand(get.GetCommandRequest{AiID: uuid.Must(uuid.FromString(request.AiID)), Name: request.CommandName}, aiRepository, logger)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.ErrorMessage}).Info("Execute Command")
		return nil, dbErr
	}

	executeCtx := context.Background()
	body := bytes.NewBuffer(request.Raw)
	headers := make(map[string]string)

	if ai.Command.PayloadType == string(m.FormData) {
		newBody, boundary, err := validateFormDataPayload(request.ContentType, body, ai.Command.Payload)

		if err != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.ErrorMessage}).Info("Execute Command")
			return nil, err
		}

		headers["Content-Type"] = *boundary
		headers[ai.AuthHeaderName] = ai.AuthHeaderContent
		body = newBody
	}

	if ai.Command.PayloadType == string(m.Json) {
		if err := validateJSONPayload(body, ai.Command.Payload); err != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.ErrorMessage}).Info("Execute Command")
			return nil, err
		}

		headers["Content-Type"] = "application/json"
		headers[ai.AuthHeaderName] = ai.AuthHeaderContent
	}

	reqResponse, reqErr := makeHTTPRequest(executeCtx, ai.Command.URL, string(ai.Command.RequestType), headers, body)

	if reqErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": reqErr.ErrorMessage}).Info("Execute Command")
		return nil, reqErr
	}

	cmdResponse, responseHeaders, responseStatus, decErr := decodeHTTPResponse(reqResponse, ai.Command.OutputType)

	if decErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": decErr.ErrorMessage}).Info("Execute Command")
		return nil, decErr
	}

	return &ExecuteCommandResponse{
		Raw:     cmdResponse,
		Headers: *responseHeaders,
		Status:  *responseStatus,
	}, nil
}

func updateUsageCount(existAi *m.AiProduct, ai d.AiInterface) *e.ErrorResponse {
	if err := ai.Update(existAi, map[string]interface{}{"used": existAi.Used + 1}); err != nil {
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}

func makeHTTPRequest(executeCtx context.Context, fullUrl string, httpMethod string, headers map[string]string, body *bytes.Buffer) (*http.Response, *e.ErrorResponse) {
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

// по дефолту возвращаем заголово - ответ формата JSON
func decodeHTTPResponse(response *http.Response, outputType string) (*bytes.Buffer, *map[string]string, *int, *e.ErrorResponse) {
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
