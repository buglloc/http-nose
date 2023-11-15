package httpfeature

import "fmt"

type HeaderValueIgnoreSymbols struct {
	BaseFeature
	Begin  []byte
	Middle []byte
	End    []byte
}

func (f *HeaderValueIgnoreSymbols) Name() string {
	return "Header value ignored symbols"
}

func (f *HeaderValueIgnoreSymbols) Export() interface{} {
	return map[string][]byte{
		"Begin":  f.Begin,
		"Middle": f.Middle,
		"End":    f.End,
	}
}

func (f *HeaderValueIgnoreSymbols) String() string {
	return fmt.Sprintf("begin(%s), middle(%s), end(%s)",
		PrintableSymbols(f.Begin), PrintableSymbols(f.Middle), PrintableSymbols(f.End))
}

func (f *HeaderValueIgnoreSymbols) Collect() error {
	f.Begin, _ = f.checkHeaderSymbols(f.BaseRequest, "x-testsym", "x-testsym", "{sym}test", "test")
	f.Middle, _ = f.checkHeaderSymbols(f.BaseRequest, "x-testsym", "x-testsym", "te{sym}st", "test")
	f.End, _ = f.checkHeaderSymbols(f.BaseRequest, "x-testsym", "x-testsym", "test{sym}", "test")
	return nil
}
