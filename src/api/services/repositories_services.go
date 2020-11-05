package services

import (
	"github.com/Komdosh/golang-microservices/src/api/config"
	"github.com/Komdosh/golang-microservices/src/api/domain/github"
	"github.com/Komdosh/golang-microservices/src/api/domain/repositories"
	"github.com/Komdosh/golang-microservices/src/api/providers/github_providers"
	"github.com/Komdosh/golang-microservices/src/api/utils/errors"
	"strings"
)

type repoService struct {
}

type reposServicesInterface interface {
	CreateRepo(input repositories.CreateRepoRequest, authorizationHeader string) (*repositories.CreateRepoResponse, errors.ApiError)
}

var (
	RepositoryService reposServicesInterface
)

func init() {
	RepositoryService = &repoService{}
}

func (s *repoService) CreateRepo(input repositories.CreateRepoRequest, authorizationHeader string) (*repositories.CreateRepoResponse, errors.ApiError) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, errors.NewBadRequestError("invalid repository name")
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	token := config.GetGithubAccessToken()
	if strings.TrimSpace(authorizationHeader) != "" {
		token = authorizationHeader
	}

	response, err := github_providers.CreateRepo(token, request)

	if err != nil {
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}

	return &result, nil
}
