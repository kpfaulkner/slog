package processor

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RoundFactor int

const (
	RoundSecond RoundFactor = iota
	RoundMinute
	RoundHour


	// misc regexs
	TIMESTAMPREGEX string = "(19|20\\d\\d)-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])T(00|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9]):(0[0-9]|[0-5][0-9])"
)

type GraphPoint struct {
	timestamp  time.Time
	errorCount int
}

type LogProcessor struct {

	// files to process.
	filenames []string

	// terms to search all log lines for. (for now just substrings, but will change to regexs
	terms []string

	// map of terms...  links to another map where time (rounded to nearest second/hour/whatever) and count of the
	// term hit.
	termDict map[string]map[time.Time]int
}

var timeStampRegexComp *regexp.Regexp

func NewLogProcessor(terms []string) *LogProcessor {
	lp := LogProcessor{}
	lp.termDict = make(map[string]map[time.Time]int)
	lp.terms = terms

	timeStampRegexComp ,_ = regexp.Compile(TIMESTAMPREGEX)
	return &lp
}

func convertTimeStamp(stringTime string) (*time.Time, error) {
	layout := "2006-01-02.15-04-05"
	t, err := time.Parse(layout, stringTime)

	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (lp *LogProcessor) findMatchedTerms(line string) []string {

	matchedTerms := []string{}
	lowerLine := strings.ToLower(line)
	for _, term := range lp.terms {
		if strings.Contains(lowerLine, term) {
			matchedTerms = append(matchedTerms, term)
		}
	}

	return matchedTerms
}

func (lp *LogProcessor) determineTimeStampForLine(line string) (*time.Time, error) {

	res := timeStampRegexComp.FindStringSubmatch(line)
	if res != nil {
		fmt.Printf("res is %v\n", res)

		year,err := strconv.Atoi(res[1])
		if err != nil {
			return nil, err
		}
		month,err := strconv.Atoi(res[2])
		if err != nil {
			return nil, err
		}
		day, err := strconv.Atoi(res[3])
		if err != nil {
			return nil, err
		}
		hour, err := strconv.Atoi(res[4])
		if err != nil {
			return nil, err
		}
		minute, err := strconv.Atoi(res[5])
		if err != nil {
			return nil, err
		}
		second,err := strconv.Atoi(res[6])
		if err != nil {
			return nil, err
		}

		t := time.Date(year, time.Month(month),day,hour,minute,second, 0, time.UTC)
		return &t, nil
	}

	return nil, errors.New("no time stamp")
}

func (lp *LogProcessor) roundTimeStamp(timeStamp time.Time, rounding RoundFactor) time.Time {

	switch rounding {
	case RoundSecond:
		return timeStamp.Round(time.Second)

	case RoundMinute:
		return timeStamp.Round(time.Minute)

	case RoundHour:
		return timeStamp.Round(time.Hour)
	}

	return timeStamp
}

func (lp *LogProcessor) ReadData(filePath string, rounding RoundFactor) (map[string]map[time.Time]int,error) {

	termDict := make(map[string]map[time.Time]int)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()

		matchedTerms := lp.findMatchedTerms(l)
		if len(matchedTerms) == 0 {
			continue
		}

		timeStamp, err := lp.determineTimeStampForLine(l)
		if err != nil {
			// unable to find timestamp...  what do we do? ditch the line?
			// throw panic for moment... want to see if this appears much
			log.Panicf("Unable to detect timestamp for line... blowing up")
		}

		roundedTimeStamp := lp.roundTimeStamp(*timeStamp, rounding)

		for _,mt := range matchedTerms {
			timeDict, ok := termDict[mt]
			if !ok  {
				timeDict = make(map[time.Time]int)
				lp.termDict[mt] = timeDict
			}

			timeDict[roundedTimeStamp]++
			termDict[mt] = timeDict
		}
	}

	return termDict, nil
}
