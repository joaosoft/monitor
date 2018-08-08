package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joaosoft/errors"
)

type CreateResponse struct {
	Index   string `json:"_index,omitempty"`
	Type    string `json:"_type,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int64  `json:"_version,omitempty"`
	Found   bool   `json:"found"`
	Result  string `json:"result"`
	Created bool   `json:"created"`
}

type CreateService struct {
	client *Elastic
	index  string
	typ    string
	id     string
	body   []byte
	method string
}

func NewCreateService(e *Elastic) *CreateService {
	return &CreateService{
		client: e,
		method: http.MethodPost,
	}
}

func (e *CreateService) Index(index string) *CreateService {
	e.index = index
	return e
}

func (e *CreateService) Type(typ string) *CreateService {
	e.typ = typ
	return e
}

func (e *CreateService) Id(id string) *CreateService {
	e.id = id
	return e
}

func (e *CreateService) Body(body interface{}) *CreateService {
	switch v := body.(type) {
	case []byte:
		e.body = v
	default:
		e.body, _ = json.Marshal(v)
	}
	return e
}

func (e *CreateService) Execute() (string, error) {

	// create data on elastic
	reader := bytes.NewReader(e.body)

	var query string

	if e.id != "" {
		query += fmt.Sprintf("/%s", e.id)
	}

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s/%s%s", e.client.config.Endpoint, e.index, e.typ, query), reader)
	if err != nil {
		return "", errors.New("0","0", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", errors.New("0","0", err)
	}
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("0","0", err)
	}

	elasticResponse := CreateResponse{}
	if err = json.Unmarshal(body, &elasticResponse); err != nil {
		return "", errors.New("0","0", err)
	}

	if !elasticResponse.Created && elasticResponse.Result != "updated" {
		return "", errors.New("0","couldn't create the resource")
	}

	return elasticResponse.ID, nil
}
