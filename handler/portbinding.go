package handler

import (
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
	err = binding.Start()
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	w.Write([]byte("{\"success\": true}"))
}

func PortBindingsIndex(w http.ResponseWriter, r *http.Request) {

}

func PortBindingDelete(w http.ResponseWriter, r *http.Request) {

}
