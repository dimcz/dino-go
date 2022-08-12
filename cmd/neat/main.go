package main

import (
	"flag"
	"log"
	"os"

	"github.com/yaricom/goNEAT/v3/experiment"
	"github.com/yaricom/goNEAT/v3/neat"
	"github.com/yaricom/goNEAT/v3/neat/genetics"
)

const NEAT_CONFIG = "assets/neat/neat.yaml"
const NEAT_GENOME = "assets/neat/genome.config"
const GENOME_DUMP = "assets/neat/last.config"

func main() {
	var genomePath string
	flag.StringVar(&genomePath, "genome", "", "genome file")
	flag.Parse()

	if len(genomePath) == 0 {
		genomePath = NEAT_GENOME
	}

	fOpts, err := os.Open(NEAT_CONFIG)
	if err != nil {
		log.Fatal(err)
	}

	defer func(fOpts *os.File) {
		_ = fOpts.Close()
	}(fOpts)

	opts, err := neat.LoadYAMLOptions(fOpts)
	if err != nil {
		log.Fatal(err)
	}

	fGenome, err := os.Open(NEAT_GENOME)
	if err != nil {
		log.Fatal(err)
	}

	defer func(fGenome *os.File) {
		_ = fGenome.Close()
	}(fGenome)

	reader, err := genetics.NewGenomeReaderFromFile(genomePath)
	if err != nil {
		log.Fatal(err)
	}

	genome, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	exp := experiment.Experiment{
		Id:     0,
		Trials: make(experiment.Trials, opts.NumRuns),
	}

	err = exp.Execute(opts.NeatContext(), genome, NewDinoEvaluator(GENOME_DUMP), nil)
	if err != nil {
		log.Fatal(err)
	}
}
