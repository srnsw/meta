package meta

import (
	"path/filepath"
	"testing"
)

func TestMetadata(t *testing.T) {
	m := NewMetadata(0, "Business Name Registration - Duntryleague Country Club")
	m.AddType("http://schema.org/Movie")
	m.Created = NewDate("1902-01-01")
	m.Creator = []Agent{MakeAgency("Office of Fair Trading", 0), MakePerson("Michael Bruce Baird", 288)}
	m.Source = "https://twitter.com/"
	m.IsPartOf = "Exhibits 50. The trials of Mel Gibson"
	m.Series = ToSeries(21404)
	m.Consignment = ToConsignment(189087)
	m.DisposalRule = DisposalRule{Authority: "DA48", Class: "1.1.1.2"}
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

func TestMetadataConst(t *testing.T) {
	m := NewMetadata(0, "State Records Act 1998 No 17")
	m.AddType("http://schema.org/Legislation")
	m.Description = "An Act to make provision for the creation, management and protection of the records of public offices of the State and to provide for public access to those records, to establish the State Archives and Records Authority; and for other purposes."
	m.Created = NewDate("1998")
	m.Creator = []Agent{MakeAgency("Department of Premier and Cabinet", 10)}
	m.Series = ToSeries(21404)
	m.Consignment = ToConsignment(189087)
	m.DisposalRule = DisposalRule{Authority: "DA48", Class: "1.1.1.2"}
	m.Language = "en"
	ctx, err := populate(metadataContext, m)
	if err != nil {
		t.Fatal(err)
	}
	m.Context = ctx
	if err = compare(m, filepath.Join("project-0", "0", "metadata.json")); err != nil {
		t.Fatal(err)
	}
}
