package processor

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type RoundFactor int

const (
	RoundSecond RoundFactor = iota
	RoundMinute
	RoundHour
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

func NewLogProcessor() *LogProcessor {
	lp := LogProcessor{}
	lp.termDict = make(map[string]map[time.Time]int)

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

	now := time.Now()

	// yeah yeah, fake it.
	return &now, nil
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

func (lp *LogProcessor) ReadData(filePath string, rounding RoundFactor) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
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
			timeDict, ok := lp.termDict[mt]
			if !ok  {
				timeDict := make(map[time.Time]int)
				lp.termDict[mt] = timeDict
			}

			timeDict[roundedTimeStamp]++
			lp.termDict[mt] = timeDict
		}
	}

	return nil
}
