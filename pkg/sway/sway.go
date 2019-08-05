package sway

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

type SwayConnection interface {
	createSocket() net.Conn
	getOutputs() []Output
}

func createSocket() net.Conn {
	swaysockpath := os.Getenv("SWAYSOCK")

	sock, err := net.Dial("unix", swaysockpath)
	if err != nil {
		panic(err)
	}
	return sock
}

func getOutputs(sock net.UnixConn) []Output {
	command := createCommand(GET_OUTPUTS, "")
	_, err = sock.Write(command)
	if err != nil {
		panic(err)
	}
	sock.Read(make([]byte, len(command)))
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
	commandbytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(commandbytes, uint32(i))
	binary.LittleEndian.PutUint32(messagelengthbytes, uint32(len(message)))
	toreturn := i3ipcbytes
	toreturn = append(toreturn, messagelengthbytes...)
	toreturn = append(toreturn, commandbytes...)
	if message != "" {
		messagebytes := []byte(message)
		toreturn = append(toreturn, messagebytes...)
		toreturn = append(toreturn, []byte("\n")...)
	}
	return toreturn
}
