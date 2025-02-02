package hw06pipelineexecution

import "fmt"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := in
	for i, stage := range stages {
		fmt.Println("stage: ", i)
		result = stageWithDone(result, done, stage)
	}
	return result
}

func stageWithDone(in In, done In, stage Stage) Out {
	stageStream := make(Bi)
	go func() {
		defer close(stageStream)
		for i := range in {
			select {
			case <-done:
				fmt.Println("done 2")
				return
			default:
			}

			select {
			case <-done:
				fmt.Println("done 3")
				return
			case stageStream <- i:
			}
		}
	}()
	return stage(stageStream)
}
