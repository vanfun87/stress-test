package stress

import "fmt"

type StressClient struct {
	Number     int
	Concurrent int
	Limitation int
}

func NewStressClient(number int, concurrent int, limitation int) *StressClient {
	return &StressClient{
		Number:     number,
		Concurrent: concurrent,
		Limitation: limitation,
	}
}

func NewStressClientWithNumber(number int) *StressClient {
	return NewStressClient(number, 1, 0)
}

func NewStressClientWithConcurrentNumber(number int, concurrent int) *StressClient {
	return NewStressClient(number, concurrent, 0)
}

func (s *StressClient) Header() {
	msg := fmt.Sprintf("%d task(s) ready to run with %d thread(s)", s.Number, s.Concurrent)

	if s.Limitation > 0 {
		msg += fmt.Sprintf("%s, with %d task(s) limitation per second", msg, s.Limitation)
	}

	fmt.Println(msg)
	fmt.Println()
}
