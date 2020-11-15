package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type operation struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Errors      []string `json:"errors,omitempty"`
}

type Endpoint struct {
	Route      string      `json:"route"`
	Operations []operation `json:"operations"`
}

type AccessControlModel struct {
	Endpoints []Endpoint `json:"endpoints"`
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

func (perms *AccessControlModel) WritePermissions() {

}

func (perms *AccessControlModel) GetChanges(changes map[string]string) {
	if len(changes) != 0 {
		for key, value := range changes {

		}
	}
}

func convertFormValueToPermissions(formValue string) {
	strings.split()
}
