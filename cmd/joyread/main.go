package main

import (
	_ "github.com/lib/pq"
	joyread "gitlab.com/joyread/ultimate"
)

func main() {
	joyread.StartServer()
}
