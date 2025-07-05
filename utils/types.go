package utils

type VM struct {
	Id        int     `json:"vmid"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	MaxMem    int     `json:"maxmem"`
	MaxCPU    int     `json:"maxcpu"`
	Mem       float64 `json:"mem"`
	CPU       float64 `json:"cpu"`
	CumNetIn  int     `json:"netin"`
	CumNetOut int     `json:"netout"`
}

type Response struct {
	Data []VM `json:"data"`
}

type VMPriority struct {
	Value    float64
	Priority int
	Weight   int
}

type VMRank struct {
	Value    float64
	Priority int
	Weight   int
	Fitness  float64
}

type VMInfo struct {
	Name string
	Load float64
	ID   int // numerik ID (1-based)
}

type ActiveRates struct {
	Rx float64 // Receive rate in bytes/sec
	Tx float64 // Transmit rate in bytes/sec
}

type KV struct {
	Key   string
	Value float64
}

type VMStats struct {
	VM       VM
	Score    float64
	Rates    ActiveRates
	BwUsage  float64
	MemUsage float64
}

type VMWeight struct {
	Name      string
	TaskTotal float64
	Weight    float64
}

type Chromosome struct {
	Genes   []int
	Fitness float64
}

type BgaEnv struct {
	APIToken       string
	PveAPIURL      string
	HAProxySock    string
	HAProxyBackend string
	VMNames        map[string]bool
	NetIfaceRate   float64
	BgaUpdater     bool
	HAProxyWeight  int
	Logger         bool
	ConsolePrint   bool
	NumTasks       int
	NumVMs         int // VMShareIdeal from NumTasks/NumVMs
	PopulationSize int
	NumElites      int
	MutationRate   float64
	FixedAlpha     float64
	GenerateDelay  int
	FetchDelay     int
	TaskLoad       float64 // TaskLoad = VMPower
	PositiveConst  float64
	Strict         bool
	UpdateNotify   bool
}
