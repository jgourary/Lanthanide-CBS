package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"strings"
)

const gaussian = true
const psi4 = true

const memory = "200 GB"
const basis = "def2-QZVPD"
const energy = "cbs, corl_wfn='mp2',corl_basis='def2-[TQ]ZVPD', delta_wfn='ccsd(t)', delta_basis='def2-[DT]ZVPD'"
const ionElement = "La"

const keyDir = "lib\\key"
const outputDir = "lib\\output"
const inputDIr = "lib\\input"

func getShiftDistances(eqDistance float64) []float64 {
	samplePts := []float64{-1.0, -0.9, -0.8, -0.7, -0.6, -0.5, -0.4, -0.3, -0.2, -0.1, 0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0, 5.5, 6.0}
	/*for i, pt := range samplePts {
		samplePts[i] = pt + eqDistance
	}*/
	return samplePts
}

/*func getCostScaling() {
	costScaling := []float64{-1.0, -0.8, -0.6, -0.4, -0.2, 0, 0.2, 0.4, 0.6, 0.8, 1.0, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0}
}*/

func main() {
	fileInfo, err := ioutil.ReadDir(inputDIr)
	if err != nil {
		fmt.Println("failed to read directory: " + inputDIr)
		log.Fatal(err)
	}
	for _, file := range fileInfo {
		name := file.Name()
		path := filepath.Join(inputDIr, name)
		if filepath.Ext(name) == ".inp" {
			fmt.Println("Reading file at " + path)
			structName, molecules := readFile(path, false)
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
			if strings.Contains(name, "bidentate") {
				unitAxis = []float64{-0.99998991, 0.004261549, -0.001418643}
				equilibDistance = 2.669451926

			} else {
				unitAxis, equilibDistance = getUnitAxis(ion, ligand, -1)
			}

			fmt.Println("\nUnit axis = [" + fmt.Sprintf("%.6f", unitAxis[0]) + " " +
				fmt.Sprintf("%.6f", unitAxis[1]) + " " + fmt.Sprintf("%.6f", unitAxis[2]) + "], len = " + fmt.Sprintf("%.6f", equilibDistance))

			ligands := generateModifiedStructures(ligand, unitAxis, equilibDistance)
			fmt.Println("New Ligands: ")
			for i, _ := range ligands {
				newUnitAxis, newEquilibDistance := getUnitAxis(ion, ligands[i], -1)
				fmt.Println("Shift Dist: " + fmt.Sprintf("%.3f", ligands[i].shift) + ": Unit axis = [" + fmt.Sprintf("%.6f", newUnitAxis[0]) + " " +
					fmt.Sprintf("%.6f", newUnitAxis[1]) + " " + fmt.Sprintf("%.6f", newUnitAxis[2]) + "], ion-ligand dist = " + fmt.Sprintf("%.6f", newEquilibDistance))
			}

			// Write XYZ, TXYZ, INP files + filelist
			writeInputs(ion, ligands, outputDir, structName, keyDir)
		}
	}

	// Assemble QM-Energy.dat from all .dat files in the directory
	// QMEnergyAssembler(outputDir)

	// Assemble an output CSV from QM-Energy.dat, result.p, filelist
	// makeCSV(outputDir)
}

func generateModifiedStructures(ligand molecule, unitAxis []float64, equilibDistance float64) []molecule {

	shiftDistances := getShiftDistances(equilibDistance)
	ligands := make([]molecule, len(shiftDistances))
	for i, _ := range shiftDistances {
		ligands[i] = copyMolecule(ligand)
	}
	// for every positioning of the ligand...
	for i, _ := range shiftDistances {
		for j := 0; j < len(ligands[i].atoms); j++ {
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
			if dist < closestDist {
				endPos = atom.pos
				closestDist = dist
			}
		}
	} else {
		startPos = ionMol.atoms[0].pos
		endPos = ligandMol.atoms[ligandIndex].pos
		closestDist = getDistance(startPos, endPos)
	}
	vector := []float64{(endPos[0] - startPos[0]) / closestDist, (endPos[1] - startPos[1]) / closestDist, (endPos[2] - startPos[2]) / closestDist}
	return vector, closestDist

}

func getDistance(pos1 []float64, pos2 []float64) float64 {
	dx2 := math.Pow(pos1[0]-pos2[0], 2)
	dy2 := math.Pow(pos1[1]-pos2[1], 2)
	dz2 := math.Pow(pos1[2]-pos2[2], 2)
	return math.Sqrt(dx2 + dy2 + dz2)
}
