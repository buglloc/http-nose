package httpfeature

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type PathNormalizations struct {
	BaseFeature
	TransformedSlashes []string
	MergedSlashes      []string
	FixCurDir          []string
	FixParentDir       []string
}

func (f *PathNormalizations) Name() string {
	return "Path normalizations symbols"
}

func (f *PathNormalizations) Export() interface{} {
	return map[string][]string{
		"TransformedSlashes": f.TransformedSlashes,
		"MergedSlashes":      f.MergedSlashes,
		"FixCurDir":          f.FixCurDir,
		"FixParentDir":       f.FixParentDir,
	}
}

func (f *PathNormalizations) String() string {
	var result strings.Builder
	_, _ = fmt.Fprintf(&result, "TransformedSlashes: %s; ", PrintableStrings(f.TransformedSlashes))
	_, _ = fmt.Fprintf(&result, "MergedSlashes: %s; ", PrintableStrings(f.MergedSlashes))
	_, _ = fmt.Fprintf(&result, "FixCurDir: %s; ", PrintableStrings(f.FixCurDir))
	_, _ = fmt.Fprintf(&result, "FixParentDir: %s; ", PrintableStrings(f.FixParentDir))

	return result.String()
}

func (f *PathNormalizations) Collect() error {
	f.TransformedSlashes = f.checkNormalizations("/{rnd}{sym}{rnd}", "/{rnd}/{rnd}")
	f.MergedSlashes = f.checkNormalizations("/{rnd}{sym}{sym}{rnd}", "/{rnd}{sym}{rnd}")
	f.FixCurDir = f.checkNormalizations("/{rnd}{sym}.{sym}{rnd}", "/{rnd}{sym}{rnd}")
	f.FixParentDir = f.checkNormalizations("/{rnd}{sym}..{sym}{rnd}", "/{rnd}")
	return nil
}

func (f *PathNormalizations) checkNormalizations(pathTmpl, expectedPath string) []string {
	randString := RandAlphanumString(8)

	buildTmpl := func(tmpl, sym string) string {
		out := strings.ReplaceAll(tmpl, "{rnd}", randString)
		return strings.ReplaceAll(out, "{sym}", sym)
	}

	allowedSymbols := make([]string, 0)
	mu := &sync.Mutex{}
	sem := make(chan bool, concurrency)
	tester := func(sym string) {
		defer func() { <-sem }()

		req := f.BaseRequest.Clone()
		req.Path += buildTmpl(pathTmpl, sym)

		resp, err := f.Client.MakeRequest(req)
		if err != nil || resp.Status != 200 {
			return
		}

		if resp.Path == "" {
			return
		}

		if resp.Path == buildTmpl(expectedPath, sym) || resp.Path == buildTmpl(expectedPath, "/") {
			mu.Lock()
			allowedSymbols = append(allowedSymbols, sym)
			mu.Unlock()
		}
	}

	for _, c := range NotAlphaNumSyms {
		// TODO(buglloc): move me from here
		switch c {
		case '?', '#', '\r', '\n', ' ':
			continue
		default:
		}

		sym := string(c)

		sem <- true
		tester(sym)

		sem <- true
		tester(url.QueryEscape(sym))
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return allowedSymbols
}
