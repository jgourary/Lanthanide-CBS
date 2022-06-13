package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func makeCSV(dir string) {
	fileList := loadTXT(filepath.Join(dir, "filelist"))
	distances := getNearestDistances(fileList)
	QMEnergy := loadTXT(filepath.Join(dir, "QM-energies.dat"))
	MMEnergy := loadTXT(filepath.Join(dir, "result.p"))

	thisFile, err := os.Create(filepath.Join(dir, "results.csv"))
	if err != nil {
		fmt.Println("Failed to create result.csv")
		log.Fatal(err)
	}
	for i := 0; i < len(fileList); i++ {
		_, _ = thisFile.WriteString(fileList[i] + ", " + fmt.Sprintf("%.6f", distances[i]) + ", " + fmt.Sprintf("%.6f", QMEnergy[i]) + ", " + fmt.Sprintf("%.6f", MMEnergy[i]))
	}
	thisFile.Close()
}


//func load
func getNearestDistances(files []string) []float64 {
	distances := make([]float64, 0)
	for _, file := range files {
		_, structure := loadTXYZ(file)
		dist := getIonLigandDist(structure)
		distances = append(distances, dist)
	}
	return distances
}

func getIonLigandDist(structure map[int]*atom) float64 {
	closestDist := 1000000.0
	ionKey := -1

	// identify ion
	for i, atom := range structure {
		if atom.element != "C" && atom.element != "O" && atom.element != "H" && atom.element != "N" && atom.element != "P" {
			ionKey = i
		}
	}

	// get closest ligand atom
	for i, atom := range structure {
		if i != ionKey {
			dist := getDistance(atom.pos, structure[ionKey].pos)
			if dist < closestDist {
				closestDist = dist
			}
		}
	}
	return closestDist
}

func loadTXYZ(filePath string) (string, map[int]*atom) {

	// Create structure to store atoms
	atoms := make(map[int]*atom)

	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open molecule file: " + filePath)
		log.Fatal(err)
	}
	structureName := strings.Split(filepath.Base(filePath),".")[0]


	// Initialize scanner
	scanner := bufio.NewScanner(file)
	// ignore first line
	scanner.Scan()
	// create line counter
	i := 1
	// iterate over all other lines
	for scanner.Scan() {
		// get next line
		line := scanner.Text()
		// split by whitespace
		tokens := strings.Fields(line)
		// check line length before proceeding
		if len(tokens) >= 6 {

			// create new atom
			var newAtom atom

			// get number of atom from file
			atomNum, err := strconv.Atoi(tokens[0])
			if err != nil {
				newErr := errors.New("Failed to convert token in position 0 on line " + strconv.Itoa(i) + " to an integer")
				log.Fatal(newErr)
			}

			// assign element
			newAtom.element = tokens[1]

			// assign positions
			pos := make([]float64,3)
			for j := 2; j < 5; j++ {
				pos[j-2], err = strconv.ParseFloat(tokens[j],64)
				if err != nil {
					newErr := errors.New("Failed to convert token in position 0 on line " + strconv.Itoa(j) + " to a float64")
					log.Fatal(newErr)
				}
			}
			newAtom.pos = pos

			// assign atomType from file
			newAtom.atomType = tokens[5]

			// assign bonds from file
			bonds := make([]string,len(tokens)-6)
			for j := 6; j < len(tokens); j++ {
				bonds[j-6] = tokens[j]
			}
			newAtom.bonds = bonds

			// add atom to map
			atoms[atomNum] = &newAtom

		} else {
			fmt.Println("Warning: line " + strconv.Itoa(i) + " has insufficient tokens. Program is skipping this " +
				"line when reading your input file.")
		}
		i++
	}

	return structureName, atoms
}

func loadTXT(filePath string) []string {
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open molecule file: " + filePath)
		log.Fatal(err)
	}
	keys := make([]string,0)
	// Initialize scanner
	scanner := bufio.NewScanner(file)
	// ignore first line
	scanner.Scan()
	// iterate over all other lines
	for scanner.Scan() {
		line := scanner.Text()
		// split by whitespace
		tokens := strings.Fields(line)
		if len(tokens) > 0 {
			keys = append(keys, tokens[0])
		}
	}
	return keys
}
