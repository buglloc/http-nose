package httpfeature

import (
	"strings"
)

type AbsoluteRequestUri struct {
	BaseFeature
	Supported bool
	Untouched  bool
	PickHost  bool
	PickPath  bool
}

func (f *AbsoluteRequestUri) Name() string {
	return "Absolute request uri support"
}

func (f *AbsoluteRequestUri) ToString() string {
	result := make([]string, 0)
	if f.Supported {
		result = append(result, "Supported")
		if f.Untouched {
			result = append(result, "Untouched")
		}
		if f.PickHost {
			result = append(result, "PickHost")
		}
		if f.PickPath {
			result = append(result, "PickPath")
		}
	} else {
		result = append(result, "Unsupported")
	}
	return strings.Join(result, ", ")
}

func (f *AbsoluteRequestUri) Collect() error {
	req := f.BaseRequest.Clone()
	req.RemoveHeader("Host")
	req.RequestURI = "http://host-from-uri/path"
	req.AddHeader("Host", "host-from-header")
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		f.Supported = false
	} else {
		f.Supported = true
		f.Untouched = resp.RequestURI == "http://host-from-uri/path"
		f.PickHost = resp.Host == "host-from-uri"
		f.PickPath = resp.Path == "/path"
	}
	return nil
}
