package cfbd

import (
   "bytes"
   "context"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type restClient struct {
   client    *http.Client
   baseURL   *url.URL
   userAgent string
   apiKey    string
}

func (c *restClient) execute(
   ctx context.Context,
   path string,
   params url.Values,
) ([]byte, error) {
   if !strings.HasPrefix(path, "/") {
      path = "/" + path
   }
   // ResolveReference preserves scheme/host.
   u := c.baseURL.ResolveReference(&url.URL{Path: path})
   u.RawQuery = params.Encode()

   req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("Accept", "application/json")
   if c.userAgent != "" {
      req.Header.Set("User-Agent", c.userAgent)
   }

   // Set Authorization header with Bearer token.
   // The API key is validated in NewClient, so it should always be present.
   req.Header.Set("Authorization", "Bearer "+c.apiKey)

   resp, err := c.client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      return nil, &APIError{StatusCode: resp.StatusCode, Body: body}
   }

   return body, nil
}

func isJSONNull(b []byte) bool {
   return bytes.Equal(bytes.TrimSpace(b), []byte("null"))
}

func setString(v url.Values, key string, val *string) {
   if val == nil {
      return
   }

   v.Set(key, *val)
}

func setInt32(v url.Values, key string, val *int32) {
   if val == nil {
      return
   }

   v.Set(key, strconv.FormatInt(int64(*val), 10))
}

func setInt(v url.Values, key string, val *int) {
   if val == nil {
      return
   }

   v.Set(key, strconv.Itoa(*val))
}

func setFloat64(v url.Values, key string, val *float64) {
   if val == nil {
      return
   }

   v.Set(key, strconv.FormatFloat(*val, 'f', -1, 64))
}

func setBool(v url.Values, key string, val *bool) {
   if val == nil {
      return
   }

   v.Set(key, strconv.FormatBool(*val))
}
