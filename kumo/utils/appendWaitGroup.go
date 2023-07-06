package utils

import "sync"

func AppendWaitGroup(wg *sync.WaitGroup, deps []*Dependency) {
	for _, d := range deps {
		d.WaitGroup = wg
	}
}