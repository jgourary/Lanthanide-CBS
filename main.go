package main

import (
	"fmt"
	"math"
	"os"
)

const memory = "4 GB"
const basis = "def2-TZVPD"
const energy = "MP2/def2-[TQ]ZVPD"
const mode = "manual"
func getShiftDistances(eqDistance float64) []float64 {
	samplePts := []float64{-0.2, -0.1, 0, 0.1, 0.2, 0.3, 0.4, 0.5}
	for i, pt := range samplePts {
		samplePts[i] = pt * eqDistance
	}
	return samplePts
}
const pointsEitherSide = 4 // 4 -> total points = 1 + 4 + 4 = 9


func main() {
	inputFile := os.Args[1]
	outputDir := os.Args[2]
	keyDir := os.Args[3]
	structName, molecules := readFile(inputFile, false)
	var ion molecule
	var ligand molecule
	if len(molecules[1].atoms) == 1 {
		ligand = molecules[0]
		ion = molecules[1]
	} else {
		ligand = molecules[1]
		ion = molecules[0]
	}
	fmt.Println("ION: ")
	printMolecule(ion)
	fmt.Println("LIGAND: ")
	printMolecule(ligand)

	var unitAxis []float64
	var equilibDistance float64
	if mode != "manual" {
		unitAxis, equilibDistance = getUnitAxis(ion, ligand, -1)
	} else {
		unitAxis = []float64{-0.99998991,0.004261549,-0.001418643}
		equilibDistance = 2.669451926
	}

	fmt.Println("\nUnit axis = [" + fmt.Sprintf("%.6f", unitAxis[0]) + " " +
		fmt.Sprintf("%.6f", unitAxis[1])  + " " + fmt.Sprintf("%.6f", unitAxis[2]) + "], len = " + fmt.Sprintf("%.6f", equilibDistance))


	ligands := generateModifiedStructures(ligand, unitAxis, equilibDistance)
	fmt.Println("New Ligands: ")
	for i, _ := range ligands {
		newUnitAxis, newEquilibDistance := getUnitAxis(ion, ligands[i], -1)
		fmt.Println(fmt.Sprintf("%.3f", ligands[i].shift) + ": Unit axis = [" + fmt.Sprintf("%.6f", newUnitAxis[0]) + " " +
			fmt.Sprintf("%.6f", newUnitAxis[1])  + " " + fmt.Sprintf("%.6f", newUnitAxis[2]) + "], len = " + fmt.Sprintf("%.6f", newEquilibDistance))
	}
	writeInputs(ion, ligands, outputDir, structName, keyDir)

}



func generateModifiedStructures(ligand molecule, unitAxis []float64, equilibDistance float64) []molecule {

	shiftDistances := getShiftDistances(equilibDistance)
	ligands := make([]molecule, len(shiftDistances))
	for i, _ := range shiftDistances {
		ligands[i] = copyMolecule(ligand)
	}
	// for every positioning of the ligand...
	for i, _ := range shiftDistances {
		for j := 0; j < len(ligands[i].atoms); j++  {
			ligands[i].atoms[j].pos[0] += unitAxis[0] * shiftDistances[i]
			ligands[i].atoms[j].pos[1] += unitAxis[1] * shiftDistances[i]
			ligands[i].atoms[j].pos[2] += unitAxis[2] * shiftDistances[i]
		}
		// save ligand
		ligands[i].shift = shiftDistances[i]
	}
	return ligands
}

func getUnitAxis(ionMol molecule, ligandMol molecule, ligandIndex int) ([]float64, float64) {
	var startPos []float64
	var endPos []float64
	closestDist := 1e100

	if ligandIndex < 0 {
		startPos = ionMol.atoms[0].pos
		endPos = make([]float64, 3)

		for _, atom := range ligandMol.atoms {
			dist := getDistance(startPos, atom.pos)
			if  dist < closestDist {
				endPos = atom.pos
				closestDist = dist
			}
		}
	} else {
		startPos = ionMol.atoms[0].pos
		endPos = ligandMol.atoms[ligandIndex].pos
		closestDist = getDistance(startPos, endPos)
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