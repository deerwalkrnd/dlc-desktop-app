package data

type Parser interface {
	ParseVideo(path string) *Video
}
