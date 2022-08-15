package main

import (
	"fmt"
	"os"
	"project/internals"
)

// Main Function starts from here

func main() {

	testCase := os.Args[1]
	readCase, _ := os.ReadFile(testCase)

	tempG := utils.HandlingInput(string(readCase))

	fmt.Printf("start room: %v \n", utils.InitialName)
	fmt.Printf("end room: %v \n", utils.FinalName)
	fmt.Printf("number of ants: %v \n", utils.NumberOfAnts)

	trail := []string{}
	utils.DepthFirstSearch(utils.InitialName, trail, tempG)

	utils.Composition(utils.FilterArray(utils.AllTracks))

	utils.ChosePath()
}
