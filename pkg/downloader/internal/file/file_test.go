package file

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestList(t *testing.T) {
	pwd, _ := os.Getwd()
	files, err := ListDir(pwd)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(files)
	for _, file := range files {
		t.Log(file)
	}
}

func TestWalk(t *testing.T) {
	pwd, _ := os.Getwd()
	pkgPath := filepath.Dir(filepath.Dir(filepath.Dir(pwd))) // should in ./pkg
	files, err := Walk(pkgPath)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(files)
	for _, file := range files {
		t.Log(file)
	}
}
