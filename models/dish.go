package models

import (
	"errors"
	"time"

	"github.com/jeffizhungry/polygon/lib/random"
)

// NOTE(Jeff): Experimenting with decoupling update / creation params
// from the actual model
//
// Inspiration from: https://github.com/stripe/stripe-go
type DishParams struct {
	Name  *string  `json:"name,omitempty"`
	Price *float64 `json:"price,omitempty"`
}

type Dish struct {
	ID    string
	Name  string
	Price float64

	Created time.Time
	Updated time.Time
}

func NewDish(params DishParams) *Dish {
	d := &Dish{
		ID:      random.SecureString(10),
		Created: time.Now(),
		Updated: time.Now(),
	}
	if params.Name != nil {
		d.Name = *params.Name
	}
	if params.Price != nil {
		d.Price = *params.Price
	}
	return d
}

func (d Dish) Validate() error {
	if d.ID == "" {
		return errors.New("Id cannot be empty string")
	}
	if d.Name == "" {
		return errors.New("Name cannot be empty string")
	}
	if d.Price == 0 {
		return errors.New("Price cannot be free")
	}
	return nil
}
