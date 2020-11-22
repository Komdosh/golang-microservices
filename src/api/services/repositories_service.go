package services

import (
	"github.com/Komdosh/golang-microservices/src/api/config"
	"github.com/Komdosh/golang-microservices/src/api/domain/github"
	"github.com/Komdosh/golang-microservices/src/api/domain/repositories"
	"github.com/Komdosh/golang-microservices/src/api/log/logger-b"
	"github.com/Komdosh/golang-microservices/src/api/providers/github_providers"
	"github.com/Komdosh/golang-microservices/src/api/utils/errors"
	"net/http"
	"strings"
	"sync"
)

type repoService struct {
}

type reposServicesInterface interface {
	CreateRepo(clientId string, input repositories.CreateRepoRequest, authorizationHeader string) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(clientId string, input []repositories.CreateRepoRequest, authorizationHeader string) (repositories.CreateReposResponse, errors.ApiError)
}

var (
	RepositoryService reposServicesInterface
)

func init() {
	RepositoryService = &repoService{}
}

func getAccessToken(authorizationHeader string) string {
	token := config.GetGithubAccessToken()
	if strings.TrimSpace(authorizationHeader) != "" {
		token = authorizationHeader
	}
	return token
}

func (s *repoService) CreateRepo(clientId string, input repositories.CreateRepoRequest, authorizationHeader string) (*repositories.CreateRepoResponse, errors.ApiError) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	token := getAccessToken(authorizationHeader)
	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	logger_b.Info("about to send request to external api", logger_b.Field("client_id:", clientId), logger_b.Field("status", "pending"))
	response, err := github_providers.CreateRepo(token, request)

	if err != nil {
		logger_b.Error("response obtained from external api", err, logger_b.Field("client_id:", clientId), logger_b.Field("status", "error"))
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}
	logger_b.Info("response obtained from external api", logger_b.Field("client_id:", clientId), logger_b.Field("status", "success"))

	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}

	return &result, nil
}

func (s *repoService) CreateRepos(clientId string, request []repositories.CreateRepoRequest, authorizationHeader string) (repositories.CreateReposResponse, errors.ApiError) {
	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	var wg sync.WaitGroup
	go s.handleRepoResults(&wg, input, output)

	for _, current := range request {
		wg.Add(1)
		go s.createRepoConcurrent(current, authorizationHeader, input)
	}

	wg.Wait()

	close(input)

	result := <-output

	successCreation := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successCreation++
		}
	}

	if successCreation == 0 {
		result.StatusCode = result.Results[0].Error.Status()
	} else if successCreation == len(request) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}

	return result, nil
}

func (s *repoService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateRepositoriesResult, output chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	for incomingEvent := range input {
		repoResult := repositories.CreateRepositoriesResult{
			Response: incomingEvent.Response,
			Error:    incomingEvent.Error,
		}
		results.Results = append(results.Results, repoResult)

		wg.Done()
	}
	output <- results
}

func (s *repoService) createRepoConcurrent(input repositories.CreateRepoRequest, authorizationHeader string, output chan repositories.CreateRepositoriesResult) {
	if err := input.Validate(); err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}

	result, err := s.CreateRepo("", input, authorizationHeader)
	if err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}

	output <- repositories.CreateRepositoriesResult{Response: result}
}
