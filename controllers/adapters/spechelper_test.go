package adapters_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/adamluzsi/frameless"
	"github.com/adamluzsi/frameless/iterators/iterateover"
	"github.com/stretchr/testify/require"
)

type mockPresenter struct {
	writer io.Writer
}

func (this *mockPresenter) Render(message interface{}) error {
	_, err := fmt.Fprint(this.writer, message)
	return err
}

func MockPresenterBuilder() func(io.Writer) frameless.Presenter {
	return func(w io.Writer) frameless.Presenter { return &mockPresenter{w} }
}

func MockIteratorBuilder() func(io.Reader) frameless.Iterator {
	return iterateover.LineByLine
}

func ControllerFor(t testing.TB, opts map[interface{}]interface{}, readBody bool, err error) frameless.Controller {
	return frameless.ControllerFunc(func(p frameless.Presenter, r frameless.Request) error {

		if opts != nil {
			for k, v := range opts {
				require.Equal(t, v, r.Context().Value(k))

				p.Render(r.Context().Value(k))
			}
		}

		if readBody {
			i := r.Data()
			defer i.Close()

			for i.Next() {
				var d string

				if err := i.Decode(&d); err != nil {
					return err
				}

				if err := p.Render(d); err != nil {
					return err
				}
			}
		}

		return err
	})
}

type Body struct {
	*bytes.Buffer
	IsClosed bool
}

func (b *Body) Close() error {
	b.IsClosed = true
	return nil
}
