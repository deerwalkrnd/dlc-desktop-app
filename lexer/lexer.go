package lexer

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

type VideoLexer struct {
	Input    string
	Position int
	Tokens   []Token
}

func NewVideoLexer(filename string) *VideoLexer {
	basename := filepath.Base(filename)
	if len(basename) > 4 && strings.HasSuffix(basename, ".mp4") {
		basename = basename[:len(basename)-4]
	}
	return &VideoLexer{
		Input:    basename,
		Position: 0,
		Tokens:   make([]Token, 0),
	}
}

// Peek returns the current rune without advancing, or 0 if we're at/after the end.
func (l *VideoLexer) Peek() rune {
	if l.Position >= len(l.Input) {
		return 0
	}
	return rune(l.Input[l.Position])
}

// PeekN peeks n runes from the current position (can be negative). Returns 0 if OOB.
func (l *VideoLexer) PeekN(n int) rune {
	pos := l.Position + n
	if pos < 0 || pos >= len(l.Input) {
		return 0
	}
	return rune(l.Input[pos])
}

// Next returns current rune and advances, or 0 if we're at/after the end.
func (l *VideoLexer) Next() rune {
	if l.Position >= len(l.Input) {
		return 0
	}
	ch := rune(l.Input[l.Position])
	l.Position++
	return ch
}

func (l *VideoLexer) SkipWhiteSpace() {
	for unicode.IsSpace(l.Peek()) {
		l.Next()
	}
}

// isDashSeparator returns true if the current rune is '-' and it is surrounded by spaces (" - ").
func (l *VideoLexer) isDashSeparator() bool {
	if l.Peek() != '-' {
		return false
	}
	prev := l.PeekN(-1)
	next := l.PeekN(1)
	return unicode.IsSpace(prev) && unicode.IsSpace(next)
}

func (l *VideoLexer) ReadUntilSeperator() string {
	start := l.Position

	for {
		ch := l.Peek()
		if ch == 0 {
			break
		}

		// Stop at " - " (space-dash-space) only
		if ch == '-' && l.isDashSeparator() {
			break
		}

		l.Next()
	}

	return strings.TrimSpace(l.Input[start:l.Position])
}

func (l *VideoLexer) Tokenize() error {
	tokenTypes := []TokenType{
		TEACHER_NAME,
		UNIT_NUMBER,
		CHAPTER_NAME,
		UNIT_NAME,
		SUBJECT_NAME,
		SUBJECT_TYPE,
		CLASS_NUMBER,
	}

	tokenIndex := 0

	for l.Position < len(l.Input) && tokenIndex < len(tokenTypes) {
		l.SkipWhiteSpace()
		if l.Position >= len(l.Input) {
			break
		}

		// Explicitly consume the " - " separators if present
		if l.Peek() == '-' && l.isDashSeparator() {
			// consume '-' and continue (skip it)
			l.Next()
			continue
		}

		value := l.ReadUntilSeperator()
		if value == "" {
			// If we didn't advance and value is empty, break to avoid infinite loop
			if l.Peek() == 0 {
				break
			}
			continue
		}

		l.Tokens = append(l.Tokens, Token{
			Type:  tokenTypes[tokenIndex],
			Value: value,
			Pos:   l.Position - len(value),
		})

		tokenIndex++
	}

	l.Tokens = append(l.Tokens, Token{
		Type:  EOF,
		Value: "",
		Pos:   l.Position,
	})

	return nil
}

func (l *VideoLexer) GetTokenByType(tokenType TokenType) (*Token, error) {
	for _, token := range l.Tokens {
		if token.Type == tokenType {
			return &token, nil
		}
	}
	return nil, fmt.Errorf("token type %d not found", tokenType)
}

func (l *VideoLexer) Validate() error {
	requiredTokens := []TokenType{
		TEACHER_NAME,
		UNIT_NUMBER,
		CHAPTER_NAME,
		UNIT_NAME,
		SUBJECT_NAME,
		SUBJECT_TYPE,
		CLASS_NUMBER,
	}
	for _, required := range requiredTokens {
		if _, err := l.GetTokenByType(required); err != nil {
			return fmt.Errorf("missing required token: %d", required)
		}
	}
	return nil
}
