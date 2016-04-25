package main

import (
	"log"
	"os"
	"time"

	"github.com/christer79/ccTray2Slack/cctray"
	"github.com/christer79/ccTray2Slack/config"
	"github.com/christer79/ccTray2Slack/webinterface"
	"github.com/codegangsta/cli"
)

//CommandLineArgs stores the arguments given on commadn line for later use
type CommandLineArgs struct {
	password   string
	username   string
	configPath string
	logFile    string
	pollTime   time.Duration
}

var commandLineArgs CommandLineArgs
var web webinterface.WebInterface

func setupLog(logPath string) {
	if logPath != "" {
		logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file:	 \"%v\" %v", commandLineArgs.logFile, err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
		log.Println("Startup with log file")
	}
}

func main() {
	var cc cctray.CcTray

	app := cli.NewApp()
	app.Name = "ccTraytoSlack"
	app.Usage = "Parse ccTray data and send upates to you slack channels"
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
			Usage:       "Path to config files drop box folder or single file",
			Destination: &commandLineArgs.configPath,
		},
		cli.StringFlag{
			Name:        "log",
			Value:       "/var/log/ccTray2Slack.log",
			Usage:       "Path to log-file",
			Destination: &commandLineArgs.logFile,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Execute the loop to retrive data and publish",
			Flags: []cli.Flag{
				cli.DurationFlag{
					Name:        "pollinterval",
					Usage:       "Set the poll inteval in seconds",
					Value:       10 * time.Second,
					Destination: &commandLineArgs.pollTime,
				},
				cli.StringFlag{
					Name:  "port",
					Usage: "Set the port to start web interface on",
					Value: "",
				},
			},
			Action: func(c *cli.Context) {

				if config, err := config.LoadConfig(commandLineArgs.configPath); err == nil {
					cc = cctray.CreateCcTray(config.Remotes[0])
					cc.Username = commandLineArgs.username
					cc.Password = commandLineArgs.password
					if c.String("port") != "" {
						web = webinterface.WebInterface{}
						go web.Start(c.String("port"))
					}
					runPollLoop(config, cc)
				} else {
					log.Fatalf("Unable to load config %v, %v ", commandLineArgs.configPath, err)
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
						if config, err := config.LoadConfig(commandLineArgs.configPath); err == nil {
							cc = cctray.CreateCcTray(config.Remotes[0])
							cc.Username = commandLineArgs.username
							cc.Password = commandLineArgs.password
							cc.ListProjects()
						} else {
							log.Fatalf("Unable to load config %v, %v ", commandLineArgs.configPath, err)
						}
					},
				},
				{
					Name:  "verify",
					Usage: "Verify all configuration files in the config folder",
					Action: func(c *cli.Context) {
						if _, err := config.LoadConfig(commandLineArgs.configPath); err != nil {
							log.Fatal("Configuration verification failed")
						}
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
						args := config.DefaultConfigArgs{RemoteURL: c.String("remote"), SlackHook: c.String("slack")}
						config.PrintDefaultConfig(args)
					},
				},
			},
		},
	}
	app.Run(os.Args)
}

func runPollLoop(conf config.Config, cc cctray.CcTray) {
	ticker := time.NewTicker(commandLineArgs.pollTime)
	go cc.GetLatest()

	for {
		select {
		case p := <-cc.Ch:
			if url, msg := conf.Process(p); url != "" {
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
		case s := <-cc.ChProjects:
			stat := conf.GroupProjects(s)
			if len(stat) != 0 {
				select {
				case web.ChStatus <- stat:
				default:
				}
			}
		case <-ticker.C:
			if config.ConfigChanged(commandLineArgs.configPath) {
				if temp, err := config.LoadConfig(commandLineArgs.configPath); err == nil {
					conf = temp
				} else {
					log.Printf("Unable to read config file %s.", commandLineArgs.configPath)
				}
			}
			log.Println("checking ...")
			go cc.GetLatest()
		}
	}
}
