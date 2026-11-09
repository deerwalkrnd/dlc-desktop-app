package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/deerwalkrnd/dlc-desktop-app/internal/library"
)

const host = "127.0.0.1:3000"

var videoBytes = []byte("0123456789abcdefghij") // stands in for an .mp4

func newTestServer(t *testing.T, ui fstest.MapFS) (*Server, string) {
	t.Helper()
	root := t.TempDir()
	files := []string{
		"Arjun Khatiwada - 1.10 - Factorization of Polynomials - Polynomials in Maths real 1.10 - Mathematics - OLD - 8 .mp4",
		"Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4",
		"Arjun Khatiwada - 1.1 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4",
		"nested/folder/Sita Rai - 2.1 - Sets - Intro #1 - Mathematics - NEW - 8.MP4",
		"Ram-Prasad Sharma - 1.1 - Light - Mid-term Review - Science - NEW - 9 .mp4",
		"badly named.mp4",
	}
	for _, f := range files {
		p := filepath.Join(root, filepath.FromSlash(f))
		os.MkdirAll(filepath.Dir(p), 0o755)
		if err := os.WriteFile(p, videoBytes, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	lib := library.New(root)
	if _, err := lib.Rescan(); err != nil {
		t.Fatal(err)
	}
	if ui == nil {
		ui = fstest.MapFS{
			"index.html":                  {Data: []byte("<html>app</html>")},
			"_app/immutable/entry/app.js": {Data: []byte("console.log(1)")},
			"_app/version.json":           {Data: []byte(`{}`)},
			"deerwalk.png":                {Data: []byte("png")},
		}
	}
	s, err := New(lib, Options{UI: ui, AllowedHosts: LoopbackHosts(host), Version: "test"})
	if err != nil {
		t.Fatal(err)
	}
	return s, root
}

func do(t *testing.T, s http.Handler, method, target string, hdr map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, nil)
	req.Host = host
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	return rec
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("decode %q: %v", rec.Body.String(), err)
	}
	return v
}

func TestClasses(t *testing.T) {
	s, _ := newTestServer(t, nil)
	rec := do(t, s, "GET", "/api/classes", nil)
	if rec.Code != 200 || !strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("%d %s", rec.Code, rec.Header())
	}
	got := decode[struct{ Classes []classSummary }](t, rec)
	if len(got.Classes) != 2 {
		t.Fatalf("%+v", got)
	}
	c8 := got.Classes[0]
	if c8.Number != 8 || c8.SubjectCount != 2 || c8.LessonCount != 4 || strings.Join(c8.Syllabi, ",") != "new,old" {
		t.Fatalf("%+v", c8)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("CORS must not be enabled")
	}
}

func TestClassDetailAndSyllabusFilter(t *testing.T) {
	s, _ := newTestServer(t, nil)

	all := decode[classDetail](t, do(t, s, "GET", "/api/classes/8", nil))
	if len(all.Subjects) != 2 {
		t.Fatalf("%+v", all)
	}

	old := decode[classDetail](t, do(t, s, "GET", "/api/classes/8?syllabus=old", nil))
	if len(old.Subjects) != 1 || old.Subjects[0].Slug != "mathematics-old" || old.Subjects[0].LessonCount != 3 {
		t.Fatalf("%+v", old)
	}

	for target, code := range map[string]int{
		"/api/classes/8?syllabus=bogus": 400,
		"/api/classes/abc":              400,
		"/api/classes/42":               404,
		"/api/classes/8/subjects/nope":  404,
		"/api/nothing":                  404,
	} {
		rec := do(t, s, "GET", target, nil)
		if rec.Code != code {
			t.Errorf("%s: %d, want %d", target, rec.Code, code)
		}
		if decode[map[string]string](t, rec)["error"] == "" {
			t.Errorf("%s: missing error message", target)
		}
	}
}

func TestSubjectDetail(t *testing.T) {
	s, _ := newTestServer(t, nil)
	rec := do(t, s, "GET", "/api/classes/8/subjects/mathematics-old", nil)
	if rec.Code != 200 {
		t.Fatal(rec.Code, rec.Body)
	}
	got := decode[subjectDetail](t, rec)
	if got.Name != "Mathematics" || got.Syllabus != "old" || len(got.Chapters) != 1 {
		t.Fatalf("%+v", got)
	}
	var nums []string
	for _, l := range got.Chapters[0].Lessons {
		nums = append(nums, l.Number)
		if l.MediaURL != "/media/"+l.ID || l.Teacher != "Arjun Khatiwada" {
			t.Fatalf("%+v", l)
		}
	}
	if strings.Join(nums, ",") != "1.1,1.2,1.10" {
		t.Fatalf("order %v", nums)
	}
}

func TestTeachersAndLibrary(t *testing.T) {
	s, _ := newTestServer(t, nil)
	teachers := decode[struct{ Teachers []teacherDTO }](t, do(t, s, "GET", "/api/teachers", nil))
	if len(teachers.Teachers) != 3 || teachers.Teachers[1].Name != "Ram-Prasad Sharma" {
		t.Fatalf("%+v", teachers)
	}

	lib := decode[libraryDTO](t, do(t, s, "GET", "/api/library", nil))
	if !lib.Exists || lib.LessonCount != 5 || len(lib.Skipped) != 1 || lib.Skipped[0].Path != "badly named.mp4" || lib.Version != "test" {
		t.Fatalf("%+v", lib)
	}
}

func TestRescan(t *testing.T) {
	s, root := newTestServer(t, nil)
	os.WriteFile(filepath.Join(root, "T - 1 - C - L - Nepali - NEW - 10.mp4"), videoBytes, 0o644)

	rec := do(t, s, "POST", "/api/library/rescan", map[string]string{"Origin": "http://" + host, "Sec-Fetch-Site": "same-origin"})
	if rec.Code != 200 || decode[libraryDTO](t, rec).LessonCount != 6 {
		t.Fatal(rec.Code, rec.Body)
	}

	for _, hdr := range []map[string]string{
		{"Origin": "https://evil.example"},
		{"Sec-Fetch-Site": "cross-site"},
	} {
		if rec := do(t, s, "POST", "/api/library/rescan", hdr); rec.Code != 403 {
			t.Errorf("%v: %d, want 403", hdr, rec.Code)
		}
	}
	if rec := do(t, s, "GET", "/api/library/rescan", nil); rec.Code == 200 {
		t.Error("GET must not trigger a rescan")
	}
}

func TestMediaStreaming(t *testing.T) {
	s, root := newTestServer(t, nil)
	subj := decode[subjectDetail](t, do(t, s, "GET", "/api/classes/8/subjects/mathematics-new", nil))
	lesson := subj.Chapters[0].Lessons[0] // lives in a nested folder, name contains '#'

	rec := do(t, s, "GET", lesson.MediaURL, nil)
	if rec.Code != 200 || rec.Body.String() != string(videoBytes) || rec.Header().Get("Content-Type") != "video/mp4" {
		t.Fatalf("%d %q %s", rec.Code, rec.Body, rec.Header())
	}
	if rec.Header().Get("Accept-Ranges") != "bytes" {
		t.Fatal("range requests not advertised")
	}

	rec = do(t, s, "GET", lesson.MediaURL, map[string]string{"Range": "bytes=10-14"})
	if rec.Code != http.StatusPartialContent || rec.Body.String() != "abcde" {
		t.Fatalf("range: %d %q", rec.Code, rec.Body)
	}

	if rec := do(t, s, "GET", "/media/doesnotexist", nil); rec.Code != 404 {
		t.Fatalf("unknown id: %d", rec.Code)
	}

	// File deleted after the scan: clean 404, not a 500 or a panic.
	os.Remove(filepath.Join(root, "nested/folder/Sita Rai - 2.1 - Sets - Intro #1 - Mathematics - NEW - 8.MP4"))
	if rec := do(t, s, "GET", lesson.MediaURL, nil); rec.Code != 404 {
		t.Fatalf("missing file: %d", rec.Code)
	}
}

func TestNoRawFilesystemAccess(t *testing.T) {
	s, _ := newTestServer(t, nil)
	for _, target := range []string{"/videos/", "/media/", "/media/..%2f..%2fetc%2fpasswd", "/DLC/"} {
		rec := do(t, s, "GET", target, nil)
		if strings.Contains(rec.Body.String(), "Arjun") || strings.Contains(rec.Body.String(), "<pre>") {
			t.Errorf("%s leaked a file listing: %q", target, rec.Body)
		}
	}
}

func TestHostGuard(t *testing.T) {
	s, _ := newTestServer(t, nil)
	for h, want := range map[string]int{
		"127.0.0.1:3000":        200,
		"localhost:3000":        200,
		"LOCALHOST:3000":        200,
		"[::1]:3000":            200,
		"attacker.example":      421,
		"attacker.example:3000": 421,
		"127.0.0.1.nip.io:3000": 421,
	} {
		req := httptest.NewRequest("GET", "/api/classes", nil)
		req.Host = h
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
		if rec.Code != want {
			t.Errorf("Host %s: %d, want %d", h, rec.Code, want)
		}
	}
}

func TestSecurityHeaders(t *testing.T) {
	s, _ := newTestServer(t, nil)
	rec := do(t, s, "GET", "/", nil)
	for _, h := range []string{"X-Content-Type-Options", "X-Frame-Options", "Referrer-Policy"} {
		if rec.Header().Get(h) == "" {
			t.Errorf("missing %s", h)
		}
	}
}

func TestSPA(t *testing.T) {
	s, _ := newTestServer(t, nil)
	for target, want := range map[string]string{
		"/":                                "<html>app</html>",
		"/index.html":                      "<html>app</html>",
		"/class/8":                         "<html>app</html>",
		"/class/8/subject/mathematics-old": "<html>app</html>",
		"/teachers":                        "<html>app</html>",
		"/deerwalk.png":                    "png",
		"/_app/immutable/entry/app.js":     "console.log(1)",
	} {
		rec := do(t, s, "GET", target, nil)
		if rec.Code != 200 || rec.Body.String() != want {
			t.Errorf("%s: %d %q", target, rec.Code, rec.Body)
		}
	}
	rec := do(t, s, "GET", "/_app/immutable/entry/app.js", nil)
	if !strings.Contains(rec.Header().Get("Cache-Control"), "immutable") {
		t.Error("hashed assets should be cached as immutable")
	}
	if rec := do(t, s, "GET", "/", nil); rec.Header().Get("Cache-Control") != "no-cache" {
		t.Error("index.html must not be cached")
	}
	if rec := do(t, s, "GET", "/_app/immutable/missing.js", nil); rec.Code != 404 {
		t.Errorf("missing chunk: %d", rec.Code)
	}
	if rec := do(t, s, "GET", "/_app/", nil); rec.Code != 404 || strings.Contains(rec.Body.String(), "version.json") {
		t.Errorf("directory listing: %d %q", rec.Code, rec.Body)
	}
}

func TestPlaceholderWithoutUI(t *testing.T) {
	s, _ := newTestServer(t, fstest.MapFS{".gitkeep": {}})
	rec := do(t, s, "GET", "/class/8", nil)
	body, _ := io.ReadAll(rec.Body)
	if rec.Code != 200 || !strings.Contains(string(body), "make build") {
		t.Fatalf("%d %q", rec.Code, body)
	}
	if rec := do(t, s, "GET", "/api/classes", nil); rec.Code != 200 {
		t.Fatal("API must work without the UI")
	}
}

func TestHealthz(t *testing.T) {
	s, _ := newTestServer(t, nil)
	got := decode[map[string]any](t, do(t, s, "GET", "/healthz", nil))
	if got["app"] != "dlc-desktop-app" || got["ok"] != true {
		t.Fatalf("%v", got)
	}
}

func TestLoopbackHosts(t *testing.T) {
	if LoopbackHosts("0.0.0.0:3000") != nil || LoopbackHosts("192.168.1.5:3000") != nil {
		t.Fatal("non-loopback addresses must disable the host guard")
	}
	if len(LoopbackHosts("127.0.0.1:8080")) != 3 || len(LoopbackHosts("localhost:3000")) != 3 {
		t.Fatal("loopback addresses must enable the host guard")
	}
}

func TestSearchEndpoint(t *testing.T) {
	s, _ := newTestServer(t, nil)
	got := decode[searchDTO](t, do(t, s, "GET", "/api/search?q=mid+term", nil))
	if got.Total != 1 || len(got.Results) != 1 {
		t.Fatalf("%+v", got)
	}
	r := got.Results[0]
	if r.Class != 9 || r.SubjectSlug != "science-new" || r.Lesson.Title != "Mid-term Review" || r.Lesson.MediaURL != "/media/"+r.Lesson.ID || r.ChapterTitle != "Light" {
		t.Fatalf("%+v", r)
	}

	empty := decode[searchDTO](t, do(t, s, "GET", "/api/search?q=", nil))
	if empty.Total != 0 || empty.Results == nil {
		t.Fatalf("empty query: %+v (results must be [] not null)", empty)
	}

	limited := decode[searchDTO](t, do(t, s, "GET", "/api/search?q=mathematics&limit=1", nil))
	if limited.Total != 4 || len(limited.Results) != 1 {
		t.Fatalf("limit: %+v", limited)
	}

	for target, code := range map[string]int{
		"/api/search?q=x&limit=0":                   400,
		"/api/search?q=x&limit=abc":                 400,
		"/api/search?q=" + strings.Repeat("a", 201): 400,
	} {
		if rec := do(t, s, "GET", target, nil); rec.Code != code {
			t.Errorf("%s: %d, want %d", target, rec.Code, code)
		}
	}
}
