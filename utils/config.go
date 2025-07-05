package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func parseVMMap(env string) map[string]bool {
	result := make(map[string]bool)
	for _, vm := range strings.Split(env, ",") {
		if trimmed := strings.TrimSpace(vm); trimmed != "" {
			result[trimmed] = true
		}
	}
	return result
}

func LoadBgaEnv() BgaEnv {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	bgaUpdater, err := strconv.ParseBool(os.Getenv("BGA_UPDATER"))
	if err != nil {
		fmt.Println("1. Error parsing boolean:", err)
		bgaUpdater = false
	}

	netIfaceRate, err := strconv.ParseFloat(os.Getenv("NET_IFACE_RATE"), 64)
	if err != nil {
		netIfaceRate = 12500000
	}

	haproxyWeight, err := strconv.Atoi(os.Getenv("HAPROXY_WEIGHT"))
	if err != nil {
		fmt.Println("2. Error parsing boolean:", err)
		haproxyWeight = 256
	}

	logger, err := strconv.ParseBool(os.Getenv("LOGGER"))
	if err != nil {
		fmt.Println("3. Error parsing boolean:", err)
		logger = false
	}

	consolePrint, err := strconv.ParseBool(os.Getenv("CONSOLE_PRINT"))
	if err != nil {
		fmt.Println("4. Error parsing boolean:", err)
		consolePrint = false
	}

	numTasks, err := strconv.Atoi(os.Getenv("NUM_TASKS"))
	if err != nil {
		fmt.Println("5. Error parsing boolean:", err)
		numTasks = 100000
	}

	numVMs, err := strconv.Atoi(os.Getenv("NUM_VMS"))
	if err != nil {
		fmt.Println("6. Error parsing boolean:", err)
		numVMs = 0
	}

	populationSize, err := strconv.Atoi(os.Getenv("POPULATION_SIZE"))
	if err != nil {
		fmt.Println("7. Error parsing boolean", err)
		populationSize = 100
	}

	numElites, err := strconv.Atoi(os.Getenv("NUM_ELITES"))
	if err != nil {
		fmt.Println("8. Error parsing boolean", err)
		numElites = 2
	}

	mutationRate, err := strconv.ParseFloat(os.Getenv("MUTATION_RATE"), 64)
	if err != nil {
		fmt.Println("9. Error parsing boolean", err)
		mutationRate = 0.5
	}

	fixedAlpha, err := strconv.ParseFloat(os.Getenv("FIXED_ALPHA"), 64)
	if err != nil {
		fmt.Println("10. Error parsing boolean", err)
		fixedAlpha = 0.2
	}

	generateDelay, err := strconv.Atoi(os.Getenv("GENERATE_DELAY"))
	if err != nil {
		fmt.Println("11. Error parsing boolean", err)
		generateDelay = 1000
	}

	taskLoad, err := strconv.ParseFloat(os.Getenv("TASK_LOAD"), 64)
	if err != nil {
		fmt.Println("12. Error parsing boolean", err)
		taskLoad = 1.0
	}

	positiveConst, err := strconv.ParseFloat(os.Getenv("POSITIVE_CONST"), 64)
	if err != nil {
		fmt.Println("13. Error parsing boolean", err)
		positiveConst = 0.00001
	}

	strict, err := strconv.ParseBool(os.Getenv("STRICT"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		strict = false
	}

	updateNotify, err := strconv.ParseBool(os.Getenv("UPDATE_NOTIFY"))
	if err != nil {
		fmt.Println("Error parsing boolean:", err)
		updateNotify = false
	}

	return BgaEnv{
		APIToken:       os.Getenv("API_TOKEN"),   // for logging purpose
		PveAPIURL:      os.Getenv("PVE_API_URL"), // for logging purpose
		HAProxySock:    os.Getenv("HAPROXY_SOCK"),
		HAProxyBackend: os.Getenv("HAPROXY_BACKEND"),
		VMNames:        parseVMMap(os.Getenv("VM_NAMES")),
		NetIfaceRate:   netIfaceRate,
		BgaUpdater:     bgaUpdater,
		UpdateNotify:   updateNotify,
		HAProxyWeight:  haproxyWeight,
		Logger:         logger,
		ConsolePrint:   consolePrint,
		NumTasks:       numTasks,
		NumVMs:         numVMs,
		PopulationSize: populationSize,
		NumElites:      numElites,
		MutationRate:   mutationRate,
		FixedAlpha:     fixedAlpha,
		GenerateDelay:  generateDelay,
		TaskLoad:       taskLoad,
		PositiveConst:  positiveConst,
		Strict:         strict,
	}
}
