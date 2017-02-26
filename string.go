package main

import "strings"

// StringService provides string operation
type StringService interface {
	ToUpper(string) (string, error)
	ToLower(string) (string, error)
	Length(string) int
}

func NewStringService() StringService {
	return &stringResource{}
}

// stringResource implements the string service
type stringResource struct{}

type ToUpperRequest struct {
	S string `json:"s"`
}

type ToUpperResponse struct {
	S   string `json:"s"`
	Err string `json:"error,omitempty"`
}

func (r *stringResource) ToUpper(s string) (string, error) {
	return strings.ToUpper(s), nil
}

type ToLowerRequest struct {
	S string `json:"s"`
}

type ToLowerResponse struct {
	S   string `json:"s"`
	Err string `json:"error,omitempty"`
}

func (r *stringResource) ToLower(s string) (string, error) {
	return strings.ToLower(s), nil
}

type LengthRequest struct {
	S string `json:"s"`
}

type LengthResponse struct {
	Length int    `json:"length"`
	Err    string `json:"error,omitempty"`
}

func (r *stringResource) Length(s string) int {
	return len(s)
}
