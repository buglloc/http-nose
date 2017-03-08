package httpfeature

type HeaderValueSymbols struct {
	BaseFeature
	Symbols []rune
}

func (f *HeaderValueSymbols) Name() string {
	return "Header value symbols"
}

func (f *HeaderValueSymbols) Export() interface{} {
	return f.Symbols
}

func (f *HeaderValueSymbols) String() string {
	return PrintableRunes(f.Symbols)
}

func (f *HeaderValueSymbols) Collect() error {
	symbols, _ := f.checkHeaderSymbols(f.BaseRequest, "x-test", "x-test", "test%csym", "test%csym")
	f.Symbols = append(symbols, AlphaNumSyms...)
	return nil
}
