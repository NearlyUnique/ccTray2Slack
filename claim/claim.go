package claim

import (
	"log"

	"github.com/hashicorp/consul/api"
)

//Action Used to send messages to claim and unclaim pipelines
type Action struct {
	Pipeline string
	Action   bool
	Comment  string
	User     string
}

//Claim Struct to clam and unclame pieplines
type Claim struct {
	consulURL        string
	consulDataCenter string
	ChClaim          chan Action
}

func CreateClaim(URL, dataCenter string) Claim {
	return Claim{
		consulURL:        URL,
		consulDataCenter: dataCenter,
		ChClaim:          make(chan Action),
	}
}

func (c *Claim) Start() {

	config := api.DefaultConfig()
	config.Datacenter = c.consulDataCenter
	config.Address = c.consulURL
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	kv := client.KV()

	for {
		select {
		case a := <-c.ChClaim:
			if a.Action {
				log.Printf("Claiming pipeline \"%s\" for user: \"%s\" with comment \"%s\"\n", a.Pipeline, a.User, a.Comment)
				p := &api.KVPair{Key: "claimed/" + a.Pipeline, Value: []byte(a.User + ":" + a.Comment)}
				_, err = kv.Put(p, nil)
				if err != nil {
					panic(err)
				}
			} else {
				log.Printf("Unclaiming pipeline \"%s\" for user: \"%s\" with comment \"%s\"\n", a.Pipeline, a.User, a.Comment)
				_, err = kv.Delete("claimed/"+a.Pipeline, nil)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
