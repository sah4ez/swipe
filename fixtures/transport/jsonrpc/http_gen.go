//+build !swipe

// Code generated by Swipe v1.22.4. DO NOT EDIT.

//go:generate swipe
package jsonrpc

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/l-vitaly/go-kit/transport/http/jsonrpc"
	"github.com/swipe-io/swipe/fixtures/user"
	"net/http"
)

type httpError struct {
	code int
}

func (e httpError) Error() string {
	return http.StatusText(e.code)
}
func (e httpError) StatusCode() int {
	return e.code
}
func ErrorDecode(code int) (_ error) {
	switch code {
	default:
		return httpError{code: code}
	case -32002:
		return user.ErrForbidden{}
	case -32001:
		return user.ErrUnauthorized{}
	}
}

func middlewareChain(middlewares []endpoint.Middleware) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		if len(middlewares) == 0 {
			return next
		}
		outer := middlewares[0]
		others := middlewares[1:]
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}

type SwipeServerOption func(*serverSwipeOpts)
type serverSwipeOpts struct {
	genericServerOption           []jsonrpc.ServerOption
	genericEndpointMiddleware     []endpoint.Middleware
	createServerOption            []jsonrpc.ServerOption
	createEndpointMiddleware      []endpoint.Middleware
	deleteServerOption            []jsonrpc.ServerOption
	deleteEndpointMiddleware      []endpoint.Middleware
	getServerOption               []jsonrpc.ServerOption
	getEndpointMiddleware         []endpoint.Middleware
	getAllServerOption            []jsonrpc.ServerOption
	getAllEndpointMiddleware      []endpoint.Middleware
	testMethodServerOption        []jsonrpc.ServerOption
	testMethodEndpointMiddleware  []endpoint.Middleware
	testMethod2ServerOption       []jsonrpc.ServerOption
	testMethod2EndpointMiddleware []endpoint.Middleware
}

func SwipeGenericServerOptions(v ...jsonrpc.ServerOption) (_ SwipeServerOption) {
	return func(o *serverSwipeOpts) { o.genericServerOption = v }
}

func SwipeGenericServerEndpointMiddlewares(v ...endpoint.Middleware) (_ SwipeServerOption) {
	return func(o *serverSwipeOpts) { o.genericEndpointMiddleware = v }
}

func SwipeCreateServerOptions(opt ...jsonrpc.ServerOption) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.createServerOption = opt }
}

func SwipeCreateServerEndpointMiddlewares(opt ...endpoint.Middleware) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.createEndpointMiddleware = opt }
}

func SwipeDeleteServerOptions(opt ...jsonrpc.ServerOption) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.deleteServerOption = opt }
}

func SwipeDeleteServerEndpointMiddlewares(opt ...endpoint.Middleware) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.deleteEndpointMiddleware = opt }
}

func SwipeGetServerOptions(opt ...jsonrpc.ServerOption) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.getServerOption = opt }
}

func SwipeGetServerEndpointMiddlewares(opt ...endpoint.Middleware) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.getEndpointMiddleware = opt }
}

func SwipeGetAllServerOptions(opt ...jsonrpc.ServerOption) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.getAllServerOption = opt }
}

func SwipeGetAllServerEndpointMiddlewares(opt ...endpoint.Middleware) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.getAllEndpointMiddleware = opt }
}

func SwipeTestMethodServerOptions(opt ...jsonrpc.ServerOption) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.testMethodServerOption = opt }
}

func SwipeTestMethodServerEndpointMiddlewares(opt ...endpoint.Middleware) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.testMethodEndpointMiddleware = opt }
}

func SwipeTestMethod2ServerOptions(opt ...jsonrpc.ServerOption) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.testMethod2ServerOption = opt }
}

func SwipeTestMethod2ServerEndpointMiddlewares(opt ...endpoint.Middleware) (_ SwipeServerOption) {
	return func(c *serverSwipeOpts) { c.testMethod2EndpointMiddleware = opt }
}
