package utils

import (
	"fmt"
	"reflect"
)

// CompareAndPrintDiff will check if the individual has changed after balancing
func PrintDiffMark(before Chromosome, after Chromosome, label string) {
	if reflect.DeepEqual(before.Genes, after.Genes) {
		fmt.Printf("%s: ❌ TIDAK BERUBAH\n", label) // Tidak berubah
		// fmt.Println("Before :", before)
		// fmt.Println("After  :", after)
	} else {
		fmt.Printf("%s: ✅ BERUBAH\n", label) // Berubah
		// fmt.Println("Before :", before)
		// fmt.Println("After  :", after)
	}
}
