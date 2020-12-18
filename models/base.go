package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

var cache map[string]map[string]map[string]*Base
var clients map[string]*http.Client
var SkipTLS bool = false

func init() {
	cache = map[string]map[string]map[string]*Base{}
	clients = map[string]*http.Client{}
}

type restConfig struct {
	objectType    string
	restKeyField  string
	tfIDField     string
	restInterface string
}

var RestConfigs = map[string]restConfig{
	"backup_restore": {
		objectType:    "backup_restore",
		restKeyField:  "_key",
		tfIDField:     "_key",
		restInterface: "backup_restore_interface",
	},
	"base_service_template": {
		objectType:    "base_service_template",
		restKeyField:  "_key",
		tfIDField:     "title",
		restInterface: "itoa_interface",
	},
	"correlation_search": {
		objectType:    "correlation_search",
		restKeyField:  "name",
		tfIDField:     "name",
		restInterface: "event_management_interface",
	},
	"entity": {
		objectType:    "entity",
		restKeyField:  "_key",
		tfIDField:     "_key",
		restInterface: "itoa_interface",
	},
	"glass_table": {
		objectType:    "glass_table",
		restKeyField:  "_key",
		tfIDField:     "_key",
		restInterface: "itoa_interface",
	},
	"kpi_base_search": {
		objectType:    "kpi_base_search",
		restKeyField:  "_key",
		tfIDField:     "title",
		restInterface: "itoa_interface",
	},
	"kpi_template": {
		objectType:    "kpi_template",
		restKeyField:  "_key",
		tfIDField:     "title",
		restInterface: "itoa_interface",
	},
	"kpi_threshold_template": {
		objectType:    "kpi_threshold_template",
		restKeyField:  "_key",
		tfIDField:     "title",
		restInterface: "itoa_interface",
	},
	"service": {
		objectType:    "service",
		restKeyField:  "_key",
		tfIDField:     "title",
		restInterface: "itoa_interface",
	},
}

type Base struct {
	restConfig
	// key used to collect this resource via the REST API
	RESTKey string
	// Terraform Identifier
	TFID    string
	RawJson json.RawMessage
	fields  []string
}

func NewBase(key, id, objectType string) *Base {
	if _, ok := RestConfigs[objectType]; !ok {
		panic(fmt.Sprintf("invalid objectype %s!", objectType))
	}
	b := &Base{
		restConfig: RestConfigs[objectType],
		RESTKey:    key,
		TFID:       id,
	}
	return b
}

// func (b *Base) GetValue(field string, out interface{}) error {
// 	if v, ok := b.RawJson[field]; ok {
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
// 	b.RawJson[field] = raw
// 	return nil
// }

func (b *Base) Clone() *Base {
	b_ := &Base{
		restConfig: b.restConfig,
		RawJson:    b.RawJson,
	}
	return b_
}

func (b *Base) NewClient(host string) *http.Client {
	if c, ok := clients[host]; ok {
		return c
	}
	tr := (http.DefaultTransport.(*http.Transport)).Clone()
	if host == "localhost" || SkipTLS {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	c := &http.Client{Transport: tr}
	return c
}

func (b *Base) Create(user, password, host string, port int) (*Base, error) {
	reqBody, err := json.Marshal(b.RawJson)
	if err != nil {
		return nil, err
	}
	client := b.NewClient(host)
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s", host, port, b.restInterface, b.objectType)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	if Verbose {
		err = logRequest(req)
		if err != nil {
			return nil, err
		}
	}
	req.SetBasicAuth(user, password)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if Verbose {
		err = logResponse(resp)
		if err != nil {
			return nil, err
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("create error: %v\n", resp.Status)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("create error: %v \n%s\n", resp.Status, body)
	}
	var r map[string]string
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	b.RESTKey = r[b.restConfig.restKeyField]
	b.storeCache()
	return b, nil
}

func (b *Base) Read(user, password, host string, port int) (*Base, error) {
	client := b.NewClient(host)
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s/%[5]s", host, port, b.restInterface, b.objectType, b.RESTKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, password)
	if Verbose {
		err = logRequest(req)
		if err != nil {
			return nil, err
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if Verbose {
		err = logResponse(resp)
		if err != nil {
			return nil, err
		}
	}
	if resp.StatusCode != 200 {
		return nil, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var raw json.RawMessage
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}
	base := b.Clone()
	err = base.Populate(raw)
	if err != nil {
		return nil, err
	}
	base.storeCache()
	return base, nil
}

func (b *Base) Update(user, password, host string, port int) error {
	reqBody, err := json.Marshal(b.RawJson)
	if err != nil {
		return err
	}

	client := b.NewClient(host)
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s/%[5]s", host, port, b.restInterface, b.objectType, b.RESTKey)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	if Verbose {
		err = logRequest(req)
		if err != nil {
			return err
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if Verbose {
		err = logResponse(resp)
		if err != nil {
			return err
		}
	}
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
	b.deleteCache()
	client := b.NewClient(host)
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s/%[5]s", host, port, b.restInterface, b.objectType, b.RESTKey)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	if Verbose {
		err = logRequest(req)
		if err != nil {
			return err
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if Verbose {
		err = logResponse(resp)
		if err != nil {
			return err
		}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("delete error: %v\n", resp.Status)
		}

		return fmt.Errorf("delete error: %v \n%s\n", resp.Status, body)
	}
	return nil
}

func (b *Base) storeCache() {
	if b.TFID == "" {
		return
	}
	if _, ok := cache[b.restInterface]; !ok {
		cache[b.restInterface] = map[string]map[string]*Base{}
	}
	if _, ok := cache[b.restInterface][b.objectType]; !ok {
		cache[b.restInterface][b.objectType] = map[string]*Base{}
	}
	cache[b.restInterface][b.objectType][b.TFID] = b
}

func (b *Base) getCache() *Base {
	if b.TFID == "" {
		return nil
	}
	if _, ok := cache[b.restInterface]; !ok {
		return nil
	}
	if _, ok := cache[b.restInterface][b.objectType]; !ok {
		return nil
	}
	if b_, ok := cache[b.restInterface][b.objectType][b.TFID]; ok {
		return b_
	} else {
		return nil
	}
}

func (b *Base) deleteCache() {
	if b.TFID == "" {
		return
	}
	if _, ok := cache[b.restInterface]; !ok {
		return
	}
	if _, ok := cache[b.restInterface][b.objectType]; !ok {
		return
	}
	delete(cache[b.restInterface][b.objectType], b.TFID)
}

func (b *Base) Find(user, password, host string, port int) (*Base, error) {
	b_ := b.getCache()
	if b_ != nil {
		return b_, nil
	}
	_, err := b.Dump(user, password, host, port)
	if err != nil {
		return nil, err
	}
	b_ = b.getCache()
	if b_ != nil {
		return b_, nil
	}
	return b_, nil
}

func (b *Base) Dump(user, password, host string, port int) ([]*Base, error) {
	client := b.NewClient(host)
	url := fmt.Sprintf("https://%[1]s:%[2]d/servicesNS/nobody/SA-ITOA/%[3]s/%[4]s", host, port, b.restInterface, b.objectType)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, password)
	if Verbose {
		err = logRequest(req)
		if err != nil {
			return nil, err
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if Verbose {
		err = logResponse(resp)
		if err != nil {
			return nil, err
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var raw []json.RawMessage
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}
	res := []*Base{}
	for _, r := range raw {
		b_ := b.Clone()
		err = b_.Populate(r)
		if err != nil {
			return nil, err
		}
		b_.storeCache()
		res = append(res, b_)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].RESTKey < res[j].RESTKey
	})
	return res, err
}

func (b *Base) Populate(raw []byte) error {
	err := json.Unmarshal(raw, &b.RawJson)
	if err != nil {
		return err
	}
	var fieldsMap map[string]*json.RawMessage
	err = json.Unmarshal(raw, &fieldsMap)
	if err != nil {
		return err
	}
	key := b.restKeyField
	if _, ok := fieldsMap[key]; !ok {
		return fmt.Errorf("missing %s RESTKey field for %s", key, b.objectType)
	}
	keyBytes, err := fieldsMap[key].MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(keyBytes, &b.RESTKey)
	if err != nil {
		return err
	}
	id := b.tfIDField
	if _, ok := fieldsMap[id]; !ok {
		return fmt.Errorf("missing %s TFID field for %s", id, b.objectType)
	}
	idBytes, err := fieldsMap[id].MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(idBytes, &b.TFID)
	if err != nil {
		return err
	}
	b.fields = []string{}
	for field := range fieldsMap {
		b.fields = append(b.fields, field)
	}
	sort.Strings(b.fields)
	return nil
}

func (b *Base) AuditLog(items []*Base, auditList []string) error {
	err := os.MkdirAll("dump", os.ModePerm)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("dump/%s.yaml", b.objectType)
	objects := []interface{}{}
	auditMap := map[string]bool{}
	for _, f := range auditList {
		auditMap[f] = true
	}
	for _, item := range items {
		by, err := json.Marshal(item.RawJson)
		if err != nil {
			return err
		}
		var auditBody map[string]*json.RawMessage
		err = json.Unmarshal(by, &auditBody)
		if err != nil {
			return err
		}
		if len(auditMap) > 0 {
			for f := range auditBody {
				if _, ok := auditMap[f]; !ok {
					delete(auditBody, f)
				}
			}
		}
		by, err = json.Marshal(auditBody)
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

func (b *Base) AuditFields(items []*Base) error {
	err := os.MkdirAll("fields", os.ModePerm)
	if err != nil {
		return err
	}
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
