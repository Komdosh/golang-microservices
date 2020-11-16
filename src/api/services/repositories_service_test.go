package services

import (
	"github.com/Komdosh/golang-microservices/src/api/clients/restclient"
	"github.com/Komdosh/golang-microservices/src/api/domain/repositories"
	"github.com/Komdosh/golang-microservices/src/api/utils/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestRepoService_CreateRepoInvalidInputName(t *testing.T) {
	request := repositories.CreateRepoRequest{}

	result, err := RepositoryService.CreateRepo("client_id", request, "")

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid repository name", err.Message())
}

func TestRepoService_CreateRepoErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication","documentation_url": "https://developer.github.com/docs"}`)),
		},
	})
	request := repositories.CreateRepoRequest{Name: "testing"}

	result, err := RepositoryService.CreateRepo("client_id", request, "")
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())
}

func TestRepoService_CreateRepoNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "testing", "owner": {"login": "Komdosh"}}`)),
		},
	})
	request := repositories.CreateRepoRequest{Name: "testing"}

	result, err := RepositoryService.CreateRepo("client_id", request, "")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 123, result.Id)
	assert.EqualValues(t, "testing", result.Name)
	assert.EqualValues(t, "Komdosh", result.Owner)
}

func TestRepoService_CreateRepoNoErrorAuthorizationProvided(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "testing", "owner": {"login": "Komdosh"}}`)),
		},
	})
	request := repositories.CreateRepoRequest{Name: "testing"}

	result, err := RepositoryService.CreateRepo("client_id", request, "token")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 123, result.Id)
	assert.EqualValues(t, "testing", result.Name)
	assert.EqualValues(t, "Komdosh", result.Owner)
}

func TestRepoService_CreateReposConcurrentInvalidRequest(t *testing.T) {
	request := repositories.CreateRepoRequest{}

	output := make(chan repositories.CreateRepositoriesResult)

	service := repoService{}
	go service.createRepoConcurrent(request, "", output)

	result := <-output

	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusBadRequest, result.Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Error.Message())
}

func TestRepoService_CreateReposConcurrentErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication","documentation_url": "https://developer.github.com/docs"}`)),
		},
	})
	request := repositories.CreateRepoRequest{Name: "testing"}

	output := make(chan repositories.CreateRepositoriesResult)

	service := repoService{}
	go service.createRepoConcurrent(request, "", output)

	result := <-output

	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusUnauthorized, result.Error.Status())
	assert.EqualValues(t, "Requires authentication", result.Error.Message())
}

func TestRepoService_CreateReposConcurrentNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "testing", "owner": {"login": "Komdosh"}}`)),
		},
	})
	request := repositories.CreateRepoRequest{Name: "testing"}

	output := make(chan repositories.CreateRepositoriesResult)

	service := repoService{}
	go service.createRepoConcurrent(request, "", output)

	result := <-output

	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Response)
	assert.EqualValues(t, 123, result.Response.Id)
	assert.EqualValues(t, "testing", result.Response.Name)
	assert.EqualValues(t, "Komdosh", result.Response.Owner)
}

func TestRepoService_HandleResults(t *testing.T) {
	input := make(chan repositories.CreateRepositoriesResult)

	output := make(chan repositories.CreateReposResponse)

	var wg sync.WaitGroup

	service := repoService{}
	go service.handleRepoResults(&wg, input, output)

	wg.Add(1)

	go func() {
		input <- repositories.CreateRepositoriesResult{
			Error: errors.NewBadRequestError("invalid repository name"),
		}
	}()

	wg.Wait()
	close(input)

	result := <-output

	assert.NotNil(t, result)
	assert.EqualValues(t, 0, result.StatusCode)

	assert.EqualValues(t, 1, len(result.Results))
	assert.NotNil(t, result.Results[0].Error)
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.Message())
	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.Status())
}

func TestRepoService_CreateReposInvalidRequest(t *testing.T) {
	requests := []repositories.CreateRepoRequest{
		{},
		{Name: ""},
	}

	result, err := RepositoryService.CreateRepos("", requests, "")

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, http.StatusBadRequest, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))
	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.Message())

	assert.EqualValues(t, http.StatusBadRequest, result.Results[1].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[1].Error.Message())
}

func TestRepoService_CreateReposOneSuccessOneFailRequest(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "testing", "owner": {"login": "Komdosh"}}`)),
		},
	})
	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "testing"},
	}

	result, err := RepositoryService.CreateRepos("", requests, "")

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, http.StatusPartialContent, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))

	for _, result := range result.Results {
		if result.Error != nil {
			assert.EqualValues(t, http.StatusBadRequest, result.Error.Status())
			assert.EqualValues(t, "invalid repository name", result.Error.Message())
			continue
		}

		assert.EqualValues(t, 123, result.Response.Id)
		assert.EqualValues(t, "testing", result.Response.Name)
		assert.EqualValues(t, "Komdosh", result.Response.Owner)
	}
}

func TestRepoService_CreateRepoAllSuccessRequest(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "testing", "owner": {"login": "Komdosh"}}`)),
		},
	})
	requests := []repositories.CreateRepoRequest{
		{Name: "testing"},
		{Name: "testing"},
	}

	result, err := RepositoryService.CreateRepos("", requests, "")

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, http.StatusCreated, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))

	assert.Nil(t, result.Results[0].Error)
	assert.EqualValues(t, 123, result.Results[0].Response.Id)
	assert.EqualValues(t, "testing", result.Results[0].Response.Name)
	assert.EqualValues(t, "Komdosh", result.Results[0].Response.Owner)

	assert.Nil(t, result.Results[1].Error)
	assert.EqualValues(t, 123, result.Results[1].Response.Id)
	assert.EqualValues(t, "testing", result.Results[1].Response.Name)
	assert.EqualValues(t, "Komdosh", result.Results[1].Response.Owner)
}
