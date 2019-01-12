package main

import (
	_ "github.com/lib/pq"
	joyread "github.com/joyread/server"
)

func main() {
	joyread.StartServer()
}
