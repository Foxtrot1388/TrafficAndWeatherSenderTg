package traffic

import (
	"context"
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

func (t *trafficGeter) Info(ctx context.Context, url string, w io.Writer) {

	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			_ = t.toImage(ctx, string(r.Body), w)
		},
	}).Start()

}

func (t *trafficGeter) toImage(ctx context.Context, htmlString string, w io.Writer) error {

	if err := ctx.Err(); err != nil {
		return err
	}

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
