package models

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"gopkg.in/yaml.v3"
)

type Crud interface {
	Create()
	Read()
	Update()
	Delete()
}

type Base struct {
	key           string
	restInterface string
	objectType    string
	rawJson       map[string]*json.RawMessage
	auditList     []string
	fields        []string
}

func NewBase(key, restInterface, objectType string, auditList []string) *Base {
	b := &Base{
		key:           key,
		restInterface: restInterface,
		objectType:    objectType,
		auditList:     auditList,
	}
	return b
}

func (b *Base) Key() string {
	return "_key"
}

func (b *Base) Clone() *Base {
	b_ := &Base{
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

func (b *Base) Read(user, password, host string, port int) error {
	fmt.Printf("objectype: %s\n", b.objectType)
	client := &http.Client{Transport: b.Transport(host)}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s/%[5]s", host, port, b.restInterface, b.objectType, b.key), nil)
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
	return b.Populate()
}

func (b *Base) Dump(user, password, host string, port int) ([]*Base, error) {
	client := &http.Client{Transport: b.Transport(host)}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s", host, port, b.restInterface, b.objectType), nil)
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
		res = append(res, b_)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].key < res[j].key
	})
	return res, err
}

func (b *Base) Populate() error {
	key := b.Key()
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
	b.fields = []string{}
	for field := range b.rawJson {
		b.fields = append(b.fields, field)
	}
	sort.Strings(b.fields)
	return nil
}

func (b *Base) auditLog(items []*Base) error {
	filename := fmt.Sprintf("dump/%s.yaml", b.objectType)
	objects := []interface{}{}
	auditMap := map[string]bool{}
	for _, f := range b.auditList {
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

// func dump(restInterface string, objectType string, objects interface{}) (fields []string, err error) {
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	client := &http.Client{}
// 	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://localhost:18089/servicesNS/nobody/SA-ITOA/%[1]s/%[2]s", restInterface, objectType), nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.SetBasicAuth(ITSI_REST_USER, ITSI_REST_PWD)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = json.Unmarshal(body, objects)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var fieldsDump = []map[string]*json.RawMessage{}
// 	err = json.Unmarshal(body, &fieldsDump)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fieldsMap := map[string]bool{}
// 	for _, fmap := range fieldsDump {
// 		for field := range fmap {
// 			fieldsMap[field] = true
// 		}
// 	}
// 	for field := range fieldsMap {
// 		fields = append(fields, field)
// 	}
// 	sort.Strings(fields)
// 	return fields, nil
// }

// func writeObjects(objectType string, objects interface{}) error {
// 	filename := fmt.Sprintf("dump/%s.yaml", objectType)
// 	b, err := yaml.Marshal(objects)
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(filename, b, 0644)
// }

// func writeFields(objectType string, fields []string) error {
// 	filename := fmt.Sprintf("fields/%s.yaml", objectType)
// 	b, err := yaml.Marshal(fields)
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(filename, b, 0644)
// }
