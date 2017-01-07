# slackmine

Bot that posts newly created Redmine issues to Slack

### Build

```
git clone git@github.com:TheDen/slackmine.git
cd slackmine
go get bytes github.com/mattn/go-redmine github.com/nickschuch/go-slack gopkg.in/alecthomas/kingpin.v2 strconv time
go build slackmine.go
```

`upx` (https://github.com/pwaller/goupx) can also be used to shrink the binary

### Run

```
usage: slackmine --rmid=RMID --rmkey=RMKEY --slack-webhook=SLACK-WEBHOOK [<flags>]

Bot that posts newly created Redmine issues to Slack.

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --rmid=RMID                    Redmine ID. Passing the option "me" uses the redmine ID of the API key owner
  --rmkey=RMKEY                  Redmine API key
  --slack-webhook=SLACK-WEBHOOK  Slack webhooks url
  --channel="#test_bots"         Which slack channel to post to. Also accepts DMs, with @username
  --url=URL                      redmine URL, e.g., https://redmine.yourcompany.com.au
  --interval=5                   interval check time for last ticket/issue (in minutes)
  --botname="redminebot"         Name of the bot whne it posts on Slack
  --emoji=":monkey:"             Which emoji to set for the bot on Slack
  --version                      Show application version.
```
### Example
To have slackmine check every 5 minutes for a new redmine issue for redmine user `3` on the slack channel #sys-ops, simply run:

```./slackmine --rmkey=xxxxxxxxxx --rmid=3 --slack-webhook=https://hooks.slack.com/services/xxxxxx/xxxxxx/xxxxxx --channel=#sys-ops --interval=5```

If a new issue is found, slackmine will post it on the `#sys-ops` channel. An issue can also be posted as a direct message on Slack by passing `@username` in the `--channel` arg.
