package log

import (
	"strings"
	"testing"
)

func TestLevels(t *testing.T) {
	w := new(strings.Builder)
	// QUIET
	log := New(WithWriter(w), WithLevel(Quiet))
	log.Error("Test")
	if w.Len() > 0 {
		t.Errorf("When level is set to QUIET no error should be displayed.")
	}
	// DEFAULT = ERROR
	log = New(WithWriter(w))
	w.Reset()
	log.Error("Test")
	if w.Len() == 0 {
		t.Errorf("In the default level all errors should be displayed.")
	}
	w.Reset()
	log.Info("Test")
	if w.Len() > 0 {
		t.Errorf("In the default level no info should be displayed.")
	}
	// VERBOSE
	log = New(WithWriter(w), WithLevel(InfoLevel))
	w.Reset()
	log.Info("Test")
	if w.Len() == 0 {
		t.Errorf("In the info level all infos should be displayed.")
	}
	w.Reset()
	log.Debug("Test")
	if w.Len() > 0 {
		t.Errorf("In the info level no debug should be displayed.")
	}
	// DEBUG
	log = New(WithWriter(w), WithLevel(DebugLevel))
	w.Reset()
	log.Debug("Test")
	if w.Len() == 0 {
		t.Errorf("In the debug level all debugs should be displayed.")
	}
}
