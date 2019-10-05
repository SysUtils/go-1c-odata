package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) Get(entity IEntity, fields []string) error {
	entity.SetClient__(c)
	typename, err := getEntityName(entity)
	if err != nil {
		return err
	}
	uri := "/" + url.PathEscape(typename)
	uri += fmt.Sprintf("(%s)", keyToQuery(entity.Key__())) // Unique key
	uri += "?$format=json"
	if len(fields) > 0 {
		uri += fmt.Sprintf("&$select=%s", url.PathEscape(strings.Join(fields, ","))) // Fields
	}
	data, err := c.get(uri)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, entity)
}

func (c *Client) GetMany(entity interface{}, where Where) error {
	typename, err := getEntityName(entity)
	if err != nil {
		return err
	}
	uri := "/" + url.PathEscape(typename)
	uri += "?$format=json&"
	uri += where.Serialize()

	data, err := c.get(uri)
	if err != nil {
		return err
	}
	type ReturnObj struct {
		Value json.RawMessage `json:"value"`
	}
	outer := ReturnObj{}
	err = json.Unmarshal(data, &outer)
	if err != nil {
		return err
	}
	err = json.Unmarshal(outer.Value, &entity)
	if err != nil {
		return err
	}

	setClientToSlice(entity, c)
	return nil
}

func (c *Client) GetNavigaion(entity IEntity, navigation string, result interface{}) error {
	typename, err := getEntityName(entity)
	if err != nil {
		return err
	}
	uri := "/" + url.PathEscape(typename)
	uri += fmt.Sprintf("(%s)", keyToQuery(entity.Key__())) // Unique key
	uri += "/" + url.PathEscape(navigation)
	uri += "?$format=json&"

	body, err := c.get(uri)
	if err != nil {
		if err.Error() == "404 Not found\nBody:\n" {
			return nil
		}
	}

	return json.Unmarshal(body, result)
}

func (c *Client) ExecuteMethod(entity IEntity, method string, params interface{}, result interface{}) error {
	typename, err := getEntityName(entity)
	if err != nil {
		return err
	}

	uri := "/" + url.PathEscape(typename)
	uri += fmt.Sprintf("(%s)", url.PathEscape(keyToQuery(entity.Key__()))) // Unique key
	uri += "/" + method
	uri += "?$format=json&"
	uri += paramsToQuery(params)
	data, err := c.post(uri, nil)
	if err != nil {
		return err
	}

	if result != nil {
		err = json.Unmarshal(data, result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Update(entity IEntity) error {
	entity.SetClient__(c)
	typename, err := getEntityName(entity)
	if err != nil {
		return err
	}
	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	uri := "/" + url.PathEscape(typename)

	uri += fmt.Sprintf("(%s)", url.PathEscape(keyToQuery(entity.Key__()))) // Unique key
	uri += "?$format=json"
	data, err = c.patch(uri, data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, entity)
}

func (c *Client) Delete(entity IEntity) error {
	entity.SetClient__(c)
	typename, err := getEntityName(entity)
	if err != nil {
		return err
	}
	uri := "/" + url.PathEscape(typename)

	uri += fmt.Sprintf("(%s)", url.PathEscape(keyToQuery(entity.Key__()))) // Unique key
	uri += "?$format=json"

	return c.delete(uri)
}

func (c *Client) Create(entity IEntity) error {
	entity.SetClient__(c)
	typename, err := getEntityName(entity)
	if err != nil {
		return err
	}
	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	uri := "/" + url.PathEscape(typename) + "?$format=json"

	data, err = c.post(uri, data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, entity)
}
