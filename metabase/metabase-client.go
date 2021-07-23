package metabase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	Id string `json:"id"`
}

// NewClient -
func NewClient(url, username, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	if url != nil {
		c.HostURL = *url
	}

	if (username != nil) && (password != nil) && (url != nil) {
		// form request body
		rb, err := json.Marshal(AuthStruct{
			Username: *username,
			Password: *password,
		})
		if err != nil {
			return nil, err
		}

		// authenticate
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/session", c.HostURL), strings.NewReader(string(rb)))
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req)

		// parse response body
		ar := AuthResponse{}
		err = json.Unmarshal(body, &ar)
		if err != nil {
			return nil, err
		}

		c.Token = ar.Id
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
