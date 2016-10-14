package web

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
	"io"
	"net/http"
	"strings"
	"time"
	"trie"
)

type Server struct {
	store map[string]*trie.Trie
}

var (
	log *logging.Logger = logging.MustGetLogger("suggest")
)

func NewServer() *Server {
	server := new(Server)
	server.store = make(map[string]*trie.Trie)
	router := mux.NewRouter()
	router.HandleFunc("/suggest/{collection}", server.suggestHandler).Methods("GET")
	router.HandleFunc("/index/{collection}", server.indexHandler).Methods("POST")
	http.Handle("/", router)
	return server
}

func parseData(line string) (string, string, error) {
	line = strings.TrimRight(line, "\n ")
	tokens := strings.Split(line, " ")
	if len(tokens) != 2 {
		tokens = strings.Split(line, "\t")
		if len(tokens) != 2 {
			return "", "", errors.New("Invalid data")
		}
	}
	return tokens[0], tokens[1], nil
}

func (s *Server) suggestHandler(w http.ResponseWriter, r *http.Request) {
	var results []string

	collection := mux.Vars(r)["collection"]

	w.Header().Set("Content-Type", "text/plain")
	r.ParseForm()
	query := r.Form.Get("query")

	if len(query) < 2 {
		results = make([]string, 0)
	} else {
		if store, found := s.store[collection]; found {
			results = store.GetAllWithPrefix(query)
		} else {
			http.Error(w, fmt.Sprintf("Collection <%s> not found", collection), 404)
			return
		}
	}
	for _, item := range results {
		io.WriteString(w, item+"\n")
	}
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {

	collection := mux.Vars(r)["collection"]
	newTrie := trie.NewTrie()
	dt := time.Now()

	count := 0
	total := 0
	bodyReader := bufio.NewReader(r.Body)

	for {
		bodyLine, err := bodyReader.ReadString('\n')
		if len(bodyLine) > 0 {
			total++

			key, value, err := parseData(bodyLine)
			if err != nil {
				log.Debugf("Can't parse data in request body: %s", bodyLine)
			} else {
				count++
				newTrie.Insert(key, value)
			}
		}
		if err != nil {
			break
		}
	}

	dur := time.Now().Sub(dt)

	if count > 0 {
		s.store[collection] = newTrie
		log.Noticef("%d/%d lines indexed. Index updated in %s. Memory used: %d", count, total, dur.String(), newTrie.Size())
		io.WriteString(w, "OK")
	} else {
		http.Error(w, "No valid data", 400)
	}

}

func (s *Server) Start(host string, port int) {
	listenAddr := fmt.Sprintf("%s:%d", host, port)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		panic(err)
	}
}
