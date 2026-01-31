package main

import (
	"bytes"
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

var (
	start        time.Time
	output       = io.Writer(os.Stdout)
	stationStats = make(map[string]*[4]float64)
)

func main() {
	start = time.Now()

	output = io.Writer(os.Stdout)

	f, err := os.Open("./utils/measurements_10.txt")
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer f.Close()

	var reader io.Reader = f

	buf := make([]byte, 4096)

	left_over := 0
	// [a b c d e \n  f   g]
	// 		       i i+1
	// [f g h i j k l \n m n]
	for {
		n, readErr := reader.Read(buf[:left_over])
		if n == 0 && readErr != nil {
			break
		}
		i := bytes.LastIndexByte(buf, '\n')
		fmt.Println(n, left_over, i)
		if i != -1 {
			// carry_over_len := len(buf) - (i + 1)

			// process that array to the cut()
			process(buf[:i])

			// put carry-over in the begining
			copy(buf, buf[i+1:])

			// find the carry-over start index
			left_over = len(buf) - i - 1
		} else {
			left_over = len(buf)
		}
	}

	// print()
}

func process(buf []byte) {
	_processStation := func(station string, temp string) {
		temp64 := toFloat64(temp)
		s := stationStats[station]
		if s == nil {
			stationStats[station] = &[4]float64{temp64, temp64, temp64, 1}
		} else {
			var smin, smax = &s[MIN_IDX], &s[MAX_IDX]
			s[MIN_IDX] = min(*smin, temp64)
			s[MAX_IDX] = min(*smax, temp64)
			s[SUM_IDX] += temp64
			s[COUNT_IDX] += 1
		}
	}

	last_newline := 0

	for i, b := range buf {
		var station string
		var temp string
		if b == '\n' {
			station, temp = cut(buf[last_newline:i])
			last_newline = i
		}
		// else if i == len(buf)-1 {
		// 	station, temp = cut(buf[last_newline+1 : i])
		// }
		_processStation(station, temp)
	}

}

func cut(buf []byte) (string, string) {
	i := bytes.LastIndexByte(buf, ';')
	return string(buf[:i-1]), string(buf[i+1:])
}

// String vers
func toFloat64(value string) float64 {
	if value == "" {
		return 0.0
	}

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

func print() {
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
