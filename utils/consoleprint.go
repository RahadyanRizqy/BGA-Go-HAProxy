package utils

import (
	"fmt"
	"sort"
)

var prevWeights map[string]float64

func ConsolePrint(result map[string][3]float64, iter int, cfg BgaEnv) ([]VMWeight, bool) {

	// Ambil weight saat ini
	currentWeights := make(map[string]float64)
	for name, values := range result {
		currentWeights[name] = values[1]
	}

	// // Bandingkan dengan iterasi sebelumnya
	// if prevWeights != nil {
	// 	for name, weight := range currentWeights {
	// 		if prev, exists := prevWeights[name]; exists && prev == weight {
	// 			// Jika ada satu weight sama, return tanpa mencetak
	// 			return nil, false
	// 		}
	// 	}
	// }

	weightChanged := false

	if prevWeights != nil {
		for name, weight := range currentWeights {
			if prev, exists := prevWeights[name]; !exists || prev != weight {
				weightChanged = true
				break
			}
		}

		if !weightChanged {
			return nil, false
		}
	}

	// Simpan weight saat ini untuk iterasi selanjutnya
	prevWeights = currentWeights

	// Cetak iterasi dan fitness
	if cfg.ConsolePrint {
		fmt.Printf("\n--- Iterasi %d ---\n", iter+1)

		var fitness float64
		for _, values := range result {
			fitness = values[2]
			break
		}
		fmt.Printf("Fitness: %.2f\n", fitness)
	}

	// Buat slice VMWeight
	var vms []VMWeight
	for name, values := range result {
		vms = append(vms, VMWeight{
			Name:      name,
			TaskTotal: values[0],
			Weight:    values[1],
		})
	}

	// Urutkan berdasarkan nama VM
	sort.Slice(vms, func(i, j int) bool {
		return vms[i].Name < vms[j].Name
	})

	// Cetak tabel
	if cfg.ConsolePrint {
		fmt.Println("Pemetaan Tugas ke VM:")
		fmt.Println("----------------------------------------------------------------")
		fmt.Printf("| %-7s | %-30s | %-6s |\n", "ID VM", "Total tugas yang bisa ditangani", "Weight")
		fmt.Println("----------------------------------------------------------------")
		for _, vm := range vms {
			fmt.Printf("| %-5s | %-30.0f | %-6.0f |\n", vm.Name, vm.TaskTotal, vm.Weight)
		}
		fmt.Println("----------------------------------------------------------------")
	}

	return vms, true
}
