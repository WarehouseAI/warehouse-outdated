package httputils

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	pg "warehouse/src/internal/database/postgresdb"
)

func MakeHTTPRequest(fullUrl string, httpMethod string, headers map[string]string, queryParameters url.Values, body io.Reader) (io.ReadCloser, *ErrorResponse) {
	client := http.Client{}

	url, err := url.Parse(fullUrl)
	if err != nil {
		return nil, NewErrorResponse(InternalError, err.Error())
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
		return nil, NewErrorResponse(InternalError, err.Error())
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, NewErrorResponse(InternalError, err.Error())
	}

	if res == nil {
		return nil, NewErrorResponse(InternalError, err.Error())
	}

	return res.Body, nil
}

func DecodeHTTPResponse(response io.ReadCloser, outputType pg.IOType) (*bytes.Buffer, *ErrorResponse) {
	if outputType == pg.Image {
		var buffer bytes.Buffer
		img, _, err := image.Decode(response)

		if err != nil {
			return nil, NewErrorResponse(InternalError, err.Error())
		}

		if err := jpeg.Encode(&buffer, img, nil); err != nil {
			return nil, NewErrorResponse(InternalError, err.Error())
		}

		return &buffer, nil
	} else {
		var buffer bytes.Buffer
		json, err := ioutil.ReadAll(response)

		if err != nil {
			return nil, NewErrorResponse(InternalError, err.Error())
		}

		if _, err := buffer.Write(json); err != nil {
			return nil, NewErrorResponse(InternalError, err.Error())
		}

		return &buffer, nil
	}
}
