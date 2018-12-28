package main

import (
	_ "github.com/lib/pq"
	joyread "gitlab.com/joyread/server"
)

func main() {
	joyread.StartServer()
}
