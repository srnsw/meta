package meta

import "time"

const w3cdtf = "2006-01-02"

type W3CDate struct {
	time.Time
}

func (d W3CDate) MarshalJSON() ([]byte, error) {
	b := make([]byte, 1, len(w3cdtf)+2)
	b[0] = '"'
	b = d.AppendFormat(b, w3cdtf)
	b = append(b, '"')
	return b, nil
}

// ParseDate makes a Date from a W3C style date string
func ParseDate(d string) (W3CDate, error) {
	t, err := time.Parse(w3cdtf, d)
	return W3CDate{t}, err
}

// ParseDateTime makes a time.Time from a W3C style datetime string
func ParseDateTime(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, t)
}
