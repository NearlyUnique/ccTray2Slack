Poll ccTray endpoint (Cruise Control xml schema), find changes and post to yourcompany.slack.com

Initially done as a learning exercise for me with golang

edit config.json to point to a slack integration point and alter the project regex to find interesting projects, changes are then sent to slack

TODO:
- Use color to make failing and fixing more obvious
- deal with http authentication to get ccTray xml
- make poll time configurable
- make the intersting map property only store interesting things according to the process func
- refactor to make the responsibilities for types clearer
- try and fond someone willing to review it to make it 'idomatic go'
- add per status/avtivity slack messages
