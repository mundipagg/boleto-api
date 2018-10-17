/*

decentcopy provides a copy files for humans

Usage

	package main

	import "github.com/hugocarreira/go-decent-copy"

	func main() {
		execPath, _ := os.Getwd()

		fileOrigin := "/testCopy.txt"
		fileDestiny := fmt.Sprintf(execPath + "/folder/" + fileOrigin)
		err := Copy(fileOrigin, fileDestiny)
		if err != nil {
			fmt.Printf("Error in copy file : %#v ", err.Error())
		}
	}

*/

package decentcopy

import (
	"io"
	"os"
)

func Copy(filepathOrigin, filepathDestiny string) error {
	execPath, _ := os.Getwd()

	srcFile, err := os.Open(execPath + filepathOrigin)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(filepathDestiny)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
