package httpfeature

type HeaderNameSymbols struct {
	BaseFeature
	Symbols []rune
}

func (f *HeaderNameSymbols) Name() string {
	return "Header name symbols"
}

func (f *HeaderNameSymbols) ToString() string {
	return PrintableRunes(f.Symbols)
}

func (f *HeaderNameSymbols) Collect() error {
	symbols, _ := f.checkHeaderSymbols(f.BaseRequest, "x-test%csym", "x-test%csym","test", "test")
	f.Symbols = append(symbols, AlphaNumSyms...)
	return nil
}
