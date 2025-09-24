package response

import (
	"fmt"

	"github.com/felixsolom/http-from-tcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprint(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}
