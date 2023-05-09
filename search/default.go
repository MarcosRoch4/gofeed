package search

// defaultMatcher implementa o matcher default
type defaultMatcher struct{}

// init registra o matcher default junto ao programa
func init() {
	var matcher defaultMatcher
	Register("default", matcher)
}

// Search implementa o comportamento do matcher default
func (m defaultMatcher) Search(feed *Feed, searchTerm string) ([]*Result, error) {
	return nil, nil
}
