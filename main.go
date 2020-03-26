package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// // https://100daysofcode.com
// var day0Re = regexp.MustCompile(`(?i)I'm publicly committing to the #?100DaysOfCode`)

// (?i)                    -- case-insensitive
// ^                       -- start            *required
// 	(?:#100daysofcode\W+)? -- #100DaysOfCode
// 	(?:r(?:ound)?\W*\d+)?  -- round ? OR r ?
// 	\W*                    -- separator
// 	d(?:ay)?\W*(\d+)       -- day ? OR d ?     *required
// (?:\W+|$)               -- separator or EOF *required
//
// https://regex101.com/r/VCM8l4/3
var re = regexp.MustCompile(`(?i)^(?:#100DaysOfCode\W+)?(?:r(?:ound)?\W*\d+)?\W*d(?:ay)?\W*(\d+)(?:\W+|$)`)

type Account struct{ *twitter.Client }

// Authenticates an account.
func Auth(ConsumerKey, ConsumerSecret, AccessKey, AccessSecret string) *Account {
	config := oauth1.NewConfig(ConsumerKey, ConsumerSecret)
	httpClient := config.Client(oauth1.NoContext, oauth1.NewToken(AccessKey, AccessSecret))
	client := twitter.NewClient(httpClient)
	return &Account{client}
}

// Streams tweets based on search terms.
func (a *Account) MustStream(terms []string) <-chan *twitter.Tweet {
	params := &twitter.StreamFilterParams{Track: terms}
	ch, err := a.Streams.Filter(params)
	must(err)
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

// Gets a status URL for a tweet.
func getStatusURL(tweet *twitter.Tweet) string {
	return "https://twitter.com/" + tweet.User.ScreenName + "/status/" + fmt.Sprint(tweet.ID)
}

func must(err error) {
	if err == nil {
		// No-op
		return
	}
	panic(err)
}

func (a *Account) Retweet(tweet *twitter.Tweet) error {
	if tweet.Retweeted {
		return nil
	}
	_, _, err := a.Statuses.Retweet(tweet.ID, nil)
	return err
}

func (a *Account) Like(tweet *twitter.Tweet) error {
	if tweet.Favorited {
		return nil
	}
	params := &twitter.FavoriteCreateParams{ID: tweet.ID}
	_, _, err := a.Favorites.Create(params)
	return err
}

func (a *Account) Follow(tweet *twitter.Tweet) error {
	if tweet.User.Following {
		return nil
	}
	params := &twitter.FriendshipCreateParams{UserID: tweet.User.ID}
	_, _, err := a.Friendships.Create(params)
	return err
}

func main() {
	log.Printf("starting…")
	var (
		CONSUMER_KEY    = os.Getenv("CONSUMER_KEY")
		CONSUMER_SECRET = os.Getenv("CONSUMER_SECRET")
		ACCESS_KEY      = os.Getenv("ACCESS_KEY")
		ACCESS_SECRET   = os.Getenv("ACCESS_SECRET")
	)
	if CONSUMER_KEY == "" {
		log.Fatal("env CONSUMER_KEY cannot be empty")
	} else if CONSUMER_SECRET == "" {
		log.Fatal("env CONSUMER_SECRET cannot be empty")
	} else if ACCESS_KEY == "" {
		log.Fatal("env ACCESS_KEY cannot be empty")
	} else if ACCESS_SECRET == "" {
		log.Fatal("env ACCESS_SECRET cannot be empty")
	}
	user := Auth(CONSUMER_KEY, CONSUMER_SECRET, ACCESS_KEY, ACCESS_SECRET)
	log.Printf("…started")
	for tweet := range user.MustStream([]string{"#100DaysOfCode"}) {
		statusURL := getStatusURL(tweet)
		if !strings.HasPrefix(tweet.Text, "I'm publicly committing to the 100DaysOfCode") && !re.MatchString(tweet.Text) {
			// No-op
			continue
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := user.Like(tweet)
			if err != nil {
				log.Print(err)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := user.Retweet(tweet)
			if err != nil {
				log.Print(err)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := user.Follow(tweet)
			if err != nil {
				log.Print(err)
			}
		}()
		wg.Wait()
		log.Printf("retweeted %s", statusURL)
	}
}
