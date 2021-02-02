// Package api provides functions for handling http requests on the control server.
// This is used to respond to requests from pixiecore, to serve files (kernel, initramfs, disk images)
// and to communicate with machines running the management os.
package api

import (
	"fmt"
	"github.com/baas-project/baas/pkg/database"
	"github.com/baas-project/baas/pkg/httplog"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func getHandler(machineStore database.Store, staticDir string, diskpath string) http.Handler {
	// Api for communicating with the management os
	api := NewApi(machineStore, diskpath)

	r := mux.NewRouter()

	r.StrictSlash(true)
	r.Use(logging)

	// Applications (in particular, the management OS) can send logs here to be logged on the control server.
	r.HandleFunc("/log", httplog.CreateLogHandler(log.StandardLogger()))

	// TODO: we may want to split this up, especially the disk images part
	// TODO: isn't this already the case?
	// Serve static files (kernel, initramfs, disk images)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	r.HandleFunc("/machine/{mac}", api.GetMachine).Methods(http.MethodGet)
	r.HandleFunc("/machines", api.GetMachines).Methods(http.MethodGet)
	r.HandleFunc("/machine", api.UpdateMachine).Methods(http.MethodPut)
	r.HandleFunc("/machine", api.UpdateMachine).Methods(http.MethodPut)

	r.HandleFunc("/users", api.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/user", api.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/user/{name}", api.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/user/{name}/createimage", api.CreateImage).Methods(http.MethodPost)

	// info about an image
	r.HandleFunc("/image/{uuid}", api.GetImage).Methods(http.MethodGet)
	r.HandleFunc("/image/{uuid}/{version}/download", api.DownloadImage).Methods(http.MethodGet)
	r.HandleFunc("/image/{uuid}/upload", api.UploadImage).Methods(http.MethodPost)

	mmosr := r.PathPrefix("/mmos").Subrouter()

	// Serve boot configurations to pixiecore (this url is hardcoded in pixiecore)
	r.HandleFunc("/v1/boot/{mac}", api.ServeBootConfigurations)

	mmosr.HandleFunc("/inform", api.BootInform).Methods(http.MethodPost)
	mmosr.HandleFunc("/disk/{uuid}", api.UploadDiskImage).Methods(http.MethodPost)
	mmosr.HandleFunc("/disk/{uuid}", api.DownloadDiskImage).Methods(http.MethodGet)

	return r
}

// StartServer defines all routes and1 then starts listening for HTTP requests.
// TODO: Config struct
func StartServer(machineStore database.Store, staticDir string, diskpath string, address string, port int) {
	srv := http.Server{
		Handler: getHandler(machineStore, staticDir, diskpath),
		Addr:    fmt.Sprintf("%s:%d", address, port),
	}
	log.Fatal(srv.ListenAndServe())
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We don't want to log the fact that we are logging.
		if r.URL.Path != "/log" {
			log.Debugf("%s request on %s", r.Method, r.URL)
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}