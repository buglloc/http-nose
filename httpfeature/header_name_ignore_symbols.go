package httpfeature

import "fmt"

type HeaderNameIgnoreSymbols struct {
	BaseFeature
	Begin  []rune
	Middle []rune
	End    []rune
}

func (f *HeaderNameIgnoreSymbols) Name() string {
	return "Header name ignored symbols"
}

func (f *HeaderNameIgnoreSymbols) ToString() string {
	return fmt.Sprintf("begin(%s), middle(%s), end(%s)",
		PrintableRunes(f.Begin), PrintableRunes(f.Middle), PrintableRunes(f.End))
}

func (f *HeaderNameIgnoreSymbols) Collect() error {
	f.Begin, _ = f.checkHeaderSymbols(f.BaseRequest, "%cx-testsym", "x-testsym","test", "test")
	f.Middle, _ = f.checkHeaderSymbols(f.BaseRequest, "x-test%csym", "x-testsym","test", "test")
	f.End, _ = f.checkHeaderSymbols(f.BaseRequest, "x-testsym%c", "x-testsym","test", "test")
	return nil
}
