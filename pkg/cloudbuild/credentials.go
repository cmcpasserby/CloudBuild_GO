package cloudbuild

import (
	"fmt"
	"github.com/cmcpasserby/CloudBuild_GO/pkg/cloudbuild/responses"
	"io"
	"net/http"
	"os"
	"strings"
)

type CredentialsService struct {
	*client
}

func NewCredentialsService(apiKey, orgId string) *CredentialsService {
	return &CredentialsService{
		client: newClient(apiKey, orgId),
	}
}

func (c *CredentialsService) GetIOS(projectId, credId string) (*responses.IOSCred, error) {
	path := fmt.Sprintf("api/v1/orgs/%s/projects/%s/credentials/signing/ios/%s", c.OrgId, projectId, credId)

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var credential responses.IOSCred
	resp, err := c.do(req, &credential)
	if err != nil {
		return nil, err
	}

	fmt.Printf("status: %s\n", resp.Status)

	return &credential, nil
}

func (c *CredentialsService) GetAllIOS(projectId string) ([]responses.IOSCred, error) {
	path := fmt.Sprintf("api/v1/orgs/%s/projects/%s/credentials/signing/ios", c.OrgId, projectId)

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var credentials []responses.IOSCred
	resp, err := c.do(req, &credentials)
	if err != nil {
		return nil, err
	}

	fmt.Printf("status: %s\n", resp.Status)

	return credentials, nil
}

func (c *CredentialsService) UpdateIOS(projectId, certId, label, certPath, profilePath, certPass string) (*responses.IOSCred, error) {
	path := fmt.Sprintf("api/v1/orgs/%s/projects/%s/credentials/signing/ios/%s", c.OrgId, projectId, certId)

	formData := map[string]io.Reader{
		"label":                   strings.NewReader(label),
		"fileCertificate":         mustOpen(certPath),
		"fileProvisioningProfile": mustOpen(profilePath),
		"certificatePass":         strings.NewReader(certPass),
	}

	req, err := c.newFormRequest("PUT", path, formData)
	if err != nil {
		return nil, err
	}

	var respData responses.IOSCred
	resp, err := c.do(req, &respData)
	if err != nil {
		return nil, err
	}

	fmt.Printf("status: %s\n", resp.Status)

	return &respData, nil
}

func (c *CredentialsService) UploadIOS(projectId, label, certPath, profilePath, certPass string) (*responses.IOSCred, error) {
	path := fmt.Sprintf("api/v1/orgs/%s/projects/%s/credentials/signing/ios", c.OrgId, projectId)

	formData := map[string]io.Reader{
		"label":                   strings.NewReader(label),
		"fileCertificate":         mustOpen(certPath),
		"fileProvisioningProfile": mustOpen(profilePath),
		"certificatePass":         strings.NewReader(certPass),
	}

	req, err := c.newFormRequest("POST", path, formData)
	if err != nil {
		return nil, err
	}

	var respData responses.IOSCred
	resp, err := c.do(req, &respData)
	if err != nil {
		return nil, err
	}

	fmt.Printf("status %s\n", resp.Status)

	return &respData, nil
}

func (c *CredentialsService) DeleteIOS(projectId, certId string) (*http.Response, error) {
	path := fmt.Sprintf("api/v1/orgs/%s/projects/%s/credentials/signing/ios/%s", c.OrgId, projectId, certId)

	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("status %s\n", resp.Status)

	return resp, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
