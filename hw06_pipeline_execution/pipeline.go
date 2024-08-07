package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// transferBetweenCh перегнать данные из in в out и остановить процесс, если done.
func transferBetweenCh(in In, out Bi, done In) {
	defer close(out)
	for {
		select {
		// приоритизация случая с ранним выходом
		case <-done:
			return
		default:
		}
		select {
		case <-done:
			return
		case v, ok := <-in:
			if !ok {
				return
			}
			out <- v
		}
	}
}

// ExecutePipeline запускает конкурентный пайплайн.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		select {
		// приоритизация случая с ранним выходом
		case <-done:
			break
		default:
		}
		processed := make(Bi)
		go transferBetweenCh(stage(out), processed, done)
		out = processed
	}
	return out
}
