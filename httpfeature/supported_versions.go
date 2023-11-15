package httpfeature

const (
	HttpVersionNone = 1 << 0
	HttpVersion09   = 1 << 1
	HttpVersion10   = 1 << 2
	HttpVersion11   = 1 << 3
	HttpVersionAny  = 1 << 4
)

type SupportedVersions struct {
	BaseFeature
	Supported int
}

func (f *SupportedVersions) Name() string {
	return "Supported HTTP versions"
}

func (f *SupportedVersions) Versions() []string {
	result := make([]string, 0)
	if (f.Supported & HttpVersionNone) != 0 {
		result = append(result, "None")
	}
	if (f.Supported & HttpVersion09) != 0 {
		result = append(result, "0.9")
	}
	if (f.Supported & HttpVersion10) != 0 {
		result = append(result, "1.0")
	}
	if (f.Supported & HttpVersion11) != 0 {
		result = append(result, "1.1")
	}
	if (f.Supported & HttpVersionAny) != 0 {
		result = append(result, "xxx")
	}

	return result
}

func (f *SupportedVersions) Export() interface{} {
	return f.Versions()
}

func (f *SupportedVersions) String() string {
	return PrintableStrings(f.Versions())
}

func (f *SupportedVersions) Collect() error {
	if f.checkVersion("") {
		f.Supported |= HttpVersionNone
	}
	if f.checkVersion("HTTP/0.9") {
		f.Supported |= HttpVersion09
	}
	if f.checkVersion("HTTP/1.0") {
		f.Supported |= HttpVersion10
	}
	if f.checkVersion("HTTP/1.1") {
		f.Supported |= HttpVersion11
	}
	if f.checkVersion("HTTP/333") {
		f.Supported |= HttpVersionAny
	}
	return nil
}

func (f *SupportedVersions) checkVersion(version string) bool {
	req := f.BaseRequest.Clone()
	req.Proto = version
	resp, err := f.Client.MakeRequest(req)
	return err == nil && resp.Status == 200
}
