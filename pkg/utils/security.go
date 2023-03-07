package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	securityTagKey = "security"
)

type securityTag struct {
	Option  bool
	Scheme  bool
	Name    string
	Type    string
	SubType string
}

type SecurityClient struct {
	client      HTTPClient
	headers     map[string]string
	queryParams map[string]string
}

func newSecurityClient(client HTTPClient) *SecurityClient {
	return &SecurityClient{
		client:      client,
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
	}
}

func (c *SecurityClient) Do(req *http.Request) (*http.Response, error) {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	queryParams := req.URL.Query()

	for k, v := range c.queryParams {
		queryParams.Set(k, v)
	}

	req.URL.RawQuery = queryParams.Encode()

	return c.client.Do(req)
}

func ConfigureSecurityClient(c HTTPClient, security interface{}) *SecurityClient {
	client := parseSecurityStruct(c, security)
	if client != nil {
		return client
	}

	return newSecurityClient(c)
}

func parseSecurityStruct(c HTTPClient, security interface{}) *SecurityClient {
	securityStructType := reflect.TypeOf(security)
	securityValType := reflect.ValueOf(security)

	if securityStructType.Kind() == reflect.Ptr {
		if securityValType.IsNil() {
			return nil
		}

		securityStructType = securityStructType.Elem()
		securityValType = securityValType.Elem()
	}

	client := newSecurityClient(c)

	for i := 0; i < securityStructType.NumField(); i++ {
		fieldType := securityStructType.Field(i)
		valType := securityValType.Field(i)

		kind := valType.Kind()

		if fieldType.Type.Kind() == reflect.Pointer {
			if valType.IsNil() {
				continue
			}

			kind = valType.Elem().Kind()
		}

		secTag := parseSecurityTag(fieldType)
		if secTag != nil {
			if secTag.Option {
				return parseSecurityOption(c, valType.Interface())
			} else if secTag.Scheme {
				// Special case for basic auth which could be a flattened struct
				if secTag.SubType == "basic" && kind != reflect.Struct {
					parseSecurityScheme(client, secTag, security)
					return client
				} else {
					parseSecurityScheme(client, secTag, valType.Interface())
				}
			}
		}
	}

	return client
}

func parseSecurityOption(c HTTPClient, option interface{}) *SecurityClient {
	optionStructType := reflect.TypeOf(option)
	optionValType := reflect.ValueOf(option)

	if optionStructType.Kind() == reflect.Ptr {
		if optionValType.IsNil() {
			return nil
		}

		optionStructType = optionStructType.Elem()
		optionValType = optionValType.Elem()
	}

	client := newSecurityClient(c)

	for i := 0; i < optionStructType.NumField(); i++ {
		fieldType := optionStructType.Field(i)
		valType := optionValType.Field(i)

		secTag := parseSecurityTag(fieldType)
		if secTag != nil && secTag.Scheme {
			parseSecurityScheme(client, secTag, valType.Interface())
		}
	}

	return client
}

func parseSecurityScheme(client *SecurityClient, schemeTag *securityTag, scheme interface{}) {
	schemeType := reflect.TypeOf(scheme)
	schemeVal := reflect.ValueOf(scheme)

	if schemeType.Kind() == reflect.Ptr {
		if schemeVal.IsNil() {
			return
		}

		schemeType = schemeType.Elem()
		schemeVal = schemeVal.Elem()
	}

	if schemeType.Kind() == reflect.Struct {
		if schemeTag.Type == "http" && schemeTag.SubType == "basic" {
			parseBasicAuthScheme(client, schemeVal.Interface())
			return
		}

		for i := 0; i < schemeType.NumField(); i++ {
			fieldType := schemeType.Field(i)
			valType := schemeVal.Field(i)

			if fieldType.Type.Kind() == reflect.Ptr {
				if valType.IsNil() {
					continue
				}

				valType = valType.Elem()
			}

			secTag := parseSecurityTag(fieldType)
			if secTag == nil || secTag.Name == "" {
				return
			}

			parseSecuritySchemeValue(client, schemeTag, secTag, valType.Interface())
		}
	} else {
		parseSecuritySchemeValue(client, schemeTag, schemeTag, schemeVal.Interface())
	}
}

func parseSecuritySchemeValue(client *SecurityClient, schemeTag *securityTag, secTag *securityTag, val interface{}) {
	switch schemeTag.Type {
	case "apiKey":
		switch schemeTag.SubType {
		case "header":
			client.headers[secTag.Name] = valToString(val)
		case "query":
			client.queryParams[secTag.Name] = valToString(val)
		case "cookie":
			client.headers["Cookie"] = fmt.Sprintf("%s=%s", secTag.Name, valToString(val))
		default:
			panic("not supported")
		}
	case "openIdConnect":
		client.headers[secTag.Name] = valToString(val)
	case "oauth2":
		client.headers[secTag.Name] = valToString(val)
	case "http":
		switch schemeTag.SubType {
		case "bearer":
			client.headers[secTag.Name] = valToString(val)
		default:
			panic("not supported")
		}
	default:
		panic("not supported")
	}
}

func parseBasicAuthScheme(client *SecurityClient, scheme interface{}) {
	schemeStructType := reflect.TypeOf(scheme)
	schemeValType := reflect.ValueOf(scheme)

	var username, password string

	for i := 0; i < schemeStructType.NumField(); i++ {
		fieldType := schemeStructType.Field(i)
		valType := schemeValType.Field(i)

		secTag := parseSecurityTag(fieldType)
		if secTag == nil || secTag.Name == "" {
			continue
		}

		switch secTag.Name {
		case "username":
			username = valType.String()
		case "password":
			password = valType.String()
		}
	}

	client.headers["Authorization"] = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))
}

func parseSecurityTag(field reflect.StructField) *securityTag {
	tag := field.Tag.Get(securityTagKey)
	if tag == "" {
		return nil
	}

	option := false
	scheme := false
	name := ""
	securityType := ""
	securitySubType := ""

	options := strings.Split(tag, ",")
	for _, optionConf := range options {
		parts := strings.Split(optionConf, "=")
		if len(parts) < 1 || len(parts) > 2 {
			continue
		}

		switch parts[0] {
		case "name":
			name = parts[1]
		case "type":
			securityType = parts[1]
		case "subtype":
			securitySubType = parts[1]
		case "option":
			option = true
		case "scheme":
			scheme = true
		}
	}

	// TODO: validate tag?

	return &securityTag{
		Option:  option,
		Scheme:  scheme,
		Name:    name,
		Type:    securityType,
		SubType: securitySubType,
	}
}