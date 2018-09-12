package uimon

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestUimon(t *testing.T) {
	ioutil.WriteFile("x.txt", []byte("x"), 0644)
	Start(func() {
		ioutil.WriteFile("x.txt", []byte("xx"), 0644)
		time.Sleep(1 * time.Second)
		os.Remove("x.txt")
		os.Exit(0)
	})
}
