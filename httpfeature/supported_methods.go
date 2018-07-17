package httpfeature

import (
	"github.com/buglloc/http-nose/httpclient"
	"strings"
	"sync"
)

type SupportedMethods struct {
	BaseFeature
	Any     bool
	Methods []string
}

func (f *SupportedMethods) Name() string {
	return "Supported HTTP methods"
}

func (f *SupportedMethods) Export() interface{} {
	if f.Any {
		return []string{"Any method"}
	} else if len(f.Methods) == 0 {
		return []string{"None"}
	}

	return f.Methods
}

func (f *SupportedMethods) String() string {
	if f.Any {
		return "Any method"
	} else if len(f.Methods) == 0 {
		return "None"
	}

	return PrintableStrings(f.Methods)
}

func (f *SupportedMethods) Collect() error {
	if f.isAnyMethodSupported() {
		f.Any = true
	} else {
		f.Methods = f.getMethods()
	}
	return nil
}

func (f *SupportedMethods) checkMethod(baseRequest httpclient.Request, method string) bool {
	req := baseRequest.Clone()
	req.Method = method
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return false
	}

	if method == "OPTIONS" || method == "HEAD" {
		return true
	}

	return strings.ToLower(method) == strings.ToLower(resp.Method)
}

func (f *SupportedMethods) isAnyMethodSupported() bool {
	return f.checkMethod(f.BaseRequest, "SOMEMETHOD")
}

func (f *SupportedMethods) getMethods() []string {
	methods := make([]string, 0)
	mu := &sync.Mutex{}
	sem := make(chan bool, concurrency)
	for _, m := range HttpMethods {
		sem <- true
		go func(method string) {
			defer func() { <-sem }()

			if f.checkMethod(f.BaseRequest, method) {
				mu.Lock()
				methods = append(methods, method)
				mu.Unlock()
			}
		}(m)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return methods
}
