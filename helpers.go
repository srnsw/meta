package meta

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

// ToID is a helper function that turns an integer identififier into a string @id
// It is useful for places where we use integer IDs e.g. series numbers but want an IRI for the @id
func ToID(i int, pat string) string {
	if i <= 0 {
		return ""
	}
	return pat + strconv.Itoa(i)
}

// ToPUID is a helper function that turns a short PUID (e.g. fmt/1) into a fully qualified PUID
// (e.g. http://www.nationalarchives.gov.uk/pronom/fmt/1)
func ToPUID(puid string) string {
	if !strings.HasPrefix(puid, "fmt/") && !strings.HasPrefix(puid, "x-fmt/") {
		return puid
	}
	return "http://www.nationalarchives.gov.uk/pronom/" + puid
}

// ToRef is a helper function that turns an integer identififier into a string ref:I reference.
// It is used for placedholder IDs like log:1 that get swapped out by the Migrate tool.
func ToRef(i int, ref string) string {
	if i < 0 {
		return ""
	}
	return ref + ":" + strconv.Itoa(i)
}

// ToSz is a helper function that turns a string int into an int64
func ToSz(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// ReadAll is a helper function that opens a file at path and reads as a CSV.
// Provide an optional lazy quote bool if you'd like lazy quotes.
// Returns a slice of string slices and an error.
func ReadAll(path string, lq ...bool) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	if len(lq) > 0 {
		reader.LazyQuotes = lq[0]
	}
	return reader.ReadAll()
}

// VarStr represents a single string or a slice of strings
type VarStr interface{}

func copyVarStr(v VarStr) VarStr {
	if v == nil {
		return nil
	}
	if str, ok := v.(string); ok {
		return str
	}
	strs := v.([]string)
	ret := make([]string, len(strs))
	copy(ret, strs)
	return ret
}

// MakeAgent returns an Agent with the given name, @id and @type.
// If @id and @type aren't given, the Agent is a simple string e.g. "Richard Lehane"
func MakeAgent(name, id, typ string) Agent {
	if id == "" && typ == "" {
		return name
	}
	return Obj{
		ID:   id,
		Typ:  typ,
		Name: name,
	}
}

// MakeSDOPerson creates an Agent that is of @type schema.org/Person. Does not set an @id.
func MakeSDOPerson(name string) Agent {
	return MakeAgent(name, "", "http://schema.org/Person")
}

// MakePerson creates an Agent that is of @type schema.org/Organization.
func MakeOrganization(name string) Agent {
	return MakeAgent(name, "", "http://schema.org/Organization")
}

// MakePerson creates an Agent that is of @type terms/Person. Sets the @id to the supplied value.
func MakePerson(name string, id int) Agent {
	return MakeAgent(name, ToID(id, "http://records.nsw.gov.au/persons/"), "http://records.nsw.gov.au/terms/Person")
}

// MakePerson creates an Agent that is of @type terms/Agency. Sets the @id to the supplied value.
func MakeAgency(name string, id int) Agent {
	return MakeAgent(name, ToID(id, "http://records.nsw.gov.au/agencies/"), "http://records.nsw.gov.au/terms/Agency")
}

// ToSeries turns a series number into an IRI @id
func ToSeries(i int) string {
	return ToID(i, "http://records.nsw.gov.au/series/")
}

// ToConsignment turns a consignment number into an IRI @id
func ToConsignment(i int) string {
	return ToID(i, "http://records.nsw.gov.au/consignments/")
}

// MakeBusiness returns a Thing of @type schema.org/Organization and sets the supplied fields.
// A variable number of proprietors can be supplied and these will be set as a slice of Organizations.
func MakeBusiness(legalName, commencedTrading, ceasedTrading, renewalDueDate, registrationNumber, abn string, proprietors ...string) Thing {
	var props Agent
	if len(proprietors) > 0 {
		if len(proprietors) > 1 {
			props := make([]Agent, len(proprietors))
			for i, v := range proprietors {
				props[i] = MakeOrganization(v)
			}
		} else {
			props = MakeOrganization(proprietors[0])
		}
	}
	return Business{
		Typ:                "http://schema.org/Organization",
		LegalName:          legalName,
		CommencedTrading:   NewDate(commencedTrading),
		CeasedTrading:      NewDate(ceasedTrading),
		RenewalDueDate:     NewDate(renewalDueDate),
		RegistrationNumber: registrationNumber,
		ABN:                abn,
		Proprietor:         props,
	}

}

// ReferenceObject makes a temporary reference to another object in the consignment.
// This reference is swapped for a UUID by the migrate tool.
func ReferenceObject(i int) string {
	return ToRef(i, "obj")
}

// ReferenceMigration makes a temporary reference to another migration (that doesn't yet have a UUID).
// This reference is swapped for a UUID by the migrate tool.
func ReferenceMigration(i int) string {
	return ToRef(i, "mig")
}
