package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pkgHelm "github.com/masahiro331/phandler/pkg/helm"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
)

const (
	Namespace           = "default"
	InternalServerError = `{"status": 500, "message":"Insternal server error"}`
)

var Config *PhandlerConfig

type Deployment struct {
	ReleaseName string `json:"release_name"`
	Chart       string `json:"chart"`
	// Replica     string `json:"replica"`
	// Endpoint    string `json:"endpoint"`
	// Port        string `json:"port"`
}

type ResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func main() {
	Config = NewPhandlerConfig(Namespace)
	if Config == nil {
		log.Fatalf("[FATAL] %+v", Config)
	}
	fmt.Print(Config)
	if err := run(); err != nil {
		log.Fatalf("[FATAL] %+v", err)
	}
}

func run() error {
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

	var deployment Deployment
	if err := json.NewDecoder(r.Body).Decode(&deployment); err != nil {
		HandleError(w, http.StatusBadRequest, err)
		return
	}
	installOpts := []helm.InstallOption{
		helm.InstallDisableHooks(false),
		helm.InstallTimeout(30),
		// helm.ReleaseName(""),
	}
	chartRequested, err := chartutil.Load(deployment.Chart)
	result, err := helmClient.InstallReleaseFromChart(chartRequested, "default", installOpts...)
	fmt.Printf("%+v", result)

	if err != nil {
		log.Printf("%+v", errors.Wrapf(err, "[ERROR] Failed to install %s", deployment.Chart))
		HandleInternalServerError(w)
		return
	}

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
