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

func readFile(filePath string) (string, []molecule) {
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

	var molecules []molecule

	// create dummy molecule
	var newMolecule molecule
	newMolAtoms := make(map[int]*atom)
	newMolecule.atoms = newMolAtoms
	atomCounter := 0

	// iterate over all other lines
	insideMoleculeBlock := false

	for scanner.Scan() {
		// get next line
		line := scanner.Text()

		// split by whitespace
		tokens := strings.Fields(line)

		if tokens[0] == "molecule" || tokens[0] == "molecule{" {
			insideMoleculeBlock = true
		} else if tokens[0] == "}" {
			insideMoleculeBlock = false
		} else if len(tokens) == 2 && insideMoleculeBlock {
			// save previous molecule
			if len(newMolecule.atoms) > 0 {
				molecules = append(molecules, newMolecule)
			}
			// create new molecule
			var newMolecule molecule
			newMolAtoms := make(map[int]*atom)
			newMolecule.atoms = newMolAtoms
			newMolecule.charge = tokens[0]
			newMolecule.multiplicity = tokens[1]
			atomCounter = 0
		} else if len(tokens) > 4 && insideMoleculeBlock {
			// create new atom
			var newAtom atom

			newAtom.element = tokens[0]
			pos := make([]float64,3)
			for j := 1; j < 4; j++ {
				pos[j-1], err = strconv.ParseFloat(tokens[j],64)
				if err != nil {
					newErr := errors.New("Failed to convert \"" + tokens[j] + "\" in position 0 on line " + strconv.Itoa(i) + " to a float64")
					log.Fatal(newErr)
				}
			}
			newAtom.pos = pos

			// add atom to map
			newMolecule.atoms[atomCounter] = &newAtom
		}
	}

	return structureName, molecules
}
