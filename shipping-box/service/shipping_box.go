package service

import (
	"github.com/Komdosh/golang-microservices/shipping-box/domain"
	"sync"
)

var (
	min = int(^uint(0) >> 1)
)

func GetBestBox(availableBoxes []domain.Measured, products []domain.Measured) domain.BestBoxResult {
	productVolume := calculateVolumeInArray(products)

	var wg sync.WaitGroup
	calculatedBoxFitnessChannel := make(chan domain.CalculatedBoxFitness)
	bestBoxChannel := make(chan domain.BestBoxResult)

	go handleResults(&wg, calculatedBoxFitnessChannel, bestBoxChannel)

	for _, box := range availableBoxes {
		wg.Add(1)

		go processBox(productVolume, box, calculatedBoxFitnessChannel)
	}

	wg.Wait()
	close(calculatedBoxFitnessChannel)

	return <-bestBoxChannel
}

func processBox(volume int, box domain.Measured, input chan domain.CalculatedBoxFitness) {
	boxVolume := calculateVolume(box)

	input <- domain.CalculatedBoxFitness{
		Box:     box.(domain.Box),
		Fitness: boxVolume - volume,
	}
}

func handleResults(wg *sync.WaitGroup, calculatedBoxFitnesses chan domain.CalculatedBoxFitness, boxes chan domain.BestBoxResult) {
	var best domain.Box
	hasResult := false
	for result := range calculatedBoxFitnesses {
		if result.Fitness >= 0 && min > result.Fitness {
			min = result.Fitness
			best = result.Box
			hasResult = true
		}
		wg.Done()
	}

	if hasResult {
		boxes <- domain.BestBoxResult{
			Box: &best,
		}
	} else {
		boxes <- domain.BestBoxResult{
			Error: &domain.Error{
				Message: "No available boxes",
			},
		}
	}
}

func calculateVolumeInArray(MeasuredArray []domain.Measured) int {
	volume := 0
	for _, Measured := range MeasuredArray {
		volume += calculateVolume(Measured)
	}
	return volume
}

func calculateVolume(Measured domain.Measured) int {
	return Measured.GetHeight() * Measured.GetLength() * Measured.GetWidth()
}
