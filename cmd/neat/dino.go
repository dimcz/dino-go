package main

import (
	"context"

	"github.com/yaricom/goNEAT/v3/experiment"
	"github.com/yaricom/goNEAT/v3/neat/genetics"
)

type DinoEvaluator struct {
}

func NewDinoEvaluator() *DinoEvaluator {
	return &DinoEvaluator{}
}

func (d *DinoEvaluator) GenerationEvaluate(ctx context.Context, pop *genetics.Population, epoch *experiment.Generation) error {
	panic("implement me")
}
