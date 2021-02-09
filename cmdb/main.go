package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

// AllocatedName holds metadata of a name allocated to a deployed resource
// within one of our theoretical deployment environments
type AllocatedName struct {
	Name string
	Type string
	Region string
}

var names map[string]*AllocatedName
var apiVersion = "v1"
var port = 8990

func main() {
	names = make(map[string]*AllocatedName)
	log.Println(fmt.Sprintf("Booting up server %s on %d", apiVersion, port))
	http.HandleFunc(fmt.Sprintf("/%s/name", apiVersion), allocateName)
	http.HandleFunc(fmt.Sprintf("/%s/details", apiVersion), getName)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// retrieves theM8081TfooRba named query param requested from the received HTTP request
func requestDetails(r *http.Request, key string) (val string, err error) {
	q := r.URL.Query()
	if val := q.Get(key); val == "" {
		return "", fmt.Errorf("Unable to fetch key %s from request", key)
	} else {
		return val, nil
	}
}

func getAllocatedName(name string) *AllocatedName {
	if v, ok := names[name]; ok {
		return v
	}
	return nil
}

func putAllocatedName(name, resource_type, region string) *AllocatedName {
	allocatedName := &AllocatedName{
		Name:   name,
		Type:   resource_type,
		Region: region,
	}
	names[name] = allocatedName
	return allocatedName
}

func generateName(resource_type, region string) string {
	newName := fmt.Sprint("M", rand.Intn(10000), "T", resource_type, "R", region[0:2])
	if _, exists := names[newName]; exists {
		return generateName(resource_type, region)
	} else {
		return newName
	}
}

func writeResponse(w http.ResponseWriter, a *AllocatedName) {
	resp, err := json.Marshal(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// HTTP PUT - /name?resource_type=${type}&region=${region}
func allocateName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resource_type, err := requestDetails(r, "resource_type")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	region, err := requestDetails(r, "region")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := generateName(resource_type, region)
	allocatedName := putAllocatedName(name, resource_type, region)
	writeResponse(w, allocatedName)
}

// HTTP GET - /details?name=${name}
func getName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name, err := requestDetails(r, "name")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	o := getAllocatedName(name)
	if o == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeResponse(w, o)
}