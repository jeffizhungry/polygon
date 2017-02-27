// +build integration
package dishes

import (
	"context"
	"fmt"
	"polymail-api/lib/utils"
	"reflect"
	"testing"

	"github.com/jeffizhungry/polygon/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeString(s string) *string {
	return &s
}

func makeFloat64(v float64) *float64 {
	return &v
}

func TestIntegrationDishesCRUD(t *testing.T) {
	testcases := map[string]struct {
		params        models.DishParams
		updateParams  models.DishParams
		errorExpected bool
	}{
		"the raised the prices!! 😨": {
			params: models.DishParams{
				Name:  makeString("Pasta"),
				Price: makeFloat64(10.0),
			},
			updateParams: models.DishParams{
				Price: makeFloat64(20.0),
			},
		},
		"the unknown dish": {
			params: models.DishParams{
				Price: makeFloat64(10.0),
			},
			errorExpected: true,
		},
	}

	s := NewService()
	for msg, tc := range testcases {

		// Create
		dish, err := s.CreateDish(context.TODO(), tc.params)
		if tc.errorExpected {
			require.Error(t, err, msg)
			continue
		} else {
			require.NoError(t, err, msg)
		}
		require.NotNil(t, dish, msg)

		// Get
		actual, err := s.GetDish(context.TODO(), dish.ID)
		require.NoError(t, err, msg)
		if !assert.True(t, reflect.DeepEqual(dish, actual), msg) {
			fmt.Println("== EXPECTED ==")
			fmt.Printf("%+v\n", dish)
			fmt.Println("== ACTUAL ==")
			fmt.Printf("%+v\n", actual)
		}

		// Update
		updatedDish, err := s.UpdateDish(context.TODO(), dish.ID, tc.updateParams)
		require.NoError(t, err, msg)

		// Get
		actual, err = s.GetDish(context.TODO(), dish.ID)
		require.NoError(t, err, msg)
		if !assert.True(t, reflect.DeepEqual(updatedDish, actual), msg) {
			fmt.Println("== EXPECTED ==")
			fmt.Printf("%+v\n", updatedDish)
			fmt.Println("== ACTUAL ==")
			fmt.Printf("%+v\n", actual)
		}

		// Delete
		err = s.DeleteDish(context.TODO(), dish.ID)
		require.NoError(t, err, msg)

		// Get
		_, err = s.GetDish(context.TODO(), dish.ID)
		require.Equal(t, models.ErrNotFound, err, msg)

		// Delete
		err = s.DeleteDish(context.TODO(), dish.ID)
		require.Equal(t, models.ErrNotFound, err, msg)
	}
}

func TestIntegrationDishesList(t *testing.T) {
	t.Skip()

	// TODO(Jeff): Fix this later 😛
	testcases := map[string]struct {
		count int
	}{
		"max page size - 1": {
			count: maxPageSize - 1,
		},
		"max page size": {
			count: maxPageSize,
		},
		"max page size + 1": {
			count: maxPageSize + 1,
		},
	}

	s := NewService()
	for msg, tc := range testcases {

		// Store set
		var expected []models.Dish
		for i := 0; i < tc.count; i++ {
			dish, err := s.CreateDish(context.TODO(), models.DishParams{
				Name:  makeString("Pasta"),
				Price: makeFloat64(10.0),
			})
			require.NoError(t, err, msg)
			expected = append(expected, *dish)
		}

		// Test different page sizes
		for pagesize := 0; pagesize < maxPageSize; pagesize++ {

			// Paginate
			var actual []models.Dish
			var page []models.Dish
			var offset string

			for once := true; once || len(page) > 0; once = false {
				page, err := s.ListDishes(context.TODO(), offset, pagesize)
				require.NoError(t, err, msg)
				if len(page) > 0 {
					actual = append(actual, page...)
					offset = page[len(page)-1].ID
				}
			}

			// Compare
			if !assert.Equal(t, len(expected), len(actual), msg) {
				continue
			}

			for i := range expected {
				if !assert.True(t, reflect.DeepEqual(expected[i], actual[i]), msg) {
					fmt.Println("== EXPECTED ==")
					fmt.Printf("%+v\n", expected[i])
					fmt.Println("== ACTUAL ==")
					fmt.Printf("%+v\n", actual[i])
					break

					utils.PPrintln("expected: ", expected)
					utils.PPrintln("actual: ", actual)
				}
			}
		}
	}
}
