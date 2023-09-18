package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func ZipFolder(folderPath, zipFilePath string) error {
	// Create a new zip file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the folder and add files to the zip archive
	err = filepath.Walk(folderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create a file header for the zip entry
		zipHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set the name of the zip entry to be relative to the folder
		zipHeader.Name, err = filepath.Rel(folderPath, filePath)
		if err != nil {
			return err
		}

		// Create a zip entry in the archive
		zipEntry, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}

		// Copy the file content to the zip entry
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	log.Info("Folder zipped successfully to:", zipFilePath)
	return nil
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
		zipPath := zippedRoot + "/" + filepath.Base(folder.Name()) + ".zip"
		log.Info("Zipping folder: ", folder, " to: ", zipPath)
		ZipFolder(rootPath+folder.Name(), zipPath)
	}
}
