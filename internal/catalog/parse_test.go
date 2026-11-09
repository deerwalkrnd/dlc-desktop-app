package catalog

import (
	"strings"
	"testing"
)

// compatCorpus holds filenames in every shape we have seen or can imagine
// content teams producing. Whatever the legacy parser accepted must parse
// identically now.
var compatCorpus = []string{
	// From the docs and the sample library.
	"Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4",
	"Arjun Khatiwada - 1.1 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4",
	"Arjun Khatiwada - 1.10 - Factorization of Polynomials - Polynomials in Maths real 1.10 - Mathematics - OLD - 8 .mp4",
	"Arjun Khatiwada - 1.3 - Factorization of Polynomials - Polynomials in Maths v2 - Mathematics - OLD - 8 .mp4",
	// Spacing variations the legacy split tolerated.
	"Arjun Khatiwada-1.2-Factorization-Polynomials-Mathematics-OLD-8.mp4",
	"Arjun Khatiwada -1.2- Factorization - Polynomials -Mathematics- OLD -8.mp4",
	"  Arjun Khatiwada  -  2  -  Sets  -  Intro  -  Maths  -  NEW  -  10  .mp4",
	"Arjun Khatiwada - 01.05 - Sets - Intro - Maths - NEW - 09 .mp4",
	"Arjun Khatiwada - 3 - Sets - Intro - Maths - NEW - 9.mp4",
	"Arjun Khatiwada - +1.5 - Sets - Intro - Maths - NEW - +9.mp4",
	"Arjun Khatiwada - 1e1 - Sets - Intro - Maths - NEW - 9.mp4",
	// Syllabus values the legacy parser mapped to NEW.
	"Sita Rai - 2.1 - Light - Reflection - Science - NEW - 9 .mp4",
	"Sita Rai - 2.1 - Light - Reflection - Science - new - 9 .mp4",
	"Sita Rai - 2.1 - Light - Reflection - Science - old - 9 .mp4",
	"Sita Rai - 2.1 - Light - Reflection - Science - NEWW - 9 .mp4",
	"Sita Rai - 2.1 - Light - Reflection - Science -  - 9 .mp4",
	// Extra trailing parts (legacy ignored everything after the class).
	"Sita Rai - 2.1 - Light - Reflection - Science - NEW - 9 - final.mp4",
	"Sita Rai - 2.1 - Light - Reflection - Science - NEW - 9 - final - v2.mp4",
	// Empty fields.
	" - 2.1 - Light - Reflection - Science - NEW - 9 .mp4",
	"Sita Rai - 2.1 -  -  - Science - NEW - 9 .mp4",
	// Unicode.
	"राम शर्मा - 4.2 - व्याकरण - सन्धि - नेपाली - NEW - 7 .mp4",
	"Ram Sharma - 4.2 - Vyakaran - Sandhi - नेपाली - OLD - 7 .mp4",
	// Things the legacy parser rejected.
	"Ram-Prasad Sharma - 2.1 - Light - Mid-term Review - Science - NEW - 9 .mp4",
	"Sita Rai - 2.1 - Light - Reflection - Science - NEW .mp4",
	"Sita Rai - two - Light - Reflection - Science - NEW - 9 .mp4",
	"Sita Rai - 2.1 - Light - Reflection - Science - NEW - nine .mp4",
	"Sita Rai - 2.1.3 - Light - Reflection - Science - NEW - 9 .mp4",
	"random video.mp4",
	"a.mp4",
}

func TestLegacyCompatibilityCorpus(t *testing.T) {
	for _, name := range compatCorpus {
		checkCompat(t, name)
	}
}

func checkCompat(t *testing.T, name string) {
	t.Helper()
	legacy := legacyParse(name)
	got, err := Parse(name)
	if legacy == nil {
		return // new parser may accept more; that is the point
	}
	if !legacyNumberSane(strings.Split(name[:len(name)-4], "-")[1]) {
		return // documented deliberate change: NaN/Inf/negative lesson numbers
	}
	if err != nil {
		t.Errorf("%q: legacy accepted it but new parser failed: %v", name, err)
		return
	}
	if !sameAsLegacy(got, legacy) {
		t.Errorf("%q:\n  legacy %+v\n  new    %+v", name, *legacy, got)
	}
}

// FuzzLegacyCompatibility generates arbitrary names; whenever the legacy
// parser accepted one, the new parser must agree.
func FuzzLegacyCompatibility(f *testing.F) {
	for _, name := range compatCorpus {
		f.Add(strings.TrimSuffix(name, ".mp4"))
	}
	f.Fuzz(func(t *testing.T, s string) {
		if strings.ContainsAny(s, `/\`) || s == "" {
			return // a path separator changes what "base name" means
		}
		checkCompat(t, s+".mp4")
	})
}

func TestParseHyphenatedFields(t *testing.T) {
	got, err := Parse("Ram-Prasad Sharma - 2.1 - Light - Mid-term Review - Science - NEW - 9 .mp4")
	if err != nil {
		t.Fatal(err)
	}
	want := Entry{
		Teacher: "Ram-Prasad Sharma", LessonNumber: "2.1", ChapterNumber: 2,
		ChapterTitle: "Light", LessonTitle: "Mid-term Review", Subject: "Science",
		Syllabus: SyllabusNew, Class: 9,
	}
	if !equalEntry(got, want) {
		t.Fatalf("got %+v\nwant %+v", got, want)
	}
}

func TestParseFieldOrder(t *testing.T) {
	// The third part is the chapter, the fourth is the lesson title.
	got, err := Parse("Arjun Khatiwada - 1.10 - Factorization of Polynomials - Polynomials in Maths real 1.10 - Mathematics - OLD - 8 .mp4")
	if err != nil {
		t.Fatal(err)
	}
	want := Entry{
		Teacher: "Arjun Khatiwada", LessonNumber: "1.10", ChapterNumber: 1,
		ChapterTitle: "Factorization of Polynomials", LessonTitle: "Polynomials in Maths real 1.10",
		Subject: "Mathematics", Syllabus: SyllabusOld, Class: 8,
	}
	if !equalEntry(got, want) {
		t.Fatalf("got %+v\nwant %+v", got, want)
	}
}

func TestParseWarnings(t *testing.T) {
	cases := map[string]string{
		"T - 1 - C - L - S - old - 8.mp4":       "",
		"T - 1 - C - L - S - New - 8.mp4":       "",
		"T - 1 - C - L - S - NEWW - 8.mp4":      `syllabus "NEWW" is not OLD or NEW`,
		"T - 1 - C - L - S - NEW - 8 - x.mp4":   "extra part",
		"T - 1 -  - L - S - NEW - 8.mp4":        "chapter title is empty",
		"T - 1 - C - L - S - OLD - 8.mp4":       "",
		"T - 1 - C - L - S - NEW - 8.mp4":       "",
		"Ram-P - 1 - C - L-x - S - NEW - 8.mp4": "",
	}
	for name, want := range cases {
		e, err := Parse(name)
		if err != nil {
			t.Errorf("%q: %v", name, err)
			continue
		}
		joined := strings.Join(e.Warnings, "; ")
		if want == "" && joined != "" {
			t.Errorf("%q: unexpected warnings %q", name, joined)
		}
		if want != "" && !strings.Contains(joined, want) {
			t.Errorf("%q: warnings %q, want substring %q", name, joined, want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	cases := map[string]string{
		"notes.txt":                         "not an .mp4",
		".mp4":                              "not an .mp4",
		"random video.mp4":                  "expected 7 parts",
		"T - x - C - L - S - NEW - 8.mp4":   "lesson number",
		"T - 1 - C - L - S - NEW - y.mp4":   "class",
		"T - NaN - C - L - S - NEW - 8.mp4": "lesson number",
	}
	for name, want := range cases {
		_, err := Parse(name)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: error %v, want substring %q", name, err, want)
		}
	}
}

func TestSyllabusIsCaseInsensitive(t *testing.T) {
	for value, want := range map[string]Syllabus{
		"OLD": SyllabusOld, "old": SyllabusOld, "Old": SyllabusOld,
		"NEW": SyllabusNew, "new": SyllabusNew, "New": SyllabusNew,
		"NEWW": SyllabusNew, "": SyllabusNew,
	} {
		e, err := Parse("T - 1 - C - L - S - " + value + " - 8.mp4")
		if err != nil {
			t.Fatal(err)
		}
		if e.Syllabus != want {
			t.Errorf("%q: syllabus %v, want %v", value, e.Syllabus, want)
		}
	}
}

func TestIsVideo(t *testing.T) {
	for name, want := range map[string]bool{
		"a - b.mp4": true, "A.MP4": true, "dir/x.Mp4": true,
		"x.mp3": false, "xmp4": false, "x.mp4.txt": false, ".mp4": false,
	} {
		if got := IsVideo(name); got != want {
			t.Errorf("IsVideo(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestCompareLessonNumbers(t *testing.T) {
	ordered := []string{"1", "1.1", "1.2", "1.10", "2", "2.01", "10.1", "1.2.3", "abc"}
	for i := range ordered {
		for j := range ordered {
			got := CompareLessonNumbers(ordered[i], ordered[j])
			want := cmpInt(i, j)
			if got != want {
				t.Errorf("Compare(%q, %q) = %d, want %d", ordered[i], ordered[j], got, want)
			}
		}
	}
}

func equalEntry(a, b Entry) bool {
	return a.Teacher == b.Teacher && a.LessonNumber == b.LessonNumber &&
		a.ChapterNumber == b.ChapterNumber && a.ChapterTitle == b.ChapterTitle &&
		a.LessonTitle == b.LessonTitle && a.Subject == b.Subject &&
		a.Syllabus == b.Syllabus && a.Class == b.Class
}
