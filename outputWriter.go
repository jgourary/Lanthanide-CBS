package main

import (
	"bufio"
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

	if psi4 {
		// Write INP
		outPath := filepath.Join(outDir, structName+"_"+index2suffix(i)+".inp")
		_ = os.MkdirAll(outDir, 0755)
		thisFile, err := os.Create(outPath)
		if err != nil {
			fmt.Println("Failed to create new fragment file: " + outPath)
			log.Fatal(err)
		}
		_, _ = thisFile.WriteString("memory " + memory + "\n")
		_, _ = thisFile.WriteString("set {\n")
		_, _ = thisFile.WriteString("\tbasis " + basis + "\n")
		_, _ = thisFile.WriteString("\tsoscf true\n")
		_, _ = thisFile.WriteString("\tfail_on_maxiter false\n")
		_, _ = thisFile.WriteString("\trelativistic dkh\n")
		_, _ = thisFile.WriteString("\tdkh_order 2\n")
		_, _ = thisFile.WriteString("}\n\n")
		_, _ = thisFile.WriteString("molecule{\n")
		_, _ = thisFile.WriteString("\t" + ion.charge + " " + ion.multiplicity + "\n")
		for _, atom := range ion.atoms {
			_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
				fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
		}
		_, _ = thisFile.WriteString("\t--\n")
		_, _ = thisFile.WriteString("\t" + ligand.charge + " " + ligand.multiplicity + "\n")
		for _, atom := range ligandSlice {
			_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
				fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
		}
		_, _ = thisFile.WriteString("}\n\n")
		_, _ = thisFile.WriteString("energy(" + energy + ")\n")
	}

	// Write XYZ
	outPath := filepath.Join(outDir, structName+"_"+index2suffix(i)+".xyz")
	thisFile, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Failed to create new fragment file: " + outPath)
		log.Fatal(err)
	}
	_, _ = thisFile.WriteString(strconv.Itoa(len(ligand.atoms)+1) + "\n\n")
	for _, atom := range ion.atoms {
		_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
			fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
	}
	for _, atom := range ligandSlice {
		_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
			fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
	}
	if gaussian {
		// Write Counterpoise GJF
		/*outPath = filepath.Join(outDir, structName+"_"+index2suffix(i)+"_cp.gjf")
		thisFile, err = os.Create(outPath)
		if err != nil {
			fmt.Println("Failed to create new fragment file: " + outPath)
			log.Fatal(err)
		}
		_, _ = thisFile.WriteString("%Mem=16GB\n")
		_, _ = thisFile.WriteString("%NProc=8\n")
		_, _ = thisFile.WriteString("%Chk=" + structName + "_" + index2suffix(i) + "_cp.chk\n\n")
		_, _ = thisFile.WriteString("#p MP2/Def2QZVP wB97XD nosymm MaxDisk=1000GB Counterpoise=2\n\n")
		_, _ = thisFile.WriteString(structName + "_" + index2suffix(i) + "\n\n")

		atomCharge, err := strconv.Atoi(ion.charge)
		ligandCharge, err := strconv.Atoi(ligand.charge)
		_, _ = thisFile.WriteString("\t" + strconv.Itoa(atomCharge+ligandCharge) + "," + ligand.multiplicity + " " + ion.charge + "," + ion.multiplicity + " " + ligand.charge + "," + ligand.multiplicity + "\n")
		for _, atom := range ion.atoms {
			_, _ = thisFile.WriteString("\t" + atom.element + "(Fragment=1) " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
				fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
		}
		for _, atom := range ligandSlice {
			_, _ = thisFile.WriteString("\t" + atom.element + "(Fragment=2) " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
				fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
		}
		_, _ = thisFile.WriteString("\n\n")

		// Write Guess GJF
		outPath = filepath.Join(outDir, structName+"_"+index2suffix(i)+"_g.gjf")
		thisFile, err = os.Create(outPath)
		if err != nil {
			fmt.Println("Failed to create new fragment file: " + outPath)
			log.Fatal(err)
		}
		_, _ = thisFile.WriteString("%Mem=16GB\n")
		_, _ = thisFile.WriteString("%NProc=8\n")
		_, _ = thisFile.WriteString("%Chk=" + structName + "_" + index2suffix(i) + ".chk\n\n")
		_, _ = thisFile.WriteString("#MP2/Def2SVP SCF(maxcycle=100) Polar Density=MP2 MaxDisk=100GB\n\n")
		_, _ = thisFile.WriteString(structName + "_" + index2suffix(i) + "\n\n")

		atomCharge, err = strconv.Atoi(ion.charge)
		ligandCharge, err = strconv.Atoi(ligand.charge)
		_, _ = thisFile.WriteString("\t" + strconv.Itoa(atomCharge+ligandCharge) + "," + ligand.multiplicity + "\n")
		for _, atom := range ion.atoms {
			_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
				fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
		}
		for _, atom := range ligandSlice {
			_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
				fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
		}
		_, _ = thisFile.WriteString("\n\n")*/

		// Write GJF
		outPath = filepath.Join(outDir, structName+"_"+index2suffix(i)+".gjf")
		thisFile, err = os.Create(outPath)
		if err != nil {
			fmt.Println("Failed to create new fragment file: " + outPath)
			log.Fatal(err)
		}
		_, _ = thisFile.WriteString("%Mem=16GB\n")
		_, _ = thisFile.WriteString("%NProc=8\n")
		_, _ = thisFile.WriteString("%Chk=" + structName + "_" + index2suffix(i) + ".chk\n\n")
		_, _ = thisFile.WriteString("#MP2/Def2QZVP GUESS=READ SCF(maxcycle=100) Polar Density=MP2 MaxDisk=100GB\n\n")
		_, _ = thisFile.WriteString(structName + "_" + index2suffix(i) + "\n\n")

		atomCharge, _ := strconv.Atoi(ion.charge)
		ligandCharge, _ := strconv.Atoi(ligand.charge)
		_, _ = thisFile.WriteString("\t" + strconv.Itoa(atomCharge+ligandCharge) + "," + ligand.multiplicity + "\n")
		for _, atom := range ion.atoms {
			_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
				fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
		}
		for _, atom := range ligandSlice {
			_, _ = thisFile.WriteString("\t" + atom.element + " " + fmt.Sprintf("%.6f", atom.pos[0]) + " " +
				fmt.Sprintf("%.6f", atom.pos[1]) + " " + fmt.Sprintf("%.6f", atom.pos[2]) + "\n")
		}
		_, _ = thisFile.WriteString("\n\n")
		file, _ := os.Open(filepath.Join("lib", "basis_sets", "gaussian"))
		// fmt.Println("Reading file at " + filePath)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// get next line
			line := scanner.Text()
			_, _ = thisFile.WriteString(line)
		}
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
