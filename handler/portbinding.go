package handler

import (
	"encoding/json"
	"github.com/romanoff/vanguard/portbinding"
	"net/http"
)

func PortBindingCreate(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	host := r.FormValue("host")
	hostPort := r.FormValue("host_port")
	binding, err := portbinding.New(port, host, hostPort)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	go binding.Start()
	content, err := json.Marshal(binding)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	w.Write(content)
}

func PortBindingsIndex(w http.ResponseWriter, r *http.Request) {

}

func PortBindingDelete(w http.ResponseWriter, r *http.Request) {

}
