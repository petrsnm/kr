package main

/*
* CLI to control krssh-agent
 */

import (
	"bitbucket.org/kryptco/krssh"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

func connectToAgent() (conn net.Conn, err error) {
	agentSockName := os.Getenv(krssh.KRSSH_CTL_SOCK_ENV)
	conn, err = net.Dial("unix", agentSockName)
	return
}

func pairCommand(c *cli.Context) (err error) {
	agentConn, err := connectToAgent()
	if err != nil {
		log.Fatal(err)
	}

	pairingSecret, err := krssh.GeneratePairingSecret()
	if err != nil {
		log.Fatal(err)
	}

	pairingSecretJson, err := json.Marshal(pairingSecret)
	if err != nil {
		log.Fatal(err)
	}

	pairRequest, err := http.NewRequest("PUT", "/pair", bytes.NewReader(pairingSecretJson))
	if err != nil {
		log.Fatal(err)
	}

	err = pairRequest.Write(agentConn)
	if err != nil {
		log.Fatal(err)
	}

	qr, err := QREncode(pairingSecretJson)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Scan this QR Code with the krSSH Mobile App to connect it with this workstation.")
	fmt.Println()
	fmt.Println(qr.Terminal)

	bufReader := bufio.NewReader(agentConn)
	response, err := http.ReadResponse(bufReader, pairRequest)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()
	return
}

func main() {
	app := cli.NewApp()
	app.Name = "kr"
	app.Usage = "communicate with krssh-agent and krssh-iOS"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "pair",
			Aliases: []string{"p"},
			Action:  pairCommand,
		},
		cli.Command{
			Name:   "me",
			Action: pairCommand,
		},
		cli.Command{
			Name:    "list",
			Aliases: []string{"ls"},
			Action:  pairCommand,
		},
	}
	app.Run(os.Args)
}
