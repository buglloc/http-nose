package httpfeature

const (
	HTTP_VERSION_NONE = 1 << 0
	HTTP_VERSION_09   = 1 << 1
	HTTP_VERSION_10   = 1 << 1
	HTTP_VERSION_11   = 1 << 1
	HTTP_VERSION_ANY  = 1 << 1
)

type SupportedVersions struct {
	BaseFeature
	Supported int
}

func (f *SupportedVersions) Name() string {
	return "Supported HTTP versions"
}

func (f *SupportedVersions) ToString() string {
	result := make([]string, 0)
	if (f.Supported & HTTP_VERSION_NONE) != 0 {
		result = append(result, "None")
	}
	if (f.Supported & HTTP_VERSION_09) != 0 {
		result = append(result, "0.9")
	}
	if (f.Supported & HTTP_VERSION_10) != 0 {
		result = append(result, "1.0")
	}
	if (f.Supported & HTTP_VERSION_11) != 0 {
		result = append(result, "1.1")
	}
	if (f.Supported & HTTP_VERSION_ANY) != 0 {
		result = append(result, "xxx")
	}

	return PrintableStrings(result)
}

func (f *SupportedVersions) Collect() error {
	if f.checkVersion("") {
		f.Supported |= HTTP_VERSION_NONE
	}
	if f.checkVersion("HTTP/0.9") {
		f.Supported |= HTTP_VERSION_09
	}
	if f.checkVersion("HTTP/1.0") {
		f.Supported |= HTTP_VERSION_10
	}
	if f.checkVersion("HTTP/1.1") {
		f.Supported |= HTTP_VERSION_11
	}
	if f.checkVersion("HTTP/333") {
		f.Supported |= HTTP_VERSION_ANY
	}
	return nil
}

func (f *SupportedVersions) checkVersion(version string) bool {
	req := f.BaseRequest.Clone()
	req.Proto = version
	resp, err := f.Client.MakeRequest(req)
	return err == nil && resp.Status == 200
}