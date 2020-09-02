package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"gopkg.in/yaml.v3"
)

var cache map[string]map[string]map[string]*Base

func init() {
	cache = map[string]map[string]map[string]*Base{}
}

type Base struct {
	// key used to collect this resource via the REST API
	key          string
	RESTKeyField func() string
	// Terraform Identifier
	id            string
	TFIDField     func() string
	restInterface string
	objectType    string
	rawJson       map[string]*json.RawMessage
	auditList     []string
	fields        []string
}

func NewBase(key, id, restInterface, objectType string) *Base {
	b := &Base{
		key:           key,
		id:            id,
		restInterface: restInterface,
		objectType:    objectType,
		rawJson:       map[string]*json.RawMessage{},
	}
	b.RESTKeyField = func() string {
		return "_key"
	}
	b.TFIDField = func() string {
		return "name"
	}
	return b
}

// func (b *Base) GetValue(field string, out interface{}) error {
// 	if v, ok := b.rawJson[field]; ok {
// 		by, err := v.MarshalJSON()
// 		if err != nil {
// 			return err
// 		}
// 		return json.Unmarshal(by, out)
// 	} else {
// 		return fmt.Errorf("no value for field %s", field)
// 	}
// }

// func (b *Base) SetValue(field string, value interface{}) error {
// 	by, err := json.Marshal(value)
// 	if err != nil {
// 		return err
// 	}
// 	var raw *json.RawMessage
// 	err = json.Unmarshal(by, raw)
// 	if err != nil {
// 		return err
// 	}
// 	b.rawJson[field] = raw
// 	return nil
// }

func (b *Base) Clone() *Base {
	b_ := &Base{
		RESTKeyField:  b.RESTKeyField,
		TFIDField:     b.TFIDField,
		auditList:     b.auditList,
		restInterface: b.restInterface,
		objectType:    b.objectType,
		rawJson:       b.rawJson,
	}
	return b_
}

func (b *Base) Transport(host string) *http.Transport {
	tr := (http.DefaultTransport.(*http.Transport)).Clone()
	if host == "localhost" {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return tr
}

func (b *Base) Create(user, password, host string, port int) error {
	reqBody, err := json.Marshal(b.rawJson)
	if err != nil {
		return err
	}

	client := &http.Client{Transport: b.Transport(host)}
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s", host, port, b.restInterface, b.objectType)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("create error: %v\n", resp.Status)
		}

		return fmt.Errorf("create error: %v \n%s\n", resp.Status, body)
	}
	b.storeCache()
	return nil
}

func (b *Base) Read(user, password, host string, port int) error {
	client := &http.Client{Transport: b.Transport(host)}
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s/%[5]s", host, port, b.restInterface, b.objectType, b.key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &b.rawJson)
	if err != nil {
		return err
	}
	err = b.Populate()
	if err != nil {
		return err
	}
	b.storeCache()
	return nil
}

func (b *Base) Update(user, password, host string, port int) error {

	reqBody, err := json.Marshal(b.rawJson)
	if err != nil {
		return err
	}

	client := &http.Client{Transport: b.Transport(host)}
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s/%[5]s", host, port, b.restInterface, b.objectType, b.key)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("update error: %v\n", resp.Status)
		}

		return fmt.Errorf("update error: %v \n%s\n", resp.Status, body)
	}
	b.storeCache()
	return nil
}

func (b *Base) Delete(user, password, host string, port int) error {
	client := &http.Client{Transport: b.Transport(host)}
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s/%[5]s", host, port, b.restInterface, b.objectType, b.key)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("delete error: %v\n", resp.Status)
		}

		return fmt.Errorf("delete error: %v \n%s\n", resp.Status, body)
	}
	b.storeDelete()
	return nil
}

func (b *Base) storeCache() {
	if b.id == "" {
		return
	}
	if _, ok := cache[b.restInterface]; !ok {
		cache[b.restInterface] = map[string]map[string]*Base{}
	}
	if _, ok := cache[b.restInterface][b.objectType]; !ok {
		cache[b.restInterface][b.objectType] = map[string]*Base{}
	}
	cache[b.restInterface][b.objectType][b.id] = b
}

func (b *Base) getCache() *Base {
	if b.id == "" {
		return nil
	}
	if _, ok := cache[b.restInterface]; !ok {
		return nil
	}
	if _, ok := cache[b.restInterface][b.objectType]; !ok {
		return nil
	}
	if b_, ok := cache[b.restInterface][b.objectType][b.id]; ok {
		return b_
	} else {
		return nil
	}
}

func (b *Base) storeDelete() {
	if b.id == "" {
		return
	}
	if _, ok := cache[b.restInterface]; !ok {
		return
	}
	if _, ok := cache[b.restInterface][b.objectType]; !ok {
		return
	}
	delete(cache[b.restInterface][b.objectType], b.id)
}

func (b *Base) Find(user, password, host string, port int) (*Base, error) {
	if b.key != "" {
		err := b.Read(user, password, host, port)
		return b, err
	}
	b_ := b.getCache()
	if b_ != nil {
		return b_, nil
	}
	_, err := b.Dump(user, password, host, port)
	if err != nil {
		return nil, err
	}
	return b.getCache(), nil
}

func (b *Base) Dump(user, password, host string, port int) ([]*Base, error) {
	client := &http.Client{Transport: b.Transport(host)}
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s", host, port, b.restInterface, b.objectType)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var raw []map[string]*json.RawMessage
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}
	res := []*Base{}
	for _, r := range raw {
		b_ := b.Clone()
		b_.rawJson = r
		err = b_.Populate()
		if err != nil {
			return nil, err
		}
		b_.storeCache()
		res = append(res, b_)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].key < res[j].key
	})
	return res, err
}

func (b *Base) Populate() error {
	key := b.RESTKeyField()
	if _, ok := b.rawJson[key]; !ok {
		fmt.Println(b.rawJson)
		return fmt.Errorf("missing %s field for %s", key, b.objectType)
	}
	keyBytes, err := b.rawJson[key].MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(keyBytes, &b.key)
	if err != nil {
		return err
	}
	id := b.TFIDField()
	if _, ok := b.rawJson[id]; !ok {
		fmt.Println(b.rawJson)
		return fmt.Errorf("missing %s field for %s", id, b.objectType)
	}
	idBytes, err := b.rawJson[id].MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(idBytes, &b.id)
	if err != nil {
		return err
	}
	b.fields = []string{}
	for field := range b.rawJson {
		b.fields = append(b.fields, field)
	}
	sort.Strings(b.fields)
	return nil
}

func (b *Base) auditLog(items []*Base, auditList []string) error {
	filename := fmt.Sprintf("dump/%s.yaml", b.objectType)
	objects := []interface{}{}
	auditMap := map[string]bool{}
	for _, f := range auditList {
		auditMap[f] = true
	}
	for _, item := range items {
		rawJson := item.rawJson
		if len(auditMap) > 0 {
			for f := range rawJson {
				if _, ok := auditMap[f]; !ok {
					delete(rawJson, f)
				}
			}
		}
		by, err := json.Marshal(rawJson)
		if err != nil {
			return err
		}
		var i interface{}
		err = json.Unmarshal(by, &i)
		if err != nil {
			return err
		}
		objects = append(objects, i)
	}
	by, err := yaml.Marshal(objects)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, by, 0644)
}

func (b *Base) auditFields(items []*Base) error {
	filename := fmt.Sprintf("fields/%s.yaml", b.objectType)
	fieldsMap := map[string]bool{}
	for _, item := range items {
		for _, f := range item.fields {
			fieldsMap[f] = true
		}
	}
	fields := []string{}
	for field := range fieldsMap {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	by, err := yaml.Marshal(fields)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, by, 0644)
}
