package httpfeature

import (
	"fmt"
	"sync"
)

type PathIgnoreSymbols struct {
	BaseFeature
	Begin []byte
	End   []byte
}

func (f *PathIgnoreSymbols) Name() string {
	return "Path ignored symbols"
}

func (f *PathIgnoreSymbols) Export() interface{} {
	return map[string][]byte{
		"Begin": f.Begin,
		"End":   f.End,
	}
}

func (f *PathIgnoreSymbols) String() string {
	return fmt.Sprintf("begin(%s), end(%s)",
		PrintableSymbols(f.Begin), PrintableSymbols(f.End))
}

func (f *PathIgnoreSymbols) Collect() error {
	testPath := fmt.Sprintf("/%s", RandAlphanumString(8))
	f.Begin = f.checkPathSymbols("{sym}"+testPath, testPath)
	f.End = f.checkPathSymbols(testPath+"{sym}", testPath)
	return nil
}

func (f *PathIgnoreSymbols) checkPathSymbols(pathTmpl, expectedPath string) []byte {
	var allowedSymbols []byte
	mu := &sync.Mutex{}
	sem := make(chan bool, concurrency)
	for _, c := range NotAlphaNumSyms {
		// TODO(buglloc): move me from here
		switch c {
		case '/', '?', '#', '\r', '\n', ' ':
			continue
		default:
		}

		sem <- true
		go func(sym byte) {
			defer func() { <-sem }()

			req := f.BaseRequest.Clone()
			req.Path = Symf(pathTmpl, sym)

			resp, err := f.Client.MakeRequest(req)
			if err != nil || resp.Status != 200 {
				return
			}

			if resp.Path != "" && expectedPath == resp.Path {
				mu.Lock()
				allowedSymbols = append(allowedSymbols, sym)
				mu.Unlock()
			}
		}(c)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return allowedSymbols
}
