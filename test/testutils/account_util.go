package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payments-backend-app/pkg/models"
	"payments-backend-app/pkg/server"
)

func (ta *TestApp) CallGetAccount(accountID int) (int, *models.Account, error) {
	var err error
	url := ta.baseUrl + fmt.Sprintf("/accounts/%d", accountID)

	httpresp, err := http.Get(url)
	if err != nil {
		return 0, nil, err
	}

	status := httpresp.StatusCode
	if status != http.StatusOK {
		return status, nil, nil
	}

	ba, err := io.ReadAll(httpresp.Body)
	if err != nil {
		return status, nil, fmt.Errorf("unable to read response body [%s]", err.Error())
	}

	resp := models.Account{}
	if err := json.Unmarshal(ba, &resp); err != nil {
		return status, nil, fmt.Errorf("unable to unmarshal response [%s]", err.Error())
	}

	return status, &resp, nil
}

func (ta *TestApp) CallCreateAccount(req *server.CreateAccountRequest) (int, *server.GetAccountResponse, error) {
	var err error
	url := ta.baseUrl + "/accounts"
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

	resp := server.GetAccountResponse{}
	if err := json.Unmarshal(ba, &resp); err != nil {
		return status, nil, fmt.Errorf("unable to unmarshal response [%s]", err.Error())
	}

	return status, &resp, nil
}

func (ta *TestApp) CallCreateAccountWithoutBody() (int, *models.Account, error) {
	var err error
	url := ta.baseUrl + "/accounts"

	httpresp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return 0, nil, err
	}

	status := httpresp.StatusCode

	ba, err := io.ReadAll(httpresp.Body)
	if err != nil {
		return status, nil, fmt.Errorf("unable to read response body [%s]", err.Error())
	}

	resp := models.Account{}
	if err := json.Unmarshal(ba, &resp); err != nil {
		return status, nil, fmt.Errorf("unable to unmarshal response [%s]", err.Error())
	}

	return status, &resp, nil
}
