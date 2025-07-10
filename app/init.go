package app

import (
	"bga_go_haproxy/funcs"
	"bga_go_haproxy/utils"
	"crypto/tls"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"
)

var (
	cfg = utils.LoadBgaEnv()
	// vmShareIdeal   = float64(cfg.NumTasks / cfg.NumVMs)
	prevStats      = make(map[string]utils.VM)
	prevScores     = make(map[string]float64)
	prevWeights    = make(map[string]int)
	activeRates    = make(map[string]utils.ActiveRates)
	lastValidRates = make(map[string]utils.ActiveRates)
	client         *http.Client
	fetchCount     int
	updateCount    int
	logLine        int = 1
	validate       bool
	mode           string
)

func InitClient() {
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func Start() {
	fmt.Println("BGA Started!")

	/*
		Random Seed Initialization
	*/
	funcs.SeedInit()
	population := funcs.PopulationInit(cfg)

	/*
		Fitness Calculation
	*/
	for i := 0; i < cfg.PopulationSize; i++ {
		funcs.FitnessCalc(&population[i], cfg)
	}

	/*
		Initialization of HTTP Client for FetchStats()
	*/
	InitClient()
	csvFileName := utils.InitCSV(cfg)
	prevTime := time.Now()

	iter := 1
	for {
		time.Sleep(time.Duration(cfg.GenerateDelay) * time.Millisecond)
		now := time.Now()
		delta := now.Sub(prevTime).Seconds()
		fetchCount++

		/*
			FetchStats() to fetch VM stats from Proxmox VE API for logging ONLY
		*/
		stats, err := funcs.FetchStats(cfg, client)
		if err != nil {
			fmt.Printf("Polling error: %v\n", err)
			continue
		}

		/*
			CSV Logging only not related to the main algorithm
		*/
		validVMs := make(map[string]bool)
		for _, vm := range cfg.VMDetails {
			validVMs[vm.Name] = true
		}

		currentStats := make(map[string]utils.VMStats)
		for _, vm := range stats {
			if !validVMs[vm.Name] {
				continue
			}
			currentStats[vm.Name] = funcs.PreviousStats(vm, delta, cfg.NetIfaceRate, lastValidRates, prevStats, activeRates)
		}

		/*
			Sort Generated Population by its Fitness
		*/
		sort.Slice(population, func(i, j int) bool { return population[i].Fitness < population[j].Fitness })

		/*
			Store its current best
		*/
		currentBest := utils.Chromosome{Genes: make([]int, cfg.NumTasks), Fitness: population[0].Fitness}
		copy(currentBest.Genes, population[0].Genes) // <- COPY

		/*
			Current Result for Weight Assignment
		*/
		currentRes := funcs.CalcPriorityWeight(currentBest, cfg)

		/*
			Strict or Loose
		*/
		if cfg.Strict {
			validate = funcs.AllWeightValidation(currentRes, prevWeights)
			mode = "STRICT"
		} else {
			validate = funcs.SomeWeightValidation(currentRes, prevWeights)
			mode = "LOOSE"
		}

		if validate {
			updateCount++
			if cfg.UpdateNotify {
				fmt.Printf("✅ [%s] UPDATE COUNT %d ITER COUNT %d\n", mode, updateCount, iter)
			}
			funcs.SetWeight(currentRes, cfg)
			utils.ConsolePrint(currentRes, cfg)
			for name, info := range currentRes {
				prevWeights[name] = info.Weight // update previous
			}
		}

		/*
			Logger
		*/
		utils.StoreCSV(
			cfg,
			csvFileName,
			&logLine,
			fetchCount,
			updateCount,
			now.Unix(),
			now.Format("2006-01-02 15:04:05"),
			currentStats,
			currentRes,
			cfg.NetIfaceRate)

		/*
			Update Previous VM State for logging purpose not related to the main algorithm
		*/
		funcs.UpdatePreviousState(prevStats, prevScores, currentStats)

		/* Looping Mechanism */
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

			parent1 := funcs.Selection(population, cfg)
			parent2 := funcs.Selection(population, cfg)

			var child1, child2 utils.Chromosome
			if opCount < numSinglePointOps {
				child1, child2 = funcs.CrossoverSinglePoint(parent1, parent2, cfg)
			} else {
				child1, child2 = funcs.CrossoverTwoPoint(parent1, parent2, cfg)
			}
			funcs.Mutation(&child1, cfg)
			funcs.Mutation(&child2, cfg)

			if cfg.Balancer {
				funcs.Balancer(&child1, cfg)
				funcs.Balancer(&child2, cfg)
			}

			funcs.FitnessCalc(&child1, cfg)
			funcs.FitnessCalc(&child2, cfg)

			newPopulation[newChildIndex] = child1
			newChildIndex++
			newPopulation[newChildIndex] = child2
			newChildIndex++
		}
		fmt.Println(population)
		population = newPopulation // Population modified to be used later again as the currentBest
		prevTime = now
		iter++
	}
}
