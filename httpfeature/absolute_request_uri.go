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

func (f *AbsoluteRequestUri) String() string {
	return strings.Join(f.features(), ", ")
}

func (f *AbsoluteRequestUri) Export() interface{} {
	return f.features()
}

func (f *AbsoluteRequestUri) features() []string {
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

	return result
}

func (f *AbsoluteRequestUri) Collect() error {
	req := f.BaseRequest.Clone()

	targetPath := req.Path + "test"
	uriHost := "host-from-uri." + req.Host
	headerHost := "host-from-header." + req.Host
	targetRequestUri := "http://" + uriHost + targetPath

	req.RequestURI = targetRequestUri
	req.Host = headerHost

	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		f.Supported = false
	} else {
		f.Supported = true
		f.Untouched = resp.RequestURI == targetRequestUri
		f.PickHost = resp.Host == uriHost
		f.PickPath = resp.Path == targetPath
	}
	return nil
}
