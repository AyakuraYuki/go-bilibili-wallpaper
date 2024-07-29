package file

import (
	"os"
	"testing"

	cjson "github.com/AyakuraYuki/bilibili-wallpaper/internal/json"
)

func TestListDir(t *testing.T) {
	home, _ := os.Getwd()
	list, err := ListDir(home)
	if err != nil {
		t.Fatal(err)
	}
	j, _ := cjson.JSON.MarshalIndent(list, "", "    ")
	t.Log(string(j))
}

func TestWalkDir(t *testing.T) {
	home, _ := os.Getwd()
	list := make([]string, 0)
	WalkDir(home, &list)
	j, _ := cjson.JSON.MarshalIndent(list, "", "    ")
	t.Log(string(j))
}

func TestWriteLines(t *testing.T) {
	filename := "test.json"
	defer func(name string) { _ = os.Remove(name) }(filename)

	WriteLines(filename, []string{"[]"})
	c := ReadFile(filename)
	t.Logf("phase 1: %v", c)
	WriteLines(filename, []string{"{}"})
	c = ReadFile(filename)
	t.Logf("phase 2: %v", c)
}
