package tasks

import (
	"go-micro.dev/v4"
)

type Task interface {
	Run()
}

func ExecuteAll(s micro.Service) {
	task := getTickers{Service: s}

	go task.Run()
}
