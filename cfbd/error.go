package cfbd

import (
   "fmt"
   "strings"
)

// APIError represents a non-2xx response.
type APIError struct {
   StatusCode int
   Body       []byte
   Endpoint   string
}

// Error returns a human readable error message detailing the API error.
func (e *APIError) Error() string {
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
