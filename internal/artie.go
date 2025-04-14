package internal

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ArtieClient struct {
	apiKey string
	apiURL string
}

func NewArtieClient(apiKey, apiURL string) ArtieClient {
	return ArtieClient{
		apiKey: apiKey,
		apiURL: cmp.Or(apiURL, "https://api.artie.com"),
	}
}

func (a ArtieClient) doRequest(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s%s", a.apiURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.apiKey))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return out, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	return out, nil
}

func (a ArtieClient) ListDeployments(ctx context.Context) error {
	out, err := a.doRequest(ctx, http.MethodGet, "/deployments", nil)
	if err != nil {
		return err
	}

	var resp ListDeploymentsResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return err
	}

	fmt.Println("--------------------------------")
	for _, deployment := range resp.Items {
		out, err := json.Marshal(deployment)
		if err != nil {
			return err
		}

		fmt.Println(string(out))
	}
	fmt.Println("--------------------------------")
	return nil
}
