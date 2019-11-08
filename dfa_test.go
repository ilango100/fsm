package fsm

import "testing"

func TestDFATable(t *testing.T) {
	table := map[string]map[byte]string{
		">*f": {'0': "f", '1': "n"},
		"n":   {'0': "f", '1': "n"},
	}

	d, err := DFATable(table)
	if err != nil {
		t.Fatalf("Error building DFA: %v\n", err)
	}

	t.Logf("%s\n", d.Info())

	t.Logf("Started: %s\n", d.Reset())
	str := "0100101"

	s := d.Next(str[0])
	for i := 0; i < len(str); i++ {
		s = d.Next(str[i])
		t.Logf("Next: %c %s\n", str[i], s)
	}
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}

	t.Logf("Accepted: %v\n", d.IsAccepted())
	if d.IsAccepted() {
		t.Fatalf("Expected Unaccepted")
	}

	d.Next('0')
	t.Logf("After 0, Accepted: %v\n", d.IsAccepted())
	if !d.IsAccepted() {
		t.Fatalf("Expected Accepted")
	}

}
