package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, s := range stages {
		out := make(Bi)
		res := s(in)

		go func() {
			defer close(out)

			for {
				select {
				case <-done:
					go drain(res)
					return
				case v, ok := <-res:
					if !ok {
						return
					}

					out <- v
				}
			}
		}()

		in = out
	}

	return in
}

func drain(ch <-chan interface{}) {
	for v := range ch {
		_ = v
	}
}
