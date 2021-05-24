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
	c.FireDurationFull = 200   //15
	c.FireDurationSmall = 15   //15
	c.WastelandDuration = 150  //300
	c.Lightnings = 20          //1
	c.LightningStartsFire = 30 //10

	//each round new trees grow with a probability of 100.000 / CreateNewTree
	c.CreateNewTree = 2   //100              //4
	c.PausePerRound = 200 //50

}
