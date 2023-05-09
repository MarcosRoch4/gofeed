package search

import (
	"fmt"
	"log"
)

// Result contém o resultado de uma pesquisa
type Result struct {
	Field   string
	Content string
}

// Matcher define o comportamento exigido pelos tipos que querem
// implementar um novo tipo de pesquisa
type Matcher interface {
	Search(feed *Feed, searchTerm string) ([]*Result, error)
}

// Match é iniciada como uma goroutine para cada feed individual a fim de executar
// pesquisas de forma concorrente
func Match(matcher Matcher, feed *Feed, searchTerm string, results chan<- *Result) {
	// Faz a pesquisa para o matcher especificado
	seachResults, err := matcher.Search(feed, searchTerm)
	if err != nil {
		log.Println(err)
		return
	}

	// Escreve os resultados no canal
	for _, result := range seachResults {
		results <- result
	}
}

// Display escreve os resultados na janela do terminal à medida que
// são recebidos pelas goroutines individuais
func Display(results chan *Result) {
	// O canal fica bloqueado até que um resultada seja escrito nele
	// Depois que o canal é fechado, o laço for termina
	for result := range results {
		fmt.Printf("%s:\n%s\n\n", result.Field, result.Content)
	}
}
