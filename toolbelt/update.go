package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func hkURL() string {
	return "https://s3.amazonaws.com/dickeyxxx_dev/releases/hk_" + runtime.GOOS + "_" + runtime.GOARCH + ".gz"
}

func updateHk(hkPath string) {
	fmt.Println("Downloading hk...")
	fmt.Println(filepath.Dir(hkPath))
	Must(os.MkdirAll(filepath.Dir(hkPath), 0777))
	out, err := os.Create(hkPath)
	Must(err)
	defer out.Close()
	if runtime.GOOS != "windows" {
		Must(out.Chmod(0777))
	}
	resp, err := http.Get(hkURL())
	Must(err)
	defer resp.Body.Close()
	uncompressed, err := gzip.NewReader(resp.Body)
	Must(err)
	_, err = io.Copy(out, uncompressed)
	Must(err)
	fmt.Println("Downloaded hk...")
}
