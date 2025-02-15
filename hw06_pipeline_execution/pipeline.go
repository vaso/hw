package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := in
	for _, stage := range stages {
		result = stageWithDone(result, done, stage)
	}
	return result
}

func stageWithDone(in In, done In, stage Stage) Out {
	stageStream := make(Bi)
	go func() {
		defer close(stageStream)

		for {
			select {
			case <-done:
				go func() {
					for range in {
					}
				}()
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				stageStream <- val
			}
		}
	}()

	return stage(stageStream)
}
