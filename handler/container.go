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
	label := r.FormValue("label")
	variables := make(map[string]string)
	dns := make([]string, 0, 0)
	volumes := make([]string, 0, 0)
	for key, values := range r.Form {
		if key != "label" && key != "name" && key != "tag" && key != "image_id" && key != "dns" && key != "volumes" {
			variables[key] = values[0]
		}
		if key == "dns" {
			for _, value := range values {
				dns = append(dns, value)
			}
		}
		if key == "volumes" {
			for _, volume := range values {
				volumes = append(volumes, volume)
			}
		}
	}
	c := &container.Container{
		Label:     label,
		Name:      name,
		Tag:       tag,
		ImageId:   imageId,
		Variables: variables,
		DNS:       dns,
		Volumes:   volumes,
	}
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
