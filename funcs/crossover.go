package funcs

import (
	"bga_go_haproxy/utils"
	"math/rand"
)

func CrossoverSinglePoint(p1, p2 utils.Chromosome, cfg utils.BgaEnv) (utils.Chromosome, utils.Chromosome) {
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

func CrossoverTwoPoint(p1, p2 utils.Chromosome, cfg utils.BgaEnv) (utils.Chromosome, utils.Chromosome) {
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
