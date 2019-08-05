package main

import (
	"github.com/ragon000/srandr/pkg/sway"
	"github.com/ragon000/srandr/internal/swaytui"
	//"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
        sock, err := sway.CreateSocket()
        if err != nil {
          panic(err)
        }
        defer sock.Close()
        swaytui.Start(sock)


}
