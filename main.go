package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
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

type Worker struct {
	PartitionSize int64
	StartByte     int64
}

func main() {
	start = time.Now()

	output = io.Writer(os.Stdout)

	f, err := os.Open("../measurements.txt")
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer f.Close()

	var numWorkers = runtime.NumCPU() - 8 // 24
	info, _ := f.Stat()
	partitionSize := info.Size() / int64(numWorkers)

	var workers []Worker
	for i := 0; i < numWorkers; i++ {
		workers = append(workers, Worker{
			StartByte: partitionSize * int64(i),
		})
	}

	done := make(chan struct{})

	for i, worker := range workers {
		go func(i int) {
			var tempBuff []byte
			var buff []byte

			if worker.StartByte == 0 {
				tempBuff = make([]byte, 40)
				f.ReadAt(tempBuff, partitionSize-1)
				idx := int64(bytes.IndexByte(tempBuff, '\n')) + 1

				buff = make([]byte, partitionSize+idx)
				f.ReadAt(buff, 0)
			} else {
				tempBuff = make([]byte, 40)
				f.ReadAt(tempBuff, worker.StartByte)
				startIdx := int64(bytes.IndexByte(tempBuff, '\n')) + 1

				tempBuff = make([]byte, 40)
				f.ReadAt(tempBuff, worker.StartByte+partitionSize)
				lastIdx := int64(bytes.IndexByte(tempBuff, '\n')) + 1

				buff = make([]byte, -startIdx+partitionSize+lastIdx)
				f.ReadAt(buff, worker.StartByte+startIdx)
			}

			processRecords(buff)

			done <- struct{}{}
		}(i)
	}

	for range workers {
		<-done
	}

	print()
}

func processRecords(buf []byte) {
	_processStation := func(station string, temp string) {
		if station == "" || temp == "" {
			return
		}
		temp64 := toFloat64(temp)
		s := stationStats[station]
		if s == nil {
			stationStats[station] = &[4]float64{temp64, temp64, temp64, 1}
		} else {
			smin, smax := &s[MIN_IDX], &s[MAX_IDX]
			s[MIN_IDX] = min(*smin, temp64)
			s[MAX_IDX] = max(*smax, temp64)
			s[SUM_IDX] += temp64
			s[COUNT_IDX] += 1
		}
	}

	lastNewline := 0

	for i, b := range buf {
		var station string
		var temp string
		if b == '\n' {
			station, temp = cut(buf[lastNewline:i])
			lastNewline = i + 1
		}
		_processStation(station, temp)
	}

	if lastNewline < len(buf) {
		station, temp := cut(buf[lastNewline:])
		_processStation(station, temp)
	}
}

func cut(buf []byte) (string, string) {
	if len(buf) < 0 {
		return "", ""
	}
	i := bytes.LastIndexByte(buf, ';')
	if i == -1 {
		return "", ""
	}
	return string(buf[:i]), string(buf[i+1:])
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
