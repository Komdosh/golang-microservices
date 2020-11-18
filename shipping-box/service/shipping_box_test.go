package service

import (
	"github.com/Komdosh/golang-microservices/shipping-box/domain"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGetBestBox_NoBoxes(t *testing.T) {
	availableBoxes := make([]domain.Measured, 0)

	products := []domain.Measured{
		domain.Product{Name: "CPU", Length: 4, Width: 4, Height: 7},
		domain.Product{Name: "GPU", Length: 2, Width: 2, Height: 1},
		domain.Product{Name: "Mother Board", Length: 1, Width: 1, Height: 2},
		domain.Product{Name: "RAM", Length: 1, Width: 5, Height: 4},
	}

	result := GetBestBox(availableBoxes, products)
	assert.Nil(t, result.Box)
	assert.NotNil(t, result.Error)

	assert.EqualValues(t, "No available boxes", result.Error.Message)
}

func TestGetBestBox_NoProductsNoBoxes(t *testing.T) {
	availableBoxes := make([]domain.Measured, 0)
	products := make([]domain.Measured, 0)

	result := GetBestBox(availableBoxes, products)
	assert.Nil(t, result.Box)
	assert.NotNil(t, result.Error)

	assert.EqualValues(t, "No available boxes", result.Error.Message)
}

func TestGetBestBox_NoProducts(t *testing.T) {
	availableBoxes := []domain.Measured{
		domain.Box{Length: 1, Width: 1, Height: 8},
	}
	products := make([]domain.Measured, 0)

	result := GetBestBox(availableBoxes, products)
	assert.NotNil(t, result.Box)
	assert.Nil(t, result.Error)

	assert.EqualValues(t, 1, result.Box.Length)
	assert.EqualValues(t, 1, result.Box.Width)
	assert.EqualValues(t, 8, result.Box.Height)
}

func TestGetBestBox_NoAvailableBoxes(t *testing.T) {
	availableBoxes := []domain.Measured{
		domain.Box{Length: 1, Width: 1, Height: 8},
	}
	products := []domain.Measured{
		domain.Product{Name: "CPU", Length: 4, Width: 4, Height: 7},
		domain.Product{Name: "GPU", Length: 2, Width: 2, Height: 1},
		domain.Product{Name: "Mother Board", Length: 1, Width: 1, Height: 2},
		domain.Product{Name: "RAM", Length: 1, Width: 5, Height: 4},
	}

	result := GetBestBox(availableBoxes, products)
	assert.Nil(t, result.Box)
	assert.NotNil(t, result.Error)

	assert.EqualValues(t, "No available boxes", result.Error.Message)
}

func TestGetBestBox_NoErrors(t *testing.T) {
	availableBoxes := []domain.Measured{
		domain.Box{Length: 5, Width: 9, Height: 8},
		domain.Box{Length: 2, Width: 8, Height: 7},
		domain.Box{Length: 3, Width: 16, Height: 3},
		domain.Box{Length: 4, Width: 4, Height: 8},
		domain.Box{Length: 15, Width: 3, Height: 8},
		domain.Box{Length: 6, Width: 2, Height: 7},
		domain.Box{Length: 7, Width: 6, Height: 5},
		domain.Box{Length: 8, Width: 23, Height: 4},
	}
	products := []domain.Measured{
		domain.Product{Name: "CPU", Length: 4, Width: 4, Height: 7},
		domain.Product{Name: "GPU", Length: 2, Width: 2, Height: 1},
		domain.Product{Name: "Mother Board", Length: 1, Width: 1, Height: 2},
		domain.Product{Name: "RAM", Length: 1, Width: 5, Height: 4},
	}

	result := GetBestBox(availableBoxes, products)

	assert.NotNil(t, result.Box)
	assert.Nil(t, result.Error)

	assert.EqualValues(t, 3, result.Box.Length)
	assert.EqualValues(t, 16, result.Box.Width)
	assert.EqualValues(t, 3, result.Box.Height)
}

func BenchmarkGetBestBox(b *testing.B) {

	availableBoxes := make([]domain.Measured, 0)
	for i := 0; i < 10_000; i++ {
		availableBoxes = append(availableBoxes, domain.Box{Length: rand.Intn(20), Width: rand.Intn(15), Height: rand.Intn(35)})
	}

	products := []domain.Measured{
		domain.Product{Name: "CPU", Length: 4, Width: 4, Height: 7},
		domain.Product{Name: "GPU", Length: 2, Width: 2, Height: 1},
		domain.Product{Name: "Mother Board", Length: 1, Width: 1, Height: 2},
		domain.Product{Name: "RAM", Length: 1, Width: 5, Height: 4},
	}

	for i := 0; i < b.N; i++ {
		GetBestBox(availableBoxes, products)
	}

}
