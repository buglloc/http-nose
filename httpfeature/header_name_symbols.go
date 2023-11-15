package httpfeature

type HeaderNameSymbols struct {
	BaseFeature
	Symbols []byte
}

func (f *HeaderNameSymbols) Name() string {
	return "Header name symbols"
}

func (f *HeaderNameSymbols) Export() interface{} {
	return f.Symbols
}

func (f *HeaderNameSymbols) String() string {
	return PrintableSymbols(f.Symbols)
}

func (f *HeaderNameSymbols) Collect() error {
	symbols, _ := f.checkHeaderSymbols(f.BaseRequest, "x-test{sym}sym", "x-test{sym}sym", "test", "test")
	f.Symbols = append(symbols, AlphaNumSyms...)
	return nil
}
