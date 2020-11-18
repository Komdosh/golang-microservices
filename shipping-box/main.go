package main

import (
	"fmt"
	"github.com/Komdosh/golang-microservices/shipping-box/domain"
	"github.com/Komdosh/golang-microservices/shipping-box/service"
)

func main() {
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

	result := service.GetBestBox(availableBoxes, products)

	fmt.Printf("Best box is: %v", result)
}
