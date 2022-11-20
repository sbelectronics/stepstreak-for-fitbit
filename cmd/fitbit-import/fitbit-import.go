package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	dateFormat = "1/2/2006"
)

type Activity struct {
	Date             time.Time
	CaloriesBurned   uint32
	Steps            uint32
	Distance         float32
	Floors           uint32
	MinuteSedentary  uint32
	MinuteLight      uint32
	MinuteFair       uint32
	MinuteVery       uint32
	ActivityCalories uint32
}

var Location *time.Location
var Activities map[time.Time]Activity

type jsonEntry struct {
	Date  string `json:"dateTime"`
	Value string `json:"value"`
}

type jsonData []jsonEntry

func readJsonFile(fn string) error {
	content, err := os.ReadFile(fn)
	if err != nil {
		return err
	}

	var payload jsonData
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}

	for _, entry := range payload {
		date, err := time.Parse("1/2/06 15:04:05", entry.Date)
		if err != nil {
			return err
		}
		date = date.In(Location)

		// strip off the time of daty
		// (I'm sure there's a more efficient way)
		dateStr := date.Format("2006-01-02")
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}

		act, okay := Activities[date]
		if !okay {
			act = Activity{}
		}

		v, err := strconv.ParseUint(entry.Value, 10, 32)
		if err != nil {
			return err
		}

		act.Steps = act.Steps + uint32(v)

		Activities[date] = act
	}

	return nil
}

type TimeSlice []time.Time

func (p TimeSlice) Len() int {
	return len(p)
}

func (p TimeSlice) Less(i, j int) bool {
	return p[i].Before(p[j])
}

func (p TimeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	var err error

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	log.Print("Convert fitbit dump to csv, by Scott Baker, http://www.smbaker.com/")

	Location, err = time.LoadLocation("America/Los_Angeles") // FIXME
	if err != nil {
		log.Fatal("Error setting timezone %s", err)
	}

	Activities = map[time.Time]Activity{}

	for _, fn := range flag.Args() {
		err := readJsonFile(fn)
		if err != nil {
			log.Fatal("Error: %s", err)
		}
	}

	dates := TimeSlice{}
	for k, _ := range Activities {
		dates = append(dates, k)
	}

	sort.Sort(TimeSlice(dates))

	writer := csv.NewWriter(os.Stdout)
	writer.Write([]string{"Activities"})
	for _, k := range dates {
		act := Activities[k]
		rec := []string{
			k.Format("2006-01-02"),
			strconv.FormatUint(uint64(act.CaloriesBurned), 10),
			strconv.FormatUint(uint64(act.Steps), 10),
			fmt.Sprintf("%0.2f", act.Distance),
			strconv.FormatUint(uint64(act.Floors), 10),
			strconv.FormatUint(uint64(act.MinuteSedentary), 10),
			strconv.FormatUint(uint64(act.MinuteLight), 10),
			strconv.FormatUint(uint64(act.MinuteFair), 10),
			strconv.FormatUint(uint64(act.MinuteVery), 10),
			strconv.FormatUint(uint64(act.ActivityCalories), 10),
		}
		writer.Write(rec)
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
}
