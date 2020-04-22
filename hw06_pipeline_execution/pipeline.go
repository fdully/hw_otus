package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here
	if len(stages) < 1 {
		return nil
	}

	generateDoneStage := func(done In) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for {
					select {
					case v, ok := <-in:
						if !ok {
							return
						}
						select {
						case out <- v:
						case <-done:
							return
						}
					case <-done:
						return
					}
				}
			}()
			return out
		}
	}

	doneStage := generateDoneStage(done)

	out := doneStage(in)
	for _, stage := range stages {
		out = stage(out)
	}
	out = doneStage(out)

	return out
}
