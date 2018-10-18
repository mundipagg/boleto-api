// Package decentcopy copy files in simple form
package decentcopy

import (
	"io"
	"os"
)

// Copy use to copy files
func Copy(filepathOrigin, filepathDestiny string) error {
	srcFile, _ := os.Open(filepathOrigin)
	defer srcFile.Close()

	destFile, _ := os.Create(filepathDestiny)
	defer destFile.Close()

	io.Copy(destFile, srcFile)

	err := destFile.Sync()

	return err
}
