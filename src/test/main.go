package main

// import "fmt"
import "gofs"
import "fmt"

func main() {
  proc := gofs.InitProc()

  fmt.Println("------first time through-------\n")
  fd, err := proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
  if err != nil { fmt.Println(err) }

  file := proc.GetFile(fd)
  file.Write([]byte("Hello, world!"))
  file.Write([]byte("Hi, hi!"))

  buffer := make([]byte, 11)
  file.Seek(0, gofs.SEEK_SET)
  file.Read(buffer)
  fmt.Println(string(buffer))

  buffer = make([]byte, 11)
  file.Read(buffer)
  fmt.Println(string(buffer))
  file.Close()

  fmt.Println("\n\n------second time through-------\n")
  fd, err = proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
  if err != nil { fmt.Println(err) }

  file = proc.GetFile(fd)
  buffer = make([]byte, 25)
  file.Read(buffer)
  fmt.Println(string(buffer))
}
