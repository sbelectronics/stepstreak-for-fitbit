package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "1/2/2006"
)

var (
	stepGoal = flag.Uint("step_goal", 10000, "Number of steps in goal")
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

var Activities map[time.Time]Activity

func removeComma(src string) string {
	dest := ""
	for _, c := range src {
		if c != ',' {
			dest = dest + string(c)
		}
	}
	return dest
}

func readUint32(s string) (uint32, error) {
	i, err := strconv.ParseUint(removeComma(s), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

func readFloat32(s string) (float32, error) {
	i, err := strconv.ParseFloat(removeComma(s), 32)
	if err != nil {
		return 0, err
	}
	return float32(i), nil
}

func readCsvFile(fn string) {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatalf("Unable to open %s: %s", fn, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	recs, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to read %s: %s", fn, err)
	}
	processing := false
	for _, rec := range recs {
		if (len(rec) == 1) && processing {
			return
		}
		if rec[0] == "Activities" {
			processing = true
			continue
		}
		if rec[0] == "Date" {
			continue
		}

		if processing {

			activity := Activity{}
			activity.Date, err = time.Parse("2006-1-2", rec[0])
			if err != nil {
				log.Printf("Failed to parse date %s: %s", rec[0], err)
				continue
			}
			activity.CaloriesBurned, err = readUint32(rec[1])
			activity.Steps, err = readUint32(rec[2])
			activity.Distance, err = readFloat32(rec[3])
			activity.Floors, err = readUint32(rec[4])
			activity.MinuteSedentary, err = readUint32(rec[5])
			activity.MinuteLight, err = readUint32(rec[6])
			activity.MinuteFair, err = readUint32(rec[7])
			activity.MinuteVery, err = readUint32(rec[8])
			activity.ActivityCalories, err = readUint32(rec[9])

			Activities[activity.Date] = activity
		}
	}
}

/*
type jsonEntry struct {
	date string `json:"name"`
	value uint32 `json:"value"`
}

type jsonData []jsonEntry

func readJsonFile(fn string) error {
	content, err := io.ReadFile(fn)
	if err != nil {
		return err
	}

	var payload jsonData
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}

	for _, entry := range jsonData {
	    date, err := time.Parse("1/2/06", jsonData.)
	}
}
*/

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

// because windows is lame.
func dirExpand(fnparam string) []string {
	result := []string{}
	if !strings.HasSuffix(fnparam, "\\") && !strings.HasSuffix(fnparam, "/") {
		result = append(result, fnparam)
		return result
	}

	entries, err := os.ReadDir(fnparam + ".")
	if err != nil {
		log.Fatal("Failed in directory read of %s: %v", fnparam, err)
	}
	for _, entry := range entries {
		if entry.Name()[0] == '.' {
			continue
		}
		if !strings.Contains(strings.ToLower(entry.Name()), ".csv") {
			continue
		}
		result = append(result, fnparam+entry.Name())
	}

	return result
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	log.Print("Stepstreak calculator for fitbit, by Scott Baker, http://www.smbaker.com/")

	Activities = map[time.Time]Activity{}

	// read the imported archive, MyFitbitData, first.
	for _, fnparam := range flag.Args() {
		for _, fn := range dirExpand(fnparam) {
			if strings.Contains(fn, "MyFitbitData") {
				readCsvFile(fn)
			}
		}
	}

	for _, fnparam := range flag.Args() {
		for _, fn := range dirExpand(fnparam) {
			if !strings.Contains(fn, "MyFitbitData") {
				readCsvFile(fn)
			}
		}
	}

	dates := TimeSlice{}
	for k, _ := range Activities {
		dates = append(dates, k)
	}

	sort.Sort(sort.Reverse(TimeSlice(dates)))

	days := 0
	cur := dates[0]
	cur = cur.Add(-time.Hour * 24) // skip today
	for {
		act, okay := Activities[cur]
		if !okay {
			log.Printf("No data on %s", cur.Format(dateFormat))
			break
		}

		if act.Steps < uint32(*stepGoal) {
			log.Printf("Step count of %d is below goal on %s", act.Steps, cur.Format(dateFormat))
			break
		}

		cur = cur.Add(-time.Hour * 24)
		days += 1
	}

	fmt.Printf("You're on a %d day step streak!\n", days)
}
