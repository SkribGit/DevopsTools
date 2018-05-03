package main

import (
  "bytes"
  "bufio"
  "flag"
  "fmt"
  "log"
  "os"
  "os/exec"
  "strings"
)

func main() {
  var memoryLimit int
  var runMode string
  var testFilename string

  flag.IntVar(&memoryLimit, "limit", 500, "worker memory limit")
  flag.StringVar(&runMode, "mode", "dryrun", "run mode")
  flag.StringVar(&testFilename, "testFilename", "test.txt", "Test file")

  flag.Parse()

  // read input from test file if it was supplied
  if runMode == "test" && testFilename != "" {
    fmt.Println("Running in test mode.")
    fmt.Printf("Reading input from %s\n", testFilename)

    // the scanner block below is from https://stackoverflow.com/a/16615559
    file, err := os.Open(testFilename)
    if err != nil {
      log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      if strings.Contains(scanner.Text(), "RackApp") {
        // extract the size of the worker
        // extract the PID
        fmt.Println(scanner.Text())
      }
    }
    if err := scanner.Err(); err != nil {
      log.Fatal(err)
    }
  } else {


  // run mode
  // dry run - do nothing, but print PIDs that will be terminated
  // live - actually terminate the PIDs

  // run passenger-memory-stats and parse the output of the command
    cmd := exec.Command("passenger-memory-stats")
    cmdReader, err := cmd.Output()
    if err != nil {
      log.Fatal(err)
    }
    fmt.Printf("Terminating workers that exceed the %dMB limit\n", memoryLimit)
    scanner := bufio.NewScanner(bytes.NewReader(cmdReader))
    for scanner.Scan() {
      if strings.Contains(scanner.Text(), "RackApp") {
        fmt.Println(scanner.Text())
      }
    }
  }

}
