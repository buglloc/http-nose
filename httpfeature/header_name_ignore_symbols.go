package httpfeature

import "fmt"

type HeaderNameIgnoreSymbols struct {
	BaseFeature
	Begin  []byte
	Middle []byte
	End    []byte
}

func (f *HeaderNameIgnoreSymbols) Name() string {
	return "Header name ignored symbols"
}

func (f *HeaderNameIgnoreSymbols) Export() interface{} {
	return map[string][]byte{
		"Begin":  f.Begin,
		"Middle": f.Middle,
		"End":    f.End,
	}
}

func (f *HeaderNameIgnoreSymbols) String() string {
	return fmt.Sprintf("begin(%s), middle(%s), end(%s)",
		PrintableSymbols(f.Begin), PrintableSymbols(f.Middle), PrintableSymbols(f.End))
}

func (f *HeaderNameIgnoreSymbols) Collect() error {
	f.Begin, _ = f.checkHeaderSymbols(f.BaseRequest, "{sym}x-testsym", "x-testsym", "test", "test")
	f.Middle, _ = f.checkHeaderSymbols(f.BaseRequest, "x-test{sym}sym", "x-testsym", "test", "test")
	f.End, _ = f.checkHeaderSymbols(f.BaseRequest, "x-testsym{sym}", "x-testsym", "test", "test")
	return nil
}
