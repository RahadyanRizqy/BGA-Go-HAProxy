package funcs

import (
	"bga_go_haproxy/utils"
	"fmt"
	"os/exec"
)

func ChangeWeight(result []utils.VMWeight, status bool, cfg utils.BgaEnv) {
	if !status {
		return // Tidak perlu ubah jika tidak ada perubahan
	}

	for _, vm := range result {
		// Format perintah
		cmdStr := fmt.Sprintf(`echo "set weight %s/%s %d" | socat stdio %s`,
			cfg.HAProxyBackend, vm.Name, int(vm.Weight), cfg.HAProxySock)

		// Jalankan command dengan shell
		if cfg.BGAUpdater {
			cmd := exec.Command("bash", "-c", cmdStr) // Jika di Windows pakai "cmd", "/C"
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Gagal set weight untuk %s: %v\n", vm.Name, err)
			} else {
				fmt.Printf("Set weight VM %s ke %d: %s\n", vm.Name, int(vm.Weight), string(output))
			}
		}
	}
}
