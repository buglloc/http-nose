package httpfeature

import (
	"fmt"
	"strings"
)

const (
	HeaderCountOverflowActionNA         = 0
	HeaderCountOverflowActionDisallowed = 1
	HeaderCountOverflowActionCut        = 2
	HeaderCountOverflowActionBody       = 3
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
	if f.Action == HeaderCountOverflowActionNA {
		return "N/A"
	}

	if f.Action == HeaderCountOverflowActionDisallowed {
		return "Disallowed"
	}

	if f.Action == HeaderCountOverflowActionCut {
		return "Dropped"
	}

	if f.Action == HeaderCountOverflowActionBody {
		return "Leak to body"
	}
	return "Unknown"
}

func (f *HeaderCountOverflowAction) Collect() error {
	max := f.Features.GetMaximumHeadersCount().Count
	if max == MaxHeadersCount {
		f.Action = HeaderCountOverflowActionNA
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
		return HeaderCountOverflowActionDisallowed
	}

	if strings.Contains(resp.Body, rand) {
		return HeaderCountOverflowActionBody
	}

	for len(resp.HeadersSlice("X-Overflow")) > 0 {
		return HeaderCountOverflowActionNA
	}

	return HeaderCountOverflowActionCut
}
