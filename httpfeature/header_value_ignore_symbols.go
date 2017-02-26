package httpfeature

import "fmt"

type HeaderValueIgnoreSymbols struct {
	BaseFeature
	Begin  []rune
	Middle []rune
	End    []rune
}

func (f *HeaderValueIgnoreSymbols) Name() string {
	return "Header value ignored symbols"
}

func (f *HeaderValueIgnoreSymbols) ToString() string {
	return fmt.Sprintf("begin(%s), middle(%s), end(%s)",
		PrintableRunes(f.Begin), PrintableRunes(f.Middle), PrintableRunes(f.End))
}

func (f *HeaderValueIgnoreSymbols) Collect() error {
	f.Begin, _ = f.checkHeaderSymbols(f.BaseRequest,  "x-testsym", "x-testsym","%ctest", "test")
	f.Middle, _ = f.checkHeaderSymbols(f.BaseRequest, "x-testsym", "x-testsym","te%cst", "test")
	f.End, _ = f.checkHeaderSymbols(f.BaseRequest, "x-testsym", "x-testsym","test%c", "test")
	return nil
}
