package presenters_test

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/adamluzsi/frameless/presenters"
	"github.com/stretchr/testify/require"
)

var (
	a = template.Must(template.New("a").Parse(`<body>{{.}}</body>`))
	b = template.Must(template.New("b").Parse(`<h1>{{.}}</h1>`))
	c = template.Must(template.New("c").Parse(`<p>{{.Data}}</p>`))
)

type content struct {
	Data string
}

func TestHTMLPresenterBuilderFor_PageTeplateGiven(t *testing.T) {
	t.Parallel()

	builder := presenters.ForHTMLTemplate(a, b, c)

	w := bytes.NewBuffer([]byte{})

	presenter := builder(w)

	presenter.Render(content{Data: "Hello, World!"})

	require.Equal(t, `<body><h1><p>Hello, World!</p></h1></body>`, w.String())
}
