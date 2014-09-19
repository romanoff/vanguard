package handler

import (
	"encoding/json"
	"github.com/romanoff/vanguard/container"
	"net/http"
)

func ContainerCreate(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	tag := r.FormValue("tag")
	imageId := r.FormValue("image_id")
	variables := make(map[string]string)
	for key, values := range r.Form {
		if key != "name" && key != "tag" && key != "image_id" {
			variables[key] = values[0]
		}
	}
	c := &container.Container{Name: name, Tag: tag, ImageId: imageId, Variables: variables}
	err := c.Run()
	if err == nil {
		content, err := json.Marshal(c)
		if err != nil {
			w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
			return
		}
		w.Write(content)
	} else {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
	}
}

func ContainerUpdate(w http.ResponseWriter, r *http.Request) {
}

func ContainerDelete(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	containerId := params.Get(":container_id")
	c, err := container.GetByContainerId(containerId)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	err = c.Stop()
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	w.Write([]byte("{\"success\": true}"))
}

func ContainersIndex(w http.ResponseWriter, r *http.Request) {
	check := r.FormValue("check")
	containers, err := container.GetContainers(check != "")
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	containersJson, err := json.Marshal(containers)
	if err != nil {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		return
	}
	w.Write(containersJson)
}
