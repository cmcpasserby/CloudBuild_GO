package cloudbuild

import (
	"fmt"
	"github.com/cmcpasserby/ucb/pkg/cloudbuild/responses"
)

type ProjectsService struct {
	*client
}

func NewProjectsService(apiKey, orgId string) *ProjectsService {
	return &ProjectsService{
		client: newClient(apiKey, orgId),
	}
}

func (c *ProjectsService) ListAll() ([]responses.Project, error) {
	path := fmt.Sprintf("api/v1/orgs/%s/projects", c.OrgId)

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var projects []responses.Project
	resp, err := c.do(req, &projects)
	if err != nil {
		return nil, err
	}

	fmt.Printf("status: %s\n", resp.Status)

	return projects, nil
}
