# ccTrayToSlack

Poll ccTray endpoint (Cruise Control xml schema), find changes and post to yourcompany.slack.com

## Getting Started

### Slack web-hook

Create a custom slack web-hook in your company slack and find out the URL to use to post to the channels.

### Configuration

Output an example configuration with the command (replacing the remote argument with your cctray URL and slack with the custom web-hook to your slack)

````
./ccTrayToSlack config default --remote http://localhost:8153/go/cctray.xml --slack http://slack.com > config.d/config
````

Once your configuration have a remote specified you can get the complete list of projects in your ccTray use the command:

````
./ccTray2Slack --config config.d/ --username <username> --password <password> config projects
````

Use the list to update the watches to monitor the projects you want to get notifications about.

You can adapt the way the slack messages look in your configuration file

The following keywords in Text and Titles will be replaced by data from your ccTray output:
* %time%
* %project%
* %status%
* %url%
* %label%

Verify the configuration

````
./ccTray2Slack --config config.d/ --username <username> --password <password> config verify
````

### Running

#### Command line

````
./ccTray2Slack --config config.d/ --username <username> --password <password> start
````

#### Docker

Drop the configuration in the folder docker/config/

````
cd docker
docker-compose build
docker-compose up -d
````

## Running the tests

````
go test
````

# TODO:

- Make it possible to specify multiple remotes
- make poll time configurable
- refactor to make the responsibilities for types clearer
