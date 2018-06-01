package meta

import "testing"

func TestLog(t *testing.T) {
	l := NewLog(0, MigrationEvent)
	l.Start, l.End = NewDateTime("2015-04-20T17:41:48+10:00"), NewDateTime("2015-04-20T17:42:00+10:00")
	l.Detail = "Manually created CSV using MS Excel 2013"
	l.Agent = MakeSDOPerson("Richard Lehane")
	ctx, err := populate(logContext, l)
	if err != nil {
		t.Fatal(err)
	}
	l.Context = ctx
	if err = compare(l, "log.json"); err != nil {
		t.Fatal(err)
	}
}
