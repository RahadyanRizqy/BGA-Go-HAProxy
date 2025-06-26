package app

import (
	"bga_go_haproxy/funcs"
	"bga_go_haproxy/utils"
	"crypto/tls"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

var (
	cfg            = utils.LoadBgaEnv()
	vmShareIdeal   = float64(cfg.NumTasks / cfg.NumVMs)
	prevStats      = make(map[string]utils.VM)
	prevScores     = make(map[string]float64)
	activeRates    = make(map[string]utils.ActiveRates)
	lastValidRates = make(map[string]utils.ActiveRates)
	client         *http.Client
	fetchCount     int
	changeCount    int = 0
	logLine        int = 1
)

func InitClient() {
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func seedInit() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandomChromosome() utils.Chromosome {
	genes := make([]int, cfg.NumTasks)
	for i := range genes {
		genes[i] = rand.Intn(cfg.NumVMs) + 1
	}
	return utils.Chromosome{Genes: genes}
}

func populationInit() []utils.Chromosome {
	population := make([]utils.Chromosome, cfg.PopulationSize)
	for i := range population {
		population[i] = generateRandomChromosome()
	}
	return population
}

func fitnessCalc(k *utils.Chromosome) {
	load := make([]float64, cfg.NumVMs)
	for _, vm := range k.Genes {
		load[vm-1] += cfg.TaskLoad
	}

	makespan := 0.0
	for _, l := range load {
		et := l / cfg.TaskLoad
		if et > makespan {
			makespan = et
		}
	}

	psuTotal := 0.0
	for _, l := range load {
		psuRaw := (l / float64(vmShareIdeal)) * 100.0
		if l == 0 && vmShareIdeal == 0 {
			psuRaw = 100
		} else if vmShareIdeal == 0 {
			psuRaw = math.Inf(1)
		}
		psuNorm := psuRaw / 100.0
		if psuRaw > 100.0 {
			psuNorm = (100.0 - (psuRaw - 100.0)) / 100.0
		}
		if psuNorm < 0 {
			psuNorm = 0
		}
		psuTotal += psuNorm
	}

	k.Fitness = makespan + (1.0 - (psuTotal / float64(cfg.NumVMs)))
}

func proportionalSelection(pop []utils.Chromosome) utils.Chromosome {
	fBest := pop[0].Fitness
	for _, k := range pop {
		if k.Fitness < fBest {
			fBest = k.Fitness
		}
	}

	sf := make([]float64, len(pop))
	total := 0.0
	for i, k := range pop {
		s := 1.0 / (cfg.PositiveConst + (k.Fitness - fBest))
		sf[i] = s
		total += s
	}

	r := rand.Float64() * total
	c := 0.0
	for i, s := range sf {
		c += s
		if c >= r {
			res := pop[i]
			copyGenes := make([]int, cfg.NumTasks)
			copy(copyGenes, res.Genes)
			return utils.Chromosome{Genes: copyGenes, Fitness: res.Fitness}
		}
	}
	return pop[len(pop)-1]
}

func crossoverSinglePoint(p1, p2 utils.Chromosome) (utils.Chromosome, utils.Chromosome) {
	a1 := utils.Chromosome{Genes: make([]int, cfg.NumTasks)}
	a2 := utils.Chromosome{Genes: make([]int, cfg.NumTasks)}
	copy(a1.Genes, p1.Genes)
	copy(a2.Genes, p2.Genes)
	point := rand.Intn(cfg.NumTasks-1) + 1
	for i := point; i < cfg.NumTasks; i++ {
		a1.Genes[i], a2.Genes[i] = a2.Genes[i], a1.Genes[i]
	}
	return a1, a2
}

func crossoverTwoPoint(p1, p2 utils.Chromosome) (utils.Chromosome, utils.Chromosome) {
	point1 := rand.Intn(cfg.NumTasks - 1)
	point2 := rand.Intn(cfg.NumTasks-point1-1) + point1 + 1
	a1 := utils.Chromosome{Genes: make([]int, cfg.NumTasks)}
	a2 := utils.Chromosome{Genes: make([]int, cfg.NumTasks)}
	copy(a1.Genes, p1.Genes)
	copy(a2.Genes, p2.Genes)
	for i := point1; i < point2; i++ {
		a1.Genes[i], a2.Genes[i] = a2.Genes[i], a1.Genes[i]
	}
	return a1, a2
}

func mutation(k *utils.Chromosome) {
	for i := range k.Genes {
		if rand.Float64() < cfg.MutationRate {
			k.Genes[i] = rand.Intn(cfg.NumVMs) + 1
		}
	}
}

func Start() {
	fmt.Println("BGA Started!")
	seedInit()
	population := populationInit()
	for i := 0; i < cfg.PopulationSize; i++ {
		fitnessCalc(&population[i])
	}

	// FETCH API //
	InitClient()
	csvFileName := utils.InitCSV(cfg)
	prevTime := time.Now()
	// FETCH API //

	iter := 0
	for {
		fetchCount++
		now := time.Now()
		delta := now.Sub(prevTime).Seconds()
		sort.Slice(population, func(i, j int) bool { return population[i].Fitness < population[j].Fitness })

		currentBest := utils.Chromosome{Genes: make([]int, cfg.NumTasks), Fitness: population[0].Fitness}
		copy(currentBest.Genes, population[0].Genes)

		// utils.ConsolePrint(currentBest, iter, cfg)
		// changeCount, changeStatus := funcs.CheckWeightChange(currentBest, iter, cfg)
		// funcs.WeightChange(currentBest, cfg, changeStatus)

		// fmt.Println(currentBest)
		newRes := funcs.AssignWeightByTaskGenes(currentBest, cfg)
		// for id, values := range newRes {
		// 	fmt.Printf("VM%d, %.2f, %.0f, %.2f\n", id, values[0], values[1], values[2])
		// }
		weightedVMs, status := utils.ConsolePrint(newRes, iter, cfg)
		if status {
			changeCount++
		}
		funcs.ChangeWeight(weightedVMs, status, cfg)

		// -- LOGGER STARTS -- //
		// Process VM stats and calculate metrics
		stats, err := funcs.FetchVMs(cfg, client)
		if err != nil {
			fmt.Printf("Polling error: %v\n", err)
			continue
		}

		// Process VM stats and calculate metrics
		currentStats := make(map[string]utils.VMStats)
		for _, vm := range stats {
			if !cfg.VMNames[vm.Name] {
				continue
			}

			stats := funcs.PreviousStats(vm, delta, cfg.NetIfaceRate, lastValidRates, prevStats, activeRates)
			currentStats[vm.Name] = stats
		}

		rankedVMs := funcs.ConvertRanked(newRes)
		utils.StoreCSV(
			cfg,
			csvFileName,
			&logLine,
			fetchCount,
			changeCount,
			now.Unix(),
			now.Format("2006-01-02 15:04:05"),
			currentStats,
			rankedVMs,
			cfg.NetIfaceRate)

		funcs.UpdatePreviousState(prevStats, prevScores, currentStats)
		prevTime = now

		newPopulation := make([]utils.Chromosome, cfg.PopulationSize)
		for i := 0; i < cfg.NumElites; i++ {
			newPopulation[i].Genes = make([]int, cfg.NumTasks)
			copy(newPopulation[i].Genes, population[i].Genes)
			newPopulation[i].Fitness = population[i].Fitness
		}

		s := (cfg.PopulationSize - cfg.NumElites) / 2
		numSinglePointOps := 0
		if s > 0 {
			numSinglePointOps = int(math.Round(float64(s) * cfg.FixedAlpha))
		}

		newChildIndex := cfg.NumElites
		for opCount := 0; opCount < s; opCount++ {
			if newChildIndex+1 >= cfg.PopulationSize {
				break
			}

			parent1 := proportionalSelection(population)
			parent2 := proportionalSelection(population)

			var child1, child2 utils.Chromosome
			if opCount < numSinglePointOps {
				child1, child2 = crossoverSinglePoint(parent1, parent2)
			} else {
				child1, child2 = crossoverTwoPoint(parent1, parent2)
			}
			mutation(&child1)
			mutation(&child2)
			funcs.ResultBalancer(&child1, cfg, vmShareIdeal)
			funcs.ResultBalancer(&child2, cfg, vmShareIdeal)
			fitnessCalc(&child1)
			fitnessCalc(&child2)
			newPopulation[newChildIndex] = child1
			newChildIndex++
			newPopulation[newChildIndex] = child2
			newChildIndex++
		}
		population = newPopulation
		iter++
		time.Sleep(time.Duration(cfg.GenerateDelay) * time.Millisecond)
	}
}
