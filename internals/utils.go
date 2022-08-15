package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	NumberOfAnts  int
	InitialName   string
	FinalName     string
	AllTracks     [][]string
	GodTracks     [][][]string
	TransientPath int
	TheIndex      int
	BGr           int
)

type Block struct {
	NameOfRoom           string
	Links                map[string]*Block
	HorizontalCoordinate int
	VerticalCoordinate   int
	Travelled            bool
}

func AllocateRoom(name string, x int, y int) *Block {
	return &Block{
		NameOfRoom:           name,
		Links:                map[string]*Block{},
		HorizontalCoordinate: x,
		VerticalCoordinate:   y,
		Travelled:            false,
	}
}

func (bl *Block) CreateLink(link *Block) {
	if _, linked := bl.Links[link.NameOfRoom]; !linked {
		bl.Links[link.NameOfRoom] = link
	}
}

func MoveAnts(path string) {
	if path != "" {
		res, err := strconv.Atoi(path)
		if err != nil {
			ErrorHandler("Please Input Whole Ants")
		}
		if res < 1 {
			ErrorHandler("There Are No Ants to Move")
		} else {
			NumberOfAnts = res
			return
		}
		ErrorHandler("Incorrect Formatting Of First Line")
	}
}

//Checked till here

func HandlingInput(data string) map[string]*Block {
	chartcreated := make(map[string]*Block)
	linescreated := strings.Split(data, "\n")

	for index, singleline := range linescreated {

		if singleline != "" && singleline[0] == 35 {
			if singleline == "##start" { //Dont change this line
				d0 := strings.Split(linescreated[index+1], " ")
				InitialName = d0[0]
			}
			if singleline == "##end" { //Dont change this line
				f0 := strings.Split(linescreated[index+1], " ")
				FinalName = f0[0]
			}
			continue
		}

		if index == 0 {
			MoveAnts(singleline)
		}

		singleroom := strings.Split(singleline, " ")
		// check the room formatting
		if len(singleroom) == 3 {
			if _, cache := chartcreated[singleroom[0]]; !cache {
				h, err1 := strconv.Atoi(singleroom[1])
				v, err2 := strconv.Atoi(singleroom[2])
				if err1 != nil || err2 != nil {
					fmt.Println("coordinate formatted wrong, " + singleroom[0] + " skipped")
					continue
				}
				chartcreated[singleroom[0]] = AllocateRoom(singleroom[0], h, v)
			} else {
				fmt.Println("already exists,  " + singleroom[0] + " skipped")
				continue
			}
		}

		connections := strings.Split(singleline, "-")
		// tunnel format checking
		if len(connections) == 2 {
			if _, err1 := chartcreated[connections[0]]; !err1 {
				ErrorHandler("Not pointed from a Block " + connections[0])
			}
			if _, err2 := chartcreated[connections[1]]; !err2 {
				ErrorHandler("Not pointed to a Block " + connections[1])
			}
			if connections[0] != connections[1] {
				chartcreated[connections[0]].CreateLink(chartcreated[connections[1]])
				chartcreated[connections[1]].CreateLink(chartcreated[connections[0]])
			} else {
				ErrorHandler("Connected to room " + connections[0])
			}
		}
	}
	return chartcreated
}

// tested till here ---->

func FilterArray(data [][]string) [][]string {
	var array []string
	for l := range data {
		array = append(array, strings.Join(data[l], " . "))
	}
	for item := range array {
		for inner := range array {
			if inner == item {
				continue
			} else if strings.Contains(array[inner], array[item]) {
				array[inner] = "pathnotfeasible"
			}
		}
	}

	var arr2 []string
	for item1 := range array {
		if array[item1] != "pathnotfeasible" {
			arr2 = append(arr2, array[item1])
		}
	}

	for l := 0; l < len(arr2)-1; l++ {
		if len(arr2[l]) > len(arr2[l+1]) {
			arr2[l], arr2[l+1] = arr2[l+1], arr2[l]
			l = -1
		}
	}

	var arr3 [][]string
	for l := range arr2 {
		arr3 = append(arr3, strings.Split(arr2[l], " . "))
	}
	return arr3
}

func ChosePath() {
	singleline := make([][]int, 0)

	for i, last := range GodTracks {

		var antLine [][]int
		for range last {
			antLine = append(antLine, make([]int, 0))
		}

		for l := 1; l <= NumberOfAnts; l++ {
			k := 10000
			for z := range last {
				token := len(last[z]) + len(antLine[z])
				if k > token {
					k = token
					BGr = z
				} else {
					continue
				}
			}
			antLine[BGr] = append(antLine[BGr], l)
		}
		move := 0
		for l := range antLine {
			token := len(last[l]) + len(antLine[l])
			if move == 0 || move < token {
				move = token
			}
		}
		if TransientPath == 0 || move < TransientPath {
			TransientPath = move
			TheIndex = i
			singleline = antLine
		}
	}
	HandlingFinalOutput(GodTracks[TheIndex], TransientPath, singleline)
}

//Checked till here

func Composition(segments [][]string) {
	for i := range segments {
		var combos [][]string
		combos = append(combos, segments[i])
		for k := range segments {
			if i == k || CrashHandler(combos, segments[k]) {
				continue
			} else {
				combos = append(combos, segments[k])
			}
		}
		GodTracks = append(GodTracks, combos)
	}
}

func CrashHandler(data [][]string, H []string) bool {
	for l := range data {
		for x := range H {
			for ii := range data[l] {
				if H[x] == data[l][ii] {
					return true
				}
			}
		}
	}
	return false
}

func DepthFirstSearch(que string, track []string, chart map[string]*Block) {
	chart[que].Travelled = true
	if que != InitialName && que != FinalName {
		track = append(track, chart[que].NameOfRoom)
	}
	if que == FinalName {
		temp := make([]string, len(track))
		copy(temp, track)
		AllTracks = append(AllTracks, temp)
	} else {
		for k, cont := range chart[que].Links {
			if !cont.Travelled {
				DepthFirstSearch(k, track, chart)
			}
		}
	}
	chart[que].Travelled = false
}

func HandlingFinalOutput(results [][]string, turns int, antID [][]int) {
	set := make([][]string, turns)
	copy(set, results)
	if results[0][0] == "" {
		var last string
		for a := range antID[0] {
			last += "L" + strconv.Itoa(antID[0][a]) + "-" + FinalName + " "
		}
		fmt.Println(last)
		return
	}
	for l := range results {
		set[l] = append(set[l], FinalName)
	}

	for r := 0; r < turns; r++ {
		var last string
		for d := range results {
			for a := range antID[d] {
				plc := r - a
				if plc < 0 || plc > len(results[d]) {
					continue
				} else {
					last += "L" + strconv.Itoa(antID[d][a]) + "-" + set[d][plc] + " "
				}
			}
		}
		fmt.Println(last)
	}
}

func ErrorHandler(message string) {
	fmt.Printf("Error: %s\n", message)
	os.Exit(1)
}
