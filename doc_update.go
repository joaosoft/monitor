package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joaosoft/errors"
)

type UpdateResponse struct {
	Index   string `json:"_index,omitempty"`
	Type    string `json:"_type,omitempty"`
	ID      string `json:"_id,omitempty"`
	Version int64  `json:"_version,omitempty"`
	Found   bool   `json:"found"`
	Result  string `json:"result"`
	Created bool   `json:"created"`
}

type UpdateService struct {
	client *Elastic
	index  string
	typ    string
	id     string
	body   []byte
	method string
}

func NewUpdateService(e *Elastic) *UpdateService {
	return &UpdateService{
		client: e,
		method: http.MethodPut,
	}
}

func (e *UpdateService) Index(index string) *UpdateService {
	e.index = index
	return e
}

func (e *UpdateService) Type(typ string) *UpdateService {
	e.typ = typ
	return e
}

func (e *UpdateService) Id(id string) *UpdateService {
	e.id = id
	return e
}

func (e *UpdateService) Body(body interface{}) *UpdateService {
	e.body, _ = json.Marshal(body)
	return e
}

func (e *UpdateService) Execute() (string, error) {

	// create data on elastic
	reader := bytes.NewReader(e.body)

	var query string

	if e.id != "" {
		query += fmt.Sprintf("/%s", e.id)
	}

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s/%s%s", e.client.config.Endpoint, e.index, e.typ, query), reader)
	if err != nil {
		return "", err
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

	elasticResponse := UpdateResponse{}
	if err := json.Unmarshal(body, &elasticResponse); err != nil {
		return "", errors.New("0","0", err)
	}

	if elasticResponse.Result != "created" && elasticResponse.Result != "updated" {
		return "", errors.New("0","couldn't update the resource")
	}

	return elasticResponse.ID, nil
}
