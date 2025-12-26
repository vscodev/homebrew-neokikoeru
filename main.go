package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	_assetNameMacosArm64 = "neokikoeru-macos-arm64.tar.gz"
	_assetNameMacosAmd64 = "neokikoeru-macos-amd64.tar.gz"
	_assetNameLinuxArm64 = "neokikoeru-linux-arm64.tar.gz"
	_assetNameLinuxAmd64 = "neokikoeru-linux-amd64.tar.gz"
)

var (
	_versionRegex = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
)

type Formula struct {
	Version string

	DownloadUrlMacosArm64 string
	Sha256MacosArm64      string

	DownloadUrlMacosAmd64 string
	Sha256MacosAmd64      string

	DownloadUrlLinuxArm64 string
	Sha256LinuxArm64      string

	DownloadUrlLinuxAmd64 string
	Sha256LinuxAmd64      string
}

type Asset struct {
	Name               string `json:"name"`
	Digest             string `json:"digest"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type Release struct {
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

type Error struct {
	Message string `json:"message"`
}

func fetchRelease(ctx context.Context, version string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/vscodev/neokikoeru/releases/tags/v%s", version)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		v := new(Error)
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return nil, err
		}

		return nil, errors.New(v.Message)
	}

	v := new(Release)
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}

	return v, nil
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		os.Exit(1)
	}

	version := flag.Arg(0)
	if !_versionRegex.MatchString(version) {
		os.Exit(1)
	}

	version = version[1:]
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	release, err := fetchRelease(ctx, version)
	if err != nil {
		log.Fatal(err)
	}

	formula := &Formula{
		Version: version,
	}
	for _, asset := range release.Assets {
		switch asset.Name {
		case _assetNameMacosArm64:
			formula.DownloadUrlMacosArm64 = asset.BrowserDownloadUrl
			formula.Sha256MacosArm64 = strings.TrimPrefix(asset.Digest, "sha256:")
		case _assetNameMacosAmd64:
			formula.DownloadUrlMacosAmd64 = asset.BrowserDownloadUrl
			formula.Sha256MacosAmd64 = strings.TrimPrefix(asset.Digest, "sha256:")
		case _assetNameLinuxArm64:
			formula.DownloadUrlLinuxArm64 = asset.BrowserDownloadUrl
			formula.Sha256LinuxArm64 = strings.TrimPrefix(asset.Digest, "sha256:")
		case _assetNameLinuxAmd64:
			formula.DownloadUrlLinuxAmd64 = asset.BrowserDownloadUrl
			formula.Sha256LinuxAmd64 = strings.TrimPrefix(asset.Digest, "sha256:")
		}
	}

	tmpl := template.Must(template.ParseFiles("./neokikoeru.rb.tmpl"))
	formulaFile, err := os.OpenFile("./Formula/neokikoeru.rb", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer formulaFile.Close()

	if err = tmpl.Execute(formulaFile, formula); err != nil {
		log.Fatal(err)
	}
}
