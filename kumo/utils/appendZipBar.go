package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func AppendZipBar(p *mpb.Progress, d *DownloadResult) {
	zipSize, err := getZipSize(d.Dependency.ZipPath)
	if err != nil {
		log.Printf("there was an error while getting the zip size: %v", err)
		return
	}

	barName := fmt.Sprintf("%s.zip", d.Dependency.Name)

	// Config the bar
	zipBar := p.AddBar(zipSize,
		mpb.BarQueueAfter(d.Dependency.DownloadBar),
		mpb.BarFillerClearOnComplete(),
		mpb.BarRemoveOnComplete(),
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
	d.Dependency.ZipBar = zipBar

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