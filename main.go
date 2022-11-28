package main

import (
	"encoding/csv"
	"fmt"
	"go/types"
	"os"
	"strconv"
	"strings"
)

type Serie struct {
	dataType types.Type
	values   []interface{}
}

type DataFrame struct {
	header []interface{}
	data   [][]interface{}
}

func (d *DataFrame) categorizeCols(columns []interface{}) map[string][]interface{} {
	result := make(map[string][]interface{})
	for _, column := range columns {
		columnStr, ok := column.(string)
		if !ok {
			columnInt, ok := column.(int)
			if !ok || columnInt < 0 || columnInt >= len(d.header) {
				continue
			}
			columnStr = d.header[columnInt].(string)
		}
		result[columnStr] = make([]interface{}, 0)
	}
	for _, row := range d.data {
		for index, value := range row {
			valueStr, ok := value.(string)
			if ok && valueStr == "" {
				row[index] = nil
				continue
			}
			if _, ok := result[d.header[index].(string)]; ok {
				label := SliceIndex(result[d.header[index].(string)], value)
				if label == -1 {
					result[d.header[index].(string)] = append(result[d.header[index].(string)], value)
					label = len(result[d.header[index].(string)]) - 1
				}
				row[index] = label
			}
		}
	}
	return result
}

func (d *DataFrame) getSerie(name string) []interface{} {
	index := SliceIndex(d.header, name)
	if index != -1 {
		return d.data[index]
	} else {
		return nil
	}

}

func (d *DataFrame) toCSV(fname string) {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		println(err.Error())
		return
	}
	csvW := csv.NewWriter(f)
	rec := make([]string, 0)
	for _, header := range d.header {
		rec = append(rec, header.(string))
	}
	csvW.Write(rec)
	for _, record := range d.data {
		rec := make([]string, 0)
		for _, col := range record {
			if col == nil {
				rec = append(rec, "")
				continue
			}
			str, ok := col.(string)
			if ok {
				rec = append(rec, str)
			} else {
				rec = append(rec, strconv.Itoa(col.(int)))
			}
		}
		err := csvW.Write(rec)
		if err != nil {
			println(err.Error())
		}
	}
	csvW.Flush()
	f.Close()
}

func SliceIndex(slice []interface{}, v interface{}) int {
	for index, value := range slice {
		if value == v {
			return index
		}
	}
	return -1
}

func main() {
	wd, _ := os.Getwd()
	trainFile, _ := os.OpenFile(wd+"/db/train.csv", os.O_RDONLY, os.ModePerm)
	frame := DataFrame{
		data:   make([][]interface{}, 0),
		header: make([]interface{}, 0),
	}
	reader := csv.NewReader(trainFile)
	records, _ := reader.ReadAll()
	for index, record := range records {
		if index != 0 {
			frame.data = append(frame.data, make([]interface{}, 0))
		}
		passGroup := 0
		cabinNum := 0
		cabinSize := ""
		for colNum, field := range record {
			if index == 0 {
				frame.header = append(frame.header, field)
			} else {
				if frame.header[colNum].(string) == "PassengerId" {
					parts := strings.Split(field, "_")
					passGroup, _ = strconv.Atoi(parts[0])
					v, _ := strconv.Atoi(parts[1])
					field = strconv.Itoa(v)
				} else if frame.header[colNum].(string) == "Cabin" {
					if field == "" {
						cabinNum = 0
						cabinSize = ""
						continue
					}
					parts := strings.Split(field, "/")
					field = parts[0]
					cabinNum, _ = strconv.Atoi(parts[1])
					cabinSize = parts[2]
				}
				frame.data[index-1] = append(frame.data[index-1], field)
			}
		}
		if index != 0 {
			frame.data[index-1] = append(frame.data[index-1], passGroup)
			frame.data[index-1] = append(frame.data[index-1], cabinNum)
			frame.data[index-1] = append(frame.data[index-1], cabinSize)
		} else {
			frame.header = append(frame.header, "Passenger_group")
			frame.header = append(frame.header, "Cabin_num")
			frame.header = append(frame.header, "Cabin_size")
		}
	}
	frame.toCSV("out.csv")
	labels := frame.categorizeCols([]interface{}{"HomePlanet", "CryoSleep", "Destination", "VIP", "Transported", "Cabin", "Cabin_size"})
	for column, lbls := range labels {
		println(column)
		for index, value := range lbls {
			fmt.Printf("%d:%s\n", index, value)
		}
	}
	frame.toCSV("out2.csv")
}
