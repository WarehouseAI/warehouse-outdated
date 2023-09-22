package httputils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"warehouse/src/internal/dto"
)

func MakeHTTPRequest[T any](fullUrl string, httpMethod string, headers map[string]string, queryParameters url.Values, body io.Reader, responseType T) (*T, error) {
	client := http.Client{}

	url, err := url.Parse(fullUrl)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, dto.EmptyResponse
	}

	responseData, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, dto.BadRequestError
	}

	var responseObject T
	if err := json.Unmarshal(responseData, &responseObject); err != nil {
		return nil, err
	}

	return &responseObject, nil
}
