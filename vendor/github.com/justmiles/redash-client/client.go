package redash

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"fmt"
)

// Client is the API client that performs all operations
// against a redash server.
type Client struct {
	BaseURL      *url.URL
	UserAgent    string
	httpClient   *http.Client
	apiKey       string
	DebugEnabled bool
}

// NewClient creates a redash client
func NewClient(uri, apiKey string) (*Client, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, err
	}
	client := Client{
		BaseURL:    u,
		UserAgent:  "redash go", // TODO: add go version and lib version
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
	}
	return &client, nil
}

func (c *Client) newRequest(method, path string, body interface{}, options *map[string]string) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

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
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", "Key "+c.apiKey)

	if options != nil {
		q := req.URL.Query()
		for key, value := range *options {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

func (c *Client) download(req *http.Request, filepath string) (*http.Response, error) {

	out, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return resp, err
}

func logRawHTTPBody(body io.ReadCloser) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	newStr := buf.String()
	fmt.Println(newStr)
}
