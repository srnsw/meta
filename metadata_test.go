package meta

import "testing"

func TestMetadata(t *testing.T) {
	m := NewMetadata("Business Name Registration - Duntryleague Country Club", "1902-01-01", 21404)
	m.Typ = append(m.Typ, "http://schema.org/Movie")
	m.Creators = []Agent{MakeAgency("Office of Fair Trading", 0), MakePerson("Michael Bruce Baird", 288)}
	m.Source = "https://twitter.com/"
	m.Consignment = ToID(189087, "http://www.records.nsw.gov.au/consignment/")
	m.DisposalRules = []DisposalRule{{Authority: "DA48", Class: "1.1.1.2"}}
	m.Duration = "13:47:30"
	m.Language = "en"
	m.Subtitles = "zh"
	m.Director = "Mel Gibson"
	m.Actor = "Mel Gibson"
	m.ProductionCompany = "Icon Films"
	m.About = MakeBusiness("Duntryleague Country Club", "1000-01-01", "1999-03-29", "1998-11-30", "A0369711", "Unknown", "The Orange Golf Club Ltd")
	ctx, err := populate(metadataContext, m)
	if err != nil {
		t.Fatal(err)
	}
	m.Context = ctx
	if err = compare(m, "metadata.json"); err != nil {
		t.Fatal(err)
	}
}
