// Package api provides low-level functions to access the RailData API.
// You shouldn't use these directly; use the [raildata.Client] interface instead.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"

	"github.com/jtarrio/raildata/errors"
)

// MethodDefinition defines a RailData API method with the types of its request and response objects.
type MethodDefinition[I any, O any] struct {
	Name   string
	Input  reflect.Type
	Output reflect.Type
}

func method[I any, O any](name string) MethodDefinition[I, O] {
	return MethodDefinition[I, O]{Name: name, Input: reflect.TypeFor[I](), Output: reflect.TypeFor[O]()}
}

// WantsToken returns true if this method's request object contains a "Token" field.
func (m MethodDefinition[I, O]) WantsToken() bool {
	_, found := m.Input.FieldByName("Token")
	return found
}

// SetToken sets the value of the request object's "Token" field, if it exists.
func (m MethodDefinition[I, O]) SetToken(input *I, token string) {
	inputType := reflect.TypeOf(input).Elem()
	if tokenField, found := inputType.FieldByName("Token"); found {
		v := reflect.ValueOf(input).Elem()
		tv := v.FieldByIndex(tokenField.Index)
		tv.Set(reflect.ValueOf(token))
	}
}

// Request sends an HTTP request via the given client, to the server at the provided API base URL.
func (m MethodDefinition[I, O]) Request(ctx context.Context, client *http.Client, apiBase url.URL, input *I) (*O, error) {
	req, err := m.createRequest(ctx, apiBase, input)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error issuing request for method '%s': %w", m.Name, err)
	}
	return m.ParseResponse(resp)
}

func (m MethodDefinition[I, O]) createRequest(ctx context.Context, apiBase url.URL, input *I) (*http.Request, error) {
	fields, err := objToMap(input)
	if err != nil {
		return nil, fmt.Errorf("could not convert a request object for method '%s': %w", m.Name, err)
	}

	body := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&body)
	for k, v := range fields {
		fw, err := bodyWriter.CreateFormField(k)
		if err != nil {
			return nil, fmt.Errorf("could not create a request body for method '%s': %w", m.Name, err)
		}
		fw.Write([]byte(v))
	}
	bodyWriter.Close()

	url := apiBase.JoinPath(m.Name)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), &body)
	if err != nil {
		return nil, fmt.Errorf("error creating request for method '%s': %w", m.Name, err)
	}
	request.Header.Add("accept", "text/plain")
	request.Header.Add("content-type", bodyWriter.FormDataContentType())
	return request, nil
}

func (m MethodDefinition[I, O]) ParseResponse(response *http.Response) (*O, error) {
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, m.parseErrorResponse(response)
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, errors.MissingCredentialsError
	}

	var output O
	err = json.Unmarshal(b, &output)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response for %s: %w", m.Name, err)
	}
	return &output, nil
}

func (m MethodDefinition[I, O]) parseErrorResponse(response *http.Response) error {
	var errResp struct {
		Message string `json:"errorMessage"`
	}
	ret := fmt.Errorf("received error status code for %s: %s", m.Name, response.Status)
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return ret
	}
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&errResp)
	if err != nil {
		return ret
	}
	if errResp.Message == "Invalid token." {
		return errors.InvalidTokenError
	}
	return errors.NewRailDataError(errResp.Message)
}

func objToMap(i any) (map[string]string, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var out map[string]string
	err = json.Unmarshal(b, &out)
	return out, err
}
