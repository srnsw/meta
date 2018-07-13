// Copyright 2018 State of New South Wales through the State Archives and Records Authority of NSW
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package meta

import (
	"errors"
	"time"
)

const (
	w3cymd   = "2006-01-02"
	w3cym    = "2006-01"
	w3cy     = "2006"
	SlashDMY = "2/01/2006"
	FBDTF    = "2006-01-02T15:04:05+0700"
)

// W3CDate contains a time.Time but marshals to json in form yyyy-mm-dd
type W3CDate struct {
	precision int
	time.Time
}

// MarshalJSON makes W3CDate a json Marshaller with yyyy-mm-dd output
func (d W3CDate) MarshalJSON() ([]byte, error) {
	fstr := w3cymd
	switch d.precision {
	case 1:
		fstr = w3cym
	case 2:
		fstr = w3cy
	}
	b := make([]byte, 1, len(fstr)+2)
	b[0] = '"'
	b = d.AppendFormat(b, fstr)
	b = append(b, '"')
	return b, nil
}

// NewDate returns a reference to W3CDate from a W3C style date string.
// If the string provided is an invalid date, a nil reference is returned.
func NewDate(d string) *W3CDate {
	var date *W3CDate
	if pd, err := ParseDate(d); err == nil {
		date = &pd
	}
	return date
}

// NewDate returns a reference to W3CDate from a time layout string and a
// date in that layout string.
// If the string provided is an invalid date, a nil reference is returned.
func NewDateLayout(layout, d string) *W3CDate {
	var date *W3CDate
	if t, err := time.Parse(layout, d); err == nil {
		date = WrapDate(t)
	}
	return date
}

// ParseDateLayout returns a reference to W3CDate from a time layout string and a
// date in that layout string.
// If the string provided is an invalid date, an error is returned
func ParseDateLayout(layout, d string) (*W3CDate, error) {
	t, err := time.Parse(layout, d)
	return WrapDate(t), err
}

// WrapDate allows you to create a *W3CDate (with YMD precision) when you already have a *time.Time
func WrapDate(t time.Time) *W3CDate {
	return &W3CDate{0, t}
}

// ParseDate makes a W3CDate from a W3C style date string
func ParseDate(d string) (W3CDate, error) {
	var (
		t   time.Time
		p   int
		err error
	)
	switch len(d) {
	case len(w3cymd):
		t, err = time.Parse(w3cymd, d)
		p = 0
	case len(w3cym):
		t, err = time.Parse(w3cym, d)
		p = 1
	case len(w3cy):
		t, err = time.Parse(w3cy, d)
		p = 2
	default:
		err = errors.New("Meta: datetime error, invalid length for provided date " + d)
	}
	return W3CDate{p, t}, err
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
