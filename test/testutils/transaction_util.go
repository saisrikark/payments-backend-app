package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payments-backend-app/pkg/server"
)

func (ta *TestApp) CallCreateTransaction(req *server.CreateTransactionRequest) (int, *server.CreateTransactionResponse, error) {
	var err error
	url := ta.baseUrl + "/transactions"
	body := &bytes.Buffer{}

	if req != nil {
		ba, err := json.Marshal(*req)
		if err != nil {
			return 0, nil, fmt.Errorf("unable to marshal [%s]", err)
		}
		body = bytes.NewBuffer(ba)
	}

	httpresp, err := http.Post(url, "application/json", body)
	if err != nil {
		return 0, nil, err
	}

	status := httpresp.StatusCode
	if status != http.StatusCreated {
		return status, nil, nil
	}

	ba, err := io.ReadAll(httpresp.Body)
	if err != nil {
		return status, nil, fmt.Errorf("unable to read response body [%s]", err.Error())
	}

	resp := server.CreateTransactionResponse{}
	if err := json.Unmarshal(ba, &resp); err != nil {
		return status, nil, fmt.Errorf("unable to unmarshal response [%s]", err.Error())
	}

	return status, &resp, nil
}
