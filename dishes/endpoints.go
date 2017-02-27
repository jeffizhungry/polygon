// Endpoint creates endpoints mapping requests and responses to service argument
// and return values.
package dishes

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/jeffizhungry/polygon/models"
)

// Endpoints
// Mimicing this: https://github.com/go-kit/kit/blob/master/examples/profilesvc/endpoints.go
//
// Mainly a helper struct for aggregating various endpoints
type Endpoints struct {
	CreateDishEndpoint endpoint.Endpoint
	UpdateDishEndpoint endpoint.Endpoint
	DeleteDishEndpoint endpoint.Endpoint
	GetDishEndpoint    endpoint.Endpoint
	ListDishesEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateDishEndpoint: MakeCreateDishEndpoint(s),
		UpdateDishEndpoint: MakeUpdateDishEndpoint(s),
		DeleteDishEndpoint: MakeDeleteDishEndpoint(s),
		GetDishEndpoint:    MakeGetDishEndpoint(s),
		ListDishesEndpoint: MakeListDishesEndpoint(s),
	}
}

// Translate request payloads to service arguments and
// services return values into response payloads.

type createDishRequest struct {
	models.DishParams
}

type createDishResponse struct {
	*models.Dish
	Err error `json:"error,omitempty"`
}

func MakeCreateDishEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(createDishRequest)
		if !ok {
			return nil, errors.New("programmer error")
		}
		dish, err := s.CreateDish(ctx, req.DishParams)
		resp := createDishResponse{Dish: dish, Err: err}
		return resp, nil
	}
}

type getDishRequest struct {
	ID string `json:"id"`
}

type getDishResponse struct {
	*models.Dish
	Err error `json:"error,omitempty"`
}

func MakeGetDishEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(getDishRequest)
		if !ok {
			return nil, errors.New("programmer error")
		}
		dish, err := s.GetDish(ctx, req.ID)
		resp := getDishResponse{Dish: dish, Err: err}
		return resp, nil
	}
}

type updateDishRequest struct {
	ID string `json:"id"`
	models.DishParams
}

type updateDishResponse struct {
	*models.Dish
	Err error `json:"error,omitempty"`
}

func MakeUpdateDishEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(updateDishRequest)
		if !ok {
			return nil, errors.New("programmer error")
		}
		dish, err := s.UpdateDish(ctx, req.ID, req.DishParams)
		resp := updateDishResponse{Dish: dish, Err: err}
		return resp, nil
	}
}

type deleteDishRequest struct {
	ID string `json:"id"`
}

type deleteDishResponse struct {
	Err error `json:"error,omitempty"`
}

func MakeDeleteDishEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(deleteDishRequest)
		if !ok {
			return nil, errors.New("programmer error")
		}
		err = s.DeleteDish(ctx, req.ID)
		resp := deleteDishResponse{Err: err}
		return resp, nil
	}
}

type listDishesRequest struct {
	Offset   string `json:"offset"`
	PageSize int    `json:"pageSize"`
}

type listDishesResponse struct {
	Dishes []models.Dish `json:"values"`
	Err    error         `json:"error,omitempty"`
}

func MakeListDishesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(listDishesRequest)
		if !ok {
			return nil, errors.New("programmer error")
		}
		dishes, err := s.ListDishes(ctx, req.Offset, req.PageSize)
		resp := listDishesResponse{Dishes: dishes, Err: err}
		return resp, nil
	}
}
