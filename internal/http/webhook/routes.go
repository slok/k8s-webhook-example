package webhook

import (
	"net/http"
)

// routes wires the routes to handlers on a specific router.
func (h handler) routes(router *http.ServeMux) error {
	allmark, err := h.allMark()
	if err != nil {
		return err
	}
	router.Handle("/wh/mutating/allmark", allmark)

	return nil
}
