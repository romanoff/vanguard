package handler

import (
	"encoding/json"
	"github.com/romanoff/vanguard/portbinding"
	"log"
	"net/http"
)

var bindingsList []*portbinding.PortBinding = []*portbinding.PortBinding{}

func getPortBinding(port string) (*portbinding.PortBinding, error) {
	var binding *portbinding.PortBinding
	for _, b := range bindingsList {
		if b.Port == port {
			binding = b
		}
	}
	if binding == nil {
		binding, err := portbinding.New(port)
		if err != nil {
			return nil, err
		}
		bindingsList = append(bindingsList, binding)
	}
	return binding, nil

}

func PortBindingCreate(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	host := r.FormValue("host")
	hostPort := r.FormValue("host_port")
	binding, err := getPortBinding(port)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	binding.AddBackend(host, hostPort)
	go binding.Start()
	content, err := json.Marshal(binding)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	log.Println("exposed " + binding.String())
	w.Write(content)
}

func PortBindingsIndex(w http.ResponseWriter, r *http.Request) {
	content, err := json.Marshal(bindingsList)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	w.Write(content)
}

func PortBindingDelete(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	port := params.Get(":port")

	i := -1
	for j, binding := range bindingsList {
		if binding.Port == port {
			err := binding.Stop()
			if err != nil {
				w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
				return
			}
			i = j
			log.Println("stopped " + binding.String())
		}
	}
	if i == -1 {
		w.Write([]byte("{}"))
		return
	}
	if i != -1 {
		binding := bindingsList[i]
		content, err := json.Marshal(binding)
		if err != nil {
			w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
			return
		}
		w.Write(content)
		copy(bindingsList[i:], bindingsList[i+1:])
		bindingsList[len(bindingsList)-1] = nil
		bindingsList = bindingsList[:len(bindingsList)-1]
	}
}
