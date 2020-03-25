package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// https://twitter.com/year_progress
var progress = []string{
	"░░░░░░░░░░░░░░░ 0%",
	"░░░░░░░░░░░░░░░ 1%",
	"░░░░░░░░░░░░░░░ 2%",
	"░░░░░░░░░░░░░░░ 3%",
	"▓░░░░░░░░░░░░░░ 4%",
	"▓░░░░░░░░░░░░░░ 5%",
	"▓░░░░░░░░░░░░░░ 6%",
	"▓░░░░░░░░░░░░░░ 7%",
	"▓░░░░░░░░░░░░░░ 8%",
	"▓░░░░░░░░░░░░░░ 9%",
	"▓▓░░░░░░░░░░░░░ 10%",
	"▓▓░░░░░░░░░░░░░ 11%",
	"▓▓░░░░░░░░░░░░░ 12%",
	"▓▓░░░░░░░░░░░░░ 13%",
	"▓▓░░░░░░░░░░░░░ 14%",
	"▓▓░░░░░░░░░░░░░ 15%",
	"▓▓░░░░░░░░░░░░░ 16%",
	"▓▓▓░░░░░░░░░░░░ 17%",
	"▓▓▓░░░░░░░░░░░░ 18%",
	"▓▓▓░░░░░░░░░░░░ 19%",
	"▓▓▓░░░░░░░░░░░░ 20%",
	"▓▓▓░░░░░░░░░░░░ 21%",
	"▓▓▓░░░░░░░░░░░░ 22%",
	"▓▓▓░░░░░░░░░░░░ 23%",
	"▓▓▓▓░░░░░░░░░░░ 24%",
	"▓▓▓▓░░░░░░░░░░░ 25%",
	"▓▓▓▓░░░░░░░░░░░ 26%",
	"▓▓▓▓░░░░░░░░░░░ 27%",
	"▓▓▓▓░░░░░░░░░░░ 28%",
	"▓▓▓▓░░░░░░░░░░░ 29%",
	"▓▓▓▓▓░░░░░░░░░░ 30%",
	"▓▓▓▓▓░░░░░░░░░░ 31%",
	"▓▓▓▓▓░░░░░░░░░░ 32%",
	"▓▓▓▓▓░░░░░░░░░░ 33%",
	"▓▓▓▓▓░░░░░░░░░░ 34%",
	"▓▓▓▓▓░░░░░░░░░░ 35%",
	"▓▓▓▓▓░░░░░░░░░░ 36%",
	"▓▓▓▓▓▓░░░░░░░░░ 37%",
	"▓▓▓▓▓▓░░░░░░░░░ 38%",
	"▓▓▓▓▓▓░░░░░░░░░ 39%",
	"▓▓▓▓▓▓░░░░░░░░░ 40%",
	"▓▓▓▓▓▓░░░░░░░░░ 41%",
	"▓▓▓▓▓▓░░░░░░░░░ 42%",
	"▓▓▓▓▓▓░░░░░░░░░ 43%",
	"▓▓▓▓▓▓▓░░░░░░░░ 44%",
	"▓▓▓▓▓▓▓░░░░░░░░ 45%",
	"▓▓▓▓▓▓▓░░░░░░░░ 46%",
	"▓▓▓▓▓▓▓░░░░░░░░ 47%",
	"▓▓▓▓▓▓▓░░░░░░░░ 48%",
	"▓▓▓▓▓▓▓░░░░░░░░ 49%",
	"▓▓▓▓▓▓▓▓░░░░░░░ 50%",
	"▓▓▓▓▓▓▓▓░░░░░░░ 51%",
	"▓▓▓▓▓▓▓▓░░░░░░░ 52%",
	"▓▓▓▓▓▓▓▓░░░░░░░ 53%",
	"▓▓▓▓▓▓▓▓░░░░░░░ 54%",
	"▓▓▓▓▓▓▓▓░░░░░░░ 55%",
	"▓▓▓▓▓▓▓▓░░░░░░░ 56%",
	"▓▓▓▓▓▓▓▓▓░░░░░░ 57%",
	"▓▓▓▓▓▓▓▓▓░░░░░░ 58%",
	"▓▓▓▓▓▓▓▓▓░░░░░░ 59%",
	"▓▓▓▓▓▓▓▓▓░░░░░░ 60%",
	"▓▓▓▓▓▓▓▓▓░░░░░░ 61%",
	"▓▓▓▓▓▓▓▓▓░░░░░░ 62%",
	"▓▓▓▓▓▓▓▓▓░░░░░░ 63%",
	"▓▓▓▓▓▓▓▓▓▓░░░░░ 64%",
	"▓▓▓▓▓▓▓▓▓▓░░░░░ 65%",
	"▓▓▓▓▓▓▓▓▓▓░░░░░ 66%",
	"▓▓▓▓▓▓▓▓▓▓░░░░░ 67%",
	"▓▓▓▓▓▓▓▓▓▓░░░░░ 68%",
	"▓▓▓▓▓▓▓▓▓▓░░░░░ 69%",
	"▓▓▓▓▓▓▓▓▓▓▓░░░░ 70%",
	"▓▓▓▓▓▓▓▓▓▓▓░░░░ 71%",
	"▓▓▓▓▓▓▓▓▓▓▓░░░░ 72%",
	"▓▓▓▓▓▓▓▓▓▓▓░░░░ 73%",
	"▓▓▓▓▓▓▓▓▓▓▓░░░░ 74%",
	"▓▓▓▓▓▓▓▓▓▓▓░░░░ 75%",
	"▓▓▓▓▓▓▓▓▓▓▓░░░░ 76%",
	"▓▓▓▓▓▓▓▓▓▓▓▓░░░ 77%",
	"▓▓▓▓▓▓▓▓▓▓▓▓░░░ 78%",
	"▓▓▓▓▓▓▓▓▓▓▓▓░░░ 79%",
	"▓▓▓▓▓▓▓▓▓▓▓▓░░░ 80%",
	"▓▓▓▓▓▓▓▓▓▓▓▓░░░ 81%",
	"▓▓▓▓▓▓▓▓▓▓▓▓░░░ 82%",
	"▓▓▓▓▓▓▓▓▓▓▓▓░░░ 83%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓░░ 84%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓░░ 85%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓░░ 86%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓░░ 87%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓░░ 88%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓░░ 89%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓░ 90%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓░ 91%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓░ 92%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓░ 93%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓░ 94%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓░ 95%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓░ 96%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 97%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 98%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 99%",
	"▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 100%",
}

// (?i)                   -- case-insensitive
// ^                      -- start            *required
// 	(?:#100daysofcode\W)? -- #100DaysOfCode
// 	(?:r(?:ound)?\W?\d+)? -- round ? OR r ?
// 	\W?                   -- separator
// 	d(?:ay)?\W?(\d+)      -- day ? OR d ?     *required
// (?:\W|$)               -- separator or EOF *required
//
// https://regex101.com/r/VCM8l4/1
var re = regexp.MustCompile(`(?i)^(?:#100daysofcode\W)?(?:r(?:ound)?\W?\d+)?\W?d(?:ay)?\W?(\d+)(?:\W|$)`)

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

func (a *Account) Tweet(status string) (*twitter.Tweet, error) {
	tweet, _, err := a.Statuses.Update(status, nil)
	return tweet, err
}

func (a *Account) Like(tweet *twitter.Tweet) error {
	if tweet.Favorited {
		// No-op
		return nil
	}
	params := &twitter.FavoriteCreateParams{ID: tweet.ID}
	_, _, err := a.Favorites.Create(params)
	return err
}

func (a *Account) Follow(tweet *twitter.Tweet) error {
	if tweet.User.Following {
		// No-op
		return nil
	}
	params := &twitter.FriendshipCreateParams{UserID: tweet.User.ID}
	_, _, err := a.Friendships.Create(params)
	return err
}

func main() {
	var (
		CONSUMER_KEY    = os.Getenv("CONSUMER_KEY")
		CONSUMER_SECRET = os.Getenv("CONSUMER_SECRET")
		ACCESS_KEY      = os.Getenv("ACCESS_KEY")
		ACCESS_SECRET   = os.Getenv("ACCESS_SECRET")
	)
	if CONSUMER_KEY == "" {
		log.Fatal("CONSUMER_KEY cannot be empty")
	} else if CONSUMER_SECRET == "" {
		log.Fatal("CONSUMER_SECRET cannot be empty")
	} else if ACCESS_KEY == "" {
		log.Fatal("ACCESS_KEY cannot be empty")
	} else if ACCESS_SECRET == "" {
		log.Fatal("ACCESS_SECRET cannot be empty")
	}
	user := Auth(CONSUMER_KEY, CONSUMER_SECRET, ACCESS_KEY, ACCESS_SECRET)
	for tweet := range user.MustStream([]string{"#100DaysOfCode"}) {
		// Day 0:
		statusURL := getStatusURL(tweet)
		if strings.HasPrefix(tweet.Text, "I'm publicly committing to the 100DaysOfCode") {
			var rt *twitter.Tweet
			var err error
			if err = user.Like(tweet); err != nil {
				log.Print(err)
				continue
			} else if err = user.Follow(tweet); err != nil {
				log.Print(err)
				continue
			} else if rt, err = user.Tweet(progress[0] + "\n\n" + statusURL); err != nil {
				log.Print(err)
				continue
			}
			log.Printf("retweeted: statusURL=%s tweet=%+v",
				statusURL, rt)
			continue
		}
		// Day 1-100:
		matches := re.FindAllStringSubmatch(tweet.Text, 1)
		if matches == nil || len(matches) == 0 || len(matches[0]) == 0 {
			// No-op
			continue
		}
		day, err := strconv.Atoi(matches[0][1])
		must(err)
		if day < 0 { // Ignore <0
			// No-op
			continue
		} else if day%100 != 0 { // Reset >100
			day %= 100
		}
		// NOTE: progress is zero-based:
		//
		// d0 -> progress[0]
		// d100 -> progress[100]
		//
		var rt *twitter.Tweet
		// var err error
		if err = user.Like(tweet); err != nil {
			log.Print(err)
			continue
		} else if err = user.Follow(tweet); err != nil {
			log.Print(err)
			continue
		} else if rt, err = user.Tweet(progress[day] + "\n\n" + statusURL); err != nil {
			log.Print(err)
			continue
		}
		log.Printf("retweeted: statusURL=%s tweet=%+v",
			statusURL, rt)
	}
}
