package utils

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(zipPath string) string {
	outFolderName := strings.Split(zipPath, "/")[len(strings.Split(zipPath, "/"))-1]
	outFolderName = strings.Split(outFolderName, ".")[0]

	dst := strings.Split(zipPath, ".")[0]

	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			err = errors.New("Invalid ZIP/IPA file path")
			panic(err)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return dst
}

func DeleteFolder (folder string) {
	err := os.RemoveAll(folder)
	if err != nil {
		log.Fatal("Failed while deleting :", folder, " :", err.Error())
	}
}
