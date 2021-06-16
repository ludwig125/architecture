package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TODO: 不要そう
func TestController(t *testing.T) {
	service := &mockService{}
	service.mockGetAllFunc = func() ([]Actor, error) {
		return []Actor{
			{ID: 1, Name: "Watson", Age: 24},
			{ID: 2, Name: "Depp", Age: 54},
			{ID: 3, Name: "Portman", Age: 32},
		}, nil
	}

	router := NewServer(Config{}, service)

	req := httptest.NewRequest("GET", "/getall", nil)
	resp := httptest.NewRecorder()

	router.GetALlHandler(resp, req)

	t.Log(resp.Code)
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(res))

	// h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Hello, client")
	// })
	mux := http.NewServeMux()
	mux.HandleFunc("/getall", router.GetALlHandler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	cli := &http.Client{}
	request, err := http.NewRequestWithContext(context.TODO(), "GET", ts.URL+"/getall", strings.NewReader(""))
	if err != nil {
		t.Errorf("NewRequest failed: %v", err)
	}

	// Act
	response, err := cli.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	// Assertion
	got, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("got", string(got))
}

func TestGetAllHandler(t *testing.T) {
	service := &mockService{}

	tests := map[string]struct {
		mockGetAllFunc func() ([]Actor, error)
		wantCode       int
		wantErr        string
		wants          []Actor
	}{
		"getAll": {
			mockGetAllFunc: func() ([]Actor, error) {
				return []Actor{
					{ID: 1, Name: "Watson", Age: 24},
					{ID: 2, Name: "Depp", Age: 54},
					{ID: 3, Name: "Portman", Age: 32},
				}, nil
			},
			wantCode: 200,
			wants: []Actor{
				{ID: 1, Name: "Watson", Age: 24},
				{ID: 2, Name: "Depp", Age: 54},
				{ID: 3, Name: "Portman", Age: 32},
			},
		},
		"error_in_GetAll": {
			mockGetAllFunc: func() ([]Actor, error) {
				return nil, errors.New("GetAll error")
			},
			wantCode: 500,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/getall", nil)
			resp := httptest.NewRecorder()

			service.mockGetAllFunc = tc.mockGetAllFunc
			router := NewServer(Config{}, service)
			router.GetALlHandler(resp, req)

			res, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if resp.Code != tc.wantCode {
				t.Errorf("gotCode: %d, wantCode: %d", resp.Code, tc.wantCode)
			}
			if resp.Code != http.StatusOK {
				if string(res) != tc.wantErr {
					t.Errorf("gotErr: %s, wantErr: %s", string(res), tc.wantErr)
				}
				return
			}

			var actors []Actor
			if err := json.Unmarshal(res, &actors); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.wants, actors); diff != "" {
				t.Errorf("got: %v, wants: %v\ndiff: %v", actors, tc.wants, diff)
			}
		})
	}
}

func TestFindHandler(t *testing.T) {
	service := &mockService{}

	tests := map[string]struct {
		params       string
		mockFindFunc func() ([]Actor, error)
		wantCode     int
		wantErr      string
		wants        []Actor
	}{
		"findByID": {
			params: "id=1",
			mockFindFunc: func() ([]Actor, error) {
				return []Actor{
					{ID: 1, Name: "Watson", Age: 24},
				}, nil
			},
			wantCode: 200,
			wants: []Actor{
				{ID: 1, Name: "Watson", Age: 24},
			},
		},
		"findByName": {
			params: "name=Depp",
			mockFindFunc: func() ([]Actor, error) {
				return []Actor{
					{ID: 2, Name: "Depp", Age: 54},
				}, nil
			},
			wantCode: 200,
			wants: []Actor{
				{ID: 2, Name: "Depp", Age: 54},
			},
		},
		"findByAge": {
			params: "age=32",
			mockFindFunc: func() ([]Actor, error) {
				return []Actor{
					{ID: 3, Name: "Portman", Age: 32},
				}, nil
			},
			wantCode: 200,
			wants: []Actor{
				{ID: 3, Name: "Portman", Age: 32},
			},
		},
		"invalid_condition": {
			params:   "",
			wantCode: 500,
			wantErr:  "failed to NewRequestCond: invalid request. id: '', name: '', age: ''\n",
		},
		"error_in_Find": {
			params: "id=1",
			mockFindFunc: func() ([]Actor, error) {
				return nil, errors.New("Find error")
			},
			wantCode: 500,
			wantErr:  "failed to Find: Find error\n",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/find?"+tc.params, nil)
			resp := httptest.NewRecorder()

			service.mockFindFunc = tc.mockFindFunc
			router := NewServer(Config{}, service)
			router.FindHandler(resp, req)

			res, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if resp.Code != tc.wantCode {
				t.Errorf("gotCode: %d, wantCode: %d", resp.Code, tc.wantCode)
			}
			if resp.Code != http.StatusOK {
				if string(res) != tc.wantErr {
					t.Errorf("gotErr: %s, wantErr: %s", string(res), tc.wantErr)
				}
				return
			}

			var actors []Actor
			if err := json.Unmarshal(res, &actors); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.wants, actors); diff != "" {
				t.Errorf("got: %v, wants: %v\ndiff: %v", actors, tc.wants, diff)
			}
		})
	}
}

func TestUpdateHandler(t *testing.T) {
	service := &mockService{}

	tests := map[string]struct {
		method         string
		contentType    string
		req            []byte
		mockUpdateFunc func() error
		wantCode       int
		wantErr        string
	}{
		"update_ok": {
			method:      "POST",
			contentType: "application/json",
			req: makeJSONBytes(t, &Actor{
				Name: "Watson",
				Age:  24,
			}),
			mockUpdateFunc: func() error {
				return nil
			},
			wantCode: 200,
		},
		"is_not_POST": {
			method:   "GET",
			wantCode: 400,
			wantErr:  "invalid method GET, request must be POST\n",
		},
		"is_not_json_request": {
			method:      "POST",
			contentType: "text/plain",
			wantCode:    400,
			wantErr:     "POST request must be JSON\n",
		},
		"bad_json": { // nameがstringではなくint
			method:      "POST",
			contentType: "application/json",
			req:         []byte(`{"name":1}`),
			wantCode:    400,
			wantErr: `Request body contains an invalid value for the "name" field (at position 9)
`,
		},
		"unexpected_EOF": {
			method:      "POST",
			contentType: "application/json",
			req:         []byte(""),
			wantCode:    400,
			wantErr:     "Request body must not be empty\n",
		},
		"decode_error": { // JSONが正しい形式ではない("で閉じられていない")
			method:      "POST",
			contentType: "application/json",
			req:         []byte(`{"invalid":"type json}`),
			wantCode:    500,
			wantErr:     "Internal Server Error\n",
		},
		"cannot_convert_Actor": {
			method:      "POST",
			contentType: "application/json",
			req:         []byte(`{"name2":"Jack"}`),
			wantCode:    500,
			wantErr:     "Internal Server Error\n",
		},
		"update_error": {
			method:      "POST",
			contentType: "application/json",
			req: makeJSONBytes(t, &Actor{
				Name: "Watson",
				Age:  24,
			}),
			mockUpdateFunc: func() error {
				return errors.New("error occurred in Update")
			},
			wantCode: 500,
			wantErr:  "Internal Server Error\n",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/update", bytes.NewBuffer(tc.req))
			resp := httptest.NewRecorder()

			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			service.mockUpdateFunc = tc.mockUpdateFunc
			router := NewServer(Config{}, service)
			router.UpdateHandler(resp, req)

			res, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if resp.Code != tc.wantCode {
				t.Errorf("gotCode: %d, wantCode: %d", resp.Code, tc.wantCode)
			}
			if resp.Code != http.StatusOK {
				if string(res) != tc.wantErr {
					t.Errorf("gotErr: %s, wantErr: %s", string(res), tc.wantErr)
				}
				return
			}
		})
	}
}

func makeJSONBytes(t *testing.T, a *Actor) []byte {
	s, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("failed to json Marshal: %v", err)
	}
	return s
}
