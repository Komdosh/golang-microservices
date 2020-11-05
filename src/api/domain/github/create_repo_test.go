package github

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRepoRequestAsJson(t *testing.T) {
	request := CreateRepoRequest{
		Name:        "golang introduction",
		Description: "a golang introduction repository",
		Homepage:    "https://github.com",
		Private:     true,
		HasIssues:   false,
		HasProjects: true,
		HasWiki:     false,
	}

	bytes, err := json.Marshal(request)

	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	assert.EqualValues(t, `{"name":"golang introduction","description":"a golang introduction repository","homepage":"https://github.com","private":true,"has_issues":false,"has_projects":true,"has_wiki":false}`,
		string(bytes))
	fmt.Println(string(bytes))

	var target CreateRepoRequest
	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.EqualValues(t, target.Name, request.Name)
	assert.EqualValues(t, target.Description, request.Description)
	assert.EqualValues(t, target.Homepage, request.Homepage)
	assert.EqualValues(t, target.Private, request.Private)
	assert.EqualValues(t, target.HasProjects, request.HasProjects)
	assert.EqualValues(t, target.HasIssues, request.HasIssues)
	assert.EqualValues(t, target.HasWiki, request.HasWiki)
}
