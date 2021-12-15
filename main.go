package main

import (
	"math"
	"os"
)

const memory = "4 GB"
const basis = "def2-TZVPD"
const energy = "MP2/def2-[TQ]ZVPD"
const movePercent = 0.1 // 0.1 = 10%
const pointsEitherSide = 4 // 4 -> total points = 1 + 4 + 4 = 9


func main() {
	inputFile := os.Args[0]
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

	unitAxis := getUnitAxis(ion, ligand)


}

func generateModifiedStructures(unitAxis []float64) {

}

func getUnitAxis(ionMol molecule, ligandMol molecule) []float64 {

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
	return vector
}

func getDistance(pos1 []float64, pos2 []float64) float64 {
	dx2 := math.Pow(pos1[0] - pos2[0],2)
	dy2 := math.Pow(pos1[1] - pos2[1],2)
	dz2 := math.Pow(pos1[2] - pos2[2],2)
	return math.Sqrt(dx2+dy2+dz2)
}