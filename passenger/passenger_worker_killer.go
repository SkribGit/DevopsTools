package main

import (
  "bytes"
  "bufio"
  "flag"
  "fmt"
  "log"
  "os"
  "os/exec"
  "strconv"
  "strings"
)

type worker struct {
  pid string
  memory int
}

func get_kill_string(pid string) string {
  return fmt.Sprintf("env kill -6 %s", pid)
}

func get_worker_signature(_passengerVersion string) string {
  switch _passengerVersion {
  case "5":
    return "Passenger AppPreloader"
  case "4":
    return "Passenger RackApp"
  default:
    return "Passenger RackApp"
  }
}

func get_passenger_workers(scanner *bufio.Scanner, _passengerVersion string) []worker {
  workers := []worker{}
  worker_signature := get_worker_signature(_passengerVersion)

  for scanner.Scan() {
    if strings.Contains(scanner.Text(), worker_signature) {
      s := strings.Fields(scanner.Text())
      pid := s[0]
      mem, err := strconv.ParseFloat(s[1], 64)
      w := worker{pid: pid, memory: int(mem)}
      workers = append(workers, w)
      if err != nil {
        log.Fatal(err)
      }
    }
  }
  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
  return workers
}

func main() {
  var memoryLimit int
  var runMode string
  var testFilename string
  var passengerVersion string

  var workers []worker

  flag.IntVar(&memoryLimit, "limit", 500, "worker memory limit")
  flag.StringVar(&runMode, "mode", "dryrun", "run mode")
  flag.StringVar(&testFilename, "test_filename", "test.txt", "Test file")
  flag.StringVar(&passengerVersion, "passenger_version", "5", "Passenger version")

  flag.Parse()

  if runMode == "test" && testFilename != "" {
    // Test mode
    // Read input from input file
    fmt.Println("Running in test mode.")
    fmt.Printf("Reading input from %s\n", testFilename)

    // the scanner block below is from https://stackoverflow.com/a/16615559
    file, err := os.Open(testFilename)
    if err != nil {
      log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    workers = get_passenger_workers(scanner, passengerVersion)
    for _, worker := range workers {
      if worker.memory > memoryLimit {
        fmt.Printf("Terminating worker with PID %s. Memory size: %d\n", worker.pid, worker.memory)
        killString := get_kill_string(worker.pid)
        fmt.Println(killString)
      }
    }
  } else {
    // Live mode
    // run passenger-memory-stats and parse the output of the command
    cmd := exec.Command("passenger-memory-stats")
    cmdReader, err := cmd.Output()
    if err != nil {
      log.Fatal(err)
    }
    fmt.Printf("Terminating workers that exceed the %dMB limit\n", memoryLimit)
    scanner := bufio.NewScanner(bytes.NewReader(cmdReader))
    workers = get_passenger_workers(scanner, passengerVersion)
    for _, worker := range workers {
      if worker.memory > memoryLimit {
        fmt.Printf("Terminating worker with PID %s. Memory size: %d\n", worker.pid, worker.memory)
        killString := get_kill_string(worker.pid)
        killCmd := exec.Command(killString)
        killErr := killCmd.Run()
        if killErr != nil {
          log.Fatal(killErr)
        }
      }
    }
  }
}
