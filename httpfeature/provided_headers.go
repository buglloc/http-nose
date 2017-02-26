package httpfeature

import (
	"strings"
	"fmt"
	"github.com/buglloc/http-nose/httpclient"
)

type ProvidedHeaders struct {
	BaseFeature
	Headers []httpclient.Header
}

func (f *ProvidedHeaders) Name() string {
	return "Server provided headers"
}

func (f *ProvidedHeaders) ToString() string {
	if len(f.Headers) == 0 {
		return "none"
	}

	result := make([]string, len(f.Headers))
	for i, h := range f.Headers {
		result[i] = fmt.Sprintf("%s: %q", h.Name, h.Value)
	}
	return strings.Join(result, ", ")
}

func (f *ProvidedHeaders) Collect() error {
	baseHeaders := make(map[string]bool, len(f.BaseRequest.Headers))
	for _, h := range f.BaseRequest.Headers {
		baseHeaders[strings.ToLower(h.Name)] = true
	}

	f.Headers = make([]httpclient.Header, 0)
	for _, h := range f.BaseResponse.Headers {
		_, requested := baseHeaders[strings.ToLower(h.Name)]
		if ! requested {
			f.Headers = append(f.Headers, h)
		}
	}
	return nil
}