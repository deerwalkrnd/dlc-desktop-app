package library

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// Hit is one search result with the context needed to display and open it.
type Hit struct {
	Lesson  *Lesson
	Subject *Subject
	Chapter *Chapter
	Score   int
}

// searchDoc is a lesson's precomputed, normalised searchable text.
type searchDoc struct {
	lesson  *Lesson
	subject *Subject
	chapter *Chapter
	order   int // position in display order, for stable tie-breaking

	title, chapterTitle, subjectName, teacher, number, class string
}

// buildSearch precomputes search documents in display order. Called once at
// the end of Build, after sorting.
func (x *Index) buildSearch() {
	for _, c := range x.Classes {
		class := "class " + strconv.Itoa(c.Number)
		for _, s := range c.Subjects {
			for _, ch := range s.Chapters {
				for _, l := range ch.Lessons {
					x.search = append(x.search, searchDoc{
						lesson:       l,
						subject:      s,
						chapter:      ch,
						order:        len(x.search),
						title:        normalize(l.Title),
						chapterTitle: normalize(ch.Title),
						subjectName:  normalize(s.Name),
						teacher:      normalize(l.Teacher),
						number:       l.Number,
						class:        class,
					})
				}
			}
		}
	}
}

// Search returns lessons matching every word of query, best matches first.
// Words can match the lesson title, chapter, subject, teacher, class
// ("class 8" or just "8") or lesson number ("1.2"). Matching is
// case-insensitive; a match at the start of a word ranks higher.
func (x *Index) Search(query string, limit int) (hits []Hit, total int) {
	terms := strings.Fields(normalize(query))
	if len(terms) == 0 {
		return nil, 0
	}
	// "class 8" is a common phrasing; treat "class" as a filler word when
	// followed by a number.
	terms = slices.DeleteFunc(slices.Clone(terms), func(t string) bool { return t == "class" })
	if len(terms) == 0 {
		return nil, 0
	}

	type scored struct {
		doc   *searchDoc
		score int
	}
	var matches []scored
docs:
	for i := range x.search {
		d := &x.search[i]
		score := 0
		for _, t := range terms {
			s := d.termScore(t)
			if s == 0 {
				continue docs
			}
			score += s
		}
		matches = append(matches, scored{d, score})
	}

	slices.SortStableFunc(matches, func(a, b scored) int {
		return cmp.Or(cmp.Compare(b.score, a.score), cmp.Compare(a.doc.order, b.doc.order))
	})
	total = len(matches)
	if limit > 0 && len(matches) > limit {
		matches = matches[:limit]
	}
	hits = make([]Hit, len(matches))
	for i, m := range matches {
		hits[i] = Hit{Lesson: m.doc.lesson, Subject: m.doc.subject, Chapter: m.doc.chapter, Score: m.score}
	}
	return hits, total
}

// termScore weights where a term matched; 0 means no match anywhere.
func (d *searchDoc) termScore(t string) int {
	best := 0
	try := func(field string, weight int) {
		if s := fieldScore(field, t, weight); s > best {
			best = s
		}
	}
	try(d.title, 6)
	try(d.chapterTitle, 4)
	try(d.subjectName, 3)
	try(d.teacher, 3)
	if d.number == t || strings.HasPrefix(d.number, t+".") {
		best = max(best, 5)
	}
	if strings.TrimPrefix(d.class, "class ") == t {
		best = max(best, 2)
	}
	return best
}

// fieldScore returns weight*2 for a match at a word start, weight for a
// match elsewhere, and 0 for no match.
func fieldScore(field, term string, weight int) int {
	i := strings.Index(field, term)
	if i < 0 {
		return 0
	}
	for i >= 0 {
		if i == 0 || field[i-1] == ' ' {
			return weight * 2
		}
		next := strings.Index(field[i+1:], term)
		if next < 0 {
			break
		}
		i += 1 + next
	}
	return weight
}

// normalize lower-cases s and turns punctuation into single spaces, so that
// "Mid-term Review" is found by "mid term" or "term". Dots are kept for
// lesson numbers like "1.2".
func normalize(s string) string {
	var sb strings.Builder
	space := true
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r) || r == '.' {
			sb.WriteRune(r)
			space = false
		} else if !space {
			sb.WriteByte(' ')
			space = true
		}
	}
	return strings.TrimSpace(sb.String())
}
