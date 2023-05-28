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
func DoRequest[Req any, Res any](ctx context.Context, req Req, resp *Res, urlPath string, httpMethod string) error {
	rawData, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("encode request: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	httpRequest, err := http.NewRequestWithContext(ctx, httpMethod, urlPath, bytes.NewBuffer(rawData))
	if err != nil {
		return fmt.Errorf("prepare request: %w", err)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong response status code: %w", err)
	}

	
	err = json.NewDecoder(httpResponse.Body).Decode(&resp)
	if err != nil {
		return fmt.Errorf("decode stock request: %w", err)
	}

	return nil
}
