package main

import (
    "flag"
    "fmt"
    "os"

    "k8s.io/klog"
)

func main() {
    klog.InitFlags(nil)
    //Init the command-line flags.
    flag.Parse()

    // Will be ignored as the program has exited in Fatal().
    defer func() {
        fmt.Println("Message in defer")
    }()

    // Flushes all pending log I/O.
    defer klog.Flush()

    // The temp folder for log files when --log_dir is not set.
    fmt.Printf("Temp folder for log files: %s\n", os.TempDir())

    klog.Info("Info")
    klog.V(4).Info("L4 info")
    klog.Error("Error")
    klog.Fatal("Fatal")

    // Will be ignored as the program has exited in Fatal().
    klog.Error("Error after Fatal")
}
