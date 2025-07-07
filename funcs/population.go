package funcs

import (
	"bga_go_haproxy/utils"
	"math/rand"
	"time"
)

func SeedInit() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateChromosome(cfg utils.BgaEnv) utils.Chromosome {
	genes := make([]int, cfg.NumTasks)
	for i := range genes {
		genes[i] = rand.Intn(cfg.NumVMs) + 1
	}
	return utils.Chromosome{Genes: genes}
}

func PopulationInit(cfg utils.BgaEnv) []utils.Chromosome {
	population := make([]utils.Chromosome, cfg.PopulationSize)
	for i := range population {
		population[i] = GenerateChromosome(cfg)
	}
	return population
}
