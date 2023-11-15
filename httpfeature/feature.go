package httpfeature

import (
	"strings"
	"sync"

	"github.com/buglloc/http-nose/httpclient"
)

type Feature interface {
	Name() string
	Collect() error
	String() string
	Export() interface{}
}

type BaseFeature struct {
	Feature
	Client       httpclient.Client
	BaseRequest  httpclient.Request
	BaseResponse httpclient.Response
	Features     *Features
}

func (f *BaseFeature) checkHeaderSymbols(baseRequest httpclient.Request, name, testName, value, testValue string) ([]byte, error) {
	var allowedSymbols []byte
	mu := &sync.Mutex{}
	sem := make(chan bool, concurrency)
	for _, c := range NotAlphaNumSyms {
		sem <- true
		go func(sym byte) {
			defer func() { <-sem }()
			req := baseRequest.Clone()

			headerName := Symf(name, sym)
			headerName = strings.ToLower(headerName)
			headerTestName := Symf(testName, sym)
			headerTestName = strings.ToLower(headerTestName)
			headerValue := Symf(value, sym)
			headerTestValue := Symf(testValue, sym)

			req.AddHeader(headerName, headerValue)
			resp, err := f.Client.MakeRequest(req)
			if err != nil || resp.Status != 200 {
				return
			}

			for _, h := range resp.Request.Headers {
				if strings.ToLower(h.Name) == headerTestName && h.Value == headerTestValue {
					mu.Lock()
					allowedSymbols = append(allowedSymbols, sym)
					mu.Unlock()
					break
				}
			}
		}(c)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return allowedSymbols, nil
}
