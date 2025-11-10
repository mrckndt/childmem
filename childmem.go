package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/procfs"
)

type ProcessInfo struct {
	Timestamp int64
	PID       int
	MemoryKB  uint64
	CmdLine   []string
}

func getPIDByName(name string) (int, error) {
	fs, err := procfs.NewDefaultFS()
	if err != nil {
		return 0, err
	}

	procs, err := fs.AllProcs()
	if err != nil {
		return 0, err
	}

	for _, p := range procs {
		comm, err := p.Comm()
		if err != nil {
			// process might have terminated since we got the list
			continue
		}
		if comm == name {
			return p.PID, nil
		}
	}

	return 0, fmt.Errorf("process %s not found", name)
}

func getChildrenInfo(PID int) ([]ProcessInfo, error) {
	var procsList []ProcessInfo

	fs, err := procfs.NewDefaultFS()
	if err != nil {
		return procsList, err
	}

	procs, err := fs.AllProcs()
	if err != nil {
		return procsList, err
	}

	for _, p := range procs {
		stat, err := p.Stat()
		if err != nil {
			// process might have terminated since we got the list
			continue
		}

		if stat.PPID == PID {
			cmdLine, err := p.CmdLine()
			if err != nil {
				cmdLine = []string{"<unavailable>"}
			}

			record := ProcessInfo{
				Timestamp: time.Now().Unix(),
				PID:       p.PID,
				MemoryKB:  uint64(stat.ResidentMemory()) / 1024,
				CmdLine:   cmdLine,
			}

			procsList = append(procsList, record)
		}
	}

	return procsList, nil
}

func writeRecordsCSV(path string, procsList []ProcessInfo) {
	var err error
	_, err = os.Stat(path)
	fileExists := !os.IsNotExist(err)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open or create CSV file %s: %v", path, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if !fileExists {
		header := []string{"timestamp", "pid", "memory_KB", "cmdline"}
		if err := writer.Write(header); err != nil {
			log.Fatalf("failed to write header to CSV: %v", err)
		}
	}

	for _, procInfo := range procsList {
		record := []string{
			strconv.FormatInt(procInfo.Timestamp, 10),
			strconv.Itoa(procInfo.PID),
			strconv.FormatUint(procInfo.MemoryKB, 10),
			strings.Join(procInfo.CmdLine, " "),
		}

		if err := writer.Write(record); err != nil {
			log.Printf("failed to write record to CSV: %v", err)
			break
		}

		fmt.Printf("record %v written to file '%s'\n", record, path)
	}
}

func main() {
	if runtime.GOOS != "linux" {
		log.Fatalf("This program is designed to run on Linux as it relies on the proc filesystem.")
	}

	PID, _ := getPIDByName("mattermost")
	procsList, _ := getChildrenInfo(PID)
	writeRecordsCSV("./plugin_mem.csv", procsList)
}
