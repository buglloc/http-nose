package httpfeature

import (
	"strings"
)

const (
	DUPLICATE_HEADERS_NA         = 0
	DUPLICATE_HEADERS_DISALLOWED = 1
	DUPLICATE_HEADERS_ALLOWED    = 2
	DUPLICATE_HEADERS_MERGED     = 3
	DUPLICATE_HEADERS_PICK_FIRST = 4
	DUPLICATE_HEADERS_PICK_LAST  = 5
)

type DuplicateHeaders struct {
	BaseFeature
	Action int
}

func (f *DuplicateHeaders) Name() string {
	return "Duplicate headers"
}

func (f *DuplicateHeaders) ToString() string {
	if f.Action == DUPLICATE_HEADERS_DISALLOWED {
		return "disallowed"
	} else if f.Action == DUPLICATE_HEADERS_ALLOWED {
		return "allowed"
	} else if f.Action == DUPLICATE_HEADERS_MERGED {
		return "merged"
	} else if f.Action == DUPLICATE_HEADERS_PICK_FIRST {
		return "pick first"
	} else if f.Action == DUPLICATE_HEADERS_PICK_LAST {
		return "pick last"
	} else if f.Action == DUPLICATE_HEADERS_NA {
		return "n/a"
	}
	return "unknown"
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
		return DUPLICATE_HEADERS_DISALLOWED
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
		return DUPLICATE_HEADERS_MERGED
	} else if first && second {
		return DUPLICATE_HEADERS_ALLOWED
	} else if first {
		return DUPLICATE_HEADERS_PICK_FIRST
	} else if second {
		return DUPLICATE_HEADERS_PICK_LAST
	}
	return DUPLICATE_HEADERS_NA
}
