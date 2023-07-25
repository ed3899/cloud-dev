package binaries

import (
	"os"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Removable interface {
	Remove() error
}

type Retrivable interface {
	GetName() string
	GetPath() string
}

type Downloadable interface {
	SetDownloadBar(*mpb.Progress)
	Download(chan<- int) error
	IncrementDownloadBar(int)
}

type Extractable interface {
	SetExtractionBar(*mpb.Progress, int64)
	Extract(string, chan<- int) error
	IncrementExtractionBar(int)
}

type ZipI interface {
	Retrivable
	Downloadable
	Extractable
	Removable
}

type Zip struct {
	Name          string
	Path          string
	URL           string
	ContentLength int64
	DownloadBar   *mpb.Bar
	ExtractionBar *mpb.Bar
}

func (z *Zip) GetName() string {
	return z.Name
}

func (z *Zip) GetPath() string {
	return z.Path
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

func (z *Zip) IncrementDownloadBar(downloadedBytes int) {
	z.DownloadBar.IncrBy(downloadedBytes)
	return
}

func (z *Zip) SetExtractionBar(p *mpb.Progress, zipSize int64) {
	// if utils.FileNotPresent(z.Path) {
	// 	err = errors.New("zip file not present")
	// 	return
	// }

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

func (z *Zip) IncrementExtractionBar(extractedBytes int) {
	// if z.ExtractionBar == nil {
	// 	err = errors.New("extraction bar not set")
	// 	return
	// }

	z.ExtractionBar.IncrBy(extractedBytes)

	return
}

func (z *Zip) Download(downloadedBytesChan chan<- int) (err error) {
	if err = utils.Download(z.URL, z.Path, downloadedBytesChan); err != nil {
		err = errors.Wrapf(err, "failed to download: %v", z.URL)
		return
	}
	return
}

func (z *Zip) Extract(extractToPath string, extractedBytesChan chan<- int) (err error) {
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
