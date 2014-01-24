package main

// import "fmt"
import "gofs"

func main() {
  proc := gofs.InitProc()
  file1 := proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
  proc.Close(file1)
}
