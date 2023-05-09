package search

import (
	"log"
	"sync"
)

// Um mapa dos matchers registrados para pesquisa
var matchers = make(map[string]Matcher)

// Run executa a lógica de pesquisa
func Run(searchTerm string) {
	//Obtém a lista de feeds para pesquisar
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	// Cria um canal sem buffer para receber os resultados das correspondências
	results := make(chan *Result)

	// Cria um grupo de espera para que possamos processar todos os feeds
	var waitGroup sync.WaitGroup

	// Definie o número de goroutines que precisamos esperar enquanto
	// elas processam os feeds individuais
	waitGroup.Add(len(feeds))

	// Inicia uma goroutine para cada fees a fim de obter os resultados
	for _, feed := range feeds {
		// Obtém um matcher para a pesquisa
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}

		// Inicia a goroutine para fazer a pesquisa
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			waitGroup.Done()
		}(matcher, feed)
	}

	// Inicia uma goroutine para saber quando todo o trabalho foi feito
	go func() {
		// Espera que tudo seja processado
		waitGroup.Wait()

		// Fecha o canal para sinalizar à função Display
		// que podemos encerrar o programa
		close(results)
	}()

	// Começa a exibir os resultados à medida que são disponibilizados e
	// retorna depois que o último resultado é exibido
	Display(results)

}

// Register é chamada pra registrar um matcher a ser usado pelo programa
func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Matcher already registered")
	}

	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher

}
