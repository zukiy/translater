package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"translator/model"
	"translator/providers/storage"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	v1            = "/v1"
	root          = "/"
	translatePath = root + "translate"
	savePath      = root + "save"

	defaultPerPage = 50
)

type (
	// WordStorage ...
	WordStorage interface {
		List(params model.ParamsList) ([]model.Word, error)
		Save(lang string, words []string) ([]int64, error)
		SaveRelations(source model.Source, translates model.Translates) error
	}
)

// Server http
type Server struct {
	exitChan    chan struct{}
	port        int
	storage     storage.IStorage
	wordStorage WordStorage
}

// New create and return new HTTP server
func New(port int, storage storage.IStorage, wordStorage WordStorage) *Server {
	return &Server{
		exitChan:    make(chan struct{}, 1),
		port:        port,
		storage:     storage,
		wordStorage: wordStorage,
	}
}

func (s *Server) Serve() {
	s.run()
	<-s.exitChan
}

func (s *Server) Stop() {
	close(s.exitChan)
}

func (s *Server) run() {

	router := mux.NewRouter().StrictSlash(false)

	v1 := router.PathPrefix(v1).Subrouter()
	v1.HandleFunc(root, s.list).Methods(http.MethodPost)
	v1.HandleFunc(translatePath, s.translate).Methods(http.MethodPost)
	v1.HandleFunc(savePath, s.save).Methods(http.MethodPost)
	http.Handle(root, router)

	go func() {
		addr := fmt.Sprintf(":%d", s.port)

		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}()
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	var request ListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Errorf("decode list request error: %s", err.Error())
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	list, err := s.wordStorage.List(model.ParamsList{
		Page:    request.Page,
		PerPage: request.PerPage,
	})

	if err != nil {
		log.Errorf("fetch list return error: %s", err.Error())
		writeError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	writeSuccess(w, list)
}

func (s *Server) translate(w http.ResponseWriter, r *http.Request) {
	var request TranslateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Errorf("decode translate request error: %s", err.Error())
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	// yandex
	result, err := s.storage.YandexTranslator().Translate(request.Word, request.Lang)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "external error")
	}

	response := &TranslateResponse{
		Results: result,
	}

	writeSuccess(w, response)
}

func (s *Server) save(w http.ResponseWriter, r *http.Request) {

	var request SaveRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Errorf("decode save request error: %s", err.Error())
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	// todo: transaction
	enLids, err := s.wordStorage.Save(request.Word.Lang, []string{request.Word.Text})
	if err != nil {
		log.Errorf("save word %s [%s] error: %s",
			strings.Join([]string{request.Word.Text}, ","), request.Word.Lang, err.Error())

		writeError(w, http.StatusInternalServerError, "save error "+request.Word.Text)
		return
	}

	ruLids, err := s.wordStorage.Save(request.Translate.Lang, []string{request.Translate.Text})
	if err != nil {
		log.Errorf("save word %s [%s] error: %s",
			strings.Join([]string{request.Translate.Text}, ","), request.Translate.Lang, err.Error())

		writeError(w, http.StatusInternalServerError, "save error "+request.Translate.Text)
		return
	}

	err = s.wordStorage.SaveRelations(
		model.Source{Lang: request.Word.Lang, WordID: enLids[0]},
		model.Translates{Lang: request.Translate.Lang, WordsIDs: ruLids},
	)

	if err != nil {
		log.Errorf("relations save error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeSuccess(w, "OK")
}

func writeSuccess(w http.ResponseWriter, resp interface{}) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf("can't write success response: %s", err.Error())
	}
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	if _, err := w.Write([]byte(message)); err != nil {
		log.Errorf("can't write error response: %s", err.Error())
	}
}
