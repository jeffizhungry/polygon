package main

import (
	"context"
	"encoding/json"
	"net/http"

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

func main() {

	// Initialize services and inject dependencies
	svc := NewStringService()

	// Initialize endpoints
	toLowerHandler := httptransport.NewServer(
		context.Background(),
		makeToLowerEndpoint(svc),
		decodeToLowerRequest,
		encodeResponse,
	)

	toUpperHandler := httptransport.NewServer(
		context.Background(),
		makeToUpperEndpoint(svc),
		decodeToUpperRequest,
		encodeResponse,
	)

	lengthHandler := httptransport.NewServer(
		context.Background(),
		makeLengthEndpoint(svc),
		decodeLengthRequest,
		encodeResponse,
	)

	// Register endpoints
	http.Handle("/toLower", toLowerHandler)
	http.Handle("/toUpper", toUpperHandler)
	http.Handle("/length", lengthHandler)

	// Start server
	logrus.Info("Listening on...  localhost:8008")
	logrus.Fatal(http.ListenAndServe(":8008", nil))
}
