package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func writeInputs(ion molecule, ligands []molecule, outDir string, structName string, keyDir string) {
	for i, ligand := range ligands {
		writeInput(ion, ligand, outDir, structName, i)
	}
	obabelConversion2(outDir, ".xyz", ".txyz", "no", false)
	keyChain := getKeyChain(keyDir)
	postProcessTXYZs(outDir, keyChain)
}

func writeInput(ion molecule, ligand molecule, outDir string, structName string, i int) {
	ligandSlice := make([]atom, len(ligand.atoms))
	for i, atom := range ligand.atoms {
		ligandSlice[i] = atom
	}

	// Write INP
	outPath := filepath.Join(outDir, structName + "_" + index2suffix(i) + ".inp")
	_ = os.MkdirAll(outDir, 0755)
	thisFile, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Failed to create new fragment file: " + outPath)
		log.Fatal(err)
	}
	_, _ = thisFile.WriteString("memory " + memory + "\n")
	_, _ = thisFile.WriteString("set basis " + basis + "\n")
	_, _ = thisFile.WriteString("set soscf true\n")
	_, _ = thisFile.WriteString("set fail_on_maxiter false\n")
	_, _ = thisFile.WriteString("molecule{\n")
	_, _ = thisFile.WriteString("\t" + ion.charge + " " + ion.multiplicity + "\n")
	for _, atom := range ion.atoms {
		_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
			fmt.Sprintf("%.6f", atom.pos[1])  + " " + fmt.Sprintf("%.6f", atom.pos[2])  + "\n")
	}
	_, _ = thisFile.WriteString("\t--\n")
	_, _ = thisFile.WriteString("\t" + ligand.charge + " " + ligand.multiplicity + "\n")
	for _, atom := range ligandSlice {
		_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
			fmt.Sprintf("%.6f", atom.pos[1])  + " " + fmt.Sprintf("%.6f", atom.pos[2])  + "\n")
	}
	_, _ = thisFile.WriteString("}\n\n")
	_, _ = thisFile.WriteString("energy('" + energy + "')\n")

	// Write XYZ
	outPath = filepath.Join(outDir, structName + "_" + index2suffix(i) + ".xyz")
	thisFile, err = os.Create(outPath)
	if err != nil {
		fmt.Println("Failed to create new fragment file: " + outPath)
		log.Fatal(err)
	}
	_, _ = thisFile.WriteString(strconv.Itoa(len(ligand.atoms)+1)+"\n\n")
	for _, atom := range ion.atoms {
		_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
			fmt.Sprintf("%.6f", atom.pos[1])  + " " + fmt.Sprintf("%.6f", atom.pos[2])  + "\n")
	}
	for _, atom := range ligandSlice {
		_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
			fmt.Sprintf("%.6f", atom.pos[1])  + " " + fmt.Sprintf("%.6f", atom.pos[2])  + "\n")
	}
}

func index2suffix(i int) string {
	iStr := strconv.Itoa(i)
	if i < 10 {
		return "00" + iStr
	} else if i < 100 {
		return "0" + iStr
	} else {
		return iStr
	}
}

