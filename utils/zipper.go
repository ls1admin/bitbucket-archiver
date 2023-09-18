package utils

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ZipFolder(folderPath, zipFilePath string) error {
	log.Debug("Zipping folder: ", folderPath, " to: ", zipFilePath)
	zipfile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(folderPath)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(folderPath)
	}

	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, folderPath))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func ZipAllProjects(rootPath string, zippedRoot string) {
	log.Info("Zipping all projects in: ", rootPath)

	// create zipped root folder
	err := os.MkdirAll(zippedRoot, os.ModePerm)
	if err != nil {
		log.WithError(err).Error("Unable to create zipped root folder")
	}

	// Get all directories in the current directory = Projects
	folders, err := os.ReadDir(rootPath)
	if err != nil {
		log.WithError(err).Error("Unable to list folders")
	}
	log.Debug("Found folders: ", folders)

	for _, folder := range folders {
		zipPath := path.Join(zippedRoot, folder.Name()+".zip")
		log.Info("Zipping folder: ", folder, " to: ", zipPath)
		ZipFolder(path.Join(rootPath, folder.Name()), zipPath)
	}
}
