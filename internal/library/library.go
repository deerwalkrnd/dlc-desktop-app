package library

import (
	"io/fs"
	"os"
	"sync"
	"sync/atomic"
)

// Library owns the current Index and serialises rescans.
type Library struct {
	root    string
	current atomic.Pointer[Index]
	scanMu  sync.Mutex
}

// New returns a Library for root. Call Rescan before reading the index.
func New(root string) *Library {
	l := &Library{root: root}
	l.current.Store(&Index{Root: root, classes: map[int]*Class{}, lessons: map[string]*Lesson{}})
	return l
}

// Root is the library folder.
func (l *Library) Root() string { return l.root }

// Index returns the current snapshot. It is safe for concurrent use.
func (l *Library) Index() *Index { return l.current.Load() }

// Rescan rebuilds the index from disk and swaps it in. On error the previous
// index stays in place. Concurrent calls are serialised.
func (l *Library) Rescan() (*Index, error) {
	l.scanMu.Lock()
	defer l.scanMu.Unlock()
	x, err := Build(l.root)
	if err != nil {
		return nil, err
	}
	l.current.Store(x)
	return x, nil
}

// Open opens a lesson's video file. Paths are resolved inside the library
// root, so a symlink or crafted name can never escape it.
func (l *Library) Open(lesson *Lesson) (*os.File, fs.FileInfo, error) {
	f, err := os.OpenInRoot(l.root, lesson.Path)
	if err != nil {
		return nil, nil, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	if info.IsDir() {
		f.Close()
		return nil, nil, fs.ErrNotExist
	}
	return f, info, nil
}
