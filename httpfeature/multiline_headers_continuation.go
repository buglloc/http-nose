package httpfeature

import (
	"fmt"
	"strings"
	"sync"
)

type MultilineHeadersContinuation struct {
	BaseFeature
	Symbols []rune
}

func (f *MultilineHeadersContinuation) Name() string {
	return "Multiline headers continuation symbols"
}

func (f *MultilineHeadersContinuation) Export() interface{} {
	return f.Symbols
}

func (f *MultilineHeadersContinuation) String() string {
	return PrintableRunes(f.Symbols)
}

func (f *MultilineHeadersContinuation) Collect() error {
	if !f.Features.GetMultilineHeadersSupport().Supported {
		return nil
	}

	mu := &sync.Mutex{}
	sem := make(chan bool, concurrency)
	for _, c := range NotAlphaNumSyms {
		sem <- true
		go func(sym rune) {
			defer func() { <-sem }()
			req := f.BaseRequest.Clone()

			req.AddHeader("X-Multiline-Test", fmt.Sprintf("test\r\n%cmultiline", sym))
			resp, err := f.Client.MakeRequest(req)
			if err != nil || resp.Status != 200 {
				return
			}

			for _, h := range resp.HeadersSlice("X-Multiline-Test") {
				if strings.HasPrefix(h.Value, "test") && strings.HasSuffix(h.Value, "multiline") {
					mu.Lock()
					f.Symbols = append(f.Symbols, sym)
					mu.Unlock()
					break
				}
			}
		}(c)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return nil
}
