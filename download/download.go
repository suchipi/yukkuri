//Package download provides functions that download files from URLs into a folder with the correct filename.
package download

import (
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

type Download struct {
	URL       string
	OutputDir string
	FileName  string
	FileSize  uint64
}

//New provides a new Download object given a url and output directory.
//An HTTP HEAD request will be issued to determine FileName and FileSize.
func New(targetUrl, outputdir string) (Download, error) {

	var empty Download

	resp, err := http.Head(targetUrl)
	defer resp.Body.Close()
	if err != nil {
		return empty, err
	}

	var contentLength string = resp.Header.Get("Content-Length")
	var fileSize uint64

	if contentLength == "" {
		//http HEAD didn't have a Content-Length header set, so we won't be able to give useful status.
		fileSize = 0
	} else {
		fileSize, err = strconv.ParseUint(contentLength, 10, 64)
		if err != nil {
			return empty, err
		}
	}

	var fileName string
	var contentDisposition string = resp.Header.Get("Content-Disposition")
	if contentDisposition == "" {
		//if no Content-Disposition header was present, we'll guess filename from the url.
		thisUrl, err := url.Parse(targetUrl)
		if err != nil {
			return empty, err
		}
		fileName = filepath.Base(thisUrl.Path)
	} else {
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			return empty, err
		}
		var prs bool
		fileName, prs = params["filename"]
		if prs != true {
			//if no filename parameter was present, we'll guess filename from the url.
			thisUrl, err := url.Parse(targetUrl)
			if err != nil {
				return empty, err
			}
			fileName = filepath.Base(thisUrl.Path)
		}
	}

	return Download{
		URL:       targetUrl,
		OutputDir: outputdir,
		FileName:  fileName,
		FileSize:  fileSize,
	}, nil
}

//Run carries out the download, blocking until it is completed.
func (dl *Download) Run() error {
	finalPath := filepath.Join(dl.OutputDir, dl.FileName)

	mode := os.FileMode(0755)

	err := os.MkdirAll(dl.OutputDir, mode)
	if err != nil {
		return err
	}

	out, err := os.Create(finalPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(dl.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
