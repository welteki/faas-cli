// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexellis/faas/gateway/requests"
)

// DeleteFunction delete a function from the FaaS server
func DeleteFunction(gateway string, functionName string) error {
	delReq := requests.DeleteFunctionRequest{FunctionName: functionName}
	reqBytes, marshalErr := json.Marshal(&delReq)
	if marshalErr != nil {
		return marshalErr
	}

	reader := bytes.NewReader(reqBytes)

	c := http.Client{}
	req, _ := http.NewRequest("DELETE", gateway+"/system/functions", reader)
	req.Header.Set("Content-Type", "application/json")
	delRes, delErr := c.Do(req)
	if delErr != nil {
		return fmt.Errorf("error removing existing function: %s, gateway=%s, functionName=%s", delErr.Error(), gateway, functionName)
	}

	if delRes.Body != nil {
		defer delRes.Body.Close()
	}

	switch delRes.StatusCode {
	case 200, 201, 202:
		fmt.Println("Removing old service.")
	case 404:
		fmt.Println("No existing service to remove")
	default:
		bytesOut, err := ioutil.ReadAll(delRes.Body)
		if err == nil {
			return fmt.Errorf("server returned unexpected status code: %d.\n Response: %s", delRes.StatusCode, string(bytesOut))
		}
	}
	return nil
}
