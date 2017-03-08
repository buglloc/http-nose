package httpfeature

import (
	"strings"
	"fmt"
	"sync"
)

type RequestLineTransformations struct {
	BaseFeature
	UriTransforms []string
	PathTransforms []string
}

func (f *RequestLineTransformations) Name() string {
	return "RequestLine transformations"
}

func (f *RequestLineTransformations) Export() interface{} {
	return map[string][]string {
		"Uri": f.UriTransforms,
		"Path": f.PathTransforms,
	}
}

func (f *RequestLineTransformations) String() string {
	return fmt.Sprintf("Uri: %s; Path: %s",
		f.formatTransform(f.UriTransforms), f.formatTransform(f.PathTransforms))
}

func (f *RequestLineTransformations) formatTransform(trans []string) string {
	if len(trans) == 0 {
		return "none"
	}
	return strings.Join(trans, ", ")
}

func (f *RequestLineTransformations) Collect() error {
	f.UriTransforms, f.PathTransforms = f.check([][]string{
		{"merge//slashes", ""},
		{"merge%2F%2Fescaped", ""},
		{"simple/../normalize", ""},
		{"escapeddot/%2e%2e/normalize", ""},
		{"escapedall%2f%2e%2e%2fnormalize", ""},
		{"fragment", "#fragment"},
		{"is/Case/Sensetive", "FoO=Bar"},
	})

	return nil
}

func (f *RequestLineTransformations) check(toCheck [][]string) ([]string, []string) {
	uriTrans := make([]string, 0)
	pathTrans := make([]string, 0)
	mu := &sync.Mutex{}
	sem := make(chan bool, concurrency)
	for _, ch := range toCheck {
		sem <- true
		go func(path, args string) {
			defer func() { <-sem }()

			req := f.BaseRequest.Clone()
			req.Path += path
			req.Args = path

			resp, err := f.Client.MakeRequest(req)
			if err != nil || resp.Status != 200 {
				return
			}

			requestPath := req.Path
			requestUri := req.BuildRequestUri()

			if resp.RequestURI != "" && requestUri != resp.RequestURI {
				mu.Lock()
				uriTrans = append(uriTrans, fmt.Sprintf("%q -> %q", requestUri, resp.RequestURI))
				mu.Unlock()
			}
			if resp.Path != "" && requestPath != resp.Path {
				mu.Lock()
				pathTrans = append(pathTrans, fmt.Sprintf("%q -> %q", req.Path, resp.Path))
				mu.Unlock()
			}
		}(ch[0], ch[1])
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return uriTrans, pathTrans
}
