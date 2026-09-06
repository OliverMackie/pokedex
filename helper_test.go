package main

import (
	"encoding/json"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  Hello world  ",
			expected: []string{"hello", "world"},
		},

		{
			input:    "  hello\n world  ",
			expected: []string{"hello", "world"},
		},

		{
			input:    "  hi\nthere ",
			expected: []string{"hi", "there"},
		},

		{
			input:    "this is a longer test case with\n multiple words and spaces  ",
			expected: []string{"this", "is", "a", "longer", "test", "case", "with", "multiple", "words", "and", "spaces"},
		},
	}

	for _, c := range cases {
		actual := CleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("CleanInput(%q) returned %v, expected %v", c.input, actual, c.expected)
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("CleanInput(%q) returned %v, expected %v", c.input, actual, c.expected)
			}
		}
	}
}

func TestLocationAreaListUnmarshal(t *testing.T) {
	payload := `{
		"count": 1539,
		"next": "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
		"previous": null,
		"results": [{
			"name": "canalave-city-area",
			"url": "https://pokeapi.co/api/v2/location-area/1/"
		}]
	}`

	var list LocationAreaList
	if err := json.Unmarshal([]byte(payload), &list); err != nil {
		t.Fatalf("json.Unmarshal returned an error: %v", err)
	}

	if list.Count != 1539 {
		t.Fatalf("expected count 1539, got %d", list.Count)
	}
	if list.Next == "" {
		t.Fatal("expected next URL to be populated")
	}
	if list.Previous != nil {
		t.Fatalf("expected previous to be nil, got %v", *list.Previous)
	}
	if len(list.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(list.Results))
	}
	if list.Results[0].Name != "canalave-city-area" {
		t.Fatalf("expected first result name to be canalave-city-area, got %q", list.Results[0].Name)
	}
}
