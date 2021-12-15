package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

const memory = "4 GB"
const basis = "def2-TZVPD"
const energy = "MP2/def2-[TQ]ZVPD"
func getSamplePoints() []float64 {
	return []float64{-0.2, -0.1, 0, 0.1, 0.25, 0.5, 0.75, 1} // 0 = everything in same spot, 1 = ligand moved to double equilibrium distance
}
const pointsEitherSide = 4 // 4 -> total points = 1 + 4 + 4 = 9


func main() {
	inputFile := os.Args[0]
	outputDir := os.Args[1]
	structName, molecules := readFile(inputFile)
	var ion molecule
	var ligand molecule
	if len(molecules[1].atoms) == 1 {
		ligand = molecules[0]
		ion = molecules[1]
	} else {
		ligand = molecules[1]
		ion = molecules[0]
	}

	unitAxis, equilibDistance := getUnitAxis(ion, ligand)
	ligands := generateModifiedStructures(ligand, unitAxis, equilibDistance)

}

func writeInputs(ion molecule, ligands []molecule) {
	for _, ligand := range ligands {

	}
}

func writeInput(ion molecule, ligand molecule, outDir string, structName string, i int) {
	outPath := filepath.Join(outDir, structName + "_" + strconv.Itoa(i) + ".inp")
	_ = os.MkdirAll(outDir, 0755)
	thisFile, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Failed to create new fragment file: " + outPath)
		log.Fatal(err)
	}
	_, _ = thisFile.WriteString("memory " + memory + "\n")
	_, _ = thisFile.WriteString("set basis " + basis + "\n")
	_, _ = thisFile.WriteString("molecule{\n")
	_, _ = thisFile.WriteString("\t" + ion.charge + " " + ion.multiplicity + "\n")
	for _, atom := range ion.atoms {
		_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprint("%.6f", atom.pos[0]) + " " +
		fmt.Sprint("%.6f", atom.pos[0])  + " " + fmt.Sprint("%.6f", atom.pos[2])  + "\n")
	}
	_, _ = thisFile.WriteString("\t--\n")
}

func generateModifiedStructures(ligand molecule, unitAxis []float64, equilibDistance float64) []molecule {
	var ligands []molecule
	samplePoints := getSamplePoints()
	// for every positioning of the ligand...
	for _, samplePoint := range samplePoints {
		// create new copy of ligand
		newLigand := copyMolecule(ligand)
		// move every atom to new position
		for _, thisAtom := range newLigand.atoms {
			thisAtom.pos[0] += unitAxis[0] * samplePoint * equilibDistance
			thisAtom.pos[1] += unitAxis[1] * samplePoint * equilibDistance
			thisAtom.pos[2] += unitAxis[2] * samplePoint * equilibDistance
		}
		// save ligand
		newLigand.shift = samplePoint
		ligands = append(ligands, newLigand)
	}
	return ligands
}

func getUnitAxis(ionMol molecule, ligandMol molecule) ([]float64, float64) {

	startPos := ionMol.atoms[0].pos

	endPos := make([]float64, 3)
	closestDist := 1e100
	for _, atom := range ligandMol.atoms {
		dist := getDistance(startPos, atom.pos)
		if  dist < closestDist {
			endPos = atom.pos
			closestDist = dist
		}
	}

	vector := []float64{(endPos[0]-startPos[0])/closestDist, (endPos[1]-startPos[1]) / closestDist, (endPos[2]-startPos[2])/closestDist}
	return vector, closestDist
}

func getDistance(pos1 []float64, pos2 []float64) float64 {
	dx2 := math.Pow(pos1[0] - pos2[0],2)
	dy2 := math.Pow(pos1[1] - pos2[1],2)
	dz2 := math.Pow(pos1[2] - pos2[2],2)
	return math.Sqrt(dx2+dy2+dz2)
}