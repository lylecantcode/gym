run: gym
	./gym

gym: gym.go
	go build gym.go

.PHONY: run
