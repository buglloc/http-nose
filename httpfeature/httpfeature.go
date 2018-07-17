package httpfeature

import (
	"github.com/buglloc/http-nose/httpclient"
	"log"
)

type Features struct {
	client                       httpclient.Client
	baseRequest                  httpclient.Request
	baseResponse                 httpclient.Response
	supportedMethods             *SupportedMethods
	multilineHeadersSupport      *MultilineHeadersSupport
	multilineHeadersContinuation *MultilineHeadersContinuation
	providedHeaders              *ProvidedHeaders
	providedHeadersOrder         *ProvidedHeadersOrder
	duplicateHeadersAction       *DuplicateHeaders
	duplicateHost                *DuplicateHost
	headerDelimiters             *HeaderDelimiters
	headerLineDelimiters         *HeaderLineDelimiters
	headerNameSymbols            *HeaderNameSymbols
	headerValueSymbols           *HeaderValueSymbols
	headerNameIgnoreSymbols      *HeaderNameIgnoreSymbols
	headerValueIgnoreSymbols     *HeaderValueIgnoreSymbols
	replaceProvidedHeaders       *ReplaceProvidedHeaders
	absoluteRequestUri           *AbsoluteRequestUri
	headerTransformations        *HeaderTransformations
	requestLineTransformations   *RequestLineTransformations
	supportedVersions            *SupportedVersions
	maximumHeadersCount          *MaximumHeadersCount
	maximumHeaderLen             *MaximumHeaderLen
	headerCountOverflowAction    *HeaderCountOverflowAction
	maximumDuplicateHeadersCount *MaximumDuplicateHeadersCount
}

func NewFeatures(client httpclient.Client, baseRequest httpclient.Request, baseResponse httpclient.Response) *Features {
	return &Features{
		client:       client,
		baseRequest:  baseRequest,
		baseResponse: baseResponse,
	}
}

func (f *Features) GetSupportedMethods() *SupportedMethods {
	if f.supportedMethods == nil {
		log.Print("Colecting SupportedMethods")
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
		log.Print("Colecting MultilineHeadersSupport")
		r := &MultilineHeadersSupport{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.multilineHeadersSupport = r
	}
	return f.multilineHeadersSupport
}

func (f *Features) GetMultilineHeadersContinuation() *MultilineHeadersContinuation {
	if f.multilineHeadersContinuation == nil {
		log.Print("Colecting MultilineHeadersContinuation")
		r := &MultilineHeadersContinuation{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.multilineHeadersContinuation = r
	}
	return f.multilineHeadersContinuation
}

func (f *Features) GetProvidedHeaders() *ProvidedHeaders {
	if f.providedHeaders == nil {
		log.Print("Colecting ProvidedHeaders")
		r := &ProvidedHeaders{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.providedHeaders = r
	}
	return f.providedHeaders
}

func (f *Features) GetProvidedHeadersOrder() *ProvidedHeadersOrder {
	if f.providedHeadersOrder == nil {
		log.Print("Colecting ProvidedHeadersOrder")
		r := &ProvidedHeadersOrder{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.providedHeadersOrder = r
	}
	return f.providedHeadersOrder
}

func (f *Features) GetDuplicateHeaders() *DuplicateHeaders {
	if f.duplicateHeadersAction == nil {
		log.Print("Colecting DuplicateHeaders")
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
		log.Print("Colecting DuplicateHost")
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
		log.Print("Colecting HeaderDelimiters")
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
		log.Print("Colecting HeaderLineDelimiters")
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
		log.Print("Colecting HeaderNameSymbols")
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
		log.Print("Colecting HeaderValueSymbols")
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
		log.Print("Colecting HeaderNameIgnoreSymbols")
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
		log.Print("Colecting HeaderValueIgnoreSymbols")
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
		log.Print("Colecting ReplaceProvidedHeaders")
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
		log.Print("Colecting AbsoluteRequestUri")
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
		log.Print("Colecting HeaderTransformations")
		r := &HeaderTransformations{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.headerTransformations = r
	}
	return f.headerTransformations
}

func (f *Features) GetRequestLineTransformations() *RequestLineTransformations {
	if f.requestLineTransformations == nil {
		log.Print("Colecting RequestLineTransformations")
		r := &RequestLineTransformations{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.requestLineTransformations = r
	}
	return f.requestLineTransformations
}

func (f *Features) GetSupportedVersions() *SupportedVersions {
	if f.supportedVersions == nil {
		log.Print("Colecting SupportedVersions")
		r := &SupportedVersions{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.supportedVersions = r
	}
	return f.supportedVersions
}

func (f *Features) GetHeaderCountOverflowAction() *HeaderCountOverflowAction {
	if f.headerCountOverflowAction == nil {
		log.Print("Colecting HeaderCountOverflowAction")
		r := &HeaderCountOverflowAction{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.headerCountOverflowAction = r
	}
	return f.headerCountOverflowAction
}

func (f *Features) GetMaximumHeadersCount() *MaximumHeadersCount {
	if f.maximumHeadersCount == nil {
		log.Print("Colecting MaximumHeadersCount")
		r := &MaximumHeadersCount{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.maximumHeadersCount = r
	}
	return f.maximumHeadersCount
}

func (f *Features) GetMaximumHeaderLen() *MaximumHeaderLen {
	if f.maximumHeaderLen == nil {
		log.Print("Colecting MaximumHeaderLen")
		r := &MaximumHeaderLen{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.maximumHeaderLen = r
	}
	return f.maximumHeaderLen
}

func (f *Features) GetMaximumDuplicateHeadersCount() *MaximumDuplicateHeadersCount {
	if f.maximumDuplicateHeadersCount == nil {
		log.Print("Colecting MaximumDuplicateHeadersCount")
		r := &MaximumDuplicateHeadersCount{
			BaseFeature: f.newBaseFeature(),
		}
		r.Collect()
		f.maximumDuplicateHeadersCount = r
	}
	return f.maximumDuplicateHeadersCount
}

func (f *Features) Collect() []Feature {
	return []Feature{
		f.GetSupportedVersions(),
		f.GetSupportedMethods(),
		f.GetMultilineHeadersSupport(),
		f.GetMultilineHeadersContinuation(),
		f.GetProvidedHeaders(),
		f.GetReplaceProvidedHeaders(),
		f.GetProvidedHeadersOrder(),
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
		f.GetRequestLineTransformations(),
		f.GetMaximumHeaderLen(),
		f.GetMaximumHeadersCount(),
		f.GetHeaderCountOverflowAction(),
		f.GetMaximumDuplicateHeadersCount(),
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
