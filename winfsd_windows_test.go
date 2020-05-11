package winfsd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestWinfsd(t *testing.T) {
	f, err := ioutil.TempFile("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		// If f still exists, then it must be closed before it can be removed.
		f.Close()
		os.Remove(f.Name())
	}()

	// If the file exists, del sets errorlevel to 0 even if it fails  as long as the file exists, so we have to
	commands := fmt.Sprintf("del %s && if exist %s exit 1", f.Name(), f.Name())
	out, err := exec.Command("cmd.exe", "/c", commands).CombinedOutput()
	if err != nil {
		t.Errorf("command failed: %s: %q", err, string(out))
	}
}
