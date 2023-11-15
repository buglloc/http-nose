package httpfeature

import (
	"fmt"
	"sync"
)

type PathIgnoreSymbols struct {
	BaseFeature
	Begin []rune
	End   []rune
}

func (f *PathIgnoreSymbols) Name() string {
	return "Path ignored symbols"
}

func (f *PathIgnoreSymbols) Export() interface{} {
	return map[string][]rune{
		"Begin": f.Begin,
		"End":   f.End,
	}
}

func (f *PathIgnoreSymbols) String() string {
	return fmt.Sprintf("begin(%s), end(%s)",
		PrintableRunes(f.Begin), PrintableRunes(f.End))
}

func (f *PathIgnoreSymbols) Collect() error {
	testPath := fmt.Sprintf("/%s", RandAlphanumString(8))
	f.Begin, _ = f.checkPathSymbols("%c"+testPath, testPath)
	f.End, _ = f.checkPathSymbols(testPath+"%c", testPath)
	return nil
}

func (f *PathIgnoreSymbols) checkPathSymbols(pathTmpl, expectedPath string) ([]rune, error) {
	allowedSymbols := make([]rune, 0)
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
		go func(sym rune) {
			defer func() { <-sem }()

			req := f.BaseRequest.Clone()
			req.Path, _ = TruncatingSprintf(pathTmpl, sym)

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

	return allowedSymbols, nil
}
