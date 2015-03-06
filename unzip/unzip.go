package unzip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

//Single extracts a single zip file to disk, blocking until the operation is complete.
func Single(filename, outputdir string) error {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		outputPath := filepath.Join(outputdir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(outputPath, file.Mode())
		} else {
			writer, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer writer.Close()

			thisreader, err := file.Open()
			if err != nil {
				return err
			}
			defer thisreader.Close()

			_, err = io.Copy(writer, thisreader)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

//Multiple takes multiple zip part files, concatenates them, and extracts the result to disk, blocking until done.
func Multiple(files []string, outputdir string) error {

	if len(files) == 1 {
		return Single(files[0], outputdir)
	}

	interim := files[0] + ".tmp"

	interim_writer, err := os.Create(interim)
	if err != nil {
		return err
	}

	for _, file := range files {
		reader, err := os.Open(file)
		if err != nil {
			return err
		}
		_, err = io.Copy(interim_writer, reader)
		if err != nil {
			return err
		}
	}
	interim_writer.Close()

	err = Single(interim, outputdir)
	if err != nil {
		return err
	}

	os.Remove(interim)
	return nil
}
