package main

import (
	"context"
	"log"

	"github.com/dimcz/dino-go/game"
	"github.com/yaricom/goNEAT/v3/experiment"
	"github.com/yaricom/goNEAT/v3/neat"
	"github.com/yaricom/goNEAT/v3/neat/genetics"
)

type DinoEvaluator struct {
}

func NewDinoEvaluator() *DinoEvaluator {
	return &DinoEvaluator{}
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
		if org.Fitness <= 0 {
			continue
		}

		if epoch.Champion == nil || org.Fitness > epoch.Champion.Fitness {
			epoch.Solved = true
			epoch.WinnerNodes = len(org.Genotype.Nodes)
			epoch.WinnerGenes = org.Genotype.Extrons()
			epoch.WinnerEvals = options.PopSize*epoch.Id + org.Genotype.Id
			epoch.Champion = org
		}
	}

	epoch.FillPopulationStatistics(pop)

	return nil
}
