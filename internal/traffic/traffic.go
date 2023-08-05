package traffic

import (
	"io"
	"log"

	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	wk "github.com/shezadkhan137/go-wkhtmltoimage"
)

func GetInfo(url string, w io.Writer) {

	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			toImage(string(r.Body), w)
		},
	}).Start()

}

func toImage(htmlString string, w io.Writer) {

	wk.Init()
	defer wk.Destroy()

	converter, err := wk.NewConverter(
		&wk.Config{
			Quality:          100,
			Fmt:              "png",
			EnableJavascript: false,
		})
	if err != nil {
		log.Fatal(err)
	}

	err = converter.Run(htmlString, w)
	if err != nil {
		log.Fatal(err)
	}

}
