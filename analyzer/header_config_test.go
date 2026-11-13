package analyzer

import "testing"

func TestHeaderConfigValidation(t *testing.T) {
	for _, h := range []HeaderConfig{
		{Methods: []string{"Unknown"}},
		{Methods: []string{"Title", "Title"}},
		{Required: [][]string{{}}},
		{Methods: []string{"Tags"}, Required: [][]string{{"Title"}}},
		{Required: [][]string{{"Title", "Title"}}},
	} {
		if _, err := New(Config{Style: StyleConfig{Header: h}}); err == nil {
			t.Errorf("accepted invalid header: %+v", h)
		}
	}
	h := HeaderConfig{Methods: []string{"Title", "Titlef"}, Required: [][]string{{"Title", "Titlef"}}}
	if _, err := New(Config{Style: StyleConfig{Header: h}}); err != nil {
		t.Fatal(err)
	}
	cloned := h.clone()
	cloned.Methods[0] = "Tags"
	cloned.Required[0][0] = "Tags"
	if h.Methods[0] != "Title" || h.Required[0][0] != "Title" {
		t.Fatal("configuration slices shared")
	}
}
