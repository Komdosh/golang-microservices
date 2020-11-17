package main

import (
	"bufio"
	"fmt"
	"github.com/Komdosh/golang-microservices/src/api/domain/repositories"
	"github.com/Komdosh/golang-microservices/src/api/services"
	"github.com/Komdosh/golang-microservices/src/api/utils/errors"
	"os"
	"sync"
)

var (
	success = make(map[string]string, 0)
	failed  = make(map[string]errors.ApiError, 0)
)

func getRequests() []repositories.CreateRepoRequest {
	result := make([]repositories.CreateRepoRequest, 0)

	file, err := os.Open("./concurrency/requests.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		request := repositories.CreateRepoRequest{
			Name: line,
		}
		result = append(result, request)
	}

	return result
}

func main() {
	requests := getRequests()

	fmt.Println("about to process requests", len(requests))

	input := make(chan createRepoResult)
	buffer := make(chan bool, 10)

	var wg sync.WaitGroup

	go handleResults(&wg, input)

	for _, request := range requests {
		buffer <- true
		wg.Add(1)
		go createRepo(buffer, input, request)
	}

	wg.Wait()
	close(input)

	fmt.Println("Failed")
	fmt.Println(failed)
	fmt.Println("Success")
	fmt.Println(success)
}

func handleResults(wg *sync.WaitGroup, input chan createRepoResult) {
	for result := range input {
		requestName := result.Request.Name
		if result.Error != nil {
			failed[requestName] = result.Error
		} else {
			success[requestName] = result.Result.Name
		}

		wg.Done()
	}
}

type createRepoResult struct {
	Request repositories.CreateRepoRequest
	Result  *repositories.CreateRepoResponse
	Error   errors.ApiError
}

func createRepo(buffer chan bool, output chan createRepoResult, request repositories.CreateRepoRequest) {
	result, err := services.RepositoryService.CreateRepo("", request, "")
	output <- createRepoResult{
		request,
		result,
		err,
	}
	<-buffer
}
