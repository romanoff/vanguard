package handler

import (
	"github.com/romanoff/vanguard/container"
	"net/http"
)

func ContainerCreate(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	tag := r.FormValue("tag")
	imageId := r.FormValue("image_id")
	c := &container.Container{Name: name, Tag: tag, ImageId: imageId}
	err := c.Run()
	if err == nil {
		w.Write([]byte("{\"container_id\": \"" + c.ContainerId + "\"}"))
	} else {
		w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
	}
}

func ContainerUpdate(w http.ResponseWriter, r *http.Request) {
}

func ContainerDelete(w http.ResponseWriter, r *http.Request) {
}

func ContainersIndex(w http.ResponseWriter, r *http.Request) {
}
