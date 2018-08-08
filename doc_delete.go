package elastic

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/joaosoft/errors"
)

type DeleteResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

type DeleteHit struct {
	Found  bool   `json:"found"`
	Result string `json:"result"`
}

type DeleteService struct {
	client *Elastic
	index  string
	typ    string
	id     string
	method string
}

func NewDeleteService(client *Elastic) *DeleteService {
	return &DeleteService{
		client: client,
		method: http.MethodDelete,
	}
}

func (e *DeleteService) Index(index string) *DeleteService {
	e.index = index
	return e
}

func (e *DeleteService) Type(typ string) *DeleteService {
	e.typ = typ
	return e
}

func (e *DeleteService) Id(id string) *DeleteService {
	e.id = id
	return e
}

func (e *DeleteService) Execute() error {

	// delete data from elastic
	var query string
	if e.typ != "" {
		query += fmt.Sprintf("/%s", e.typ)
	}

	if e.id != "" {
		query += fmt.Sprintf("/%s", e.id)
	}

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s%s", e.client.config.Endpoint, e.index, query), nil)
	if err != nil {
		return errors.New("0", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("0", err)
	}
	defer response.Body.Close()

	// unmarshal data
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("0", err)
	}

	if e.id != "" {
		elasticResponse := DeleteHit{}
		if err = json.Unmarshal(body, &elasticResponse); err != nil {
			return errors.New("0", err)
		}

		if !elasticResponse.Found || elasticResponse.Result != "deleted" {
			return errors.New("0","couldn't delete the resource")
		}
	} else {
		elasticResponse := DeleteResponse{}
		if err = json.Unmarshal(body, &elasticResponse); err != nil {
			return errors.New("0", err)
		}

		if !elasticResponse.Acknowledged {
			return errors.New("0","couldn't delete the resource")
		}
	}

	return nil
}
