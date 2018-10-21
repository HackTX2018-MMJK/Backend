package main

import (
	"encoding/json"
	"flag"
	"github.com/dchest/uniuri"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type DataFile struct {
	Mappings map[string]string `json:"mappings"`
}

func main() {
	/* Load Arguments */
	log.Println("Loading Arguments.")
	var rows int
	var seats int
	flag.IntVar(&rows, "rows", -1, "Number of rows to generate")
	flag.IntVar(&seats, "seats", -1, "Number of seats per row")

	var fileName string
	flag.StringVar(&fileName, "file", "data.json", "File for seat-mappings")

	flag.Parse()

	/* Load data file */
	log.Println("Attempting to load data file...")
	var dataFile DataFile
	file, err := os.Open(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		} else {
			/* Create missing file */
			log.Println("No data file. Creating one now")
			file, err = os.Create(fileName)
			checkErr(err)
		}
	} else {
		/* Load Data from File */
		log.Println("Data file found. Loading values")
		jsonData, err := ioutil.ReadAll(file)
		checkErr(err)
		err = json.Unmarshal(jsonData, &dataFile)
		checkErr(err)
	}

	/* Handle Arguments */
	if rows != -1 && seats != -1 {
		generate_mapping(rows, seats, fileName)
	} else {
		handle_requests(dataFile.Mappings)
	}
}

func generate_mapping(rows int, seats int, out string) {
	log.Printf("Generating mappings for input, rows=%d seats=%d\n", rows, seats)
	mapping := make(map[string]string)
	BASE := 64
	for row := 1; row <= rows; row++ {
		for seat := 1; seat <= seats; seat++ {
			mapping[uniuri.NewLen(6)] = strconv.Itoa(row) + string(seat+BASE)
		}
	}
	var dataFile DataFile
	dataFile.Mappings = mapping

	jsonData, err := json.Marshal(dataFile)
	checkErr(err)

	err = ioutil.WriteFile(out, jsonData, 0644)
	checkErr(err)
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
