package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// swagger:route GET / Index index
// Returns OK if there's no problem
// responses:
//	200: OK

// Index returns OK handles GET requests
func (p *APIContext) Index(rw http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	spanname := "docmanAPI.Index"
	var span opentracing.Span

	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		// The method is called without a span context in the http header.
		span = tracer.StartSpan(spanname)
	} else {
		// Create the span referring to the RPC client if available.
		// If wireContext == nil, a root span will be created.
		span = opentracing.StartSpan(spanname, ext.RPCServerOption(wireContext))
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, r.URL.RequestURI())
	ext.HTTPMethod.Set(span, r.Method)
	defer span.Finish()

	rw.WriteHeader(200)
}

// swagger:route GET /version Index version
// Returns version information
// responses:
//	200: OK

// Version returns the version info for the service by reading from a static file
func (p *APIContext) Version(rw http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	spanname := "docmanAPI.Version"
	var span opentracing.Span

	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		// The method is called without a span context in the http header.
		span = tracer.StartSpan(spanname)
	} else {
		// Create the span referring to the RPC client if available.
		// If wireContext == nil, a root span will be created.
		span = opentracing.StartSpan(spanname, ext.RPCServerOption(wireContext))
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, r.URL.RequestURI())
	ext.HTTPMethod.Set(span, r.Method)
	defer span.Finish()
	dat, err := ioutil.ReadFile("./static/version.txt")
	if err != nil {
		dat = append(dat, '0')
	}
	fmt.Fprintf(rw, "Welcome to docman API! Version:%s", dat)
}

// CorsHandler swagger:route OPTIONS /
//
// # Handler to respond to CORS preflight requests
//
// Responses:
//
//	200: OK
func CorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, lang")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, OPTIONS")
	w.WriteHeader(200)
}
