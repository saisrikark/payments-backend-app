package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"payments-backend-app/pkg/server"
)

func CallCreateAccount(port int, req server.CreateAccountRequest) (int, error) {

	ba, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("unable to marshal [%s]", err)
	}

	url := fmt.Sprintf("http://localhost:%d%s", port, server.CreateAccountExtension)
	body := bytes.NewBuffer(ba)

	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return 0, err
	}

	status := resp.StatusCode

	return status, nil
}
