package download

import (
	"os"
	"sync"

	"github.com/ed3899/kumo/binz/download/draft"
)

func RemoveDownloads(tobeRemoved *draft.Dependencies) (err error) {
	var wg sync.WaitGroup
	var errChan = make(chan error, len(*tobeRemoved))

	for _, d := range *tobeRemoved {
		wg.Add(1)
		go func(d *draft.Dependency) {
			defer wg.Done()
			err = os.Remove(d.ZipPath)
			if err != nil {
				errChan <- err
			}
		}(d)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for e := range errChan {
		if e != nil {
			return e
		}
	}

	return nil
}
