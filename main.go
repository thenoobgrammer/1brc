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

	var reader io.Reader = os.Stdin

	buf := make([]byte, 4096)

	left_over := 0

	for {
		n, readErr := reader.Read(buf[:left_over])
		if n == 0 && readErr != nil {
			break
		}

		i := bytes.LastIndexByte(buf, '\n')

		// find the carry-over start index
		carry_over_len := len(buf) - (i + 1)

		// process that array to the cut()
		process(buf[:i])

		// put carry-over in the begining
		copy(buf, buf[i+1:])

		left_over = carry_over_len
	}

	print()
}

func process(buf []byte) {
	for _, b := range buf {
		if b == '\n' {
			station, temp := cut(buf)
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
	}
}

func cut(buf []byte) (string, string) {
	i := bytes.LastIndexByte(buf, ';')
	return string(buf[:i-1]), string(buf[i+1:])
}

// String vers
func toFloat64(value string) float64 {
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
