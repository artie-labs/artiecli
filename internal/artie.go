package internal

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
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

func (a ArtieClient) CancelDeploymentBackfill(ctx context.Context, deploymentUUID string, tableUUIDs []string) error {
	request := map[string]any{
		"optionalReason": "Done through Artie CLI",
		"tableUUIDs":     tableUUIDs,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	out, err := a.doRequest(ctx, http.MethodPost, fmt.Sprintf("/deployments/%s/backfill/cancel", deploymentUUID), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to cancel deployment backfill: %w, response: %q", err, string(out))
	}

	return nil
}

func (a ArtieClient) DeploySourceReader(ctx context.Context, sourceReaderUUID uuid.UUID) error {
	_, err := a.doRequest(ctx, http.MethodPost, fmt.Sprintf("/source-readers/%s/deploy", sourceReaderUUID), nil)
	if err != nil {
		return fmt.Errorf("failed to deploy source reader: %w", err)
	}

	return nil
}
func (a ArtieClient) GetDeploymentByUUID(ctx context.Context, deploymentUUID string) error {
	out, err := a.doRequest(ctx, http.MethodGet, fmt.Sprintf("/deployments/%s", deploymentUUID), nil)
	if err != nil {
		return err
	}

	var resp GetDeploymentResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return err
	}

	fmt.Println("--------------------------------")
	fmt.Println("Deployment:")
	out, err = json.Marshal(resp.FullDeployment.Deployment)
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	fmt.Println("Tables:")
	for _, table := range resp.FullDeployment.Source.Tables {
		out, err := json.Marshal(table)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	}
	fmt.Println("--------------------------------")
	return nil
}
