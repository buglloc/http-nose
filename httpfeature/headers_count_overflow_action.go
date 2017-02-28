package httpfeature

import (
	"fmt"
	"strings"
)

const (
	HEADER_COUNT_OVERFLOW_ACTION_NA         = 0
	HEADER_COUNT_OVERFLOW_ACTION_DISALLOWED = 1
	HEADER_COUNT_OVERFLOW_ACTION_CUT        = 2
	HEADER_COUNT_OVERFLOW_ACTION_BODY       = 3
)

type HeaderCountOverflowAction struct {
	BaseFeature
	Action int
}

func (f *HeaderCountOverflowAction) Name() string {
	return "Headers count overflow action"
}

func (f *HeaderCountOverflowAction) ToString() string {
	if f.Action == HEADER_COUNT_OVERFLOW_ACTION_NA {
		return "n/a"
	}

	if f.Action == HEADER_COUNT_OVERFLOW_ACTION_DISALLOWED {
		return "disallowed"
	}

	if f.Action == HEADER_COUNT_OVERFLOW_ACTION_CUT {
		return "cut"
	}

	if f.Action == HEADER_COUNT_OVERFLOW_ACTION_BODY {
		return "body"
	}
	return "unknown"
}

func (f *HeaderCountOverflowAction) Collect() error {
	max := f.Features.GetMaximumHeadersCount().Count
	if max == MAX_HEADERS_COUNT {
		f.Action = HEADER_COUNT_OVERFLOW_ACTION_NA
	} else {
		f.Action = f.check(max)
	}
	return nil
}

func (f *HeaderCountOverflowAction) check(maximum int) int {
	req := f.BaseRequest.Clone()
	for i := 0; i < maximum; i++ {
		req.AddHeader(fmt.Sprintf("X-Foo%d", i), "a")
	}

	rand := RandAlphanumString(8)
	req.AddHeader("X-Overflow", "overflow" + rand)

	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return HEADER_COUNT_OVERFLOW_ACTION_DISALLOWED
	}

	if strings.Contains(resp.Body, rand) {
		return HEADER_COUNT_OVERFLOW_ACTION_BODY
	}

	for _, h := range resp.Headers {
		if strings.ToLower(h.Name) == "x-overflow" {
			return HEADER_COUNT_OVERFLOW_ACTION_NA
		}
	}

	return HEADER_COUNT_OVERFLOW_ACTION_CUT
}
