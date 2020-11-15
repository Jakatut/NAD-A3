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

	endpoints := []map[string]interface{}{}
	endpointToListLocation := map[string]int{}

	if len(changes) != 0 {
		updates := map[string]interface{}{}

		for key := range changes {
			formValueKeyDetails := strings.Split(key, "-")
			if len(formValueKeyDetails) == 3 {
				return errors.New("incorrect permission change value")
			}
			changedEndpoint := formValueKeyDetails[0]
			changedRole := formValueKeyDetails[1]
			changedAction := formValueKeyDetails[2]

			if index, ok := endpointToListLocation[changedEndpoint]; ok {
				endpoints[index]["route"] = changedEndpoint
				if _, ok := endpoints[index]["operations"]; ok {
					endpoints[index]["operations"] = append(endpoints[index]["operations"], changedRole)
				} else {

				}
			} else {
				endpoints = append(endpoints, updates)
				endpointToListLocation[changedEndpoint] = len(endpoints) + 1
			}
		}
	}

	return nil
}
