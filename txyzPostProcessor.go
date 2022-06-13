package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	//fmt.Println("Reading file at " + path)
	if err != nil {
		fmt.Println("Failed to open molecule file: " + path)
		log.Fatal(err)
	}
	structureName := strings.Split(filepath.Base(path),".")[0]
	//fmt.Println("Structure Name: " + structureName)
	structurePieces := strings.Split(structureName, "_")
	structurePieces = structurePieces[:len(structurePieces)-1]
	newStructureName := ""
	for i, piece := range structurePieces {
		newStructureName += piece
		if i != len(structurePieces)-1 {
			newStructureName += "_"
		}
	}
	//fmt.Println("Searching for key: " + newStructureName)


	relevantKey := keyChain[newStructureName]

	// Initialize scanner
	scanner := bufio.NewScanner(file)

	outLines := make([]string,0)

	// revise atom types
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

	//fmt.Println("file = " + path)

	atoms := make(map[string]*atom)
	// read as molecule and edit bonds
	for _, nextLine := range outLines {
		tokens := strings.Fields(nextLine)
		if len(tokens) > 4 {
			// create new atom
			var newAtom atom
			newAtom.id = tokens[0]
			newAtom.element = tokens[1]
			//fmt.Println(newAtom.element)
			pos := make([]float64,3)
			for j := 2; j < 5; j++ {
				pos[j-2], err = strconv.ParseFloat(tokens[j],64)
				if err != nil {
					newErr := errors.New("Failed to convert \"" + tokens[j] + "\" in position 0 on line " + strconv.Itoa(j) + " to a float64")
					log.Fatal(newErr)
				}
			}
			newAtom.pos = pos
			newAtom.atomType = tokens[5]
			newAtom.bonds = make([]string,0)
			if len(tokens) > 6 {
				for k := 6; k < len(tokens); k++ {
					newAtom.bonds = append(newAtom.bonds, tokens[k])
				}
			}
			atoms[newAtom.id] = &newAtom
		}
	}
	for _, atom := range atoms {
		for _, bond := range atom.bonds {
			if bond == "1" {
				disconnect(atoms,"1", atom.id)
			}
		}
	}


	/*fmt.Println("slice length = " + strconv.Itoa(len(atoms)))
	for id, atom := range atoms {
		fmt.Println("id = " + id + " element = " + atom.element)
	}*/

	atomSlice := make([]atom, len(atoms))
	for id, atom := range atoms {
		intID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Failed ID conversion for " + id)
		}
		atomSlice[intID-1] = *atom
	}

	// Write OUT
	thisFile, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to create new fragment file: " + path)
		log.Fatal(err)
	}
	_, _ = thisFile.WriteString(outLines[0] + "\n")
	for i, _ := range atomSlice {
		line := atomSlice[i].id + "\t" + atomSlice[i].element + "\t" + fmt.Sprintf("%.6f",atomSlice[i].pos[0]) + "\t" +
			fmt.Sprintf("%.6f",atomSlice[i].pos[1]) + "\t" + fmt.Sprintf("%.6f",atomSlice[i].pos[2]) + "\t" +
			atomSlice[i].atomType
		for _, bondedAtom := range atomSlice[i].bonds {
			line += "\t" + bondedAtom
		}

		_, err = thisFile.WriteString(line + "\n")
	}


}

// removes each atom from the other's list of bonded atoms
func disconnect(atoms map[string]*atom, atom1 string, atom2 string) {
	for i, bondedAtom := range atoms[atom1].bonds {
		if bondedAtom == atom2 {
			atoms[atom1].bonds = remove(atoms[atom1].bonds,i)
		}
	}
	for i, bondedAtom := range atoms[atom2].bonds {
		if bondedAtom == atom1 {
			atoms[atom2].bonds = remove(atoms[atom2].bonds,i)
		}
	}
}

// removes element from slice
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
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
	//fmt.Println("Reading file at " + filePath)
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



func QMEnergyAssembler(dir string) []float64 {

	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("failed to read directory: " + dir)
		log.Fatal(err)
	}
	QmEnergies := make([]float64,0)
	for _, file := range fileInfo {
		path := filepath.Join(dir, file.Name())
		nameFields := strings.Split(file.Name(),".")
		if nameFields[len(nameFields)-1] == "dat" {
			CBSenergy := getCBSEnergy(path)
			QmEnergies = append(QmEnergies, CBSenergy)
		}
	}

	outPath := filepath.Join(dir,"QM-energy.dat")
	thisFile, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Failed to create new fragment file: " + outPath)
		log.Fatal(err)
	}
	for _, i := range QmEnergies {
		_,_ = thisFile.WriteString(fmt.Sprintf("%.6f", i) + "\n")
	}
	_ = thisFile.Close()
	return QmEnergies
}

func getCBSEnergy(path string) float64 {
	// open file
	file, err := os.Open(path)
	//fmt.Println("Reading file at " + path)
	if err != nil {
		fmt.Println("Failed to open molecule file: " + path)
		log.Fatal(err)
	}

	// Initialize scanner
	scanner := bufio.NewScanner(file)

	e := 0.0

	for scanner.Scan() {
		// get next line
		line := scanner.Text()
		// split by whitespace
		tokens := strings.Fields(line)
		if len(tokens) == 3 && tokens[0] == "total" && tokens[1] == "CBS" {
			e, err = strconv.ParseFloat(tokens[2], 64)
		}
	}
	_ = file.Close()
	return e
}
