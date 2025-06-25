package utils

func VMTaskCalc(c Chromosome, cfg BgaEnv) []int {
	taskPerVM := make([]int, cfg.NumVMs)
	for _, vmID := range c.Genes {
		if vmID >= 1 && vmID <= cfg.NumVMs {
			taskPerVM[vmID-1]++
		}
	}
	return taskPerVM
}
