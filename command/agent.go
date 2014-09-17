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
	http.Handle("/", mux)
	log.Println("Listening on port 2728")
	http.ListenAndServe(":2728", nil)
}
