package httpfeature

import (
	"strings"
	"fmt"
	"sync"
)

type RequestUriTransformations struct {
	BaseFeature
	Transforms []string
}

func (f *RequestUriTransformations) Name() string {
	return "RequestUri transformations"
}

func (f *RequestUriTransformations) ToString() string {
	if len(f.Transforms) == 0 {
		return "none"
	}
	return strings.Join(f.Transforms, ", ")
}

func (f *RequestUriTransformations) Collect() error {
	prefix := f.BaseRequest.Path
	f.Transforms = f.check([]string{
		prefix + "/trimed?",
		prefix + "//some/merge",
		prefix + "/%2F%2Fmerge/escaped",
		prefix + "/some/../normalize",
		prefix + "/some%2f%2e%2e%2fnormalize/escaped",
		prefix + "/path?#fragment",
		prefix + "/is/Case/Sensitive?FoO=Bar",
	}, true)

	return nil
}

func (f *RequestUriTransformations) check(toCheck []string, isUri bool) []string {
	result := make([]string, 0)
	mu := &sync.Mutex{}
	sem := make(chan bool, concurrency)
	for _, ch := range toCheck {
		sem <- true
		go func(ch string) {
			defer func() { <-sem }()

			req := f.BaseRequest.Clone()
			req.RequestURI = ch
			resp, err := f.Client.MakeRequest(req)
			if err != nil || resp.Status != 200 {
				return
			}

			var same bool
			if isUri {
				same = len(resp.RequestURI) > 0 && resp.RequestURI != ch
			} else {
				same = len(resp.Path) > 0 && resp.Path != ch
			}

			if same {
				mu.Lock()
				result = append(result, fmt.Sprintf("%q -> %q", ch, resp.RequestURI))
				mu.Unlock()
			}
		}(ch)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return result
}
