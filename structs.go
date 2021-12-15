package main

type molecule struct {
	charge string
	multiplicity string
	atoms map[int]*atom

	shift float64
}

type atom struct {
	element string
	pos []float64
}

func copyAtom(oldAtom atom) atom {
	var newAtom atom
	newAtom.element = oldAtom.element
	newAtom.pos = oldAtom.pos
	return newAtom
}

func copyMolecule(oldMol molecule) molecule {
	var newMol molecule

	newMol.charge = oldMol.charge
	newMol.multiplicity = oldMol.multiplicity

	newMolAtoms := make(map[int]*atom)
	newMol.atoms = newMolAtoms
	for name, oldAtom := range oldMol.atoms {
		newAtom := copyAtom(*oldAtom)
		newMol.atoms[name] = &newAtom
	}
	return newMol
}