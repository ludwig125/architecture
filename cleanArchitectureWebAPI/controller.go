package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Server struct {
	Port    string
	Service ActorService
}

func NewServer(config Config, service ActorService) *Server {
	return &Server{
		Port:    config.Port,
		Service: service,
	}
}

func (s *Server) GetALlHandler(w http.ResponseWriter, r *http.Request) {
	ps, err := s.Service.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseActorByJSON(ps, w)
}

func (s *Server) FindHandler(w http.ResponseWriter, r *http.Request) {
	cond, err := NewRequestCond(r.URL.Query().Get("id"),
		r.URL.Query().Get("name"),
		r.URL.Query().Get("age"))
	if err != nil {
		e := fmt.Sprintf("failed to NewRequestCond: %v", err)
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	ps, err := s.Service.Find(*cond)
	if err != nil {
		e := fmt.Sprintf("failed to Find: %v", err)
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}
	responseActorByJSON(ps, w)
}

func (s *Server) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		e := fmt.Sprintf("invalid method %s, request must be POST", r.Method)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "POST request must be JSON", http.StatusBadRequest)
		return
	}

	var a Actor

	// request をdecode。失敗したらエラーを返す
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&a); err != nil {
		// 参考：https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
		default:
			// 内部のロジックをさらさないようにあえて詳細は返さない
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// 構造体のvalidation check
	if err := a.ValidateActor(); err != nil {
		// 内部のロジックをさらさないようにあえて詳細は返さない
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := s.Service.Update(a); err != nil {
		log.Println("failed to Update:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf, w.Writeはw.WriteHeaderやhttp.Errorの後に呼ばないと、
	// ステータスコード200を返してしまうので、最後に出力する
	// 例えば、fmt.Fprintfのあとにw.WriteHeader(http.StatusInternalServerError)をしても200を返すので注意
	fmt.Fprintf(w, "Actor name: '%s', age: '%d'", a.Name, a.Age)
}

func responseActorByJSON(ps []Actor, w http.ResponseWriter) {
	jsonData, err := json.Marshal(ps)
	if err != nil {
		log.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 省略可能
	fmt.Fprintln(w, string(jsonData))
}

func (s *Server) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/getall", s.GetALlHandler)
	mux.HandleFunc("/search", s.FindHandler)
	mux.HandleFunc("/update", s.UpdateHandler) //POST
	srv := &http.Server{
		Addr:    "localhost:" + s.Port,
		Handler: mux,
	}
	fmt.Println("starting http server on :", s.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln("Server closed with error:", err)
	}
}
