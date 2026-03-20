package hh

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-querystring/query"
)

const API_URL = "https://api.hh.ru"

var ErrRequest = errors.New("Request failed")

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{http.DefaultClient}
}

func (c Client) GetVacancies(reqb GetVacanciesRequest) (*GetVacanciesResponse, error) {
	v, err := query.Values(reqb)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRequest, err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/vacancies?%s", API_URL, v.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRequest, err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRequest, err)
	}

	defer res.Body.Close()

	resb, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRequest, err)
	}

	var r GetVacanciesResponse

	if err := json.Unmarshal(resb, &r); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRequest, err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: %d: %s", ErrRequest, res.StatusCode, resb)
	}

	return &r, nil
}
