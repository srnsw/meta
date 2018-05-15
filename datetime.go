package meta

import "time"

const w3cdtf = "2006-01-02"

// W3CDate contains a time.Time but marshals to json in form yyyy-mm-dd
type W3CDate struct {
	time.Time
}

// MarshalJSON makes W3CDate a json Marshaller with yyyy-mm-dd output
func (d W3CDate) MarshalJSON() ([]byte, error) {
	b := make([]byte, 1, len(w3cdtf)+2)
	b[0] = '"'
	b = d.AppendFormat(b, w3cdtf)
	b = append(b, '"')
	return b, nil
}

// NewDate returns a reference to W3CDate from a W3C style date string.
// If the string provided is an invalid date, a nil reference is returned.
func NewDate(d string) *W3CDate {
	var date *W3CDate
	if d != "" {
		if pd, err := ParseDate(d); err == nil {
			date = &pd
		}
	}
	return date
}

// ParseDate makes a W3CDate from a W3C style date string
func ParseDate(d string) (W3CDate, error) {
	t, err := time.Parse(w3cdtf, d)
	return W3CDate{t}, err
}

// NewDate returns a reference to a time.Time from a W3C style datetime string.
// If the string provided is an invalid datetime, a nil reference is returned.
func NewDateTime(t string) *time.Time {
	var ti *time.Time
	if t != "" {
		if pt, err := ParseDateTime(t); err == nil {
			ti = &pt
		}
	}
	return ti
}

// ParseDateTime makes a time.Time from a W3C style datetime string
func ParseDateTime(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, t)
}
