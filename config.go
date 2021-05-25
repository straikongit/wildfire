package main

import "time"

type Config struct {
	FireDurationFull    int
	FireDurationSmall   int
	WastelandDuration   int
	Lightnings          int
	LightningStartsFire int

	//each round new trees grow with a probability of 100.000 / CreateNewTree
	CreateNewTree int
	PausePerRound time.Duration
}

func (c *Config) init() {
	c.FireDurationFull = 40   //15
	c.FireDurationSmall = 10   //15
	c.WastelandDuration = 35  //300
	c.Lightnings = 2         //1
	c.LightningStartsFire = 10 //10

	//each round new trees grow with a probability of 100.000 / CreateNewTree
	c.CreateNewTree = 1   //100              //4
	c.PausePerRound = 200 //50

}
