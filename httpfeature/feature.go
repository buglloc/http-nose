package httpfeature

import (
	"github.com/buglloc/http-nose/httpclient"
	"strings"
	"sync"
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

func (f *BaseFeature) checkHeaderSymbols(baseRequest httpclient.Request, name, testName, value, testValue string) ([]rune, error) {
	allowed_symbols := make([]rune, 0)
	mu := &sync.Mutex{}
	sem := make(chan bool, concurrency)
	for _, c := range NotAlphaNumSyms {
		sem <- true
		go func(sym rune) {
			defer func() { <-sem }()
			req := baseRequest.Clone()

			headerName, _ := TruncatingSprintf(name, sym)
			headerName = strings.ToLower(headerName)
			headerTestName, _ := TruncatingSprintf(testName, sym)
			headerTestName = strings.ToLower(headerTestName)
			headerValue, _ := TruncatingSprintf(value, sym)
			headerTestValue, _ := TruncatingSprintf(testValue, sym)

			req.AddHeader(headerName, headerValue)
			resp, err := f.Client.MakeRequest(req)
			if err != nil || resp.Status != 200 {
				return
			}

			for _, h := range resp.Request.Headers {
				if strings.ToLower(h.Name) == headerTestName && h.Value == headerTestValue {
					mu.Lock()
					allowed_symbols = append(allowed_symbols, sym)
					mu.Unlock()
					break
				}
			}
		}(c)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return allowed_symbols, nil
}
