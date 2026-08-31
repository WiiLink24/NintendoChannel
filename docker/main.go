package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/logrusorgru/aurora/v4"
)

var (
	config Config
)

const SocketSuccess = `{"success": true}`

func checkError(err error) {
	if err != nil {
		log.Fatalf("WiiLink Nintendo Channel File Generator has encountered a fatal error! Reason: %v\n", err)
	}
}

func socketFail(err error) []byte {
	data, _ := json.Marshal(SocketFailResponse{
		Success: false,
		Error:   err.Error(),
	})
	return data
}

func main() {
	rawConfig, err := os.ReadFile("config.xml")
	checkError(err)

	err = xml.Unmarshal(rawConfig, &config)
	checkError(err)

	// Remove if it didn't gracefully exit for some reason
	err = os.Remove("/tmp/nc-gen.sock")
	if !os.IsNotExist(err) {
		checkError(err)
	}

	socket, err := net.Listen("unix", "/tmp/nc-gen.sock")
	checkError(err)

	defer func(socket net.Listener) {
		err := socket.Close()
		if err != nil {
			log.Println(aurora.Red("error closing socket:"), err)
		}
	}(socket)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err = os.Remove("/tmp/nc-gen.sock")
		checkError(err)
		os.Exit(0)
	}()

	log.Printf("%s", aurora.Green("UNIX socket connected."))
	log.Printf("%s %s\n", aurora.Green("Listening on UNIX socket:"), socket.Addr())

	// Listen forever
	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Print(aurora.Red("Socket Accept ERROR: "), err.Error(), "\n")
		}

		go func(conn net.Conn) {
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					log.Println(aurora.Red("error closing connection:"), err)
				}
			}(conn)
			buf := make([]byte, 4096)

			n, err := conn.Read(buf)
			if err != nil {
				log.Print(aurora.Red("Socket Read ERROR: "), err.Error(), "\n")
				return
			}

			reply := []byte(SocketSuccess)
			payloadRaw := strings.ReplaceAll(string(buf[:n]), "\n", "")

			var payload NintendoChannelWebPayload
			err = json.Unmarshal([]byte(payloadRaw), &payload)
			if err != nil {
				log.Print(aurora.Red("Unmarshal ERROR: "), err.Error(), "\n")
				_, err = conn.Write(append(socketFail(err), []byte("\n")...))
				if err != nil {
					log.Println(aurora.Red("error writing error (wow!):"), err)
				}
				return
			}

			if config.SocketSecret != payload.SocketSecret {
				log.Print(aurora.Red("Authentication failure."), "\n")
				_, err = conn.Write(append(socketFail(errors.New("authentication error")), []byte("\n")...))
				if err != nil {
					log.Println(aurora.Red("error writing error (wow!):"), err)
				}
				return
			}

			log.Println("Running file generator with the following arguments:", payload.Arguments)
			outRaw, err := exec.Command("./cli", payload.Arguments).Output()
			out := string(outRaw)
			lines := strings.Split(out, "\n")
			for _, line := range lines {
				if len(strings.TrimSpace(line)) > 0 {
					log.Println(line)
				}
			}

			if err != nil {
				log.Print(aurora.Red("Error running file generator: "), err.Error(), "\n")
				_, err = conn.Write(append(socketFail(err), []byte("\n")...))
				if err != nil {
					log.Println(aurora.Red("error writing error (wow!):"), err)
				}
				return
			}
			log.Println(aurora.Green("Successfully generated files with arguments:"), payload.Arguments)

			_, err = conn.Write(append(reply, []byte("\n")...))
			if err != nil {
				log.Print(aurora.Red("Socket Write ERROR: "), err.Error(), "\n")
				return
			}
		}(conn)
	}
}
