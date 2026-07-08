package brreg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const brregBaseUrl = "https://data.brreg.no"

type Client struct {
	HTTP    *http.Client
	BaseURL string
}

func NewClient() *Client {
	return &Client{
		HTTP:    &http.Client{Timeout: 30 * time.Second},
		BaseURL: brregBaseUrl,
	}
}

func (c *Client) GetHovedenhetByOrgnummer(orgnummer string) (*Hovedenhet, error) {
	endpoint := fmt.Sprintf("%s/enhetsregisteret/api/enheter/%s", c.BaseURL, orgnummer)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	log.Println("requesting ", endpoint)

	req.Header.Set("User-Agent", "SN-Utils")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var hovedenhet Hovedenhet
	if err := json.NewDecoder(resp.Body).Decode(&hovedenhet); err != nil {
		return nil, ErrUnmarshallingResponse
	}

	return &hovedenhet, nil
}
