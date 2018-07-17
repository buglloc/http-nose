package httpfeature

const (
	DUPLICATE_HOST_NA    = 0
	DUPLICATE_HOST_FIRST = 1
	DUPLICATE_HOST_LAST  = 2
)

type DuplicateHost struct {
	BaseFeature
	Action int
}

func (f *DuplicateHost) Name() string {
	return "Duplicate host"
}

func (f *DuplicateHost) Export() interface{} {
	return f.String()
}

func (f *DuplicateHost) String() string {
	if f.Action == DUPLICATE_HOST_FIRST {
		return "Pick first"
	} else if f.Action == DUPLICATE_HOST_LAST {
		return "Pick last"
	} else if f.Action == DUPLICATE_HOST_NA {
		return "N/A"
	}
	return "Unknown"
}

func (f *DuplicateHost) Collect() error {
	if f.Features.GetDuplicateHeaders().Action > DUPLICATE_HEADERS_DISALLOWED {
		f.Action = f.check()
	} else {
		f.Action = DUPLICATE_HOST_NA
	}

	return nil
}

func (f *DuplicateHost) check() int {
	req := f.BaseRequest.Clone()
	req.RemoveHeader("Host")
	req.AddHeader("Host", "foo-host-first")
	req.AddHeader("Host", "foo-host-last")
	resp, err := f.Client.MakeRequest(req)
	if err != nil || resp.Status != 200 {
		return DUPLICATE_HOST_NA
	}

	if resp.Host == "foo-host-first" {
		return DUPLICATE_HOST_FIRST
	} else if resp.Host == "foo-host-last" {
		return DUPLICATE_HOST_LAST
	} else {
		return DUPLICATE_HOST_NA
	}
}
