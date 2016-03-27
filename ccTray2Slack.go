package main

import (
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
)

type CommandLineArgs struct {
	password   string
	username   string
	configPath string
}

var commandLineArgs CommandLineArgs

func main() {
	var cc ccTray

	app := cli.NewApp()
	app.Name = "ccTraytoSlack"
	app.Usage = "Parce ccTray data and send upates to you slack channels"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "username",
			Value:       "",
			Usage:       "Username to authenticate to retrieve ccTray",
			Destination: &commandLineArgs.username,
		},
		cli.StringFlag{
			Name:        "password",
			Value:       "",
			Usage:       "Password to authenticate to retrieve ccTray",
			Destination: &commandLineArgs.password,
		},
		cli.StringFlag{
			Name:        "config",
			Value:       "config.d",
			Usage:       "Path to config files drop box folder",
			Destination: &commandLineArgs.configPath,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Execute the loop to retrive data and publish",
			Action: func(c *cli.Context) {
				if config, err := LoadConfig(commandLineArgs.configPath); err == nil {
					cc = CreateCcTray(config.Remotes[0])
					cc.Username = commandLineArgs.username
					cc.Password = commandLineArgs.password
					runPollLoop(config, cc)
				} else {
					log.Fatal("Unable to load config stoping executions")
				}
			},
		},
		{
			Name:  "config",
			Usage: "Configuration command",
			Subcommands: []cli.Command{
				{
					Name:  "projects",
					Usage: "Print all availabale projects on ccTray endpoint",
					Action: func(c *cli.Context) {
						if config, err := LoadConfig(commandLineArgs.configPath); err == nil {
							cc = CreateCcTray(config.Remotes[0])
							cc.Username = commandLineArgs.username
							cc.Password = commandLineArgs.password
							cc.ListProjects()
						} else {
							log.Fatal("Unable to load config stoping executions")
						}
					},
				},
				{
					Name:  "verify",
					Usage: "Verify all configuration files in the config folder",
					Action: func(c *cli.Context) {
						VerifyConfig(commandLineArgs.configPath)
					},
				},
				{
					Name:  "default",
					Usage: "Print a default configfile",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "slack",
							Value: "http://slack.com/",
							Usage: "slack-url to use when outputting default config",
						},
						cli.StringFlag{
							Name:  "remote",
							Value: "http://localhost:8153/go/cctray.xml",
							Usage: "RemoteURL to use when outputting default config",
						},
					},
					Action: func(c *cli.Context) {
						args := DefaultConfigArgs{RemoteURL: c.String("remote"), SlackHook: c.String("slack")}
						PrintDefaultConfig(args)
					},
				},
			},
		},
	}
	app.Run(os.Args)
}

func runPollLoop(config Config, cc ccTray) {
	ticker := time.NewTicker(10 * time.Second)
	go cc.GetLatest()

	for {
		select {
		case p := <-cc.Ch:
			if url, msg := config.Process(p); url != "" {
				log.Printf("posting for %q\n", p.Name)
				msg.UpdateMessage(p)
				msg.PostSlackMessage(url)
			} else {
				log.Printf("skipped %q\n", p.Name)
			}
		case e := <-cc.ChErr:
			if e != nil {
				log.Fatalf("Failed to get ccTray Data \n%v\n", e)
			}
			log.Println("Cycle complete")
		case <-ticker.C:
			if ConfigChanged(commandLineArgs.configPath) {
				if temp, err := LoadConfig(commandLineArgs.configPath); err == nil {
					config = temp
				} else {
					log.Printf("Unable to read config file %s.", commandLineArgs.configPath)
				}
			}
			log.Println("checking ...")
			go cc.GetLatest()
		}
	}
}
