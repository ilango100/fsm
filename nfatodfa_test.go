package fsm

import "testing"

func TestNFAtoDFA(t *testing.T) {

	table := map[string]map[rune][]string{
		">*f": {'0': []string{"f", "n"}, '1': []string{"n"}},
		"n":   {'0': []string{"f"}, '1': []string{"n"}},
	}

	n, err := NFAfromTable(table)
	if err != nil {
		t.Fatal(err)
	}

	d, err := NFAtoDFA(*n)
	if err != nil {
		t.Fatalf("Error building DFA: %v\n", err)
	}

	t.Logf("%s\n", d)

	t.Logf("Started: %s\n", d.Reset())
	str := "0100101"

	s := d.Next(rune(str[0]))
	for i := 0; i < len(str); i++ {
		s = d.Next(rune(str[i]))
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

	d.ToCSV("nfatodfa.csv")

}
