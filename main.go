package main

import (
	"context"
	"encoding/json"
	"net/http"
	"polymail-api/config"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

/**************************************
 * Endpoints
 *	- map requests to service
 *************************************/

func makeToLowerEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ToLowerRequest)
		ans, err := svc.ToLower(req.S)
		if err != nil {
			return ToLowerResponse{Err: err.Error()}, nil
		}
		return ToLowerResponse{S: ans}, nil
	}
}

func makeToUpperEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ToUpperRequest)
		ans, err := svc.ToUpper(req.S)
		if err != nil {
			return ToUpperResponse{Err: err.Error()}, nil
		}
		return ToUpperResponse{S: ans}, nil
	}
}

func makeLengthEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LengthRequest)
		ans := svc.Length(req.S)
		return LengthResponse{Length: ans}, nil
	}
}

/**************************************
 * Translate http requests to service
 * inputs and service outputs to http
 * responses.
 *************************************/

func decodeToLowerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req ToLowerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeToUpperRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req ToUpperRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLengthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req LengthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	return json.NewEncoder(w).Encode(resp)
}

/**************************************
 * Middleware
 *************************************/

func poorManMetricsMiddleware(route string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			defer func(begin time.Time) {
				logrus.WithField("route", route).Infof("Duration = %v", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func main() {

	// Initialize services and inject dependencies
	svc := NewStringService()

	// Initialize endpoints
	toLowerEndpoint := makeToLowerEndpoint(svc)
	toLowerEndpoint = poorManMetricsMiddleware(toLowerEndpoint)
	toLowerHandler := httptransport.NewServer(
		context.Background(),
		toLowerEndpoint,
		decodeToLowerRequest,
		encodeResponse,
	)

	toUpperEndpoint := makeToUpperEndpoint(svc)
	toUpperEndpoint = poorManMetricsMiddleware(toUpperEndpoint)
	toUpperHandler := httptransport.NewServer(
		context.Background(),
		toUpperEndpoint,
		decodeToUpperRequest,
		encodeResponse,
	)

	lengthEndpoint := makeLengthEndpoint(svc)
	lengthEndpoint = poorManMetricsMiddleware(lengthEndpoint)
	lengthHandler := httptransport.NewServer(
		context.Background(),
		lengthEndpoint,
		decodeLengthRequest,
		encodeResponse,
	)

	// Register endpoints
	http.Handle("/toLower", toLowerHandler)
	http.Handle("/toUpper", toUpperHandler)
	http.Handle("/length", lengthHandler)

	// Start server
	logrus.Info("Listening on...  %v", config.Server.Address())
	logrus.Fatal(http.ListenAndServe(config.Server.Address(), nil))
}
