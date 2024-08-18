package kalkan

import "log"

type Kalkan struct {
}

func New() *Kalkan {
	return &Kalkan{}
}

func (k *Kalkan) Run() error {
	log.Printf("starting kalkan")
	return nil
}