package binaries

import (
	"os"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type ZipI interface {
	Download() (err error)
	Extract() (err error)
	Remove() (err error)
}

type Zip struct {
	Name          string
	Path          string
	URL           string
	ContentLength int64
	DownloadBar   *mpb.Bar
	ExtractionBar *mpb.Bar
}

func (z *Zip) SetDownloadBar(p *mpb.Progress) {
	z.DownloadBar = p.AddBar(int64(z.ContentLength),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(z.Name),
			decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Percentage(decor.WCSyncSpace),
				"downloaded",
			),
		),
	)
}

func (z *Zip) SetExtractionBar(p *mpb.Progress) (err error) {
	if utils.FileNotPresent(z.Path) {
		err = errors.New("zip file not present")
		return
	}

	zipSize, err := utils.GetZipSize(z.Path)
	if err != nil {
		err = errors.Wrapf(err, "failed to get zip size for: %v", z.Path)
		return
	}

	z.ExtractionBar = p.AddBar(zipSize,
		mpb.BarQueueAfter(z.DownloadBar),
		mpb.BarFillerClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(z.Name),
			decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Percentage(decor.WCSyncSpace),
				"unzipped",
			),
		),
	)

	return
}

func (z *Zip) Download(downloadedBytesChan chan<- int) (err error) {
	if utils.FilePresent(z.Path) {
		err = errors.New("zip file already present")
		return
	}

	if z.DownloadBar == nil {
		err = errors.New("download bar not set")
		return
	}

	if err = utils.Download(z.URL, z.Path, downloadedBytesChan); err != nil {
		err = errors.Wrapf(err, "failed to download: %v", z.URL)
		return
	}
	return
}

func (z *Zip) Extract(extractToPath string, extractedBytesChan chan<- int) (err error) {
	if utils.FileNotPresent(z.Path) {
		err = errors.New("zip file not present")
		return
	}

	if z.ExtractionBar == nil {
		err = errors.New("extraction bar not set")
		return
	}

	if err = utils.Unzip(z.Path, extractToPath, extractedBytesChan); err != nil {
		err = errors.Wrapf(err, "failed to unzip: %v", z.Path)
		return
	}
	return
}

func (z *Zip) Remove() (err error) {
	if err = os.Remove(z.Path); err != nil {
		err = errors.Wrapf(err, "failed to remove: %v", z.Path)
		return
	}
	return
}
