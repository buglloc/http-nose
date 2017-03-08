package httpfeature

import (
	"fmt"
)

type MaximumDuplicateHeadersCount struct {
	BaseFeature
	Count int
}

func (f *MaximumDuplicateHeadersCount) Name() string {
	return "Maximum duplicate headers count"
}

func (f *MaximumDuplicateHeadersCount) Export() interface{} {
	return f.Count
}

func (f *MaximumDuplicateHeadersCount) String() string {
	if f.Count == 0 {
		return "N/A"
	}

	return fmt.Sprintf("%d", f.Count)
}

func (f *MaximumDuplicateHeadersCount) Collect() error {
	if f.Features.GetDuplicateHeaders().Action >= DUPLICATE_HEADERS_ALLOWED {
		max := f.Features.GetMaximumHeadersCount().Count + 1
		if max == 0 {
			max = MAX_HEADERS_COUNT
		}
		f.Count = f.check(max)
	} else {
		f.Count = 0
	}

	return nil
}

func (f *MaximumDuplicateHeadersCount) check(maxCount int) int {
	counts := make([]int, maxCount)
	sem := make(chan bool, concurrency)
	stop := false
	for i := 0; i < maxCount && !stop; i++ {
		sem <- true
		go func(cur int) {
			defer func() { <-sem }()

			req := f.BaseRequest.Clone()
			for k := 0; k < cur; k++ {
				req.AddHeader("X-Foo", "a")
			}

			resp, err := f.Client.MakeRequest(req)
			if err != nil || resp.Status != 200 {
				counts[cur] = 0
				stop = true
			} else {
				counts[cur] = cur
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
	return max
}
