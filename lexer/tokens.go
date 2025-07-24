package lexer

type TokenType int

const (
	TEACHER_NAME TokenType = iota
	CHAPTER_NUMBER
	CHAPTER_NAME
	UNIT_NUMBER
	UNIT_NAME
	SUBJECT_NAME
	SUBJECT_TYPE
	CLASS_NUMBER
	SEPERATOR
	EOF
	ERROR
)

type Token struct {
	Type  TokenType
	Value string
	Pos   int
}
