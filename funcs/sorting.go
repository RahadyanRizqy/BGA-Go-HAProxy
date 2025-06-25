package funcs

import (
	"bga_go_haproxy/utils"
	"math"
	"sort"
)

func ScorePriority(stats map[string]utils.VMStats) map[string]utils.VMPriority {
	var sorted []utils.KV
	for name, stat := range stats {
		sorted = append(sorted, utils.KV{Key: name, Value: stat.Score})
	}

	// Urutkan berdasarkan nilai Score dari kecil ke besar
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value < sorted[j].Value
	})

	result := make(map[string]utils.VMPriority)
	totalVMs := len(sorted)
	totalWeight := 256
	totalPrioritySum := totalVMs * (totalVMs + 1) / 2 // jumlah 1 + 2 + ... + n

	for i, item := range sorted {
		priority := i + 1
		weight := int(math.Round(float64(totalWeight) * float64(priority) / float64(totalPrioritySum)))

		result[item.Key] = utils.VMPriority{
			Value:    item.Value,
			Priority: priority,
			Weight:   weight,
		}
	}

	return result
}

func ConvertRanked(result map[string][3]float64) map[string]utils.VMPriority {
	ranked := make(map[string]utils.VMPriority)

	// Buat slice untuk sorting berdasarkan nilai task (Value)
	type kv struct {
		Name  string
		Value float64
	}
	var sorted []kv
	for name, vals := range result {
		sorted = append(sorted, kv{name, vals[0]}) // vals[0] = task total
	}

	// Urutkan ascending → prioritas rendah untuk task kecil
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value < sorted[j].Value
	})

	for i, item := range sorted {
		ranked[item.Name] = utils.VMPriority{
			Value:    item.Value,
			Priority: i + 1,                     // 1 = paling kecil
			Weight:   int(result[item.Name][1]), // dari hasil weight
		}
	}

	return ranked
}
