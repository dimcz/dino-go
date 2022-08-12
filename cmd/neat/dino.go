package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dimcz/dino-go/game"
	"github.com/yaricom/goNEAT/v3/experiment"
	"github.com/yaricom/goNEAT/v3/neat"
	"github.com/yaricom/goNEAT/v3/neat/genetics"
)

type DinoEvaluator struct {
	path string
}

func NewDinoEvaluator(path string) *DinoEvaluator {
	return &DinoEvaluator{path}
}

func (d *DinoEvaluator) GenerationEvaluate(ctx context.Context, pop *genetics.Population, epoch *experiment.Generation) error {
	options, ok := neat.FromContext(ctx)
	if !ok {
		return neat.ErrNEATOptionsNotFound
	}

	g, err := game.NewNEATGame(pop.Organisms, epoch.Id)
	if err != nil {
		log.Fatal(err)
	}
	g.Run()

	for _, org := range pop.Organisms {
		if epoch.Champion == nil || org.Fitness > epoch.Champion.Fitness {
			epoch.WinnerNodes = len(org.Genotype.Nodes)
			epoch.WinnerGenes = org.Genotype.Extrons()
			epoch.WinnerEvals = options.PopSize*epoch.Id + org.Genotype.Id
			epoch.Champion = org
		}
	}

	if g.Stop && len(d.path) > 0 {
		if file, err := os.Create(d.path); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to create file, reason: %s\n", err))
		} else if err = epoch.Champion.Genotype.Write(file); err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump champion genome, reason: %s\n", err))
		}
		epoch.Solved = true
	}

	epoch.FillPopulationStatistics(pop)

	return nil
}
