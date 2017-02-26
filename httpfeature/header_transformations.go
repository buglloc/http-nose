package httpfeature

import (
	"strings"
	"fmt"
)

type HeaderTransformations struct {
	BaseFeature
	HeaderName []string
	HeaderDelimiter []string
	HeaderLineDelimiter []string
}

func (f *HeaderTransformations) Name() string {
	return "Header transformations"
}

func (f *HeaderTransformations) ToString() string {
	result := make([]string, 0)
	if len(f.HeaderName) > 0 {
		result = append(result, "HeaderName: " + strings.Join(f.HeaderName, ", "))
	}
	if len(f.HeaderDelimiter) > 0 {
		result = append(result, "HeaderDelimiter: " + strings.Join(f.HeaderDelimiter, ", "))
	}
	if len(f.HeaderLineDelimiter) > 0 {
		result = append(result, "HeaderLineDelimiter: " + strings.Join(f.HeaderLineDelimiter, ", "))
	}
	if len(result) == 0 {
		return "none"
	}
	return strings.Join(result, "; ")
}

func (f *HeaderTransformations) Collect() error {
	f.HeaderName = f.checkHeaderName()
	f.HeaderDelimiter = f.checkHeaderDelimiter(f.Features.GetHeaderDelimiters().Symbols)
	f.HeaderLineDelimiter = f.checkHeaderLineDelimiter(f.Features.GetHeaderLineDelimiters().Delims)
	return nil
}

func (f *HeaderTransformations) checkHeaderName() []string {
	result := make([]string, 0)
	req := f.BaseRequest.Clone()
	testName := "X-Test-CASE"
	testValue := RandAlphanumString(8)
	req.AddHeader(testName, testValue)
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return result
	}

	for _, h := range resp.Headers {
		if h.Value == testValue && h.Name != testName {
			result = append(result, fmt.Sprintf("%q -> %q", testName, h.Name))
		}
	}
	return result
}

func (f *HeaderTransformations) checkHeaderLineDelimiter(delims []string) []string {
	result := make([]string, 0)
	if len(f.BaseResponse.Raw) == 0 {
		return result
	}

	for _, d := range delims {
		req := f.BaseRequest.Clone()
		rand := RandAlphanumString(8)
		req.AddHeader("X-Test", fmt.Sprintf("%s%s%s:test", rand, d, rand))
		resp, err := f.Client.MakeRequest(req)
		if err != nil || resp.Status != 200 {
			continue
		}

		if len(resp.Raw) > 0 {
			startPos := strings.Index(resp.Raw, rand)
			if startPos == -1 {
				continue
			} else {
				startPos += 8
			}
			endPos := strings.Index(resp.Raw[startPos:], rand)
			if endPos == -1 {
				continue
			} else {
				endPos += startPos
			}
			trans := resp.Raw[startPos:endPos]
			if trans != d {
				result = append(result, fmt.Sprintf("%q -> %q", d, trans))
			}
		}
	}

	return result
}

func (f *HeaderTransformations) checkHeaderDelimiter(delims []rune) []string {
	result := make([]string, 0)
	if len(f.BaseResponse.Raw) == 0 {
		return result
	}

	for _, d := range delims {
		delimiter := string(d)
		req := f.BaseRequest.Clone()
		testName := fmt.Sprintf("X-%s", RandAlphanumString(8))
		testValue := fmt.Sprintf("%s", RandAlphanumString(8))
		req.AddHeader(testName, testValue)
		resp, err := f.Client.MakeRawRequest(req.Build(nil, []byte(delimiter)))
		if err != nil || resp.Status != 200 {
			continue
		}

		if len(resp.Raw) > 0 {
			startPos := strings.Index(resp.Raw, testName)
			if startPos == -1 {
				continue
			} else {
				startPos += len(testName)
			}
			endPos := strings.Index(resp.Raw[startPos:], testValue)
			if endPos == -1 {
				continue
			} else {
				endPos += startPos
			}
			trans := resp.Raw[startPos:endPos]
			if trans != delimiter {
				result = append(result, fmt.Sprintf("%q -> %q", d, trans))
			}
		}
	}

	return result
}