package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/teambition/rrule-go"
)

type SimplifiedCalendar struct {
	Events      []SimplifiedCalendarEvent `json:"events"`
	DateCreated time.Time                 `json:"dateCreated"`
}

type SimplifiedCalendarEvent struct {
	Title string    `json:"title"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func parseICalDate(prop *ics.IANAProperty, defaultLoc *time.Location) (time.Time, error) {
	if prop == nil || prop.Value == "" {
		return time.Time{}, fmt.Errorf("missing date value")
	}

	if tzid := getTZID(prop); tzid != "" {
		return rrule.StrToDtStart("TZID="+tzid+":"+prop.Value, defaultLoc)
	}

	return rrule.StrToDtStart(prop.Value, defaultLoc)
}

func getTZID(prop *ics.IANAProperty) string {
	if prop == nil || prop.ICalParameters == nil {
		return ""
	}

	if values, ok := prop.ICalParameters["TZID"]; ok && len(values) > 0 {
		return values[0]
	}

	return ""
}

func main() {
	// set calendars
	var calendarURLs = []string{
		os.Getenv("CALENDAR_1"),
		os.Getenv("CALENDAR_2"),
		os.Getenv("CALENDAR_3"),
	}

	var allEvents []SimplifiedCalendarEvent
	// iterate through each
	for i, url := range calendarURLs {
		cal, err := ics.ParseCalendarFromUrl(url)
		if err != nil {
			log.Fatal(err)
		}

		windowStart := time.Now()
		windowEnd := time.Now().Add(7 * 24 * time.Hour)

		for _, event := range cal.Events() {
			// check each event for proximity to current date
			// if event is within 1 week, save to new format
			componentDate := event.GetProperty(ics.ComponentPropertyDtStart)
			parsedDate, err := parseICalDate(componentDate, time.Local)
			if err != nil {
				continue
			}

			duration := time.Duration(0)
			endProp := event.GetProperty(ics.ComponentPropertyDtEnd)
			if endProp != nil {
				parsedEndDate, err := parseICalDate(endProp, time.Local)
				if err == nil {
					duration = parsedEndDate.Sub(parsedDate)
				}
			}

			summaryProp := event.GetProperty(ics.ComponentPropertySummary)
			title := ""
			if summaryProp != nil {
				title = summaryProp.Value
			}

			rruleProp := event.GetProperty(ics.ComponentProperty("RRULE"))
			if rruleProp != nil {
				opt, err := rrule.StrToROptionInLocation(rruleProp.Value, time.Local)
				if err != nil {
					continue
				}
				opt.Dtstart = parsedDate
				r, err := rrule.NewRRule(*opt)
				if err != nil {
					continue
				}

				for _, occurrence := range r.Between(windowStart, windowEnd, true) {
					parsedEvent := SimplifiedCalendarEvent{
						Title: title,
						Start: occurrence,
						End:   occurrence.Add(duration),
					}
					allEvents = append(allEvents, parsedEvent)
				}
				continue
			}

			if parsedDate.Before(windowEnd) && parsedDate.After(windowStart) {
				parsedEvent := SimplifiedCalendarEvent{
					Title: title,
					Start: parsedDate,
					End:   parsedDate.Add(duration),
				}
				allEvents = append(allEvents, parsedEvent)
			}
		}
		fmt.Printf("Calendar %d has %d events\n", i, len(cal.Events()))

	}

	jsonData, err := json.Marshal(SimplifiedCalendar{Events: allEvents, DateCreated: time.Now()})
	if err != nil {
		log.Println("Error marshalling calendar:", err)
		return
	}

	// Hardcoded 128-bit AES key
	key, err := hex.DecodeString(os.Getenv("CAL_KEY"))
	if err != nil {
		log.Fatal("Error decoding key:", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("Error creating cipher:", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal("Error generating IV:", err)
	}

	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(jsonData))
	stream.XORKeyStream(ciphertext, jsonData)

	os.MkdirAll("docs", 0755)
	file, err := os.Create("docs/cal.aes")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	if _, err := file.WriteString(hex.EncodeToString(iv) + "\n"); err != nil {
		log.Fatal("Error writing IV:", err)
	}
	if _, err := file.Write(ciphertext); err != nil {
		log.Fatal("Error writing ciphertext:", err)
	}

	fmt.Printf("Successfully encrypted and saved %d events to docs/cal.aes\n", len(allEvents))
}
