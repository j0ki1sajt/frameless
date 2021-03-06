package adapters

import (
	"io"
	"net/http"

	"github.com/adamluzsi/frameless"
	fhttp "github.com/adamluzsi/frameless/controllers/adapters/integrations/net/http"
)

// NetHTTP creates an adapter http.Hander object that can be given to a http.ServerMux
func NetHTTP(
	controller frameless.Controller,
	buildPresenter func(io.Writer) frameless.Presenter,
	buildIterator func(io.Reader) frameless.Iterator,
) http.Handler {

	closer := func(c io.Closer) {
		if c != nil {
			c.Close()
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer closer(r.Body)

		presenter := buildPresenter(w)
		request := fhttp.NewRequest(r, buildIterator)

		if err := controller.Serve(presenter, request); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})
}
