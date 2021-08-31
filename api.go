package main

import (
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterAPIAuthentication struct {
	*twitter.Client
}

func newTwitterAPIAuthentication(params OAuth1AuthenticationParameters) *TwitterAPIAuthentication {
	config := oauth1.NewConfig(params.ConsumerKey, params.ConsumerSecret)
	httpClient := config.Client(oauth1.NoContext, oauth1.NewToken(params.AccessKey, params.AccessSecret))
	client := twitter.NewClient(httpClient)
	if client == nil {
		return nil
	}
	return &TwitterAPIAuthentication{client}
}

// Stream tweets based on terms e.g. a hashtag.
func (a *TwitterAPIAuthentication) MustStream(terms []string) <-chan *twitter.Tweet {
	params := &twitter.StreamFilterParams{Track: terms}
	ch, err := a.Streams.Filter(params)
	if err != nil {
		panic(err)
	}
	out := make(chan *twitter.Tweet)
	go func() {
		defer func() { ch.Stop(); close(out) }()
		for msg := range ch.Messages {
			tweet, ok := msg.(*twitter.Tweet)
			if !ok {
				log.Printf("(error) msg.(*twitter.Tweet): tweet=%+v", tweet)
				continue
			} else if tweet.RetweetedStatus != nil { // Ignore RTs
				// No-op
				continue
			}
			out <- tweet
		}
	}()
	return out
}

func (t *TwitterAPIAuthentication) Retweet(tweet *twitter.Tweet) error {
	if tweet.Retweeted {
		return nil
	}
	_, _, err := t.Statuses.Retweet(tweet.ID, nil)
	return err
}

func (a *TwitterAPIAuthentication) Like(tweet *twitter.Tweet) error {
	if tweet.Favorited {
		return nil
	}
	params := &twitter.FavoriteCreateParams{ID: tweet.ID}
	_, _, err := a.Favorites.Create(params)
	return err
}

func (a *TwitterAPIAuthentication) Follow(tweet *twitter.Tweet) error {
	if tweet.User.Following {
		return nil
	}
	params := &twitter.FriendshipCreateParams{UserID: tweet.User.ID}
	_, _, err := a.Friendships.Create(params)
	return err
}
