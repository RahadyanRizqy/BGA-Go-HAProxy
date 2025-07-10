package funcs

import (
	"bga_go_haproxy/utils"
	"math/rand"
)

func Balancer(chromosome *utils.Chromosome, cfg utils.BgaEnv) {
	numVMs := cfg.NumVMs
	taskSizes := make([]float64, len(chromosome.Genes))
	for i := range taskSizes {
		taskSizes[i] = cfg.TaskSize
	}

	// Step 1: Ambil VM GHz dari cfg.VMDetails dan hitung ShareRatios serta VMShare
	vmGHz := make([]float64, numVMs)
	for _, vm := range cfg.VMDetails {
		if vm.ID >= 1 && vm.ID <= numVMs {
			vmGHz[vm.ID-1] = vm.GHz
		}
	}

	shareRatios := CalcShareRatios(vmGHz)
	totalTaskSize := float64(len(taskSizes)) * cfg.TaskSize
	vmShares := CalcVMShare(shareRatios, totalTaskSize)

	// Step 2: Hitung beban aktual per VM (ShareUsed)
	shareUsed := CalcShareUsed(taskSizes, chromosome.Genes, numVMs)

	// Step 3: Identifikasi VM overutilized dan underutilized
	mostOverUtilizedVMIndex := -1
	maxOverUtilization := -1.0
	mostUnderUtilizedVMIndex := -1
	maxUnderUtilization := -1.0

	tasksOnVM := make([][]int, numVMs)
	for i := range tasksOnVM {
		tasksOnVM[i] = []int{}
	}
	for idx, vmID := range chromosome.Genes {
		if vmID >= 1 && vmID <= numVMs {
			tasksOnVM[vmID-1] = append(tasksOnVM[vmID-1], idx)
		}
	}

	for i := 0; i < numVMs; i++ {
		over := shareUsed[i] - vmShares[i]
		under := vmShares[i] - shareUsed[i]
		if over > 0 && (mostOverUtilizedVMIndex == -1 || over > maxOverUtilization) {
			mostOverUtilizedVMIndex = i
			maxOverUtilization = over
		}
		if under > 0 && (mostUnderUtilizedVMIndex == -1 || under > maxUnderUtilization) {
			mostUnderUtilizedVMIndex = i
			maxUnderUtilization = under
		}
	}

	// Step 4: Pindahkan satu task dari over → under
	if mostOverUtilizedVMIndex != -1 && mostUnderUtilizedVMIndex != -1 && len(tasksOnVM[mostOverUtilizedVMIndex]) > 0 {
		taskToMoveIdx := tasksOnVM[mostOverUtilizedVMIndex][rand.Intn(len(tasksOnVM[mostOverUtilizedVMIndex]))]
		chromosome.Genes[taskToMoveIdx] = mostUnderUtilizedVMIndex + 1
	}
}
