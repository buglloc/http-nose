package httpfeature

import (
	"fmt"
	"strings"
)

const (
	REPLACE_PROVIDED_HEADERS_NO     = 0
	REPLACE_PROVIDED_HEADERS_YES    = 1
	REPLACE_PROVIDED_HEADERS_MERGED = 2
)

type ReplaceProvidedHeaders struct {
	BaseFeature
	Headers map[string]int
}

func (f *ReplaceProvidedHeaders) Name() string {
	return "Replace provided headers"
}

func (f *ReplaceProvidedHeaders) HeadersAction() []string {
	if len(f.Headers) == 0 {
		return []string{"N/A"}
	}

	result := make([]string, len(f.Headers))
	count := 0
	for name, action := range f.Headers {
		var resultAction string
		if action == REPLACE_PROVIDED_HEADERS_NO {
			resultAction = "No"
		} else if action == REPLACE_PROVIDED_HEADERS_YES {
			resultAction = "Yes"
		} else if action == REPLACE_PROVIDED_HEADERS_MERGED {
			resultAction = "Merged"
		} else {
			resultAction = "Unknown"
		}

		result[count] = fmt.Sprintf("%s: %s", name, resultAction)
		count++
	}
	return result
}

func (f *ReplaceProvidedHeaders) Export() interface{} {
	return f.HeadersAction()
}

func (f *ReplaceProvidedHeaders) String() string {
	return strings.Join(f.HeadersAction(), ", ")
}

func (f *ReplaceProvidedHeaders) Collect() error {
	providedHeaders := f.Features.GetProvidedHeaders().Headers
	f.Headers = make(map[string]int, len(providedHeaders))
	for _, h := range providedHeaders {
		f.Headers[h.Name] = f.check(h.Name)
	}
	return nil
}

func (f *ReplaceProvidedHeaders) check(name string) int {
	req := f.BaseRequest.Clone()
	testValue := RandAlphanumString(8)
	req.AddHeader(name, testValue)
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return REPLACE_PROVIDED_HEADERS_NO
	}

	for _, h := range resp.HeadersSlice(name) {
		if h.Value == testValue {
			return REPLACE_PROVIDED_HEADERS_YES
		}

		if strings.Contains(h.Value, testValue) {
			return REPLACE_PROVIDED_HEADERS_MERGED
		}
	}
	return REPLACE_PROVIDED_HEADERS_NO
}
