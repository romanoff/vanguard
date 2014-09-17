package command

import (
	"github.com/bmizerany/pat"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/handler"
	"log"
	"net/http"
)

func NewAgentCommand() cli.Command {
	return cli.Command{
		Name:  "agent",
		Usage: "start agent server",
		Action: func(c *cli.Context) {
			agentCommandFunc()
		},
	}
}

func agentCommandFunc() {
	mux := pat.New()
	mux.Post("/containers", http.HandlerFunc(handler.ContainerCreate))
	mux.Put("/containers/:container_id", http.HandlerFunc(handler.ContainerUpdate))
	mux.Del("/containers/:container_id", http.HandlerFunc(handler.ContainerDelete))
	mux.Get("/containers", http.HandlerFunc(handler.ContainersIndex))
	http.Handle("/", mux)
	log.Println("Listening on port 2728")
	http.ListenAndServe(":2728", nil)
}
