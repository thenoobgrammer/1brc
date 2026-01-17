package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

const (
	MIN_IDX   = 0
	MAX_IDX   = 1
	SUM_IDX   = 2
	COUNT_IDX = 3
)

func main() {
	start := time.Now()

	output := io.Writer(os.Stdout)

	f, err := os.Open("../measurements.txt")
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer f.Close()

	stationStats := make(map[string]*[4]float64)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		station, tempStr := cut(line, ';')
		temp := parseTemp(tempStr)

		s := stationStats[station]
		if s == nil {
			stationStats[station] = &[4]float64{temp, temp, temp, 1}
		} else {
			var smin, smax = &s[MIN_IDX], &s[MAX_IDX]
			s[MIN_IDX] = min(*smin, temp)
			s[MAX_IDX] = min(*smax, temp)
			s[SUM_IDX] += temp
			s[COUNT_IDX] += 1
		}
	}

	stations := make([]string, 0, len(stationStats))
	for station := range stationStats {
		stations = append(stations, station)
	}
	sort.Strings(stations)

	fmt.Fprint(output, "{")
	for i, station := range stations {
		if i > 0 {
			fmt.Fprint(output, ", ")
		}
		s := stationStats[station]
		mean := s[SUM_IDX] / s[COUNT_IDX]
		fmt.Fprintf(output, "%s=%.1f/%.1f/%.1f", station, s[MIN_IDX], mean, s[MAX_IDX])
	}
	fmt.Fprint(output, "}\n")
	fmt.Fprintf(output, "Elapsed: %v ------------------\n", time.Since(start))
}

func parseTemp(value string) float64 {
	index := 0
	sign := 1.0

	if value[index] == '-' {
		sign = -1.0
		index++
	}

	res := float64(value[index] - '0')
	index++
	if value[index] != '.' {
		res = res*10 + float64(value[index]-'0')
		index++
	}
	index++

	res = res + ((float64(value[index]) - '0') / 10.0)

	return sign * res
}

func cut(s string, sep rune) (string, string) {
	var sepIdx int
	for i, v := range s {
		if v == sep {
			sepIdx = i
			break
		}
	}

	return s[:sepIdx-1], s[sepIdx+1:]
}
