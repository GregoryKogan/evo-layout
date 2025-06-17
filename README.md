# EvoLayout: A Go Library for Hybrid Genetic Algorithms

A high-performance, modular Go library for exploring and applying genetic algorithms, with a special focus on multi-objective graph layout optimization. This project is the artifact of the research paper "Development of a system implementing a genetic algorithm and its application to the arrangement of graph vertices on a plane."

<p align="center">
  <img alt="Animation showing the layout process for a 200-vertex planar graph using the FR-NSGA2 hybrid algorithm" src="https://github.com/GregoryKogan/GregoryKogan/blob/4b380cdd26bf43603a0521440122e2a3f614014a/readme_assets/gp-200-planar-FR-NSGA2.gif" />
  
  <i align="center">Animation showing the layout process for a 200-vertex planar graph using the FR-NSGA2 hybrid algorithm</i>
</p>

---

## Table of Contents

- [About The Project](#about-the-project)
- [Key Features](#sparkles-key-features)
- [Getting Started](#rocket-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Quick Start](#runner-quick-start)
- [Project Structure](#open_file_folder-project-structure)
- [Implemented Algorithms & Problems](#wrench-implemented-algorithms--problems)
  - [Genetic Algorithms](#genetic-algorithms)
  - [Genetic Operators](#genetic-operators)
  - [Optimization Problems](#optimization-problems)
- [Key Results & Visualizations](#bar_chart-key-results--visualizations)
- [License](#scroll-license)
- [Acknowledgments](#pray-acknowledgments)

## About The Project

Modern optimization problems are often multi-criteria and combinatorially complex, rendering traditional methods inefficient. Graph layout is a prime example: creating a clear, readable visualization of a graph is a multi-objective challenge that involves minimizing edge crossings, ensuring uniform vertex distribution, and optimizing angles.

This library was developed to tackle this challenge by exploring **hybrid genetic algorithms**. It provides a flexible, extensible, and high-performance framework for implementing, testing, and comparing various evolutionary computation techniques. While its primary focus is graph layout, the modular architecture allows it to solve other classic optimization problems like the Traveling Salesperson Problem (TSP) and the Knapsack Problem.

The core of this research demonstrates that hybrid approaches, particularly the custom **FR-NSGA2** algorithm, can significantly outperform both classic genetic algorithms and standalone force-directed methods, especially for large, complex graphs.

## :sparkles: Key Features

- **Hybrid Algorithms:** Implements novel hybrid methods that combine the strengths of force-directed placement (Fruchterman-Reingold) and multi-objective genetic algorithms (NSGA-II).
- **Rich Operator Toolkit:** Provides a comprehensive suite of **12+ specialized genetic operators**, including standard, problem-specific, and adaptive mutation/crossover functions.
- **Modular & Extensible Architecture:** Designed around clean Go interfaces (`Problem`, `Solution`). Easily add new algorithms, problems, or operators without modifying the core library.
- **Classic & Modern GAs:** Includes implementations of SGA, SSGA, NSGA-II, and SPEA2 for comparative analysis.
- **Built-in Problem Suite:** Comes with ready-to-use implementations for Graph Layout, Traveling Salesperson Problem (TSP), 0-1 Knapsack Problem, and the ZDT test suite for multi-objective optimization.
- **Performance-Oriented:** Written in Go for high performance and low memory footprint, capable of handling graphs with hundreds of vertices efficiently.

## :rocket: Getting Started

### Prerequisites

- **Go:** Version 1.23 or later.

### Installation

To add the library to your project, use `go get`:

```bash
go get github.com/GregoryKogan/genetic-algorithms
```

## :runner: Quick Start

Here's a simple example of how to use the library to solve a graph layout problem using the high-performing `FR-NSGA2` hybrid method.

```go
package main

import (
 "context"
 "fmt"
 "time"

 "github.com/GregoryKogan/genetic-algorithms/pkg/algos/nsga2"
 "github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane"
 "github.com/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/crossover"
 "github.comcom/GregoryKogan/genetic-algorithms/pkg/problems/graphplane/operators/mutation"
)

func main() {
 // 1. Define the problem: A planar graph with 50 vertices
 problem := graphplane.NewPlanarGraphPlaneProblem(50)

 // 2. Configure the hybrid algorithm FR-NSGA2
 // Phase 1: Force-Directed (Fruchterman-Reingold)
 frParams := graphplane.FDSParams{
  Steps: 2000,
  Temp:  0.005,
  K:     0.6, // Optimal K for this graph size
 }
 // Phase 2: NSGA-II
 nsga2Params := nsga2.Params{
  PopulationSize: 500,
  CrossoverFunc:  crossover.Uniform(0.4),
  MutationFunc:   mutation.ConservativeNorm(0.1),
 }
 nsga2GenerationLimit := 350

 // 3. Run the hybrid algorithm
 fmt.Println("Starting FR-NSGA2 optimization...")
 startTime := time.Now()

 // Run FR phase
 frSolver := graphplane.NewForceDirectedSolver(problem.RandomSolution(), frParams, nil)
 frSolution := frSolver.Solve()

 // Run NSGA-II phase, seeding it with the result from FR
 ga := nsga2.NewAlgorithm(problem, nsga2Params, nsga2GenerationLimit, nil)
 ga.Seed(frSolution.Solution) // Seed with the best solution from the FR phase
 
 // Set a timeout for the optimization process
 ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
 defer cancel()
 ga.Run(ctx)
 
 finalSolution := ga.GetSolution()
 elapsed := time.Since(startTime)

 // 4. Print the results
 fmt.Printf("Optimization finished in %v\n", elapsed)
 fmt.Printf("Final Solution Objectives: %v\n", finalSolution.Objectives())
}
```

## :open_file_folder: Project Structure

The library is organized into a clear, modular structure within the `pkg/` directory.

```plaintext
pkg/
├── algos/                # Core genetic algorithm implementations
│   ├── nsga2/
│   ├── sga/
│   ├── spea2/
│   └── ssga/
├── problems/             # Problem definitions and solutions
│   ├── graphplane/       # Graph Layout problem
│   │   └── operators/    # Specialized crossover and mutation operators
│   ├── knapsack/         # 0-1 Knapsack problem
│   ├── tsp/              # Traveling Salesperson Problem
│   └── zdt/              # ZDT benchmark functions
└── visual/                 # Scripts for generating visualizations
```

- **`pkg/algos`**: Contains the implementations of different genetic algorithms (SGA, NSGA-II, etc.). They all work with the generic `problems.Solution` interface.
- **`pkg/problems`**: Defines the core interfaces (`Problem`, `Solution`) and contains sub-packages for each implemented optimization problem.
- **`cmd/`**: Contains example executables for running experiments.
- **`visual/`**: Contains Python and p5.js scripts used to generate the charts and animations from the research paper.

## :wrench: Implemented Algorithms & Problems

### Genetic Algorithms

| Algorithm | Type | Key Feature |
| :--- | :--- | :--- |
| **SGA** | Single-Objective | Simple, generational model. |
| **SSGA** | Single-Objective | Steady-state model, replaces worst individuals. |
| **NSGA-II** | Multi-Objective | Fast non-dominated sorting and crowding distance. |
| **SPEA2** | Multi-Objective | Strength-based fitness and density estimation. |
| **FR-NSGA2** | Multi-Objective, Hybrid | Uses Force-Directed placement for a fast start, then NSGA-II for refinement. **(Best performer)** |
| **SSGA-FR** | Single-Objective, Hybrid | Uses SSGA for initial layout, then FR for local optimization. |
| **FR-SSGA-NSGA2** | Multi-Objective, Hybrid | A three-phase approach combining all three methods. |

### Genetic Operators

A rich set of crossover and mutation operators is provided, especially for the graph layout problem.

| Operator Type | Name | Description |
| :--- | :--- | :--- |
| **Crossover** | Uniform Crossover | Exchanges genes between parents based on a probability. |
| **Mutation** | Uniform, Normal, Mirror, Percentage | Standard operators for exploration and exploitation. |
| **Fixed Mutation** | Fixed Uniform/Normal/Percentage | **Целенаправленные:** Apply mutation only to vertices involved in edge crossings. |
| **Hybrid Mutation** | Tension Vector (TV), Fixed TV | Incorporates force-directed principles into the mutation step. |
| **Adaptive Mutation** | Adaptive Normal, Conservative Normal | Change behavior based on the state of the solution (e.g., presence of crossings). |

### Optimization Problems

| Problem | Type | Description |
| :--- | :--- | :--- |
| **Graph Layout** | Multi-Objective | Minimize edge crossings, maximize vertex/angle uniformity. |
| **Knapsack (0-1)** | Single-Objective | Maximize total value within resource constraints. |
| **TSP** | Single-Objective | Find the shortest possible route that visits each city once. |
| **ZDT Suite** | Multi-Objective | Benchmark functions (ZDT1, 2, 3, 4, 6) for testing MOEAs. |

## :bar_chart: Key Results & Visualizations

The experimental results demonstrate the significant advantage of the hybrid **FR-NSGA2** algorithm.

<p align="center">
  <img alt="Edge Crossings Comparison Chart" src="https://github.com/user-attachments/assets/cb25d202-04bf-4ec1-84a3-f70521c8e387" />
  
  <i align="center">Comparison of average edge crossings across all algorithms and graph sizes. Lower is better. The hybrid `FR-NSGA2` (dark blue) consistently outperforms others.</i>
</p>

<p align="center">
  <img alt="Untangled Solutions Comparison Chart" src="https://github.com/user-attachments/assets/36480d24-ed36-4e1c-9854-66bc45e701d1" />
  
  <i align="center">Percentage of Successfully Untangled Solutions. This chart compares the ability of different algorithms to achieve a perfect planar embedding (0 edge crossings) for planar graphs of increasing size. The results dramatically illustrate the superiority of the hybrid FR-NSGA2 algorithm, which is the only method that consistently finds planar layouts for large graphs (100 and 200 vertices).</i>
</p>

<p align="center">
  <img alt="ZDT1-SSGA" src="https://github.com/user-attachments/assets/d06c8d3e-e5e4-4124-87ea-c042aa3f1ade" width=49% />
  <img alt="ZDT1-NSGA2" src="https://github.com/user-attachments/assets/4f9745a3-9a60-4424-93a3-521576f29093" width=49% />

  <i align="center">Solution of ZDT1 problem by Single-Objective SSGA (left) and Multi-Objective NSGA-II (right)</i>
</p>

<p align="center">
  <img alt="TSP-100" src="https://github.com/user-attachments/assets/5045c41d-c2a5-4e33-8a78-4b17205a2a9d" width=49% />
  
  <i align="center">An example solution for a 100-city Traveling Salesperson Problem found by the SSGA algorithm.</i>
</p>

## :scroll: License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## :pray: Acknowledgments

This project is based on the undergraduate research paper completed at the **National Research Nuclear University MEPhI (Moscow Engineering Physics Institute)**.

- **Author:** Gregory Koganovsky
- **Supervisor:** M.A. Korotkova, Ph.D., Associate Professor

A special thanks to the faculty of the Department of Cybernetics (No. 22) for their guidance and support.
