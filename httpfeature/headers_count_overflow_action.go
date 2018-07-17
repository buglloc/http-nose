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

func (f *HeaderCountOverflowAction) Export() interface{} {
	return f.String()
}

func (f *HeaderCountOverflowAction) String() string {
	if f.Action == HEADER_COUNT_OVERFLOW_ACTION_NA {
		return "N/A"
	}

	if f.Action == HEADER_COUNT_OVERFLOW_ACTION_DISALLOWED {
		return "Disallowed"
	}

	if f.Action == HEADER_COUNT_OVERFLOW_ACTION_CUT {
		return "Dropped"
	}

	if f.Action == HEADER_COUNT_OVERFLOW_ACTION_BODY {
		return "Leak to body"
	}
	return "Unknown"
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
	req.AddHeader("X-Overflow", "overflow"+rand)

	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return HEADER_COUNT_OVERFLOW_ACTION_DISALLOWED
	}

	if strings.Contains(resp.Body, rand) {
		return HEADER_COUNT_OVERFLOW_ACTION_BODY
	}

	for len(resp.HeadersSlice("X-Overflow")) > 0 {
		return HEADER_COUNT_OVERFLOW_ACTION_NA
	}

	return HEADER_COUNT_OVERFLOW_ACTION_CUT
}
