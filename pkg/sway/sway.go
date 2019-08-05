package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type IPC_command uint32

const (
	RUN_COMMAND       IPC_command = 0
	GET_WORKSPACES    IPC_command = 1
	SUBSCRIBE         IPC_command = 2
	GET_OUTPUTS       IPC_command = 3
	GET_TREE          IPC_command = 4
	GET_MARKS         IPC_command = 5
	GET_BAR_COFIG     IPC_command = 6
	GET_VERSION       IPC_command = 7
	GET_BINDING_MODES IPC_command = 8
	GET_CONFIG        IPC_command = 9
	SEND_TICK         IPC_command = 10
	// SYNC ommited cause it's only there for i3 compability
	GET_INPUTS IPC_command = 100
	GET_SEATS  IPC_command = 101
)

type Mode struct {
	Width   int
	Height  int
	Refresh int
}

type Output struct {
	Name              string
	Make              string
	Model             string
	Serial            string
	Active            bool
	Primary           bool
	Scale             float64
	Subpixel_hinting  string
	Transform         string
	Current_workspace string
	Modes             []Mode
	Current_mode      Mode
}

var sock net.Conn = nil

func main() {
	swaysockpath := os.Getenv("SWAYSOCK")
	sock, err := net.Dial("unix", swaysockpath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Message: ", getOutputs())
	sock.Close()
	return

}

func getOutputs() []Output {
	command := createCommand(GET_OUTPUTS, "")
	sock.Write(command)
	dec := json.NewDecoder(sock)
	var o []Output
	if err := dec.Decode(&o); err != nil {
		panic(err)
	}
	return o
}

func createCommand(i IPC_command, message string) []byte {
	i3ipcbytes := []byte("i3-ipc")
	messagelengthbytes := make([]byte, 4)
	messagebytes := []byte(message)
	commandbytes := make([]byte, 4)
	binary.BigEndian.PutUint32(commandbytes, uint32(i))
	binary.BigEndian.PutUint32(messagebytes, uint32(len(message)))
	toreturn := i3ipcbytes
	toreturn = append(toreturn, messagelengthbytes...)
	toreturn = append(toreturn, commandbytes...)
	toreturn = append(toreturn, messagebytes...)
	return toreturn
}
