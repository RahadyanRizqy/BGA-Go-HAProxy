package funcs

import (
	"bga_go_haproxy/utils"
	"math"
)

func CalcShareRatios(vmGHz []float64) []float64 {
	total := 0.0
	for _, ghz := range vmGHz {
		total += ghz
	}
	shareRatios := make([]float64, len(vmGHz))
	for i, ghz := range vmGHz {
		shareRatios[i] = ghz / total
	}
	return shareRatios
}

func CalcVMShare(shareRatios []float64, totalTaskSize float64) []float64 {
	vmShares := make([]float64, len(shareRatios))
	for i, ratio := range shareRatios {
		vmShares[i] = totalTaskSize * ratio
	}
	return vmShares
}

func CalcShareUsed(taskSizes []float64, genes []int, numVMs int) []float64 {
	shareUsed := make([]float64, numVMs)
	for i, vmID := range genes {
		if vmID >= 1 && vmID <= numVMs {
			shareUsed[vmID-1] += taskSizes[i]
		}
	}
	return shareUsed
}

func FitnessCalc(k *utils.Chromosome, cfg utils.BgaEnv) {
	numVMs := cfg.NumVMs
	taskSizes := make([]float64, len(k.Genes))

	// Semua task diasumsikan 1.0 (konstan)
	for i := range taskSizes {
		taskSizes[i] = cfg.TaskSize
	}

	// Ambil list GHz dari cfg.VMDetails
	vmGHz := make([]float64, numVMs)
	for i, vm := range cfg.VMDetails {
		vmGHz[i] = vm.GHz
	}

	// Step 1: Hitung ShareRatio, VMShare, dan ShareUsed
	shareRatios := CalcShareRatios(vmGHz)
	totalTaskSize := float64(len(taskSizes)) * cfg.TaskSize
	vmShares := CalcVMShare(shareRatios, totalTaskSize)
	shareUsed := CalcShareUsed(taskSizes, k.Genes, numVMs)

	// Step 2: Hitung Makespan
	makespan := 0.0
	for _, used := range shareUsed {
		if used > makespan {
			makespan = used
		}
	}

	// Step 3: Hitung PSU dan PSU Norm
	psuTotal := 0.0
	for i := 0; i < numVMs; i++ {
		vmShare := vmShares[i]
		used := shareUsed[i]

		var psuRaw float64
		if used == 0 && vmShare == 0 {
			psuRaw = 100
		} else if vmShare == 0 {
			psuRaw = math.Inf(1)
		} else {
			psuRaw = (used / vmShare) * 100
		}

		psuNorm := psuRaw / 100
		if psuRaw > 100 {
			psuNorm = (100 - (psuRaw - 100)) / 100
		}
		if psuNorm < 0 {
			psuNorm = 0
		}
		psuTotal += psuNorm
	}

	avgPSU := psuTotal / float64(numVMs)
	k.Fitness = makespan + (1.0 - avgPSU)
}
