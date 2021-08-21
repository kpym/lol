package app

import (
	"testing"

	"github.com/kpym/lol/builder"
	"github.com/kpym/lol/log"
)

func TestGetFiles(t *testing.T) {
	// no Main, no Patterns
	params := builder.Parameters{Log: log.New()}
	files, err := GetFiles(params)

	// Check for error
	if err == nil {
		t.Errorf("The Main field is emty, there should be an error.")
	}

	// Main only (no Patterns)
	testData := []struct {
		fname string
		ok    bool
	}{
		{"app.go", true},
		{"app_test.go", false},
		{"app.md", false},
	}
	params.Main = "app.go"
	files, err = GetFiles(params)

	// Check for error
	if err != nil {
		t.Errorf("Error while converting parameters to files. %v", err)
	}
	// Check for files, no patterns main only
	for _, check := range testData {
		if _, ok := files[check.fname]; ok != check.ok {
			if check.ok {
				t.Errorf("Missing %s in files.", check.fname)
			} else {
				t.Errorf("Present %s in files, but it shouldn't be there.", check.fname)
			}
		}
	}

	// Main and Patterns
	testData = []struct {
		fname string
		ok    bool
	}{
		{"app.go", true},
		{"app_test.go", true},
		{"app.md", false},
	}
	params.Patterns = []string{"*.go"}
	files, err = GetFiles(params)

	// Check for error
	if err != nil {
		t.Errorf("Error while converting parameters to files. %v", err)
	}
	// Check for files, no patterns main only
	for _, check := range testData {
		if _, ok := files[check.fname]; ok != check.ok {
			if check.ok {
				t.Errorf("Missing %s in files.", check.fname)
			} else {
				t.Errorf("Present %s in files, but it shouldn't be there.", check.fname)
			}
		}
	}
}
