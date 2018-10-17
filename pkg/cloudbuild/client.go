package cloudbuild

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

const baseUrl = "build-api.cloud.unity3d.com"

type Client struct {
	BaseUrl    *url.URL
	ApiKey     string
	OrgId      string
	httpClient *http.Client
}

func New(apiKey, orgId string) *Client {
	return &Client{
		BaseUrl:    &url.URL{Scheme: "https", Host: baseUrl},
		ApiKey:     apiKey,
		OrgId:      orgId,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
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

func (c *Client) newFormRequest(method, path string, form map[string]io.Reader) (*http.Request, error) {
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

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
