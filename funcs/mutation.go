package funcs

import (
	"bga_go_haproxy/utils"
	"math/rand"
)

func Mutation(k *utils.Chromosome, cfg utils.BgaEnv) {
	for i := range k.Genes {
		if rand.Float64() < cfg.MutationRate {
			k.Genes[i] = rand.Intn(cfg.NumVMs) + 1
		}
	}
}
