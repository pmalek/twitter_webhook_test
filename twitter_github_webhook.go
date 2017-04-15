package main

import (
	"flag"
	"fmt"
	"log"

	webhooks "gopkg.in/go-playground/webhooks.v2"
	"gopkg.in/go-playground/webhooks.v2/github"

	twitter "github.com/ChimeraCoder/anaconda"
)

func main() {
	// Parse comannd line flags
	twitterConsumerKey := flag.String("twitterConsumerKey", "", "")
	twitterConsumerSecret := flag.String("twitterConsumerSecret", "", "")
	twitterAccessToken := flag.String("twitterAccessToken", "", "")
	twitterSecret := flag.String("twitterSecret", "", "")
	githubWebhookSecret := flag.String("githubWebhookSecret", "", "")
	flag.Parse()

	// Setup Twitter API
	twitter.SetConsumerKey(*twitterConsumerKey)
	twitter.SetConsumerSecret(*twitterConsumerSecret)
	twitterAPI := twitter.NewTwitterApi(*twitterAccessToken, *twitterSecret)

	user, err := twitterAPI.GetSelf(nil)
	if err != nil {
		log.Fatalf("Twitter VerifyCredentials failed with error: %v", err)
		return
	}
	twitterUserName := user.ScreenName
	log.Printf("Twitter VerifyCredentials for '%s' successful :)", twitterUserName)

	// Setup GitHub webhook events and Listen
	hook := github.New(&github.Config{Secret: *githubWebhookSecret})
	hook.RegisterEvents(func(payload interface{}, header webhooks.Header) {
		pl, ok := payload.(github.PushPayload)
		if !ok {
			log.Printf("Recieved non PushPayload payload! Check what types of events you have configured on github")
			return
		}

		if pl.Forced {
			log.Printf("Tweeting push forced commit in %s - ups this might be ugly...", pl.Repository.URL)
		} else {
			log.Printf("Tweeting push event in %s, %d commits", pl.Repository.URL, len(pl.Commits))
		}

		for i := 0; i < len(pl.Commits); i++ {
			post := fmt.Sprintf("new commit in %s by %s %s/commit/%s",
				pl.Repository.URL,
				pl.Commits[i].Committer.Username,
				pl.Repository.URL,
				pl.Commits[i].ID[:6])

			tweet, err := twitterAPI.PostTweet(post, nil)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Successfully tweeted '%s' at https://www.twitter.com/%s/status/%d", post, twitterUserName, tweet.Id)
			}

		}
	}, github.PushEvent)

	log.Println("Starting webhooks at :8080...")
	err = webhooks.Run(hook, fmt.Sprintf(":%d", 8080), "/")
	if err != nil {
		log.Fatal(err)
	}
}
