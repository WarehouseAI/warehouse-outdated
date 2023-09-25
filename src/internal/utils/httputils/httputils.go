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
	dbm "warehouse/src/internal/db/models"
	"warehouse/src/internal/dto"
)

func MakeHTTPRequest(fullUrl string, httpMethod string, headers map[string]string, queryParameters url.Values, body io.Reader) (io.ReadCloser, error) {
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

	return res.Body, nil
}

func DecodeHTTPResponse(response io.ReadCloser, outputType dbm.IOType) (*bytes.Buffer, error) {
	if outputType == dbm.Image {
		var buffer bytes.Buffer
		img, _, err := image.Decode(response)

		if err != nil {
			return nil, err
		}

		if err := jpeg.Encode(&buffer, img, nil); err != nil {
			return nil, err
		}

		return &buffer, nil
	} else {
		var buffer bytes.Buffer
		json, err := ioutil.ReadAll(response)

		if err != nil {
			return nil, err
		}

		if _, err := buffer.Write(json); err != nil {
			return nil, err
		}

		return &buffer, nil
	}
}
