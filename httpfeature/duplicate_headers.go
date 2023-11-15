package httpfeature

import (
	"strings"
)

const (
	DuplicateHeadersNA         = 0
	DuplicateHeadersDisallowed = 1
	DuplicateHeadersAllowed    = 2
	DuplicateHeadersMerged     = 3
	DuplicateHeadersPickFirst  = 4
	DuplicateHeadersPickLast   = 5
)

type DuplicateHeaders struct {
	BaseFeature
	Action int
}

func (f *DuplicateHeaders) Name() string {
	return "Duplicate headers"
}

func (f *DuplicateHeaders) Export() interface{} {
	return f.String()
}

func (f *DuplicateHeaders) String() string {
	if f.Action == DuplicateHeadersDisallowed {
		return "Disallowed"
	} else if f.Action == DuplicateHeadersAllowed {
		return "Allowed"
	} else if f.Action == DuplicateHeadersMerged {
		return "Merged"
	} else if f.Action == DuplicateHeadersPickFirst {
		return "Pick first"
	} else if f.Action == DuplicateHeadersPickLast {
		return "Pick last"
	} else if f.Action == DuplicateHeadersNA {
		return "N/A"
	}
	return "Unknown"
}

func (f *DuplicateHeaders) Collect() error {
	f.Action = f.check()
	return nil
}

func (f *DuplicateHeaders) check() int {
	req := f.BaseRequest.Clone()
	req.AddHeader("X-Foo", "foo1")
	req.AddHeader("X-Foo", "foo2")
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return DuplicateHeadersDisallowed
	}

	first := false
	second := false
	merged := false
	for _, h := range resp.HeadersSlice("X-Foo") {
		f := strings.Contains(h.Value, "foo1")
		s := strings.Contains(h.Value, "foo2")
		if f && s {
			merged = true
			break
		} else if f {
			first = true
		} else if s {
			second = true
		}
	}

	if merged {
		return DuplicateHeadersMerged
	} else if first && second {
		return DuplicateHeadersAllowed
	} else if first {
		return DuplicateHeadersPickFirst
	} else if second {
		return DuplicateHeadersPickLast
	}
	return DuplicateHeadersNA
}
