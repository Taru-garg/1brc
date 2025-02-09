package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
)

const (
	FILENAME = "measurements.txt"
)

type City struct {
	Min, Max, Sum float64
	TotalReading  int64
}

func main() {
	_f, err := os.Create("cpu.prof")
	if err := pprof.StartCPUProfile(_f); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()
	f, err := os.Open(FILENAME)
	if err != nil {
		log.Printf("Error while opening the file %s\n", FILENAME)
	}
	defer f.Close()

	cityStats := make(map[string]*City)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		city, temp, hasSemi := strings.Cut(line, ";")
		if !hasSemi {
			continue
		}
		ftemp, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			log.Printf("Error Parsing temp %s\n", temp)
			return
		}

		if cityData, ok := cityStats[city]; ok {
			cityData.Max = max(cityData.Max, ftemp)
			cityData.Min = min(cityData.Min, ftemp)
			cityData.Sum += ftemp
			cityData.TotalReading++
		} else {
			cityStats[city] = &City{
				ftemp,
				ftemp,
				ftemp,
				1,
			}
		}
	}
	cities := make([]string, 0, len(cityStats))
	for city := range cityStats {
		cities = append(cities, city)
	}
	sort.Strings(cities)
	fmt.Printf("{")
	for i, cityName := range cities {
		if i > 0 {
			fmt.Printf(", ")
		}
		city := cityStats[cityName]
		tempMean := city.Sum / float64(city.TotalReading)
		fmt.Printf("%s=%.1f/%.1f/%.1f", cityName, city.Min, tempMean, city.Max)
	}
	fmt.Printf("}\n")

}
