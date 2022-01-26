package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const obabel string = "C:\\Program Files\\OpenBabel-3.1.1\\obabel.exe"
const goRoutinesBatchSize int = 128

// for file structure directory > file to be converted
func obabelConversion2(directory string, ext1 string, ext2 string, addHydrogens string, addCoords bool) {
	// Read in all files in dir
	fileInfo, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("failed to read directory: " + directory)
		log.Fatal(err)
	}


	// Iterate through all items in directory
	maxFrame := min(goRoutinesBatchSize,len(fileInfo)-1)
	frame := []int{0,maxFrame}

	for frame[0] < len(fileInfo) {

		wg := sync.WaitGroup{}

		for i := frame[0]; i <= frame[1]; i++ {
			if filepath.Ext(fileInfo[i].Name()) == ext2 {
				_ = os.Remove(filepath.Join(directory, fileInfo[i].Name()))
			} else if filepath.Ext(fileInfo[i].Name()) == ext1 {
				baseName := strings.Split(fileInfo[i].Name(), ".")[0]
				convName := baseName + ext2
				basePath := filepath.Join(directory, fileInfo[i].Name())
				convPath := filepath.Join(directory, convName)

				wg.Add(1)
				go obabelWrapper(basePath, convPath, addHydrogens, addCoords, &wg)

			}
		}
		frame[0] += goRoutinesBatchSize
		frame[1] += goRoutinesBatchSize
		frame[1] = min(frame[1], len(fileInfo)-1)
		wg.Wait()
	}
}

func obabelWrapper(path1 string, path2 string, addHydrogens string, addCoords bool, wg *sync.WaitGroup) {
	cmdArgs := []string{path1, "-O", path2}
	if addHydrogens == "add" {
		cmdArgs = append(cmdArgs, "-h")
	} else if addHydrogens == "remove" {
		cmdArgs = append(cmdArgs, "-d")
	}
	if addCoords == true {
		cmdArgs = append(cmdArgs, "--gen3d")
	}

	//fmt.Println(path2)
	//cmdstring := obabel + " -i " + path1 + " -o " + path2
	out, err := exec.Command(obabel, cmdArgs...).CombinedOutput()
	//fmt.Println(string(out))
	wg.Done()
	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err)
		log.Fatal(err)
	}
}

// get min of two integers
func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}