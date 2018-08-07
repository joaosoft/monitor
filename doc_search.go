package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	errors "github.com/joaosoft/errors"
)

type SearchResponse struct {
	Hits   SearchHits   `json:"hits,omitempty"`
	Error  *SearchError `json:"error"`
	Status int          `json:"status"`
}

type SearchHits struct {
	Total int         `json:"total,omitempty"`
	Hits  []SearchHit `json:"hits,omitempty"`
}

type SearchHit struct {
	Index   string          `json:"_index,omitempty"`
	Type    string          `json:"_type,omitempty"`
	ID      string          `json:"_id,omitempty"`
	Version int64           `json:"_version,omitempty"`
	Found   bool            `json:"found"`
	Source  json.RawMessage `json:"_source,omitempty"`
}

type SearchError struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type SearchService struct {
	client *Elastic
	index  string
	typ    string
	id     string
	query  string
	object interface{}
	method string
}

func NewSearchService(client *Elastic) *SearchService {
	return &SearchService{
		client: client,
		method: http.MethodGet,
	}
}

func (e *SearchService) Index(index string) *SearchService {
	e.index = index
	return e
}

func (e *SearchService) Type(typ string) *SearchService {
	e.typ = typ
	return e
}

func (e *SearchService) Id(id string) *SearchService {
	e.id = id
	return e
}

func (e *SearchService) Query(query string) *SearchService {
	e.query = query
	return e
}

func (e *SearchService) Object(object interface{}) *SearchService {
	e.object = object
	return e
}

type SearchTemplate struct {
	Data interface{} `json:"data,omitempty"`
	From int         `json:"from,omitempty"`
	Size int         `json:"size,omitempty"`
}

func (e *SearchService) Template(path, name string, data *SearchTemplate, reload bool) *SearchService {
	key := fmt.Sprintf("%s/%s", path, name)

	var result bytes.Buffer
	var err error

	if _, found := templates[key]; !found {
		e.client.mux.Lock()
		defer e.client.mux.Unlock()
		templates[key], err = ReadFile(key, nil)
		if err != nil {
			log.Error(err)
			return e
		}
	}

	t := template.New(name)
	t, err = t.Parse(string(templates[key]))
	if err == nil {
		if err := t.ExecuteTemplate(&result, name, data); err != nil {
			log.Error(err)
			return e
		}

		e.query = result.String()
	} else {
		log.Error(err)
		return e
	}

	return e
}

func (e *SearchService) Execute() error {

	if e.query != "" {
		e.method = http.MethodPost
	}

	// get data from elastic
	reader := strings.NewReader(e.query)

	var q string
	if e.typ != "" {
		q += fmt.Sprintf("/%s", e.typ)
	}

	if e.id != "" {
		q += fmt.Sprintf("/%s", e.id)
	} else {
		q += "/_search"
	}

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s%s", e.client.config.Endpoint, e.index, q), reader)
	if err != nil {
		return errors.New("0", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error(err)
		return errors.New("0", err)
	}
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return errors.New("0", err)
	}

	var hit []byte

	if e.id != "" {
		elasticResponse := SearchHit{}
		if err := json.Unmarshal(body, &elasticResponse); err != nil {
			log.Error(err)
			return errors.New("0", err)
		}

		hit, err = json.Marshal(elasticResponse.Source)
		if err != nil {
			log.Error(err)
			return errors.New("0", err)
		}
	} else {
		elasticResponse := SearchResponse{}
		if err := json.Unmarshal(body, &elasticResponse); err != nil {
			log.Error(err)
			return errors.New("0", err)
		}

		if elasticResponse.Error != nil {
			return errors.New(fmt.Sprintf("[%s] %s", elasticResponse.Error.Type, elasticResponse.Error.Reason))
		}

		rawHits := make([]json.RawMessage, len(elasticResponse.Hits.Hits))
		for i, rawHit := range elasticResponse.Hits.Hits {
			rawHits[i] = rawHit.Source
		}

		hit, err = json.Marshal(rawHits)
		if err != nil {
			return errors.New("0", err)
		}
	}

	if err := json.Unmarshal(hit, e.object); err != nil {
		return errors.New("0", err)
	}

	return nil
}
