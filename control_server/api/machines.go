package api

import (
	"encoding/json"
	"github.com/baas-project/baas/pkg/model"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (api *Api) GetMachine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mac, ok := vars["mac"]
	if !ok || mac == "" {
		http.Error(w, "invalid mac address", http.StatusBadRequest)
		log.Error("Invalid mac address given")
		return
	}

	machine, err := api.store.GetMachineByMac(mac)
	if err != nil {
		http.Error(w, "couldn't get machine", http.StatusInternalServerError)
		log.Errorf("get machine by mac: %v", err)
		return
	}

	e := json.NewEncoder(w)
	_ = e.Encode(machine)
}

func (api *Api) GetMachines(w http.ResponseWriter, r *http.Request) {
	machines, err := api.store.GetMachines()
	if err != nil {
		http.Error(w, "couldn't get machines", http.StatusInternalServerError)
		log.Errorf("get machines: %v", err)
		return
	}

	e := json.NewEncoder(w)
	_ = e.Encode(machines)
}

func (api *Api) UpdateMachine(w http.ResponseWriter, r *http.Request) {
	var machine model.MachineModel
	err := json.NewDecoder(r.Body).Decode(&machine)
	if err != nil {
		http.Error(w, "invalid machine given", http.StatusBadRequest)
		log.Errorf("Invalid machine given: %v", err)
		return
	}

	err = api.store.UpdateMachine(&machine)
	if err != nil {
		http.Error(w, "couldn't update machine", http.StatusInternalServerError)
		log.Errorf("get update machine: %v", err)
		return
	}
}
