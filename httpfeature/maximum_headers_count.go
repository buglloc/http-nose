package httpfeature

import (
	"fmt"
)

const MAX_HEADERS_COUNT = 512

type MaximumHeadersCount struct {
	BaseFeature
	Count int
}

func (f *MaximumHeadersCount) Name() string {
	return "Maximum headers count"
}

func (f *MaximumHeadersCount) ToString() string {
	if f.Count == -1 {
		return fmt.Sprintf("%d+", MAX_HEADERS_COUNT)
	}
	return fmt.Sprintf("%d", f.Count)
}

func (f *MaximumHeadersCount) Collect() error {
	f.Count = f.check()
	return nil
}

func (f *MaximumHeadersCount) check() int {
	counts := make([]int, MAX_HEADERS_COUNT)
	sem := make(chan bool, concurrency)
	stop := false
	for i := 0; i < MAX_HEADERS_COUNT && !stop; i++ {
		sem <- true
		go func(cur int) {
			defer func() { <-sem }()

			req := f.BaseRequest.Clone()
			for k := 0; k < cur; k++ {
				req.AddHeader(fmt.Sprintf("X-Foo%d", k), "a")
			}

			resp, err := f.Client.MakeRequest(req)
			if err != nil || resp.Status != 200 {
				counts[cur] = 0
				stop = true
			} else if len(resp.Headers) >= cur {
				counts[cur] = len(resp.Headers)
			} else {
				counts[cur] = 0
				stop = true
			}
		}(i)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	max := 0
	for _, v := range counts {
		if v > max {
			max = v
		}
	}

	if max >= MAX_HEADERS_COUNT {
		return -1
	}
	return max
}
