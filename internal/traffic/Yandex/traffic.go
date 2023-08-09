package traffic

import (
	"io"

	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	wk "github.com/shezadkhan137/go-wkhtmltoimage"
)

type trafficGeter struct {
}

func New() *trafficGeter {
	return &trafficGeter{}
}

func (t *trafficGeter) Info(url string, w io.Writer) {

	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			_ = t.toImage(string(r.Body), w)
		},
	}).Start()

}

func (t *trafficGeter) toImage(htmlString string, w io.Writer) error {

	wk.Init()
	defer wk.Destroy()

	converter, err := wk.NewConverter(
		&wk.Config{
			Quality:          100,
			Fmt:              "png",
			EnableJavascript: false,
		})
	if err != nil {
		return err
	}

	err = converter.Run(htmlString, w)
	if err != nil {
		return err
	}

	return nil

}
