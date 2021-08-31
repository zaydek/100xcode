package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

var (
	CONSUMER_KEY    = os.Getenv("CONSUMER_KEY")
	CONSUMER_SECRET = os.Getenv("CONSUMER_SECRET")
	ACCESS_KEY      = os.Getenv("ACCESS_KEY")
	ACCESS_SECRET   = os.Getenv("ACCESS_SECRET")

	CHECK_FOR_BLOCKED_SCREENNAMES_INTERVAL = 15 * time.Minute
)

// (?i)                     -- case-insensitive
// ^                        -- start              *required
// 	(?:#100daysofcode\W+)?  -- #100DaysOfCode     *optional
// 	(?:r(?:ounds?)?\W*\d+)? -- round(s?) ? OR r ? *optional
// 	\W*                     -- separator          *optional
// 	d(?:ays?)?\W*(\d+)      -- day(s?) ? OR d ?   *required
// (?:\W+|$)                -- separator or EOF   *required
//
// https://regex101.com/r/VCM8l4/3
var progressRegex = regexp.MustCompile(`(?i)^(?:#100DaysOfCode\W+)?(?:r(?:ounds?)?\W*\d+)?\W*d(?:ays?)?\W*(\d+)(?:\W+|$)`)

////////////////////////////////////////////////////////////////////////////////

type BlockedService struct {
	interval                        time.Duration
	shouldRefreshBlockedScreenNames bool
	blockedScreenNamesMap           map[string]struct{}

	*TwitterOAuth1Authentication
}

func newBlockedService(twitterOAuth1Auth *TwitterOAuth1Authentication, interval time.Duration) *BlockedService {
	srv := &BlockedService{
		interval:                        interval,
		shouldRefreshBlockedScreenNames: true,
		blockedScreenNamesMap:           map[string]struct{}{},
	}
	srv.Refresh()
	go func() {
		ticker := time.NewTicker(interval)
		for ; true; <-ticker.C {
			srv.shouldRefreshBlockedScreenNames = true
		}
	}()
	return srv
}

func (b *BlockedService) Refresh() error {
	if !b.shouldRefreshBlockedScreenNames {
		return nil
	}
	screenNames, err := b.GetBlockedScreenNames()
	if err != nil {
		return err
	}
	for _, screenName := range screenNames {
		b.blockedScreenNamesMap[strings.ToLower(screenName)] = struct{}{}
	}
	b.shouldRefreshBlockedScreenNames = false
	return nil
}

func (b *BlockedService) IsBlocked(screenName string) bool {
	_, ok := b.blockedScreenNamesMap[strings.ToLower(screenName)]
	return ok
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	log.Println("initializing")

	if CONSUMER_KEY == "" {
		log.Fatal("env CONSUMER_KEY cannot be empty")
	} else if CONSUMER_SECRET == "" {
		log.Fatal("env CONSUMER_SECRET cannot be empty")
	} else if ACCESS_KEY == "" {
		log.Fatal("env ACCESS_KEY cannot be empty")
	} else if ACCESS_SECRET == "" {
		log.Fatal("env ACCESS_SECRET cannot be empty")
	}
}

func isRelevant(tweet *twitter.Tweet) bool {
	return strings.HasPrefix(tweet.Text, "I'm publicly committing to the 100DaysOfCode") || progressRegex.MatchString(tweet.Text)
}

func main() {
	// Connect to the Twitter API
	api := newTwitterAPIAuthentication(OAuth1AuthenticationParameters{
		ConsumerKey:    CONSUMER_KEY,
		ConsumerSecret: CONSUMER_SECRET,
		AccessKey:      ACCESS_KEY,
		AccessSecret:   ACCESS_SECRET,
	})
	if api == nil {
		panic("failed to authenticate twitter api")
	}
	log.Println("connected to twitter api")

	// Connect to the Twitter OAuth1 API (for checking for blocker screen names)
	twitterOAuth1Auth := newTwitterOAuth1Authentication(OAuth1AuthenticationParameters{
		ConsumerKey:    CONSUMER_KEY,
		ConsumerSecret: CONSUMER_SECRET,
		AccessKey:      ACCESS_KEY,
		AccessSecret:   ACCESS_SECRET,
	})
	if twitterOAuth1Auth == nil {
		panic("failed to authenticate twitter oauth1 api")
	}
	log.Println("connected to twitter oauth1 api")

	// Create a blocked service and start streaming `"#100DaysOfCode"` tweets
	blockedService := newBlockedService(twitterOAuth1Auth, CHECK_FOR_BLOCKED_SCREENNAMES_INTERVAL)
	for tweet := range api.MustStream([]string{"#100DaysOfCode"}) {
		var (
			username = strings.ToLower(tweet.User.ScreenName)
			url      = fmt.Sprintf("https://twitter.com/%s/status/%s", username, fmt.Sprint(tweet.ID))
		)
		if !isRelevant(tweet) {
			log.Printf("ignored irrelevant user @%s tweet %s\n",
				username, url)
			continue
		}
		if err := blockedService.Refresh(); err != nil {
			panic(fmt.Sprintf("failed to refresh blocked service; %s", err))
		}
		if blockedService.IsBlocked(tweet.User.ScreenName) {
			log.Printf("ignored blocked user @%s tweet %s\n",
				username, url)
			continue
		}
		if err := api.Retweet(tweet); err != nil {
			log.Printf("cannot retweet user @%s tweet %s; %s\n",
				username, url, err)
			continue
		}
		log.Printf("retweeted user @%s tweet %s\n",
			username, url)
	}
}
