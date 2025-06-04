package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/nsga2"
	"github.com/GregoryKogan/genetic-algorithms/pkg/algos/ssga"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/crossover"
	"github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/mutation"
)

type Test struct {
	Repeat   int
	Name     string
	Planar   bool
	Vertexes int
}

func main() {
	tests := []Test{
		{Repeat: 30, Name: "SSGA-FR P25", Planar: true, Vertexes: 25},
		{Repeat: 30, Name: "FR-NSGA2 P25", Planar: true, Vertexes: 25},
		{Repeat: 30, Name: "FR-SSGA-NSGA2 P25", Planar: true, Vertexes: 25},

		{Repeat: 9, Name: "SSGA-FR A25", Planar: false, Vertexes: 25},
		{Repeat: 9, Name: "FR-NSGA2 A25", Planar: false, Vertexes: 25},
		{Repeat: 9, Name: "FR-SSGA-NSGA2 A25", Planar: false, Vertexes: 25},

		{Repeat: 30, Name: "SSGA-FR P50", Planar: true, Vertexes: 50},
		{Repeat: 30, Name: "FR-NSGA2 P50", Planar: true, Vertexes: 50},
		{Repeat: 30, Name: "FR-SSGA-NSGA2 P50", Planar: true, Vertexes: 50},

		{Repeat: 9, Name: "SSGA-FR A50", Planar: false, Vertexes: 50},
		{Repeat: 9, Name: "FR-NSGA2 A50", Planar: false, Vertexes: 50},
		{Repeat: 9, Name: "FR-SSGA-NSGA2 A50", Planar: false, Vertexes: 50},

		{Repeat: 9, Name: "SSGA-FR P100", Planar: true, Vertexes: 100},
		{Repeat: 9, Name: "FR-NSGA2 P100", Planar: true, Vertexes: 100},
		{Repeat: 9, Name: "FR-SSGA-NSGA2 P100", Planar: true, Vertexes: 100},

		{Repeat: 9, Name: "SSGA-FR A100", Planar: false, Vertexes: 100},
		{Repeat: 9, Name: "FR-NSGA2 A100", Planar: false, Vertexes: 100},
		{Repeat: 9, Name: "FR-SSGA-NSGA2 A100", Planar: false, Vertexes: 100},

		{Repeat: 9, Name: "SSGA-FR P200", Planar: true, Vertexes: 200},
		{Repeat: 9, Name: "FR-NSGA2 P200", Planar: true, Vertexes: 200},
		{Repeat: 9, Name: "FR-SSGA-NSGA2 P200", Planar: true, Vertexes: 200},

		{Repeat: 9, Name: "SSGA-FR A200", Planar: false, Vertexes: 200},
		{Repeat: 9, Name: "FR-NSGA2 A200", Planar: false, Vertexes: 200},
		{Repeat: 9, Name: "FR-SSGA-NSGA2 A200", Planar: false, Vertexes: 200},
	}

	// initLogsDir()
	for _, test := range tests {
		totalIntersections := 0
		zeroCounter := 0
		totalFitness := 0.0
		for range test.Repeat {
			var problem problems.Problem
			if test.Planar {
				problem = graphplane.NewPlanarGraphPlaneProblem(test.Vertexes)
			} else {
				problem = graphplane.NewGraphPlaneProblem(test.Vertexes, test.Vertexes*3)
			}

			var intersections int
			var fitness float64
			algoName := strings.Split(test.Name, " ")[0]
			if algoName == "SSGA-FR" {
				intersections, fitness = ssga_fr(problem, test.Planar, test.Vertexes)
			} else if algoName == "FR-NSGA2" {
				intersections, fitness = fr_nsga2(problem, test.Planar, test.Vertexes)
			} else if algoName == "FR-SSGA-NSGA2" {
				intersections, fitness = fr_ssga_nsga2(problem, test.Planar, test.Vertexes)
			}

			totalFitness += fitness
			totalIntersections += intersections
			if intersections == 0 {
				zeroCounter++
			}
		}
		fmt.Printf(
			"%s - Avg I: %.2f, Avg F: %.2f, 0-rate: %.2f%%\n",
			test.Name,
			float64(totalIntersections)/float64(test.Repeat),
			totalFitness/float64(test.Repeat),
			float64(zeroCounter)/float64(test.Repeat)*100.0,
		)
	}
}

// func initLogger(problem problems.Problem, method string) algos.ProgressLoggerProvider {
//   logPath := filepath.Join("logs", fmt.Sprintf("%s_%s.jsonl", problem.Name(), method))
//   logger := algos.NewProgressLogger(logPath)
//   logger.InitLogging()
//   logger.LogProblem(problem)
//   return logger
// }

// func initLogsDir() {
//   os.RemoveAll("logs")
//   os.Mkdir("logs", 0755)
// }

func ssga_fr(p problems.Problem, planar bool, vertexes int) (int, float64) {
	gaParams := ssga.Params{
		PopulationSize: 500,
		CrossoverFunc:  crossover.Uniform(0.4),
		MutationFunc:   mutation.ConservativeNorm(0.1),
	}
	ga := ssga.NewAlgorithm(p, gaParams, 100000, nil)
	ga.Run(context.Background())

	fdsParams := graphplane.FDSParams{
		Steps: 2000,
		Temp:  0.005,
		K:     getKCoefficient(planar, vertexes),
	}
	fdSolver := graphplane.NewForceDirectedSolver(ga.GetSolution(), fdsParams, nil)
	fdSol, _ := fdSolver.Solve().Solution.(*graphplane.GraphPlaneSolution)
	return fdSol.CountIntersections(), fdSol.Fitness()
}

func fr_nsga2(p problems.Problem, planar bool, vertexes int) (int, float64) {
	fdsParams := graphplane.FDSParams{
		Steps: 2000,
		Temp:  0.005,
		K:     getKCoefficient(planar, vertexes),
	}
	fdSolver := graphplane.NewForceDirectedSolver(p.RandomSolution(), fdsParams, nil)
	fdSol := fdSolver.Solve()

	gaParams := nsga2.Params{
		PopulationSize: 500,
		CrossoverFunc:  crossover.Uniform(0.4),
		MutationFunc:   mutation.ConservativeNorm(0.1),
	}
	ga := nsga2.NewAlgorithm(p, gaParams, 350, nil)
	ga.Seed(fdSol.Solution)
	ga.Run(context.Background())
	gaSol, _ := ga.GetSolution().(*graphplane.GraphPlaneSolution)
	return gaSol.CountIntersections(), gaSol.Fitness()
}

func fr_ssga_nsga2(p problems.Problem, planar bool, vertexes int) (int, float64) {
	fdsParams := graphplane.FDSParams{
		Steps: 2000,
		Temp:  0.005,
		K:     getKCoefficient(planar, vertexes),
	}
	fdSolver := graphplane.NewForceDirectedSolver(p.RandomSolution(), fdsParams, nil)
	fdSol := fdSolver.Solve()

	ssgaParams := ssga.Params{
		PopulationSize: 500,
		CrossoverFunc:  crossover.Uniform(0.4),
		MutationFunc:   mutation.ConservativeNorm(0.1),
	}
	ss := ssga.NewAlgorithm(p, ssgaParams, 70000, nil)
	ss.Seed(fdSol.Solution)
	ss.Run(context.Background())

	nsga2Params := nsga2.Params{
		PopulationSize: 500,
		CrossoverFunc:  crossover.Uniform(0.4),
		MutationFunc:   mutation.ConservativeNorm(0.1),
	}
	ns2 := nsga2.NewAlgorithm(p, nsga2Params, 250, nil)
	ns2.SetPopulation(ss.GetPopulation())
	ns2.Run(context.Background())
	ns2Sol, _ := ns2.GetSolution().(*graphplane.GraphPlaneSolution)
	return ns2Sol.CountIntersections(), ns2Sol.Fitness()
}

func getKCoefficient(planar bool, vertexes int) float64 {
	if planar {
		if vertexes == 25 {
			return 0.7
		} else if vertexes == 50 {
			return 0.6
		} else if vertexes == 100 {
			return 0.5
		} else if vertexes == 200 {
			return 0.4
		}
	} else {
		if vertexes == 25 {
			return 1.0
		} else if vertexes == 50 {
			return 0.95
		} else if vertexes == 100 {
			return 0.90
		} else if vertexes == 200 {
			return 0.85
		}
	}

	panic("Bad number of vertexes")
}
