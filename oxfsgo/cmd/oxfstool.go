package main

import (
  "fmt"

  "github.com/io-core/oxfs-linux/oxfsgo"
)

func main() {
  fmt.Println("starting oxfstool")
  fmt.Println("Config:", oxfsgo.Config())
}
