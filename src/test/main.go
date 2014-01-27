package main

// import "fmt"
import "gofs"
import "fmt"

func main() {
  proc := gofs.InitProc()
  fd, err := proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
  if err != nil { fmt.Println(err) }

  file := proc.GetFile(fd)
  file.Read(make([]byte, 100))
  file.Close()
}
