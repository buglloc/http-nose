package httpclient

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Request
	Status int
}

func NewFromHttpResponse(response *http.Response) (*Response, error) {
	r := &Response{
		Status: response.StatusCode,
	}

	if response.StatusCode == 200 && response.ContentLength > 0 {
		err := json.NewDecoder(response.Body).Decode(&r.Request)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}
