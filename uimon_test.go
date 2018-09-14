package uimon

import (
	"testing"
	"time"
)

func TestResetFlag(t *testing.T) {
	shouldBeTrue := false
	resetFlag(&shouldBeTrue)
	time.Sleep(time.Second * 3)
	if !shouldBeTrue {
		t.Error("Flag was not reset within 2 seconds.")
	}
}

func TestMatchFile(t *testing.T) {
	path := `"/dir/Dir/di_r/di-r/diR/file.go": OPERATION`
	f := matchFile(path)
	if f != "file.go" {
		t.Errorf("Could not parse file, got: %s", f)
	}
}
