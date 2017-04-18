package main

import (
	"fmt"
	"strconv"
	"time"

	redmine "github.com/mattn/go-redmine"
	slack "github.com/nickschuch/go-slack"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	rmid         = kingpin.Flag("rmid", "Redmine ID. Passing the option \"me\" uses the redmine ID of the API key owner").Required().OverrideDefaultFromEnvar("RM_ID").String()
	rmkey        = kingpin.Flag("rmkey", "Redmine API key").Required().OverrideDefaultFromEnvar("RM_KEY").String()
	slackWebhook = kingpin.Flag("slack-webhook", "Slack webhooks url").Required().OverrideDefaultFromEnvar("SLACK_WEBHOOK").String()
	channel      = kingpin.Flag("channel", "Which slack channel to post to. Also accepts DMs, with @username").Default("#test_bots").OverrideDefaultFromEnvar("SLACK_CHANNEL").String()
	url          = kingpin.Flag("url", "redmine URL, e.g., https://redmine.youraccount.com.au").Default("https://redmine.com.au").OverrideDefaultFromEnvar("RM_URL").String()
	interval     = kingpin.Flag("interval", "interval check time for last ticket/issue (in minutes)").Default("5").OverrideDefaultFromEnvar("INTERVAL_CHECK").Int()
	botname      = kingpin.Flag("botname", "Name of the bot when it posts on Slack").Default("redmine-bot").OverrideDefaultFromEnvar("SLACK_BOTNAME").String()
	emoji        = kingpin.Flag("emoji", "Which emoji to set for the bot on Slack").Default(":redmine:").OverrideDefaultFromEnvar("SLACK_EMOJI").String()
)

func timediff(minute_interval int) string {
	now := time.Now()
	minutes := minute_interval
	then := now.Add(-time.Duration(minutes) * time.Minute)
	return then.Format("2006-01-02T15:04:05Z")
}

func getparams(timefrom string) string {
	return fmt.Sprintf("%s&assigned_to_id=%s&created_on=%%3E%%3D%s", *rmkey, *rmid, timefrom)
}

func sendmsg(rmsub, rmissue, priority string) {
	msg := fmt.Sprintf("*<%s/issues/%s|RM-%s>*  Priority: *%s*\n%s", *url, rmissue, rmissue, priority, rmsub)
	slack.Send(*botname, *emoji, *channel, msg, *slackWebhook)
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
			sendmsg(i.Subject, strconv.Itoa(i.Id), i.Priority.Name)
		}
	}
}

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("Denis Khoshaba")
	kingpin.CommandLine.Help = "Bot that posts newly created Redmine issues to Slack."
	kingpin.Parse()
	doEvery(time.Duration(*interval)*time.Minute, redminecheck)
}
