package main

import (
  "os/exec"
  "bytes"
  "strings"
  "fmt"
)

func main() {
  benchmarks := []string{
    "BenchmarkOC1$",
    "BenchmarkOtC$",
    "BenchmarkOC$",
    "BenchmarkOtCtU$",
    "BenchmarkOCU$",
    "BenchmarkOWsC$",
    "BenchmarkOWsCU$",
    "BenchmarkOWbC$",
    "BenchmarkOWbCU$",
    "BenchmarkOWMsC$",
    "BenchmarkOWMsCU$",
    "BenchmarkOWMbC$",
    "BenchmarkOWMbCU$",
    "BenchmarkOWbbC$",
    "BenchmarkOWbbCU$",
  }

  for _, benchRegEx := range benchmarks {
    var out bytes.Buffer
    cmd := exec.Command("go", "test", "-bench", benchRegEx, "bench")
	cmd.Stdout = &out
    cmd.Start()

    err := cmd.Wait()
    if (err != nil) {
      fmt.Println("Benchmark", benchRegEx, "failed to run:", err)
      panic("Error running benchmark.")
    }

    output := out.String()
    fmt.Println(strings.Split(output, "\n")[1])
  }
}
