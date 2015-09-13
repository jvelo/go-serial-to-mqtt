package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"

	"github.com/tarm/serial"
)

func main() {
	config := &serial.Config{Name: "/dev/tty.usbmodem1411", Baud: 57600}
	s, err := serial.OpenPort(config)
	if err != nil {
		panic(err)
	}

	opts := MQTT.NewClientOptions().AddBroker("tcp://46.101.145.61:1883")
	opts.SetClientID("meteo-studio")

	// Connect MQTT client
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Loop over serial port lines
	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "|")
		if parts[0] == "6" {
			// if command is "6" (meteo data payload), push it to the broker
			token := c.Publish("studio/meteo", 0, false, parts[1])
			token.Wait()
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
