package cfbd

import (
   "bytes"
   "context"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "reflect"
   "strconv"
   "strings"

   "google.golang.org/protobuf/proto"
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

func (c *Client) unmarshal(b []byte, out proto.Message) error {
   if out == nil {
      return fmt.Errorf("out cannot be nil")
   }
   if len(bytes.TrimSpace(b)) == 0 || isJSONNull(b) {
      return nil
   }

   if err := c.unmarshaller.Unmarshal(b, out); err != nil {
      return fmt.Errorf("")
   }

   return nil
}

func (c *Client) unmarshalList(
   b []byte, out any, prototype proto.Message,
) error {
   if len(bytes.TrimSpace(b)) == 0 || isJSONNull(b) {
      return nil
   }
   if prototype == nil {
      return fmt.Errorf("prototype cannot be nil (e.g. &pb.Drive{})")
   }

   rv := reflect.ValueOf(out)
   if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Slice {
      return fmt.Errorf("out must be pointer to slice, got %T", out)
   }

   var raws []json.RawMessage
   if err := json.Unmarshal(b, &raws); err != nil {
      return err
   }

   slice := rv.Elem()
   for _, raw := range raws {
      if isJSONNull(raw) {
         continue
      }

      msg := proto.Clone(prototype)
      if err := c.unmarshaller.Unmarshal(raw, msg); err != nil {
         return err
      }

      // Ensure msg type matches slice element type
      msgV := reflect.ValueOf(msg)
      if !msgV.Type().AssignableTo(slice.Type().Elem()) {
         return fmt.Errorf(
            "prototype type %T not assignable to slice element type %s",
            msg, slice.Type().Elem(),
         )
      }

      slice = reflect.Append(slice, msgV)
   }

   rv.Elem().Set(slice)
   return nil
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
