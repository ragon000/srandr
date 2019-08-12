// Package sway provides functions and structs to communicate with a sway IPC socket according to sway-ipc(7)
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

type SwayConnection struct {
	Outputs []Output
	sock    net.Conn
}

type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

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
	Rect              Rectangle
}

func (s *SwayConnection) CloseConnection() {
  s.sock.Close()
}


func (s *Output) IsEqualTo(o Output)  bool{
  return s.Name == o.Name && s.Serial == o.Serial
}

func UpOf(o Output, outputs []Output) Output{
  var returnOutput Output

  for i, out := range outputs {
    if i == 0 {
      returnOutput = out
    }else{
    if returnOutput.Rect.Y > out.Rect.Y && o.Rect.Y < out.Rect.Y {
      returnOutput = out
    }
    }
  }
  return returnOutput
}

func DownOf(o Output, outputs []Output) Output{
  var returnOutput Output

  for i, out := range outputs {
    if i == 0 {
      returnOutput = out
    }else{
    if returnOutput.Rect.Y < out.Rect.Y && o.Rect.Y > out.Rect.Y {
      returnOutput = out
    }
    }
  }
  return returnOutput
}

func RightOf(o Output, outputs []Output) Output{
  var returnOutput Output

  for i, out := range outputs {
    if i == 0 {
      returnOutput = out
    }else{
    if returnOutput.Rect.X > out.Rect.X && o.Rect.X < out.Rect.X {
      returnOutput = out
    }
    }
  }
  return returnOutput
}

func LeftOf(o Output, outputs []Output) Output{
  var returnOutput Output

  for i, out := range outputs {
    if i == 0 {
      returnOutput = out
    }else{
    if returnOutput.Rect.X < out.Rect.X && o.Rect.X > out.Rect.X {
      returnOutput = out
    }
    }
  }
  return returnOutput
}

func CreateSwayConnection() (SwayConnection, error) {
	s := SwayConnection{}
	err := s.createSocket()
	if err != nil {
		return SwayConnection{}, err
	}
        s.GetOutputsFromSocket()
	return s, nil
}

func (s *SwayConnection) applyingModes(o []Output) error{
    s.


}

// Creates the socket to communicate with sway
// Uses the environment variable $SWAYSOCK to find the socket path
func (s *SwayConnection) createSocket() error {
	swaysockpath := os.Getenv("SWAYSOCK")
	if swaysockpath == "" {
		return fmt.Errorf("$SWAYSOCK is empty")
	}
	var err error
	s.sock, err = net.Dial("unix", swaysockpath)
	if err != nil {
		return err
	}
	return nil
}

// Sets SwayConnection.Outputs to the outputs currently present on the Sway Socket
func (s *SwayConnection) GetOutputsFromSocket() error {
	command := createCommand(GET_OUTPUTS, "")
	_, err := s.sock.Write(command)
	if err != nil {
		return err
	}
	s.sock.Read(make([]byte, len(command))) // sway returns the sent command, but we don't want it
	dec := json.NewDecoder(s.sock)
	var o []Output
	if err := dec.Decode(&o); err != nil {
		return err
	}
	s.Outputs = o
	return nil
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
