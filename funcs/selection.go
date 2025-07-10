package funcs

import (
	"bga_go_haproxy/utils"
	"math/rand"
)

func Selection(pop []utils.Chromosome, cfg utils.BgaEnv) utils.Chromosome {
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
