package main

type molecule struct {
	charge string
	multiplicity string
	atoms map[int]*atom
}

type atom struct {
	element string
	pos []float64

}