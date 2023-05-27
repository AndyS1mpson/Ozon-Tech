// Generic function for making a request to a third-party service
package clientwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Make a request to a third-party service 
func DoRequest(ctx context.Context, req any, urlPath string, httpMethod string) (*http.Response, error) {
	rawData, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	httpRequest, err := http.NewRequestWithContext(ctx, httpMethod, urlPath, bytes.NewBuffer(rawData))
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong response status code: %w", err)
	}


	return httpResponse, nil
}
