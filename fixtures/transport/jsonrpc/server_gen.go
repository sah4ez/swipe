//+build !swipe

// Code generated by Swipe v1.24.0. DO NOT EDIT.

//go:generate swipe
package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/l-vitaly/go-kit/transport/http/jsonrpc"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/swipe-io/swipe/fixtures/service"
	"net/http"
	"strings"
)

func encodeResponseJSONRPCSwipe(_ context.Context, result interface{}) (json.RawMessage, error) {
	b, err := ffjson.Marshal(result)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func MakeSwipeEndpointCodecMap(ep EndpointSet, ns ...string) jsonrpc.EndpointCodecMap {
	var namespace = strings.Join(ns, ".")
	if len(ns) > 0 {
		namespace += "."
	}
	ecm := jsonrpc.EndpointCodecMap{}
	if ep.CreateEndpoint != nil {
		ecm[namespace+"create"] = jsonrpc.EndpointCodec{
			Endpoint: ep.CreateEndpoint,
			Decode:   DecodeUploadFileRequest,
			Encode:   encodeResponseJSONRPCSwipe,
		}
	}
	if ep.DeleteEndpoint != nil {
		ecm[namespace+"delete"] = jsonrpc.EndpointCodec{
			Endpoint: ep.DeleteEndpoint,
			Decode: func(_ context.Context, msg json.RawMessage) (interface{}, error) {
				var req deleteRequestSwipe
				err := ffjson.Unmarshal(msg, &req)
				if err != nil {
					return nil, fmt.Errorf("couldn't unmarshal body to deleteRequestSwipe: %s", err)
				}
				return req, nil
			},
			Encode: encodeResponseJSONRPCSwipe,
		}
	}
	if ep.GetEndpoint != nil {
		ecm[namespace+"get"] = jsonrpc.EndpointCodec{
			Endpoint: ep.GetEndpoint,
			Decode: func(_ context.Context, msg json.RawMessage) (interface{}, error) {
				var req getRequestSwipe
				err := ffjson.Unmarshal(msg, &req)
				if err != nil {
					return nil, fmt.Errorf("couldn't unmarshal body to getRequestSwipe: %s", err)
				}
				return req, nil
			},
			Encode: encodeResponseJSONRPCSwipe,
		}
	}
	if ep.GetAllEndpoint != nil {
		ecm[namespace+"getAll"] = jsonrpc.EndpointCodec{
			Endpoint: ep.GetAllEndpoint,
			Decode: func(_ context.Context, msg json.RawMessage) (interface{}, error) {
				return nil, nil
			},
			Encode: encodeResponseJSONRPCSwipe,
		}
	}
	if ep.TestMethodEndpoint != nil {
		ecm[namespace+"testMethod"] = jsonrpc.EndpointCodec{
			Endpoint: ep.TestMethodEndpoint,
			Decode: func(_ context.Context, msg json.RawMessage) (interface{}, error) {
				var req testMethodRequestSwipe
				err := ffjson.Unmarshal(msg, &req)
				if err != nil {
					return nil, fmt.Errorf("couldn't unmarshal body to testMethodRequestSwipe: %s", err)
				}
				return req, nil
			},
			Encode: encodeResponseJSONRPCSwipe,
		}
	}
	if ep.TestMethod2Endpoint != nil {
		ecm[namespace+"testMethod2"] = jsonrpc.EndpointCodec{
			Endpoint: ep.TestMethod2Endpoint,
			Decode: func(_ context.Context, msg json.RawMessage) (interface{}, error) {
				var req testMethod2RequestSwipe
				err := ffjson.Unmarshal(msg, &req)
				if err != nil {
					return nil, fmt.Errorf("couldn't unmarshal body to testMethod2RequestSwipe: %s", err)
				}
				return req, nil
			},
			Encode: encodeResponseJSONRPCSwipe,
		}
	}
	return ecm
}

// HTTP JSONRPC Transport
func MakeHandlerJSONRPCSwipe(s service.Interface, opts ...SwipeServerOption) (http.Handler, error) {
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
	r := mux.NewRouter()
	handler := jsonrpc.NewServer(MakeSwipeEndpointCodecMap(ep), sopt.genericServerOption...)
	r.Methods("POST").Path("/rpc/{method}").Handler(handler)
	return r, nil
}
