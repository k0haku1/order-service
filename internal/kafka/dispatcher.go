package kafka

import (
	"log"
	"sync"
)

type Event struct {
	Key   string
	Value []byte
}

type Dispatcher struct {
	producer *Producer
	events   chan Event
	wg       sync.WaitGroup
	stopCh   chan struct{}
}

func NewDispatcher(producer *Producer, bufferSize int) *Dispatcher {
	d := &Dispatcher{
		producer: producer,
		events:   make(chan Event, bufferSize),
		stopCh:   make(chan struct{}),
	}

	d.wg.Add(1)
	go d.run()

	return d
}

func (d *Dispatcher) run() {
	defer d.wg.Done()

	for {
		select {
		case e := <-d.events:
			d.producer.Send(e.Key, e.Value)

		case <-d.stopCh:
			for e := range d.events {
				d.producer.Send(e.Key, e.Value)
			}
		}
		return

	}
}
func (d *Dispatcher) Publish(key string, value []byte) {
	select {
	case d.events <- Event{Key: key, Value: value}:
	default:
		log.Println("WARNING: Event channel is full")
	}
}

func (d *Dispatcher) Close() error {
	close(d.stopCh)
	close(d.events)
	d.wg.Wait()
	return d.producer.Close()
}
