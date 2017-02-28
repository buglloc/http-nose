package httpfeature

import (
	"strings"
	"fmt"
)

type ProvidedHeadersOrder struct {
	BaseFeature
	After bool
	Before bool
}

func (f *ProvidedHeadersOrder) Name() string {
	return "Server provided headers order"
}

func (f *ProvidedHeadersOrder) ToString() string {
	if ! f.After && ! f.Before {
		return "n/a"
	}
	if f.After {
		return "after user headers"
	}
	return "before user headers"
}

func (f *ProvidedHeadersOrder) Collect() error {
	providedHeaders := f.Features.GetProvidedHeaders().Headers
	if len(providedHeaders) == 0 {
		return nil
	}

	req := f.BaseRequest.Clone()
	rand := RandAlphanumString(8)
	headerName := fmt.Sprintf("X-%s", rand)
	req.AddHeader(headerName, rand)
	resp, err := f.Client.MakeRequest(req)
	if err == nil && resp.Status == 200 {
		testName := strings.ToLower(headerName)
		user := false
		provided := false
		for _, h := range resp.Headers {
			if h.Name == providedHeaders[0].Name {
				if user {
					f.After = true
					f.Before = false
					break
				}
				provided = true
			} else if strings.ToLower(h.Name) == testName {
				if provided {
					f.After = false
					f.Before = true
					break
				}
				user = true
			}
		}
	}

	return nil
}

