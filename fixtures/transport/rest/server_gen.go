//+build !swipe

// Code generated by Swipe v1.22.4. DO NOT EDIT.

//go:generate swipe
package rest

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	fasthttp2 "github.com/l-vitaly/go-kit/transport/fasthttp"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/swipe-io/swipe/fixtures/service"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
)

type errorWrapper struct {
	Error string `json:"error"`
}

func encodeResponseHTTPSwipe(ctx context.Context, w *fasthttp.Response, response interface{}) error {
	h := w.Header
	h.Set("Content-Iface", "application/json; charset=utf-8")
	if e, ok := response.(endpoint.Failer); ok && e.Failed() != nil {
		data, err := ffjson.Marshal(errorWrapper{Error: e.Failed().Error()})
		if err != nil {
			return err
		}
		w.SetBody(data)
		return nil
	}
	data, err := ffjson.Marshal(response)
	if err != nil {
		return err
	}
	w.SetBody(data)
	return nil
}

// HTTP REST Transport
func MakeHandlerRESTSwipe(s service.Interface, opts ...SwipeServerOption) (fasthttp.RequestHandler, error) {
	sopt := &serverSwipeOpts{}
	for _, o := range opts {
		o(sopt)
	}
	ep := MakeEndpointSet(s)
	ep.CreateEndpoint = middlewareChain(append(sopt.genericEndpointMiddleware, sopt.createEndpointMiddleware...))(ep.CreateEndpoint)
	ep.DeleteEndpoint = middlewareChain(append(sopt.genericEndpointMiddleware, sopt.deleteEndpointMiddleware...))(ep.DeleteEndpoint)
	ep.GetEndpoint = middlewareChain(append(sopt.genericEndpointMiddleware, sopt.getEndpointMiddleware...))(ep.GetEndpoint)
	ep.GetAllEndpoint = middlewareChain(append(sopt.genericEndpointMiddleware, sopt.getAllEndpointMiddleware...))(ep.GetAllEndpoint)
	ep.TestMethodEndpoint = middlewareChain(append(sopt.genericEndpointMiddleware, sopt.testMethodEndpointMiddleware...))(ep.TestMethodEndpoint)
	ep.TestMethod2Endpoint = middlewareChain(append(sopt.genericEndpointMiddleware, sopt.testMethod2EndpointMiddleware...))(ep.TestMethod2Endpoint)
	r := routing.New()
	r.To(fasthttp.MethodPost, "/users", fasthttp2.NewServer(
		ep.CreateEndpoint,
		func(ctx context.Context, r *fasthttp.Request) (interface{}, error) {
			var req createRequestSwipe
			err := ffjson.Unmarshal(r.Body(), &req)
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("couldn't unmarshal body to createRequestSwipe: %s", err)
			}
			return req, nil
		},
		encodeResponseHTTPSwipe,
		append(sopt.genericServerOption, sopt.createServerOption...)...,
	).RouterHandle())
	r.To(fasthttp.MethodPost, "/delete", fasthttp2.NewServer(
		ep.DeleteEndpoint,
		func(ctx context.Context, r *fasthttp.Request) (interface{}, error) {
			var req deleteRequestSwipe
			err := ffjson.Unmarshal(r.Body(), &req)
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("couldn't unmarshal body to deleteRequestSwipe: %s", err)
			}
			return req, nil
		},
		encodeResponseHTTPSwipe,
		append(sopt.genericServerOption, sopt.deleteServerOption...)...,
	).RouterHandle())
	r.To(fasthttp.MethodGet, "/users/<id:[0-9]>/<name:[a-z]>/<fname>", fasthttp2.NewServer(
		ep.GetEndpoint,
		ServerDecodeRequestTest,
		encodeResponseHTTPSwipe,
		append(sopt.genericServerOption, sopt.getServerOption...)...,
	).RouterHandle())
	r.To(fasthttp.MethodGet, "/users", fasthttp2.NewServer(
		ep.GetAllEndpoint,
		func(ctx context.Context, r *fasthttp.Request) (interface{}, error) {
			return nil, nil
		},
		encodeResponseHTTPSwipe,
		append(sopt.genericServerOption, sopt.getAllServerOption...)...,
	).RouterHandle())
	r.To("GET", "/testMethod", fasthttp2.NewServer(
		ep.TestMethodEndpoint,
		func(ctx context.Context, r *fasthttp.Request) (interface{}, error) {
			var req testMethodRequestSwipe
			return req, nil
		},
		encodeResponseHTTPSwipe,
		append(sopt.genericServerOption, sopt.testMethodServerOption...)...,
	).RouterHandle())
	r.To(http.MethodPut, "/<ns>/auth/<utype>/<user>/<restype>/<resource>/<permission>", fasthttp2.NewServer(
		ep.TestMethod2Endpoint,
		func(ctx context.Context, r *fasthttp.Request) (interface{}, error) {
			var req testMethod2RequestSwipe
			err := ffjson.Unmarshal(r.Body(), &req)
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("couldn't unmarshal body to testMethod2RequestSwipe: %s", err)
			}
			vars, ok := ctx.Value(fasthttp2.ContextKeyRouter).(*routing.Context)
			if !ok {
				return nil, fmt.Errorf("couldn't assert fasthttp2.ContextKeyRouter to *routing.Context")
			}
			req.Ns = vars.Param("ns")
			req.Utype = vars.Param("utype")
			req.User = vars.Param("user")
			req.Restype = vars.Param("restype")
			req.Resource = vars.Param("resource")
			req.Permission = vars.Param("permission")
			return req, nil
		},
		encodeResponseHTTPSwipe,
		append(sopt.genericServerOption, sopt.testMethod2ServerOption...)...,
	).RouterHandle())
	return r.HandleRequest, nil
}
