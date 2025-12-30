package httpget

import (
   "context"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// apiError represents a non-2xx response.
type apiError struct {
   StatusCode int
   Body       []byte
   Endpoint   string
}

// Error returns a human readable error message detailing the API error.
func (e *apiError) Error() string {
   b := strings.TrimSpace(string(e.Body))
   msgCharLimit := 400
   if len(b) > msgCharLimit {
      b = b[:msgCharLimit] + "…"
   }

   if b == "" {
      return fmt.Sprintf(
         "cfbd api error for %s: status=%d", e.Endpoint, e.StatusCode,
      )
   }

   return fmt.Sprintf(
      "cfbd api error for %s: status=%d body=%s", e.Endpoint, e.StatusCode, b,
   )
}

// Client is a wrapper around http.Client which enables dependency
// injection/mocking without relying on an external resource.
type Client struct {
   HttpClient *http.Client
   BaseURL    *url.URL
   UserAgent  string
   APIKey     string
}

func (c *Client) Execute(
   ctx context.Context,
   path string,
   params url.Values,
) ([]byte, error) {
   if !strings.HasPrefix(path, "/") {
      path = "/" + path
   }

   u := c.BaseURL.ResolveReference(&url.URL{Path: path})
   u.RawQuery = params.Encode()

   req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
   if err != nil {
      return nil, fmt.Errorf("could not create request with context; %w", err)
   }

   req.Header.Set("Accept", "application/json")
   if c.UserAgent != "" {
      req.Header.Set("User-Agent", c.UserAgent)
   }

   // Set Authorization header with Bearer token.
   // The API key is validated in NewClient, so it should always be present.
   req.Header.Set("Authorization", "Bearer "+c.APIKey)

   resp, err := c.HttpClient.Do(req)
   if err != nil {
      return nil, fmt.Errorf("failed to execute request; %w", err)
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("failed to read body; %w", err)
   }

   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      return nil, &apiError{StatusCode: resp.StatusCode, Body: body}
   }

   return body, nil
}
