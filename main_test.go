package main

import (
	"fmt"
	"io"
	"testing"
)

func TestEmbedFS(t *testing.T) {
	f, err := htmlFiles.Open("assets/html/test.html")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(io.ReadAll(f))
	f.Close()
}
