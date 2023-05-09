package search

import(
	"encoding/json"
	"os"
)

const dataFile := "data/data.json"

// Feed contém informações necessárias para processar um feed
type Feed struct {
	Name string `json:"site"`
	URI string `json:"site"`
	Type string `json:"type"`
}

// RetrieveFeeds lê e faz o unmarshal do arquivo de dados de feed
func RetrieveFeeds()([]*Feed,error)  {
	//Abre o arquivo
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	// Escalona o arquivo para ser fechado depois
	// que a função retornar
	defer file.Close()

	// Decodifica o arquivo em uma fatia de ponteiros
	// para valores do tipo Feed
	var feeds []*Feed
	err = json.NewDecoder(file).Decode(&feeds)

	// Não precisamos verificar erros; Quem chama a função pode fazer isso
	return feeds, err
}