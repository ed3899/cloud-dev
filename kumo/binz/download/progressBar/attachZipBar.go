package progressBar

import (
	"fmt"
	"log"
	"os"

	"github.com/ed3899/kumo/binz/download/draft"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type DownloadResult struct {
	Dependency *draft.Dependency
	Err        error
}

func AttachZipBar(p *mpb.Progress, dr *DownloadResult) {
	zipSize, err := getZipSize(dr.Dependency.ZipPath)
	if err != nil {
		log.Printf("there was an error while getting the zip size: %v", err)
		return
	}

	barName := fmt.Sprintf("%s.zip", dr.Dependency.Name)

	// Config the bar
	zipBar := p.AddBar(zipSize,
		mpb.BarQueueAfter(dr.Dependency.DownloadBar.(*mpb.Bar)),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(barName),
			decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Percentage(decor.WCSyncSpace),
				"unzipped",
			),
		),
	)

	// Assign the bar to the dependency
	dr.Dependency.ZipBar = zipBar

}

func getZipSize(path string) (size int64, err error) {
	zipfile, err := os.Open(path)
	if err != nil {
		return
	}
	defer zipfile.Close()

	info, err := zipfile.Stat()
	if err != nil {
		return
	}
	size = info.Size()

	return
}
