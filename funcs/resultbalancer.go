package funcs

import (
	"bga_go_haproxy/utils"
	"math/rand"
)

func ResultBalancer(chromosome *utils.Chromosome, cfg utils.BgaEnv, vmShareIdeal float64) {
	actualLoadVM := make([]float64, cfg.NumVMs)
	tasksOnVM := make([][]int, cfg.NumVMs)
	for i := 0; i < cfg.NumVMs; i++ {
		tasksOnVM[i] = make([]int, 0)
	}
	for taskIdx, vmID := range chromosome.Genes {
		if vmID >= 1 && vmID <= cfg.NumVMs {
			actualLoadVM[vmID-1] += cfg.TaskSize
			tasksOnVM[vmID-1] = append(tasksOnVM[vmID-1], taskIdx)
		}
	}
	mostOverUtilizedVMIndex := -1
	maxOverUtilization := -1.0
	mostUnderUtilizedVMIndex := -1
	maxUnderUtilizationCapacity := -1.0
	vmIndicesShuffled := make([]int, cfg.NumVMs)
	for k := 0; k < cfg.NumVMs; k++ {
		vmIndicesShuffled[k] = k
	}
	rand.Shuffle(len(vmIndicesShuffled), func(a, b int) {
		vmIndicesShuffled[a], vmIndicesShuffled[b] = vmIndicesShuffled[b], vmIndicesShuffled[a]
	})
	for _, vmIdx := range vmIndicesShuffled {
		overUtilization := actualLoadVM[vmIdx] - vmShareIdeal
		if overUtilization > 0 && (mostOverUtilizedVMIndex == -1 || overUtilization > maxOverUtilization) {
			maxOverUtilization = overUtilization
			mostOverUtilizedVMIndex = vmIdx
		}
		underUtilizationCapacity := vmShareIdeal - actualLoadVM[vmIdx]
		if underUtilizationCapacity > 0 && (mostUnderUtilizedVMIndex == -1 || underUtilizationCapacity > maxUnderUtilizationCapacity) {
			maxUnderUtilizationCapacity = underUtilizationCapacity
			mostUnderUtilizedVMIndex = vmIdx
		}
	}
	if mostOverUtilizedVMIndex != -1 && mostUnderUtilizedVMIndex != -1 &&
		len(tasksOnVM[mostOverUtilizedVMIndex]) > 0 {
		taskToMoveIdxOriginal := tasksOnVM[mostOverUtilizedVMIndex][rand.Intn(len(tasksOnVM[mostOverUtilizedVMIndex]))]
		chromosome.Genes[taskToMoveIdxOriginal] = mostUnderUtilizedVMIndex + 1
	}
}
