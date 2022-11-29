package main

import (
	"encoding/csv"
	"os"
	"path"
	"strconv"
	"strings"
)

var srcCols = map[string]int{
	"PassengerId":  0,
	"HomePlanet":   1,
	"CryoSleep":    2,
	"Cabin":        3,
	"Destination":  4,
	"Age":          5,
	"VIP":          6,
	"RoomService":  7,
	"FoodCourt":    8,
	"ShoppingMall": 9,
	"Spa":          10,
	"VRDeck":       11,
	"Name":         12,
	"Transported":  13,
}
var srcColsRev = map[int]string{}

var dstCols = map[string]int{
	"GroupId":      0,
	"PassengerId":  1,
	"HomePlanet":   2,
	"CryoSleep":    3,
	"Cabin_Deck":   4,
	"Cabin_Num":    5,
	"Cabin_Side":   6,
	"Destination":  7,
	"Age":          8,
	"VIP":          9,
	"RoomService":  10,
	"FoodCourt":    11,
	"ShoppingMall": 12,
	"Spa":          13,
	"VRDeck":       14,
	"Transported":  15,
}

var data = make([][]float64, 0)

func main3() {
	makeRevs()
	readData()
	saveDataToCsv()
}

func readData() {
	dir, _ := os.Getwd()
	input, _ := os.OpenFile(path.Join(dir, "db", "test.csv"), os.O_RDONLY, os.ModePerm)
	reader := csv.NewReader(input)
	reader.Read()
	for {
		row, _ := reader.Read()
		if row == nil {
			break
		}
		destRow := make([]float64, 16)
		for index, value := range row {
			col, ok := srcColsRev[index]
			if !ok {
				continue
			}
			switch col {
			case "CryoSleep":
				setBoolCol(row, destRow, col, value)
			case "Transported":
				setBoolCol(row, destRow, col, value)
			case "VIP":
				setBoolCol(row, destRow, col, value)
			case "HomePlanet":
				setLabelCol(row, destRow, col, value)
			case "Destination":
				setLabelCol(row, destRow, col, value)
			case "PassengerId":
				setPassengerId(row, destRow, col, value)
			case "Cabin":
				setCabin(row, destRow, col, value)
			default:
				setFloatCol(row, destRow, col, value)
			}
		}
		data = append(data, destRow)
	}
}

func saveDataToCsv() {
	dir, _ := os.Getwd()
	output, _ := os.OpenFile(path.Join(dir, "db", "out2.csv"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	writer := csv.NewWriter(output)
	destRow := make([]string, len(dstCols))
	for col, index := range dstCols {
		destRow[index] = col
	}
	writer.Write(destRow)
	for _, row := range data {
		for index, col := range row {
			destRow[index] = strconv.FormatFloat(col, 'f', 2, 64)
		}
		writer.Write(destRow)
	}
	writer.Flush()
	output.Close()
}

func makeRevs() {
	for index, value := range srcCols {
		srcColsRev[value] = index
	}
}

func setCabin(row []string, destRow []float64, col string, value string) {
	if value == "" {
		destRow[dstCols["Cabin_Deck"]] = -1
		destRow[dstCols["Cabin_Num"]] = -1
		destRow[dstCols["Cabin_Side"]] = -1
		return
	}
	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		destRow[dstCols["Cabin_Deck"]] = -1
		destRow[dstCols["Cabin_Num"]] = -1
		destRow[dstCols["Cabin_Side"]] = -1
		return
	}
	labelsColection := getLabels()
	cabinDeck, ok := labelsColection["Cabin_Deck"][parts[0]]
	if !ok {
		destRow[dstCols["Cabin_Deck"]] = -1
	} else {
		destRow[dstCols["Cabin_Deck"]] = cabinDeck
	}
	cabinNum, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		destRow[dstCols["Cabin_Num"]] = -1
	} else {
		destRow[dstCols["Cabin_Num"]] = cabinNum
	}
	cabinSide, ok := labelsColection["Cabin_Side"][parts[2]]
	if !ok {
		destRow[dstCols["Cabin_Side"]] = -1
	} else {
		destRow[dstCols["Cabin_Side"]] = cabinSide
	}
}

func setPassengerId(row []string, destRow []float64, col string, value string) {
	if value == "" {
		destRow[dstCols["PassengerId"]] = -1
		destRow[dstCols["GroupId"]] = -1
		return
	}
	parts := strings.Split(value, "_")
	if len(parts) != 2 {
		destRow[dstCols["PassengerId"]] = -1
		destRow[dstCols["GroupId"]] = -1
		destRow[dstCols[col]] = -1
		return
	}
	groupId, err := strconv.Atoi(parts[0])
	if err != nil {
		destRow[dstCols["GroupId"]] = -1
	} else {
		destRow[dstCols["GroupId"]] = float64(groupId)
	}
	passengerId, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		destRow[dstCols["PassengerId"]] = -1
	} else {
		destRow[dstCols["PassengerId"]] = passengerId
	}
}

func setLabelCol(row []string, destRow []float64, col string, value string) {
	labelsColection := getLabels()
	labels, ok := labelsColection[col]
	if !ok {
		destRow[dstCols[col]] = -1
		return
	}
	label, ok := labels[value]
	if !ok {
		destRow[dstCols[col]] = -1
		return
	}
	destRow[dstCols[col]] = label
}

func setBoolCol(row []string, destRow []float64, col string, value string) {
	if value == "True" {
		destRow[dstCols[col]] = 1
	} else if value == "False" {
		destRow[dstCols[col]] = 2
	} else if value == "" {
		destRow[dstCols[col]] = -1
	}
}

func setFloatCol(row []string, destRow []float64, col string, value string) {
	if dstCols[col] == 0 {
		return
	}
	fValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		destRow[dstCols[col]] = -1
	} else {
		destRow[dstCols[col]] = fValue
	}
}

func getLabels() map[string]map[string]float64 {
	return map[string]map[string]float64{
		"HomePlanet": {
			"Earth":  1.0,
			"Europa": 2.0,
			"Mars":   3.0,
		},
		"Cabin_Deck": {
			"A": 1.0,
			"B": 2.0,
			"C": 3.0,
			"D": 4.0,
			"E": 5.0,
			"F": 6.0,
			"G": 7.0,
			"T": 8.0,
		},
		"Cabin_Side": {
			"P": 1.0,
			"S": 2.0,
		},
		"Destination": {
			"TRAPPIST-1e":   1.0,
			"PSO J318.5-22": 2.0,
			"55 Cancri e":   3.0,
		},
	}

}
