package yandex

// Format parameter
type Format string

// String format
func (f Format) String() string {
	return string(f)
}

const (
	baseURL       = "https://translate.yandex.net/api"
	jsonInterface = "tr.json"

	// Plain text without markdown, default
	Plain Format = "plain"

	// HTML text
	HTML Format = "html"
)

// TranslateResponse ...
type TranslateResponse struct {
	Code int
	Lang string
	Text []string
}

// Error response struct
type Error struct {
	Code    int
	Message string
}
