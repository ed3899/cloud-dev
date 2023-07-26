package instances

import (
	"os"

	"github.com/ed3899/kumo/utils"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Zip struct {
	Name          string
	AbsPath       string
	URL           string
	ContentLength int64
	DownloadBar   *mpb.Bar
	ExtractionBar *mpb.Bar
}

func (z *Zip) GetName() string {
	return z.Name
}

func (z *Zip) GetPath() string {
	return z.AbsPath
}

func (z *Zip) IsPresent() bool {
	return utils.FilePresent(z.AbsPath)
}

func (z *Zip) IsNotPresent() bool {
	return utils.FileNotPresent(z.AbsPath)
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
}

func (z *Zip) SetExtractionBar(p *mpb.Progress, zipSize int64) {
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
}

func (z *Zip) IncrementExtractionBar(extractedBytes int) {
	z.ExtractionBar.IncrBy(extractedBytes)
}

func (z *Zip) Download(downloadedBytesChan chan<- int) (err error) {
	if err = utils.Download(z.URL, z.AbsPath, downloadedBytesChan); err != nil {
		return errors.Wrapf(err, "failed to download: %v", z.URL)
	}
	return
}

func (z *Zip) ExtractTo(extractToPath string, extractedBytesChan chan<- int) (err error) {
	if err = utils.Unzip(z.AbsPath, extractToPath, extractedBytesChan); err != nil {
		return errors.Wrapf(err, "failed to unzip: %v", z.AbsPath)
	}
	return
}

func (z *Zip) Remove() (err error) {
	if err = os.Remove(z.AbsPath); err != nil {
		return errors.Wrapf(err, "failed to remove: %v", z.AbsPath)
	}
	return
}
