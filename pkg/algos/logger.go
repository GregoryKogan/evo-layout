package algos

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/GregoryKogan/genetic-algorithms/pkg/problems"
)

type ProgressLoggerProvider interface {
	InitLogging()
	LogProblem(problem problems.Problem)
	LogStep(step any)
}

type ProgressLogger struct {
	filepath string
}

func NewProgressLogger(filepath string) ProgressLoggerProvider {
	return &ProgressLogger{
		filepath: filepath,
	}
}

func (pl *ProgressLogger) InitLogging() {
	if pl.filepath == "" {
		pl.filepath = "progress-log.jsonl"
	}

	if _, err := os.Stat(pl.filepath); !errors.Is(err, os.ErrNotExist) {
		err = os.Remove(pl.filepath)
		if err != nil {
			panic(fmt.Sprintf("error removing old log file: %v", err))
		}
	}

	if _, err := os.Stat(pl.filepath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(pl.filepath)
		if err != nil {
			panic(fmt.Sprintf("error creating log file: %v", err))
		}
		file.Close()
	} else {
		panic(fmt.Sprintf("%s progress log file already exists", pl.filepath))
	}
}

func (pl *ProgressLogger) LogProblem(problem problems.Problem) {
	pl.log(problem)
}

func (pl *ProgressLogger) LogStep(step any) {
	pl.log(step)
}

func (pl *ProgressLogger) log(obj any) {
	file, err := os.OpenFile(pl.filepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}
