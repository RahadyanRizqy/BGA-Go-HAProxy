package funcs

import (
	"bga_go_haproxy/utils"
	"math"
	"sort"
)

func DistributeWeights(arr []int, weightTotal int) []int {
	sum := Sum(arr)
	result := make([]int, len(arr))
	for i, val := range arr {
		ratio := float64(val) / float64(sum)
		result[i] = int(math.Round(ratio * float64(weightTotal)))
	}
	return result
}

func Sum(arr []int) int {
	total := 0
	for _, v := range arr {
		total += v
	}
	return total
}

func CalcPriorityWeight(chromosome utils.Chromosome, cfg utils.BgaEnv) map[string]utils.VMRank {
	// Step 1: Hitung total tugas per VM ID (angka)
	taskCounts := make(map[int]int)
	for _, vmID := range chromosome.Genes {
		taskCounts[vmID]++
	}

	// Step 2: Ambil nama VM dan urutkan
	var vmNames []string
	for name := range cfg.VMNames {
		vmNames = append(vmNames, name)
	}
	sort.Strings(vmNames)

	// Step 3: Buat slice untuk menyimpan informasi VM
	var infos []utils.VMInfo
	for idx, name := range vmNames {
		vmID := idx + 1
		load := float64(taskCounts[vmID]) * cfg.TaskLoad
		infos = append(infos, utils.VMInfo{
			Name: name,
			Load: load,
			ID:   vmID,
		})
	}

	// Step 4: Urutkan berdasarkan Load descending (terbesar dulu)
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Load > infos[j].Load
	})

	// Step 5: Hitung distribusi bobot (descending)
	n := len(infos)
	base := make([]int, n)
	for i := range base {
		base[i] = i + 1
	}
	weights := DistributeWeights(base, cfg.HAProxyWeight)
	sort.Sort(sort.Reverse(sort.IntSlice(weights))) // urutkan dari terbesar

	// Step 6: Buat map hasil akhir
	result := make(map[string]utils.VMRank)
	for i, vm := range infos {
		priority := i + 1 // posisi dalam slice = priority
		result[vm.Name] = utils.VMRank{
			Value:    vm.Load,
			Priority: priority,
			Weight:   weights[i],
			Fitness:  chromosome.Fitness,
		}
	}

	return result
}
