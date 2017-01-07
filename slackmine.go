package main

import (
	"bytes"
	redmine "github.com/mattn/go-redmine"
	slack "github.com/nickschuch/go-slack"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
	"time"
)

var (
	rmid         = kingpin.Flag("rmid", "Redmine ID. Passing the option \"me\" uses the redmine ID of the API key owner").Required().String()
	rmkey        = kingpin.Flag("rmkey", "Redmine API key").Required().OverrideDefaultFromEnvar("RM_KEY").String()
	slackWebhook = kingpin.Flag("slack-webhook", "Slack webhooks url").Required().OverrideDefaultFromEnvar("SLACK_WEBHOOK").String()
	channel      = kingpin.Flag("channel", "Which slack channel to post to. Also accepts DMs, with @username").Default("#test_bots").String()
	url          = kingpin.Flag("url", "redmine URL, e.g., https://redmine.yourcompany.com.au").String()
	interval     = kingpin.Flag("interval", "interval check time for last ticket/issue (in minutes)").Default("5").Int()
	botname      = kingpin.Flag("botname", "Name of the bot when it posts on Slack").Default("redminebot").String()
	emoji        = kingpin.Flag("emoji", "Which emoji to set for the bot on Slack").Default(":monkey:").String()
)

func timediff(minute_interval int) string {
	now := time.Now()
	minutes := minute_interval
	then := now.Add(-time.Duration(minutes) * time.Minute)
	return then.Format("2006-01-02T15:04:05Z")
}

func getparams(timefrom string) string {
	var buffer bytes.Buffer
	buffer.WriteString(*rmkey)
	buffer.WriteString("&assigned_to_id=")
	buffer.WriteString(*rmid)
	buffer.WriteString("&created_on=%3E%3D")
	buffer.WriteString(timefrom)
	return buffer.String()
}

func sendmsg(rmsub, rmissue string) {
	var buffer bytes.Buffer
	buffer.WriteString("*<")
	buffer.WriteString(*url)
	buffer.WriteString("/issues/")
	buffer.WriteString(rmissue)
	buffer.WriteString("|RM")
	buffer.WriteString(rmissue)
	buffer.WriteString(">*\n")
	buffer.WriteString(rmsub)
	slack.Send(*botname, *emoji, *channel, buffer.String(), *slackWebhook)
	buffer.Reset()
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func redminecheck(t time.Time) {
	timefrom := timediff(*interval)
	params := getparams(timefrom)
	client := redmine.NewClient(*url, params)
	issues, err := client.Issues()

	if err != nil {
		panic(err)
	}

	for _, i := range issues {
		if len(i.Subject) > 0 {
			sendmsg(i.Subject, strconv.Itoa(i.Id))
		}
	}
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Denis Khoshaba")
	kingpin.CommandLine.Help = "Bot that posts newly created Redmine issues to Slack."
	kingpin.Parse()
	doEvery(time.Duration(*interval)*time.Minute, redminecheck)
}
