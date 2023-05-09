package matchers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/MarcosRoch4/gofeed/search"
)

type (
	//item define os campos associados à tag item
	// no documento rss
	item struct {
		XMLName     xml.Name `xml:"item"`
		PubDate     string   `xml:"pubDate"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Link        string   `xml:"link"`
		GUID        string   `xml:"guid"`
		GeoRssPoint string   `xml:"georss:point"`
	}

	// image define os campos associados à tag image
	// no documento rss

	image struct {
		XMLName xml.Name `xml:"image"`
		URL     string   `xml:"url"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
	}

	// channel define os campos associados à tag channel
	// no documento rss
	channel struct {
		XMLName        xml.Name `xml:"channel"`
		Title          string   `xml:"title"`
		Description    string   `xml:"description"`
		Link           string   `xml:"link"`
		PubDate        string   `xml:"pubDate"`
		LastBuildDate  string   `xml:"lastBuildDate"`
		TTL            string   `xml:"ttl"`
		Language       string   `xml:"language"`
		ManagingEditor string   `xml:"managingEditor"`
		WebMaster      string   `xml:"webMaster"`
		Image          image    `xml:"image"`
		Item           []item   `xml:"item"`
	}
	//  rssDocument define os campos associados ao documeno rss
	rssDocument struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)

// rssMatcher implementa a interface Matcher
type rssMatcher struct{}

// init registra o matcher junto ao programa
func init() {
	var matcher rssMatcher
	search.Register("rss", matcher)
}

// Search procura o termo de pesquisa especificado no documento
func (m rssMatcher) Search(feed *search.Feed, searchTerm string) ([]*search.Result, error) {
	var results []*search.Result

	log.Printf("Search Feed Type[%s] Site[%s] for Uri[%s]\n", feed.Type, feed.Name, feed.URI)

	// Obtém o dado para pesquisar
	document, err := m.retrieve(feed)
	if err != nil {
		return nil, err
	}

	for _, channelItem := range document.Channel.Item {
		// Verifica se o termo de pesquisa está no título
		matched, err := regexp.MatchString(searchTerm, channelItem.Title)
		if err != nil {
			return nil, err
		}

		// Se houver uma correspondência, salva o resultado
		if matched {
			results = append(results, &search.Result{
				Field:   "Title",
				Content: channelItem.Title,
			})
		}

		// Verifica se o termo de pesquisa está na descrição
		matched, err = regexp.MatchString(searchTerm, channelItem.Description)
		if err != nil {
			return nil, err
		}

		// Se houver uma correspondência, salva o resultado
		if matched {
			results = append(results, &search.Result{
				Field:   "Description",
				Content: channelItem.Description,
			})
		}
	}

	return results, nil
}

// retrieve faz uma requisição HTTP GET  para o feed rss e o decodifica
func (m rssMatcher) retrieve(feed *search.Feed) (*rssDocument, error) {
	if feed.URI == "" {
		return nil, errors.New("No rss feed URI provided")
	}

	// Obtém o documento de feed rss da web
	resp, err := http.Get(feed.URI)
	if err != nil {
		return nil, err
	}

	// Fecha a resposta depois que retornarmo da função
	defer resp.Body.Close()

	// Verifica se o código de status é 200 para saber se recebemos uma
	// resposta apropriada
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Rresponse Error %d\n", resp.StatusCode)
	}

	// Decodifica o documento de feed rss em nosso tipo estrutura
	// Não precisamos verifica erros; quem chama a função pode fazer isso
	var document rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&document)
	return &document, err
}
