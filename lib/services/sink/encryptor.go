package sink

import (
	"fmt"
	"sync"
	"telemetry-task/lib/crypto"
)

const (
	workersCount = 5
)

type EncryptManager struct {
	wg       *sync.WaitGroup
	doneCh   chan struct{}
	inputCh  chan string
	outputCh chan string
	jobsCh   chan string
}

func NewEncryptManager() *EncryptManager {
	return &EncryptManager{
		wg:       &sync.WaitGroup{},
		doneCh:   make(chan struct{}),
		inputCh:  make(chan string),
		jobsCh:   make(chan string, workersCount),
		outputCh: make(chan string, workersCount),
	}
}

func (em *EncryptManager) Run() error {
	errCh := make(chan error, workersCount)

	// warm up workers
	for range workersCount {
		worker, err := NewWorker()
		if err != nil {
			return fmt.Errorf("create worker: %v", err)
		}

		em.wg.Add(1)
		go func() {
			worker.run(em.wg, em.inputCh, em.jobsCh, errCh)
		}()
	}
	go func() {
		for {
			select {
			case <-em.doneCh:
				return
			case in := <-em.inputCh:
				em.jobsCh <- in
			case res := <-em.jobsCh:
				em.outputCh <- res
			case err := <-errCh:
				logger.Error("Error in crypto worker", "err:", err)
			}
		}
	}()
	return nil
}

func (em *EncryptManager) Stop() {
	em.doneCh <- struct{}{}
	close(em.doneCh)
	close(em.inputCh)
	em.wg.Wait()
	close(em.jobsCh)
	close(em.outputCh)
}

type Worker struct {
	encryptor crypto.Encryptor
}

func NewWorker() (*Worker, error) {
	enc, err := crypto.NewEncryptor()
	if err != nil {
		return nil, fmt.Errorf("create encryptor: %v", err)
	}
	return &Worker{
		encryptor: enc,
	}, nil
}

func (w *Worker) run(wg *sync.WaitGroup, in <-chan string, out chan<- string, errCh chan<- error) {
	defer wg.Done()
	for message := range in {
		encrypted, err := w.encryptor.EncryptMessage(message)
		if err != nil {
			errCh <- fmt.Errorf("encrypt message: %w", err)
			continue
		}
		out <- encrypted
	}
}
