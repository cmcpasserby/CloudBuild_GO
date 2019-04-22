package cloudbuild

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

const baseUrl = "build-api.cloud.unity3d.com"

type client struct {
	BaseUrl    *url.URL
	ApiKey     string
	OrgId      string
	httpClient *http.Client
}

func newClient(apiKey, orgId string) *client {
	return &client{
		BaseUrl:    &url.URL{Scheme: "https", Host: baseUrl},
		ApiKey:     apiKey,
		OrgId:      orgId,
		httpClient: http.DefaultClient,
	}
}

func (c *client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseUrl.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.ApiKey))

	return req, nil
}

func (c *client) newFormRequest(method, path string, form map[string]io.Reader) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseUrl.ResolveReference(rel)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for key, value := range form {
		var fw io.Writer
		var err error

		if x, ok := value.(io.Closer); ok {
			//noinspection GoDeferInLoop
			defer x.Close()
		}

		if x, ok := value.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return nil, err
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return nil, err
			}
		}
		if _, err := io.Copy(fw, value); err != nil {
			return nil, err
		}
	}

	w.Close()

	req, err := http.NewRequest(method, u.String(), &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.ApiKey))

	return req, nil
}

func (c *client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyString, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(bodyString))
	}

	if resp.StatusCode == 204 { // no content to decode
		return resp, nil
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return nil, err
		}
	}
	return resp, nil
}
