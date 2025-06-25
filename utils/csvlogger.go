package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func InitCSV(cfg BgaEnv) string {
	if cfg.Logger {

		if err := os.MkdirAll("data", 0755); err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return ""
		}

		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("data/vm_metrics_%d.csv", timestamp)

		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return ""
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		headers := []string{
			"no", "fetch", "update", "vm_id", "vm_name", "cpu_usage", "max_cpu",
			"mem_usage", "max_mem", "cum_netin", "cum_netout", "rate_netin",
			"rate_netout", "bw_usage", "max_bw", "score", "priority", "weight",
			"unix_timestamp", "timestamp",
		}
		if err := writer.Write(headers); err != nil {
			fmt.Printf("Error writing headers: %v\n", err)
			return ""
		}
		writer.Flush()

		return filename
	}
	return ""
}

func StoreCSV(
	cfg BgaEnv,
	filename string,
	logLine *int,
	fetchCount int,
	updateCount int,
	unix_timestamp int64,
	timestamp string,
	stats map[string]VMStats,
	ranked map[string]VMPriority,
	netIfaceRate float64,
) error {
	if cfg.Logger {
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		for name, stat := range stats {
			record := []string{
				strconv.Itoa(*logLine),
				strconv.Itoa(fetchCount),
				strconv.Itoa(updateCount),
				strconv.Itoa(stat.VM.Id),
				stat.VM.Name,
				strconv.FormatFloat(stat.VM.CPU, 'f', -1, 64),
				strconv.Itoa(stat.VM.MaxCPU),
				strconv.FormatFloat(stat.MemUsage, 'f', -1, 64),
				strconv.Itoa(stat.VM.MaxMem),
				strconv.Itoa(stat.VM.CumNetIn),
				strconv.Itoa(stat.VM.CumNetOut),
				strconv.FormatFloat(stat.Rates.Rx, 'f', -1, 64),
				strconv.FormatFloat(stat.Rates.Tx, 'f', -1, 64),
				strconv.FormatFloat(stat.BwUsage, 'f', -1, 64),
				strconv.FormatFloat(netIfaceRate, 'f', -1, 64),
				strconv.FormatFloat(stat.Score, 'f', -1, 64),
				strconv.Itoa(ranked[name].Priority),
				strconv.Itoa(ranked[name].Weight),
				strconv.FormatInt(unix_timestamp, 10),
				timestamp,
			}

			if err := writer.Write(record); err != nil {
				return fmt.Errorf("error writing record: %v", err)
			}
			*logLine++
		}
	}

	return nil
}

// func Logger(currentBest Chromosome, iter int, change int) {
// 	InitClient()
// 	cfg := LoadBgaEnv()
// 	csvFileName := InitCSV(cfg)
// 	prevTime := time.Now()

// }
