package command

import (
	"github.com/bmizerany/pat"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/handler"
	"github.com/romanoff/vanguard/host"
	"log"
	"net/http"
	"os"
)

func NewAgentCommand() cli.Command {
	return cli.Command{
		Name:  "agent",
		Usage: "start agent server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hostname",
				Usage: "server hostname",
			},
			cli.StringFlag{
				Name:  "ip",
				Usage: "server weave ip address",
			},
		},
		Action: func(c *cli.Context) {
			agentCommandFunc(c)
		},
	}
}

func agentCommandFunc(c *cli.Context) {
	mux := pat.New()
	mux.Post("/containers", http.HandlerFunc(handler.ContainerCreate))
	mux.Put("/containers/:container_id", http.HandlerFunc(handler.ContainerUpdate))
	mux.Del("/containers/:container_id", http.HandlerFunc(handler.ContainerDelete))
	mux.Get("/containers", http.HandlerFunc(handler.ContainersIndex))

	mux.Post("/portbindings", http.HandlerFunc(handler.PortBindingCreate))
	mux.Get("/portbindings", http.HandlerFunc(handler.PortBindingsIndex))
	mux.Del("/portbindings/:port", http.HandlerFunc(handler.PortBindingDelete))
	http.Handle("/", mux)
	currentHost, err := host.New(c.String("hostname"), c.String("ip"))
	if err != nil {
		log.Printf("Error getting host information: %v\n", err)
		os.Exit(1)
	}
	err = currentHost.Persist()
	if err != nil {
		log.Printf("Error persisting host information: %v\n", err)
		os.Exit(1)
	}
	log.Println("Listening on port 2728")
	http.ListenAndServe(":2728", nil)
}
