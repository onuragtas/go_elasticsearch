package v7

import (
	"encoding/json"
	"github.com/onuragtas/go_elasticsearch"
	"github.com/onuragtas/go_elasticsearch/adapters"
	"log"
)

const defaultType = "_doc"
const defaultSize = 100

type ElasticSearchV7 struct {
	Host  string
	Index string
	Type  string
	From  int
	Size  int
}

func NewElasticSearch(host, index, doc string, from, size int) go_elasticsearch.IOperation {
	return &ElasticSearchV7{
		Host:  host,
		Index: index,
		Type:  doc,
		From:  from,
		Size:  size,
	}
}

func (t *ElasticSearchV7) AddToTerm(to []map[string]interface{}, key string, value interface{}) []map[string]interface{} {
	mainTerm := map[string]interface{}{}
	termInterface := map[string]interface{}{}

	termInterface[key] = value
	mainTerm["term"] = termInterface
	to = append(to, mainTerm)
	return to
}
func (t *ElasticSearchV7) AddToExists(to []map[string]interface{}, value interface{}) []map[string]interface{} {
	mainTerm := map[string]interface{}{}
	termInterface := map[string]interface{}{}

	termInterface["exists"] = value
	mainTerm["term"] = termInterface
	to = append(to, mainTerm)
	return to
}

func (t *ElasticSearchV7) AddToRange(slice []map[string]interface{}, key string, from, to interface{}) []map[string]interface{} {
	mainTerm := map[string]interface{}{}
	rangeInterface := map[string]interface{}{}
	defInterface := map[string]interface{}{}

	if from != nil {
		defInterface["from"] = from
	}

	if to != nil {
		defInterface["to"] = to
	}

	rangeInterface[key] = defInterface
	mainTerm["range"] = rangeInterface
	slice = append(slice, mainTerm)
	return slice
}

func (t *ElasticSearchV7) Search(query go_elasticsearch.Main) (go_elasticsearch.Result, error) {
	if t.Type == "" {
		t.Type = defaultType
	}
	if query.Size == 0 {
		query.Size = defaultSize
	}

	byteJson, err := json.Marshal(query)
	if err != nil {
		log.Println(err, "search json error")
	}

	res, err := t.search(byteJson)
	return adapters.Decorate(res)

}

func (t *ElasticSearchV7) AddToTerms(slice []map[string]interface{}, key string, value ...interface{}) []map[string]interface{} {
	mainTerm := map[string]interface{}{}
	termInterface := map[string]interface{}{}

	termInterface[key] = value
	mainTerm["terms"] = termInterface
	slice = append(slice, mainTerm)
	return slice
}

func (t *ElasticSearchV7) Scroll(query go_elasticsearch.Main) (go_elasticsearch.Result, error) {

	if t.Type == "" {
		t.Type = defaultType
	}
	if query.Size == 0 {
		query.Size = defaultSize
	}

	byteJson, err := json.Marshal(query)
	if err != nil {
		log.Println(err, "scroll json error")
	}

	res, err := t.scroll(byteJson)
	return adapters.Decorate(res)
}

func (t *ElasticSearchV7) ScrollById(result go_elasticsearch.Result) (go_elasticsearch.Result, error) {

	if result.ScrollID != "" {

		scrollRequest := map[string]interface{}{}
		scrollRequest["scroll"] = "2m"
		scrollRequest["scroll_id"] = result.ScrollID
		scrollJson, _ := json.Marshal(scrollRequest)
		byteScroll, err := t.request("POST", t.Host+"/_search/scroll", scrollJson)
		if err != nil {
			panic(err)
		}

		return adapters.Decorate(byteScroll)
	}
	return go_elasticsearch.Result{}, nil
}

func (t *ElasticSearchV7) UpdateWithId(id string, source map[string]interface{}) ([]byte, error) {
	scrollJson, _ := json.Marshal(source)
	byteScroll, err := t.request("PUT", t.Host+"/"+t.Index+"/_doc/"+id, scrollJson)
	return byteScroll, err
}

func (t *ElasticSearchV7) UpdateByQuery(query go_elasticsearch.Main) ([]byte, error) {
	scrollJson, _ := json.Marshal(query)
	byteScroll, err := t.request("POST", t.Host+"/"+t.Index+"/_update_by_query", scrollJson)
	return byteScroll, err
}
