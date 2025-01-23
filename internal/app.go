package internal

import (
	"github.com/vulcand/oxy/v2/forward"
	"github.com/vulcand/oxy/v2/testutils"
	"net/http"
)

func app() {
	// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
	fwd := forward.New(false)

	redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// let us forward this request to another server
		req.URL = testutils.MustParseRequestURI("http://localhost:63450")
		fwd.ServeHTTP(w, req)
	})

	// that's it! our reverse proxy is ready!
	s := &http.Server{
		Addr:    ":8080",
		Handler: redirect,
	}
	s.ListenAndServe()
}
