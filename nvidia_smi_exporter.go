package main

import (
    "bytes"
    "encoding/csv"
    "flag"
    "fmt"
    "log"
    "net/http"

    //    "os"
    "os/exec"
    "strings"
)

// name, index, temperature.gpu, utilization.gpu,
// utilization.memory, memory.total, memory.free, memory.used, uuid

var (
    listenAddress string
    metricsPath   string
)

func metrics(response http.ResponseWriter, request *http.Request) {
    out, err := exec.Command(
        "nvidia-smi",
        "--query-gpu=name,index,uuid,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used",
        "--format=csv,noheader,nounits").Output()

    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    csvReader := csv.NewReader(bytes.NewReader(out))
    csvReader.TrimLeadingSpace = true
    records, err := csvReader.ReadAll()

    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    metricList := []string{
        "temperature.gpu", "utilization.gpu",
        "utilization.memory", "memory.total", "memory.free", "memory.used"}

    result := ""
    for _, row := range records {
        name := fmt.Sprintf("%s[%s]", row[0], row[1])
        uuid := row[2]
        for idx, value := range row[3:] {
            metric := metricList[idx]
            result = fmt.Sprintf(
                "%s%s%s{gpu=\"%s\"} %s\n", result, "nvidia.",
                metric, name, value)

            result = fmt.Sprintf(
                "%s%s%s%s{uuid=\"%s\"} %s\n", result, "nvidia.",
                metric, ".uuid", uuid, value)
        }
    }

    fmt.Fprintf(response, strings.Replace(result, ".", "_", -1))
}

func init() {
    flag.StringVar(&listenAddress, "web.listen-address", ":9101", "Address to listen on")
    flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
    flag.Parse()
}

func main() {
    //    addr := ":9101"
    //    if len(os.Args) > 1 {
    //        addr = ":" + os.Args[1]
    //    }

    http.HandleFunc(metricsPath, metrics)
    err := http.ListenAndServe(listenAddress, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
