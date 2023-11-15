package httpfeature

const (
	DuplicateHostNA    = 0
	DuplicateHostFirst = 1
	DuplicateHostLast  = 2
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
	if f.Action == DuplicateHostFirst {
		return "Pick first"
	} else if f.Action == DuplicateHostLast {
		return "Pick last"
	} else if f.Action == DuplicateHostNA {
		return "N/A"
	}
	return "Unknown"
}

func (f *DuplicateHost) Collect() error {
	if f.Features.GetDuplicateHeaders().Action > DuplicateHeadersDisallowed {
		f.Action = f.check()
	} else {
		f.Action = DuplicateHostNA
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
		return DuplicateHostNA
	}

	if resp.Host == "foo-host-first" {
		return DuplicateHostFirst
	} else if resp.Host == "foo-host-last" {
		return DuplicateHostLast
	} else {
		return DuplicateHostNA
	}
}
