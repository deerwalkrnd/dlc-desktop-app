package library

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// writeFiles creates empty files (fake videos) under root.
func writeFiles(t *testing.T, root string, names ...string) {
	t.Helper()
	for _, n := range names {
		p := filepath.Join(root, filepath.FromSlash(n))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("fake"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestBuildFakeLibrary(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root,
		// Sample library from the repository, deliberately out of order.
		"Arjun Khatiwada - 1.10 - Factorization of Polynomials - Polynomials in Maths real 1.10 - Mathematics - OLD - 8 .mp4",
		"Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4",
		"Arjun Khatiwada - 1.1 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4",
		"Arjun Khatiwada - 1.3 - Factorization of Polynomials - Polynomials in Maths v2 - Mathematics - OLD - 8 .mp4",
		// Nested folders and upper-case extension (both were broken before).
		"Class 8/Maths/Arjun Khatiwada - 2.1 - Sets - Intro to Sets - Mathematics - NEW - 8.MP4",
		// Hyphens inside fields (rejected before).
		"Class 9/Ram-Prasad Sharma - 3.1 - Light - Mid-term Review - Science - NEW - 9 .mp4",
		// Same subject name, other syllabus, other class.
		"Sita Rai - 1.1 - Matter - States of Matter - Science - OLD - 9 .mp4",
		// Junk that must be ignored or reported.
		"._Arjun Khatiwada - 1.1 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4",
		"notes.txt",
		"random clip.mp4",
		// Lower-case syllabus.
		"Sita Rai - 1.2 - Matter - Solids - Science - old - 9 .mp4",
		// Duplicate of 1.2 in another folder.
		"backup/Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4",
	)

	x, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if !x.Exists {
		t.Fatal("Exists = false")
	}

	if got := classNumbers(x); got != "8,9" {
		t.Fatalf("classes = %s", got)
	}

	c8 := x.Class(8)
	if got := subjectSlugs(c8); got != "mathematics-new,mathematics-old" {
		t.Fatalf("class 8 subjects = %s", got)
	}
	old := c8.Subject("mathematics-old")
	if len(old.Chapters) != 1 || old.Chapters[0].Title != "Factorization of Polynomials" || old.Chapters[0].Number != 1 {
		t.Fatalf("chapters = %+v", old.Chapters)
	}
	if got := lessonNumbers(old.Chapters[0]); got != "1.1,1.2,1.2,1.3,1.10" {
		t.Fatalf("lesson order = %s", got)
	}

	nested := c8.Subject("mathematics-new").Chapters[0].Lessons[0]
	if nested.Path != "Class 8/Maths/Arjun Khatiwada - 2.1 - Sets - Intro to Sets - Mathematics - NEW - 8.MP4" {
		t.Fatalf("nested path = %q", nested.Path)
	}
	if x.Lesson(nested.ID) != nested {
		t.Fatal("lesson lookup by ID failed")
	}

	c9 := x.Class(9)
	if got := subjectSlugs(c9); got != "science-new,science-old" {
		t.Fatalf("class 9 subjects = %s", got)
	}
	sciNew := c9.Subject("science-new")
	if len(sciNew.Chapters) != 1 || sciNew.Chapters[0].Number != 3 {
		t.Fatalf("science-new chapters = %+v", sciNew.Chapters)
	}
	hy := sciNew.Chapters[0].Lessons[0]
	if hy.Teacher != "Ram-Prasad Sharma" || hy.Title != "Mid-term Review" {
		t.Fatalf("hyphenated lesson = %+v", hy)
	}
	// "old" in lower case means the old syllabus, just like "OLD".
	sciOld := c9.Subject("science-old").Chapters[0]
	if got := lessonNumbers(sciOld); got != "1.1,1.2" || sciOld.Lessons[1].Title != "Solids" {
		t.Fatalf("lower-case 'old' should join the old syllabus, got %s", got)
	}

	if x.LessonCount() != 9 {
		t.Fatalf("lesson count = %d", x.LessonCount())
	}
	if len(x.Skipped) != 1 || x.Skipped[0].Path != "random clip.mp4" {
		t.Fatalf("skipped = %+v", x.Skipped)
	}
	wantWarn := map[string]string{
		"backup/Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4": "same lesson as",
	}
	if len(x.Warnings) != len(wantWarn) {
		t.Fatalf("warnings = %+v", x.Warnings)
	}
	for _, w := range x.Warnings {
		if !strings.Contains(w.Message, wantWarn[w.Path]) {
			t.Errorf("warning %+v", w)
		}
	}

	if got := teacherNames(x); got != "Arjun Khatiwada,Ram-Prasad Sharma,Sita Rai" {
		t.Fatalf("teachers = %s", got)
	}
	if x.Teachers[0].Lessons != 6 || strings.Join(x.Teachers[0].Subjects, ",") != "Mathematics" {
		t.Fatalf("teacher = %+v", x.Teachers[0])
	}
}

func TestBuildMissingRoot(t *testing.T) {
	x, err := Build(filepath.Join(t.TempDir(), "nope"))
	if err != nil {
		t.Fatal(err)
	}
	if x.Exists || len(x.Classes) != 0 {
		t.Fatalf("got %+v", x)
	}
}

func TestBuildRootIsFile(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, "file")
	if _, err := Build(filepath.Join(root, "file")); err == nil {
		t.Fatal("expected error")
	}
}

func TestRescanPicksUpChangesAndKeepsIDs(t *testing.T) {
	root := t.TempDir()
	a := "T - 1.1 - C - A - S - NEW - 5.mp4"
	writeFiles(t, root, a)
	l := New(root)
	if l.Index().LessonCount() != 0 {
		t.Fatal("index should start empty")
	}
	x1, err := l.Rescan()
	if err != nil {
		t.Fatal(err)
	}
	id := x1.Class(5).Subjects[0].Chapters[0].Lessons[0].ID

	writeFiles(t, root, "T - 1.2 - C - B - S - NEW - 5.mp4")
	if err := os.Remove(filepath.Join(root, a)); err != nil {
		t.Fatal(err)
	}
	writeFiles(t, root, a) // re-add: same path, same ID
	x2, _ := l.Rescan()
	if x2.LessonCount() != 2 || l.Index() != x2 {
		t.Fatalf("rescan: %d lessons", x2.LessonCount())
	}
	if x2.Lesson(id) == nil {
		t.Fatal("lesson ID changed across rescans")
	}
	if x1.LessonCount() != 1 {
		t.Fatal("old snapshot was mutated")
	}
}

func TestOpenStaysInsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	writeFiles(t, outside, "secret.mp4")
	if err := os.Symlink(filepath.Join(outside, "secret.mp4"), filepath.Join(root, "T - 1 - C - L - S - NEW - 1.mp4")); err != nil {
		t.Skip("symlinks unsupported:", err)
	}
	l := New(root)
	x, _ := l.Rescan()
	lesson := x.Class(1).Subjects[0].Chapters[0].Lessons[0]
	if f, _, err := l.Open(lesson); err == nil {
		f.Close()
		t.Fatal("opened a file outside the library root via symlink")
	}
	if _, _, err := l.Open(&Lesson{Path: "../" + filepath.Base(outside) + "/secret.mp4"}); err == nil {
		t.Fatal("opened a file outside the library root via ..")
	}
}

func TestSlugify(t *testing.T) {
	for in, want := range map[string]string{
		"Mathematics":       "mathematics",
		"Social Studies":    "social-studies",
		"  C++ / Computer ": "c-computer",
		"नेपाली":            "नेपाली",
		"!!!":               "subject",
	} {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSlugCollision(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root,
		"T - 1 - C - L - Maths - NEW - 1.mp4",
		"T - 1 - C - L - maths - NEW - 1.mp4",
	)
	x, _ := Build(root)
	if got := subjectSlugs(x.Class(1)); got != "maths-new,maths-new-2" {
		t.Fatalf("slugs = %s", got)
	}
}

func classNumbers(x *Index) string {
	var s []string
	for _, c := range x.Classes {
		s = append(s, strconv.Itoa(c.Number))
	}
	return strings.Join(s, ",")
}

func subjectSlugs(c *Class) string {
	var s []string
	for _, sub := range c.Subjects {
		s = append(s, sub.Slug)
	}
	return strings.Join(s, ",")
}

func lessonNumbers(ch *Chapter) string {
	var s []string
	for _, l := range ch.Lessons {
		s = append(s, l.Number)
	}
	return strings.Join(s, ",")
}

func teacherNames(x *Index) string {
	var s []string
	for _, t := range x.Teachers {
		s = append(s, t.Name)
	}
	return strings.Join(s, ",")
}
