package core

import (
	"sync"
	"testing"
)

func TestApplicationDirectoriesConcurrent(t *testing.T) {
	app := NewApp()
	app.setDirectories([]string{"a", "b", "c"})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			app.setDirectories([]string{"x", "y", "z"})
		}
	}()

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				_ = len(app.Directories())
			}
		}()
	}

	wg.Wait()
}
