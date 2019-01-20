package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pkgHelm "github.com/masahiro331/phandler/pkg/helm"
	"github.com/pkg/errors"
)

var (
	Config *PhandlerConfig
)

const (
	Namespace           = "default"
	InternalServerError = `{"status": 500, "message":"Insternal server error"}`
)

type Deployment struct {
	ImageName string `json:"image_name"`
	ImageTag  string `json:"image_tag"`
	Replica   string `json:"replica"`
	Endpoint  string `json:"endpoint"`
	Port      string `json:"port"`
}

type ResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("[FATAL] %+v", err)
	}
}

func run() error {
	Config, err := NewPhandlerConfig(Namespace)
	if err != nil {
		return errors.Wrap(err, "")
	}
	fmt.Print(Config)
	http.HandleFunc("/deployment", DeploymentHandler)
	return http.ListenAndServe(":8080", nil)
}

func DeploymentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	helmClient, err := pkgHelm.NewHelmClient(Config.kubernetesClientSet, Config.kubernetesConfig)
	if err != nil {
		HandleError(w, http.StatusBadRequest, err)
		return
	}

	fmt.Print(helmClient)
	var deployment Deployment
	if err := json.NewDecoder(r.Body).Decode(&deployment); err != nil {
		HandleError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Printf("%+v", deployment)

	rb := &ResponseBody{
		Status:  200,
		Message: "success",
	}
	if err := json.NewEncoder(w).Encode(rb); err != nil {
		log.Print(err)
		HandleInternalServerError(w)
	}
}

func HandleError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	log.Print(err)
	w.WriteHeader(status)
	rb := &ResponseBody{
		Status:  status,
		Message: err.Error(),
	}
	if err := json.NewEncoder(w).Encode(rb); err != nil {
		log.Print(err)
		HandleInternalServerError(w)
	}
}

func HandleInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, InternalServerError)
}
