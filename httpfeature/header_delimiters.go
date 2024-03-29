package httpfeature

import (
	"sync"
)

type HeaderDelimiters struct {
	BaseFeature
	Symbols []byte
}

func (f *HeaderDelimiters) Name() string {
	return "Header name/value delimitiers"
}

func (f *HeaderDelimiters) Export() interface{} {
	return f.Symbols
}

func (f *HeaderDelimiters) String() string {
	return PrintableSymbols(f.Symbols)
}

func (f *HeaderDelimiters) Collect() error {
	f.Symbols, _ = f.collectSymbols()
	return nil
}

func (f *HeaderDelimiters) collectSymbols() ([]byte, error) {
	var symbols []byte
	mu := &sync.Mutex{}
	sem := make(chan bool, concurrency)
	for _, c := range NotAlphaNumSyms {
		sem <- true
		go func(sym byte) {
			defer func() { <-sem }()

			req := f.BaseRequest.Clone()
			req.AddHeader("X-Foo", "foo")
			resp, err := f.Client.MakeRawRequest(req.Build(nil, []byte{sym}))
			if err != nil || resp.Status != 200 {
				return
			}

			for _, h := range resp.HeadersSlice("X-Foo") {
				if h.Value == "foo" {
					mu.Lock()
					symbols = append(symbols, sym)
					mu.Unlock()
					break
				}
			}
		}(c)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return symbols, nil
}
