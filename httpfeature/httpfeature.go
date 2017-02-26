package httpfeature

import (
	"github.com/buglloc/http-nose/httpclient"
)

type Features struct {
	client                   httpclient.Client
	baseRequest              httpclient.Request
	baseResponse             httpclient.Response
	supportedMethods         *SupportedMethods
	multilineHeadersSupport  *MultilineHeadersSupport
	providedHeaders          *ProvidedHeaders
	duplicateHeadersAction   *DuplicateHeaders
	duplicateHost            *DuplicateHost
	headerDelimiters         *HeaderDelimiters
	headerLineDelimiters     *HeaderLineDelimiters
	headerNameSymbols        *HeaderNameSymbols
	headerValueSymbols       *HeaderValueSymbols
	headerNameIgnoreSymbols  *HeaderNameIgnoreSymbols
	headerValueIgnoreSymbols *HeaderValueIgnoreSymbols
	replaceProvidedHeaders   *ReplaceProvidedHeaders
	absoluteRequestUri       *AbsoluteRequestUri
	headerTransformations    *HeaderTransformations
	pathTransformations      *RequestUriTransformations
	supportedVersions        *SupportedVersions
}

func NewFeatures(client httpclient.Client, baseRequest httpclient.Request, baseResponse httpclient.Response) (*Features) {
	return &Features{
		client:       client,
		baseRequest:  baseRequest,
		baseResponse: baseResponse,
	}
}

func (f *Features) GetSupportedMethods() *SupportedMethods {
	if f.supportedMethods == nil {
		r := &SupportedMethods{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.supportedMethods = r
	}
	return f.supportedMethods
}

func (f *Features) GetMultilineHeadersSupport() *MultilineHeadersSupport {
	if f.multilineHeadersSupport == nil {
		r := &MultilineHeadersSupport{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.multilineHeadersSupport = r
	}
	return f.multilineHeadersSupport
}

func (f *Features) GetProvidedHeaders() *ProvidedHeaders {
	if f.providedHeaders == nil {
		r := &ProvidedHeaders{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.providedHeaders = r
	}
	return f.providedHeaders
}

func (f *Features) GetDuplicateHeaders() *DuplicateHeaders {
	if f.duplicateHeadersAction == nil {
		r := &DuplicateHeaders{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.duplicateHeadersAction = r
	}
	return f.duplicateHeadersAction
}

func (f *Features) GetDuplicateHost() *DuplicateHost {
	if f.duplicateHost == nil {
		r := &DuplicateHost{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.duplicateHost = r
	}
	return f.duplicateHost
}

func (f *Features) GetHeaderDelimiters() *HeaderDelimiters {
	if f.headerDelimiters == nil {
		r := &HeaderDelimiters{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.headerDelimiters = r
	}
	return f.headerDelimiters
}

func (f *Features) GetHeaderLineDelimiters() *HeaderLineDelimiters {
	if f.headerLineDelimiters == nil {
		r := &HeaderLineDelimiters{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.headerLineDelimiters = r
	}
	return f.headerLineDelimiters
}

func (f *Features) GetHeaderNameSymbols() *HeaderNameSymbols {
	if f.headerNameSymbols == nil {
		r := &HeaderNameSymbols{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.headerNameSymbols = r
	}
	return f.headerNameSymbols
}

func (f *Features) GetHeaderValueSymbols() *HeaderValueSymbols {
	if f.headerValueSymbols == nil {
		r := &HeaderValueSymbols{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.headerValueSymbols = r
	}
	return f.headerValueSymbols
}

func (f *Features) GetHeaderNameIgnoreSymbols() *HeaderNameIgnoreSymbols {
	if f.headerNameIgnoreSymbols == nil {
		r := &HeaderNameIgnoreSymbols{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.headerNameIgnoreSymbols = r
	}
	return f.headerNameIgnoreSymbols
}

func (f *Features) GetHeaderValueIgnoreSymbols() *HeaderValueIgnoreSymbols {
	if f.headerValueIgnoreSymbols == nil {
		r := &HeaderValueIgnoreSymbols{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.headerValueIgnoreSymbols = r
	}
	return f.headerValueIgnoreSymbols
}

func (f *Features) GetReplaceProvidedHeaders() *ReplaceProvidedHeaders {
	if f.replaceProvidedHeaders == nil {
		r := &ReplaceProvidedHeaders{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.replaceProvidedHeaders = r
	}
	return f.replaceProvidedHeaders
}

func (f *Features) GetAbsoluteRequestUri() *AbsoluteRequestUri {
	if f.absoluteRequestUri == nil {
		r := &AbsoluteRequestUri{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.absoluteRequestUri = r
	}
	return f.absoluteRequestUri
}

func (f *Features) GetHeaderTransformations() *HeaderTransformations {
	if f.headerTransformations == nil {
		r := &HeaderTransformations{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.headerTransformations = r
	}
	return f.headerTransformations
}

func (f *Features) GetRequestUriTransformations() *RequestUriTransformations {
	if f.pathTransformations == nil {
		r := &RequestUriTransformations{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.pathTransformations = r
	}
	return f.pathTransformations
}

func (f *Features) GetSupportedVersions() *SupportedVersions {
	if f.supportedVersions == nil {
		r := &SupportedVersions{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.supportedVersions = r
	}
	return f.supportedVersions
}

func (f *Features) Collect() []Feature {
	return []Feature{
		f.GetSupportedVersions(),
		f.GetSupportedMethods(),
		f.GetMultilineHeadersSupport(),
		f.GetProvidedHeaders(),
		f.GetReplaceProvidedHeaders(),
		f.GetDuplicateHeaders(),
		f.GetDuplicateHost(),
		f.GetHeaderDelimiters(),
		f.GetHeaderLineDelimiters(),
		f.GetHeaderNameSymbols(),
		f.GetHeaderValueSymbols(),
		f.GetHeaderNameIgnoreSymbols(),
		f.GetHeaderValueIgnoreSymbols(),
		f.GetAbsoluteRequestUri(),
		f.GetHeaderTransformations(),
		f.GetRequestUriTransformations(),
	}
}

func (f *Features) newBaseFeature() BaseFeature {
	return BaseFeature{
		Client:       f.client,
		BaseRequest:  f.baseRequest,
		BaseResponse: f.baseResponse,
		Features:     f,
	}
}
