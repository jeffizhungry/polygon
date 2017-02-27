// Service implements the buisness logic for a service
package dishes

import (
	"context"
	"sync"

	"github.com/jeffizhungry/polygon/models"
	"github.com/pkg/errors"
)

const (
	maxPageSize = 5
)

// NOTE(Jeff): The goal of this interface is to provide a standard interface
// for implementing the service and building a client library for communicating
// with this service.
type Service interface {
	CreateDish(ctx context.Context, d models.DishParams) (*models.Dish, error)
	GetDish(ctx context.Context, id string) (*models.Dish, error)
	UpdateDish(ctx context.Context, id string, d models.DishParams) (*models.Dish, error)
	DeleteDish(ctx context.Context, id string) error
	ListDishes(ctx context.Context, offset string, limit int) ([]models.Dish, error)
}

func NewService() Service {
	return &resource{
		local: make(map[string]*models.Dish),
		mu:    &sync.RWMutex{},
	}
}

type resource struct {
	local          map[string]*models.Dish
	secondaryIndex []models.Dish
	mu             *sync.RWMutex
}

func (r *resource) CreateDish(ctx context.Context, d models.DishParams) (*models.Dish, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create model
	dish := models.NewDish(d)

	// Validate
	if err := dish.Validate(); err != nil {
		return nil, err
	}

	// Save model
	r.local[dish.ID] = dish
	r.secondaryIndex = append(r.secondaryIndex, *dish)
	return dish, nil
}

func (r *resource) GetDish(ctx context.Context, id string) (*models.Dish, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Get model
	dish, found := r.local[id]
	if !found {
		return nil, models.ErrNotFound
	}
	return dish, nil
}

func (r *resource) UpdateDish(ctx context.Context, id string, params models.DishParams) (*models.Dish, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get model
	dish, found := r.local[id]
	if !found {
		return nil, models.ErrNotFound
	}

	// Update model
	if params.Name != nil {
		dish.Name = *params.Name
	}
	if params.Price != nil {
		dish.Price = *params.Price
	}

	// Validate
	if err := dish.Validate(); err != nil {
		return nil, err
	}

	// Update local
	r.local[id] = dish

	// Update secondary
	for i := range r.secondaryIndex {
		if r.secondaryIndex[i].ID == id {
			r.secondaryIndex[i] = *dish
		}
	}
	return dish, nil
}

func (r *resource) DeleteDish(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if it exists
	_, found := r.local[id]
	if !found {
		return models.ErrNotFound
	}

	// Delete from local
	delete(r.local, id)

	// Delete from secondary
	for i := range r.secondaryIndex {
		if r.secondaryIndex[i].ID == id {
			r.secondaryIndex = append(r.secondaryIndex[:i], r.secondaryIndex[i+1:]...)
		}
	}
	return nil
}

func (r *resource) ListDishes(ctx context.Context, offset string, limit int) ([]models.Dish, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit > maxPageSize {
		return nil, errors.Errorf("max page size is ", maxPageSize)
	}

	var set []models.Dish
	if offset == "" {
		set = r.secondaryIndex
	} else {
		for i := range r.secondaryIndex {
			if r.secondaryIndex[i].ID == offset {
				set = r.secondaryIndex[i+1:]
				break
			}
		}
	}

	// Limit set to page size
	if len(set) > limit {
		return set[:limit], nil
	}
	return set, nil
}
