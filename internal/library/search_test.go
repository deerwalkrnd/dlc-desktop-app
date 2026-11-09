package library

import (
	"strings"
	"testing"
)

func searchFixture(t *testing.T) *Index {
	t.Helper()
	root := t.TempDir()
	writeFiles(t, root,
		"Arjun Khatiwada - 1.1 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8.mp4",
		"Arjun Khatiwada - 1.2 - Factorization of Polynomials - Common Factors - Mathematics - OLD - 8.mp4",
		"Arjun Khatiwada - 1.10 - Factorization of Polynomials - Revision - Mathematics - OLD - 8.mp4",
		"Arjun Khatiwada - 2.1 - Sets - Introduction to Sets - Mathematics - NEW - 8.mp4",
		"Sita Rai - 1.1 - Matter - States of Matter - Science - NEW - 9.mp4",
		"Ram-Prasad Sharma - 3.1 - Light - Mid-term Review - Science - NEW - 9.mp4",
		"Gita Thapa - 1.1 - व्याकरण - सन्धि - नेपाली - NEW - 9.mp4",
		"Hari Bhandari - 4.1 - Trigonometry - Heights of Polynomials - Mathematics - NEW - 10.mp4",
	)
	x, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	return x
}

func titles(hits []Hit) string {
	var s []string
	for _, h := range hits {
		s = append(s, h.Lesson.Title)
	}
	return strings.Join(s, " | ")
}

func TestSearch(t *testing.T) {
	x := searchFixture(t)
	cases := []struct {
		query, want string
	}{
		// Title matches outrank chapter matches.
		{"polynomials", "Polynomials in Maths | Heights of Polynomials | Common Factors | Revision"},
		// Every word must match, in any field.
		{"arjun sets", "Introduction to Sets"},
		{"science matter", "States of Matter"},
		// Case-insensitive; punctuation in titles is ignored.
		{"MID TERM", "Mid-term Review"},
		{"mid-term", "Mid-term Review"},
		// Teacher with a hyphen.
		{"ram-prasad", "Mid-term Review"},
		// Lesson numbers: exact, or a chapter prefix.
		{"1.10", "Revision"},
		// Class numbers, with or without the word "class".
		{"class 10", "Heights of Polynomials"},
		{"10 trig", "Heights of Polynomials"},
		// Non-Latin script.
		{"सन्धि", "सन्धि"},
		{"नेपाली", "सन्धि"},
		// No match, empty and filler-only queries.
		{"chemistry", ""},
		{"   ", ""},
		{"class", ""},
	}
	for _, c := range cases {
		hits, total := x.Search(c.query, 0)
		if got := titles(hits); got != c.want {
			t.Errorf("Search(%q) = %q, want %q", c.query, got, c.want)
		}
		if total != len(hits) {
			t.Errorf("Search(%q): total %d != %d hits", c.query, total, len(hits))
		}
	}
}

func TestSearchChapterPrefix(t *testing.T) {
	x := searchFixture(t)
	// "1" matches lesson numbers 1.x and class/chapter-free text containing 1.
	hits, _ := x.Search("1 factorization", 0)
	if got := titles(hits); got != "Polynomials in Maths | Common Factors | Revision" {
		t.Fatalf("got %q", got)
	}
}

func TestSearchLimitAndContext(t *testing.T) {
	x := searchFixture(t)
	hits, total := x.Search("mathematics", 2)
	if total != 5 || len(hits) != 2 {
		t.Fatalf("total=%d hits=%d", total, len(hits))
	}
	h := hits[0]
	// Ties keep display order: class 8, new syllabus before old.
	if h.Subject.Slug != "mathematics-new" || h.Subject.Class != 8 || h.Chapter.Title != "Sets" || h.Lesson.Title != "Introduction to Sets" {
		t.Fatalf("context = %+v %+v", h.Subject, h.Chapter)
	}
}

func TestFieldScore(t *testing.T) {
	if fieldScore("polynomials in maths", "poly", 3) != 6 {
		t.Error("word-start match should score double")
	}
	if fieldScore("common polynomials", "nomial", 3) != 3 {
		t.Error("mid-word match should score the base weight")
	}
	if fieldScore("abcab", "ab", 1) != 2 {
		t.Error("first occurrence at index 0 is a word start")
	}
	if fieldScore("xab ab", "ab", 1) != 2 {
		t.Error("a later word-start occurrence should be found")
	}
	if fieldScore("maths", "sci", 3) != 0 {
		t.Error("no match should score 0")
	}
}
