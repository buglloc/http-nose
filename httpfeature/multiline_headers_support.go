package httpfeature

import (
	"strings"
)

type MultilineHeadersSupport struct {
	BaseFeature
	Supported bool
}

func (f *MultilineHeadersSupport) Name() string {
	return "Supported multiline headers"
}

func (f *MultilineHeadersSupport) ToString() string {
	return PrintableBool(f.Supported)
}

func (f *MultilineHeadersSupport) Collect() error {
	f.Supported = f.check()
	return nil
}

func (f *MultilineHeadersSupport) check() bool {
	req := f.BaseRequest.Clone()
	req.AddHeader("X-Foo", "foo\r\n\tbar")
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return false
	}
	for _, h := range resp.Headers {
		if strings.ToLower(h.Name) == "x-foo" {
			if strings.Contains(h.Value, "foo") && strings.Contains(h.Value, "bar") {
				return true
			}
		}
	}
	return false
}
