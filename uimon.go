package uimon

import (
	"github.com/fsnotify/fsnotify"
	"go/build"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"
)

func init() {
	log.SetFlags(log.Flags() &^ log.Ldate)
}

func Starter() {
	bin, _ := exec.LookPath("go")
	env := os.Environ()
	args := []string{"go", "test", "-args", "-uimon=run"}
	syscall.Exec(bin, args, env)
}

func HotfixLoop() {
	logger("Killing %v, from %v", os.Getppid(), os.Getpid())
	syscall.Kill(os.Getppid(), syscall.SIGKILL)
	Starter()
}

func Start(exec func(), q func()) {
	c := make(chan int)
	go startWatcher(c)
	go Quit(c, q)
	exec()
	HotfixLoop()
}

func Quit(c chan int, q func()) {
	<-c
	q()
}

func startWatcher(c chan int) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Print(err)
	}
	defer watcher.Close()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	gp := build.Default.GOPATH
	logger("Serving files from: $GOPATH%s", strings.TrimPrefix(wd, gp))
	if err := watcher.Add(wd); err != nil {
		panic(err)
	}

	logger("Watching files...")
	runWatcherLoop(watcher, c)
}

func runWatcherLoop(w *fsnotify.Watcher, c chan int) {
	flag := true
	go func() {
		// file path | operation
		for {
			select {
			case ev, ok := <-w.Events:
				if !ok {
					continue
				}
				if flag {
					f := matchFile(ev.String())
					logger("%s : %v", f, ev.Op)
					resetFlag(&flag)
					c <- 1
				}
			case err, ok := <-w.Errors:
				if !ok {
					continue
				}
				panic(err)
			}
		}
	}()
	select {}
}

var r *regexp.Regexp

func matchFile(s string) string {
	if r == nil {
		r = regexp.MustCompile(`.*/(.*)"`)
	}
	m := r.FindStringSubmatch(s)
	return m[1]
}

func resetFlag(s *bool) {
	*s = false
	go func() {
		time.Sleep(time.Second * 2)
		*s = true
	}()
}

func logger(f string, args ...interface{}) {
	f = "\033[1;34m[uimon]\033[0m " + f
	log.Printf(f, args...)
}
