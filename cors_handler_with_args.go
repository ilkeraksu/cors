package cors

import (
	"context"
	"net/http"
)

// HandlerFunc provides context compatible handler
// Also takes args as request arguments
func (c *Cors) HandlerContextFuncEx(h func(ctx context.Context, w http.ResponseWriter, r *http.Request, args map[string]string)) func(ctx context.Context, w http.ResponseWriter, r *http.Request, args map[string]string) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, args map[string]string) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			c.logf("Handler: Preflight request")
			c.handlePreflight(w, r)
			// Preflight requests are standalone and should stop the chain as some other
			// middleware may not handle OPTIONS requests correctly. One typical example
			// is authentication middleware ; OPTIONS requests won't carry authentication
			// headers (see #1)
			if c.optionPassthrough {
				h(ctx, w, r, args)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			c.logf("Handler: Actual request")
			c.handleActualRequest(w, r)
			h(ctx, w, r, args)
		}
	}
}
