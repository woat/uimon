package uimon

import (
	"github.com/fsnotify/fsnotify"
	"go/build"
	"log"
	"os"
	"regexp"
	"strings"
)

func logger(s string) {
	log.Printf("\033[1;34m[uimon]\033[0m %s", s)
}

func Start(exec func()) {
	log.SetFlags(log.Flags() &^ log.Ldate)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Print(err)
	}
	defer watcher.Close()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := watcher.Add(wd); err != nil {
		panic(err)
	}

	gp := build.Default.GOPATH
	logger("Serving files from: $GOPATH" + strings.TrimPrefix(wd, gp))

	logger("Watching files...")
	go func() {
		r := regexp.MustCompile(`.*/(.*)"`)
		for {
			select {
			case event := <-watcher.Events:
				f := r.FindStringSubmatch(event.String())
				logger("File changed: " + f[1])
				exec()
			case err := <-watcher.Errors:
				panic(err)
			}
		}
	}()
	select {}
}
