package elastic

import (
	"encoding/json"

	"bytes"
	"fmt"
	"net/http"

	"github.com/joaosoft/errors"
)

type ExistsIndexService struct {
	client *Elastic
	index  string
	body   []byte
	method string
}

func NewExistsIndexService(e *Elastic) *ExistsIndexService {
	return &ExistsIndexService{
		client: e,
		method: http.MethodHead,
	}
}

func (e *ExistsIndexService) Index(index string) *ExistsIndexService {
	e.index = index
	return e
}

func (e *ExistsIndexService) Body(body interface{}) *ExistsIndexService {
	switch v := body.(type) {
	case []byte:
		e.body = v
	default:
		e.body, _ = json.Marshal(v)
	}
	return e
}

func (e *ExistsIndexService) Execute() (int, error) {

	// create data on elastic
	reader := bytes.NewReader(e.body)

	request, err := http.NewRequest(e.method, fmt.Sprintf("%s/%s", e.client.config.Endpoint, e.index), reader)
	if err != nil {
		return 0, err
	}

	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()

	if err != nil {
		return response.StatusCode, errors.New(err)
	}

	return response.StatusCode, nil
}
