package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func postProcessTXYZs(dir string, keyChain map[string]map[string]string) {
	thisFile, err := os.Create(filepath.Join(dir, "filelist"))
	if err != nil {
		fmt.Println("Failed to create new fragment file: " + dir)
		log.Fatal(err)
	}

	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("failed to read directory: " + dir)
		log.Fatal(err)
	}

	for _, file := range fileInfo {
		path := filepath.Join(dir, file.Name())
		if filepath.Ext(file.Name()) == ".txyz" {
			postProcessTXYZ(path, keyChain)
			thisFile.WriteString(filepath.Base(path)+"\n")
		}
	}
}

func postProcessTXYZ(path string, keyChain map[string]map[string]string) {
	// open file
	file, err := os.Open(path)
	fmt.Println("Reading file at " + path)
	if err != nil {
		fmt.Println("Failed to open molecule file: " + path)
		log.Fatal(err)
	}
	structureName := strings.Split(filepath.Base(path),".")[0]
	fmt.Println("Structure Name: " + structureName)
	structurePieces := strings.Split(structureName, "_")
	structurePieces = structurePieces[:len(structurePieces)-1]
	newStructureName := ""
	for i, piece := range structurePieces {
		newStructureName += piece
		if i != len(structurePieces)-1 {
			newStructureName += "_"
		}
	}
	fmt.Println("Searching for key: " + newStructureName)


	relevantKey := keyChain[newStructureName]

	// Initialize scanner
	scanner := bufio.NewScanner(file)

	outLines := make([]string,0)

	for scanner.Scan() {
		// get next line
		line := scanner.Text()
		tokens := strings.Fields(line)
		if len(tokens) < 5 {
			outLines = append(outLines, line)
		} else {
			tokens[5] = relevantKey[tokens[0]]
			newLine := ""
			for _, token := range tokens {
				newLine += token + " "
			}
			outLines = append(outLines, newLine)
		}
	}

	// Write OUT
	thisFile, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to create new fragment file: " + path)
		log.Fatal(err)
	}
	for _, line := range outLines {
		_, _ = thisFile.WriteString(line + "\n")
	}
}

func getKeyChain(dir string) map[string]map[string]string {
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("failed to read directory: " + dir)
		log.Fatal(err)
	}

	keyChain := make(map[string]map[string]string)
	for _, file := range fileInfo {
		path := filepath.Join(dir, file.Name())
		keyName, keyDict := keyReader(path)
		keyChain[keyName] = keyDict
	}
	return keyChain
}

func keyReader(filePath string) (string, map[string]string) {
	// open file
	file, err := os.Open(filePath)
	fmt.Println("Reading file at " + filePath)
	if err != nil {
		fmt.Println("Failed to open molecule file: " + filePath)
		log.Fatal(err)
	}
	structureName := strings.Split(filepath.Base(filePath),".")[0]

	// Initialize scanner
	scanner := bufio.NewScanner(file)

	params := make(map[string]string)

	for scanner.Scan() {
		// get next line
		line := scanner.Text()
		// split by whitespace
		tokens := strings.Fields(line)
		if len(tokens) > 5 {
			params[tokens[0]] = tokens[5]
		}
	}
	return structureName, params
}

