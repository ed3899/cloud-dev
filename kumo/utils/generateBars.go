package utils

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func AppendBarToDependencies(wg *sync.WaitGroup, deps []*Dependency) {
	progress := mpb.New(mpb.WithWaitGroup(wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))

	for _, d := range deps {
		// Config the bar
		bar := progress.AddBar(int64(d.ContentLength),
			mpb.PrependDecorators(
				decor.Name(d.Name),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"done",
				),

			),
		)

		// Assign the bar to the dependency
		d.Bar = bar
	}
}

func GenerateBars(progress *mpb.Progress, ze []*Dependency) []*mpb.Bar {
	// Create the bars
	bars := make([]*mpb.Bar, 0)
	for i := 0; i < len(ze); i++ {
		var bar *mpb.Bar
		url := ze[i].URL
		name := filepath.Base(ze[i].ExtractionPath)
		// Perform a HEAD request to get the content length
		resp, err := http.Head(url)
		if err != nil {
			log.Printf("Error occurred while sending HEAD request: %v\n", err)
			continue
		}

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			log.Printf("Received non-200 status code: %d\n", resp.StatusCode)
			resp.Body.Close()
			continue
		}
		defer resp.Body.Close()

		// Retrieve the Content-Length header
		contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			log.Printf("There was an error while attempting to parse the content length from '%s': '%#v'", url, err)
			resp.Body.Close()
			continue
		}

		// Assign the bar length to the content length
		bar = progress.AddBar(int64(contentLength),
			mpb.PrependDecorators(
				decor.Name(name),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"done",
				),
			),
		)

		bars = append(bars, bar)
	}

	return bars
}
