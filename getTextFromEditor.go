package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

func getTextFromEditor() (string, error) {
	tmpfile, err := ioutil.TempFile("", "tmp.*")
	if err != nil {
		return "", err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = findEditor("vim")
	}

	if editor == "" {
		editor = findEditor("nano")
	}

	if editor == "" {
		editor = findEditor("pico")
	}

	// Ugh fiiiiine :þ
	if editor == "" {
		editor = findEditor("emacs")
	}

	if editor == "" {
		return "", errors.New("No known editor could be found (including via $EDITOR env variable).")
	}

	defer os.Remove(tmpfile.Name()) // clean up

	cmd := exec.Command(editor, tmpfile.Name())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	var fileContents []byte
	fileContents, err = ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		return "", err
	}

	return string(fileContents), nil
}

func findEditor(editorName string) string {
	path, err := exec.LookPath(editorName)
	if err != nil {
		return ""
	}
	return path
}
