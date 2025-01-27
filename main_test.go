package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// resetGlobals resets all global variables to their default state.
func resetGlobals() {
	twidth = 80
	oneflag = false
	aflag = false
	bflag = false
	cflag = false
	dflag = false
	fflag = false
	lflag = false
	mflag = false
	pflag = false
	allflag = true
	ndir = 0
	printed = false
	maxwidth = 0
	lwidth = twidth - indent2
	wout = os.Stdout
	werr = os.Stderr
}

func TestLcSingleFile(t *testing.T) {
	resetGlobals()
	color.NoColor = true

	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Capture output.
	var buf bytes.Buffer
	wout = &buf

	if err := lc(tmpfile.Name()); err != nil {
		t.Errorf("lc returned error: %v", err)
	}

	got := buf.String()
	want := tmpfile.Name() + ": file\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLcDirectory(t *testing.T) {
	resetGlobals()

	tmpdir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	testfile := filepath.Join(tmpdir, "test.txt")
	if err := os.WriteFile(testfile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set flags as they would normally be.
	fflag = true
	dflag = true
	cflag = true
	bflag = true
	mflag = true
	lflag = true
	pflag = true

	var buf bytes.Buffer
	wout = &buf

	if err := lc(tmpdir); err != nil {
		t.Errorf("lc returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "Files:") {
		t.Errorf("output %q doesn't contain 'Files:'", got)
	}
	if !strings.Contains(got, "test.txt") {
		t.Errorf("output %q doesn't contain 'test.txt'", got)
	}
}

func TestLcNotFound(t *testing.T) {
	resetGlobals()

	var buf bytes.Buffer
	wout = &buf

	if err := lc("nonexistent"); err == nil {
		t.Error("lc didn't return error")
	}

	got := buf.String()
	want := "nonexistent: not found\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
