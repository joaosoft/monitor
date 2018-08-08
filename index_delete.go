package elastic

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joaosoft/errors"
)

type DeleteIndexResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

type DeleteIndexService struct {
	client *Elastic
	index  string
	typ    string
	method string
}

func NewDeleteIndexService(client *Elastic) *DeleteIndexService {
	return &DeleteIndexService{
		client: client,
		method: http.MethodDelete,
	}
}

func (e *DeleteIndexService) Index(index string) *DeleteIndexService {
	e.index = index
	return e
}

func (e *DeleteIndexService) Type(typ string) *DeleteIndexService {
	e.typ = typ
	return e
}

func (e *DeleteIndexService) Execute() error {

	// delete data from elastic
	var query string
	if e.typ != "" {
		query += fmt.Sprintf("/%s", e.typ)
	}

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s%s", e.client.config.Endpoint, e.index, query), nil)
	if err != nil {
		return errors.New("0","0", err)
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)

	elasticResponse := DeleteIndexResponse{}
	json.Unmarshal(body, &elasticResponse)

	if !elasticResponse.Acknowledged {
		return errors.New("0","couldn't delete the index")
	}

	return nil
}
