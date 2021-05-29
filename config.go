package main

import (
	"time"
	"encoding/json"
	"io/ioutil"
)


type Configs struct {
	        Configs map[string]*Config
	}

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
	c.FireDurationFull = 80   //15
	c.FireDurationSmall = 5   //15
	c.WastelandDuration = 35  //300
	c.Lightnings = 1          //1
	c.LightningStartsFire = 1 //10

	//each round new trees grow with a probability of 100.000 / CreateNewTree
	c.CreateNewTree = 10  //100              //4
	c.PausePerRound = 400 //50

}
func LoadConfig(configName string) *Config {
	file, _ := ioutil.ReadFile("data.json")
	cx := Configs{}

	_ = json.Unmarshal([]byte(file), &cx)
	if err := json.Unmarshal([]byte(file), &cx); err != nil {
		panic(err)
	}
	return (cx.Configs[configName])

}
