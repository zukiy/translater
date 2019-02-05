package providers

// TranslateProvider ...
type TranslateProvider interface {
	Translate(text, lang string) ([]string, error)
}
