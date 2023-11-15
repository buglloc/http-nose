package httpfeature

import (
	"fmt"
	"strings"
)

const MinHeaderLen = 128
const MaxHeaderLen = 10 * 1024

type MaximumHeaderLen struct {
	BaseFeature
	NameLen  int
	ValueLen int
}

func (f *MaximumHeaderLen) Name() string {
	return "Maximum header length"
}

func (f *MaximumHeaderLen) Export() interface{} {
	return map[string]int{
		"Name":  f.NameLen,
		"Value": f.ValueLen,
	}
}

func (f *MaximumHeaderLen) String() string {
	return fmt.Sprintf("Name: %s; Value: %s",
		PrintableBytes(f.NameLen), PrintableBytes(f.ValueLen))
}

func (f *MaximumHeaderLen) Collect() error {
	c := make(chan float64)
	go func() {
		c <- GoldenSectionSearch(MinHeaderLen, MaxHeaderLen, 2, f.goldenNameSize)
	}()
	go func() {
		c <- GoldenSectionSearch(MinHeaderLen, MaxHeaderLen, 2, f.golderValueSize)
	}()
	f.NameLen = int(<-c)
	f.ValueLen = int(<-c)
	return nil
}

func (f *MaximumHeaderLen) goldenNameSize(size float64) float64 {
	if size == MinHeaderLen {
		return size
	}

	req := f.BaseRequest.Clone()
	headerName := strings.Repeat("a", int(size))
	req.AddHeader(headerName, "a")
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return size - 1
	}

	return 0
}

func (f *MaximumHeaderLen) golderValueSize(size float64) float64 {
	if size == MinHeaderLen {
		return size
	}

	req := f.BaseRequest.Clone()
	headerValue := strings.Repeat("a", int(size))
	req.AddHeader("a", headerValue)
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return size - 1
	}

	return 0
}
