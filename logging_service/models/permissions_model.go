package models

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"sync"
)

type operation struct {
	Action  string   `json:"action"`
	Allowed []string `json:""`
}

type endpoint struct {
	Route      string      `json:"route"`
	Operations []operation `json:"operations"`
}

type AccessControlModel struct {
	Endpoints []endpoint `json:"endpoints"`
	Lock      sync.RWMutex
}

func (perms *AccessControlModel) GetPermissions() error {
	perms.Lock.RLock()
	defer perms.Lock.RUnlock()
	jsonFile, err := os.Open("config/permissions.json")
	defer jsonFile.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	data, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(data, perms)

	return nil
}

func (perms *AccessControlModel) WritePermissions() error {
	bytes, err := json.Marshal(perms)
	if err != nil {
		return err
	}

	perms.Lock.Lock()
	defer perms.Lock.Unlock()

	jsonFile, err := os.Open("config/permissions.json")
	if err != nil {
		return err
	}
	bufferedWriter := bufio.NewWriter(jsonFile)
	_, err = bufferedWriter.Write(bytes)

	return nil
}

func (perms *AccessControlModel) GetChanges(changes url.Values) error {

	var permissions AccessControlModel

	endpoints := []endpoint{}

	if len(changes) != 0 {
		for key := range changes {
			formValueKeyDetails := strings.Split(key, "-")
			if len(formValueKeyDetails) < 3 {
				return errors.New("incorrect permission change value")
			}
			changedRoute := formValueKeyDetails[0]
			changedRole := formValueKeyDetails[1]
			changedAction := formValueKeyDetails[2]

			if len(endpoints) > 0 {
				for _, _endpoint := range endpoints {
					if _endpoint.Route == changedRoute {
						for _, _operation := range _endpoint.Operations {
							if _operation.Action == changedAction {
								_operation.Allowed = append(_operation.Allowed, changedRole)
							}
						}
					}
				}
			} else {
				newOperation := operation{Action: changedAction, Allowed: []string{changedRole}}
				newEndpoint := endpoint{Route: changedRoute, Operations: []operation{newOperation}}
				permissions.Endpoints = append(permissions.Endpoints, newEndpoint)
			}
		}

		for _, _endpoint := range permissions.Endpoints {
			for endpointIndex, currentEndpoint := range perms.Endpoints {
				if _endpoint.Route == currentEndpoint.Route {
					for _, _operation := range _endpoint.Operations {
						for operationIndex, currentOperations := range currentEndpoint.Operations {
							if currentOperations.Action == _operation.Action {
								perms.Endpoints[endpointIndex].Operations[operationIndex].Allowed = _operation.Allowed
							}
						}
					}
				}
			}
		}

	}
	return nil
}
