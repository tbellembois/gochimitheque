package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func main() {
	csvFile, _ := os.Open("cmr-clp-atp10-FR-total2.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	m := make(map[string]string)

	r := regexp.MustCompile("([0-9]+-[0-9]+-[0-9]+)")

	i := 0
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			global.Log.Fatal(error)
		}
		fmt.Println("----")
		fmt.Println(line[0])
		fmt.Println(line[len(line)-1])
		s := r.FindAllString(line[0], -1)
		fmt.Println(s)
		for _, c := range s {
			m[c] = line[len(line)-1]
		}
		i++
	}

	fout, _ := os.Create("/tmp/dat2")
	defer fout.Close()

	for k, v := range m {
		fout.WriteString(k + "," + v + "\n")
	}
	fout.Sync()
}
