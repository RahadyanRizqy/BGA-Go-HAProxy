package funcs

import (
	"bga_go_haproxy/utils"
	"fmt"
	"os/exec"
)

func AllWeightValidation(current map[string]utils.VMRank, previous map[string]int) bool {
	for name, info := range current {
		prev, ok := previous[name]
		if !ok {
			continue // Tidak ada data sebelumnya → anggap valid
		}
		if info.Weight == prev {
			return false // Ada satu saja yang sama → tidak valid
		}
	}
	return true // Semua weight berbeda dari sebelumnya
}

func SomeWeightValidation(current map[string]utils.VMRank, previous map[string]int) bool {
	for name, info := range current {
		prev, ok := previous[name]
		if !ok || info.Weight != prev {
			// Ada satu saja VM yang belum pernah dicek (prev=0) atau weight-nya berubah
			return true
		}
	}
	return false // Semua weight sama dengan sebelumnya
}

func SetWeight(ranked map[string]utils.VMRank, cfg utils.BgaEnv) error {
	// Ambil backend dan path sock
	backend := cfg.HAProxyBackend
	sockPath := cfg.HAProxySock

	for vmName, data := range ranked {
		weight := data.Weight

		// Perintah shell untuk mengatur bobot
		cmdStr := fmt.Sprintf(`echo "set weight %s/%s %d" | socat stdio %s`, backend, vmName, weight, sockPath)

		if cfg.BgaUpdater {
			cmd := exec.Command("bash", "-c", cmdStr)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Gagal set weight untuk %s: %v\nOutput: %s\n", vmName, err, string(output))
				continue
			}
			if cfg.ConsolePrint {
				fmt.Printf("Set weight untuk %s: %d\n", vmName, weight)
			}
		}
	}

	return nil
}
