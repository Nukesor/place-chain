package app

import (
	"fmt"
	"os"

	twitter "github.com/dghubble/go-twitter/twitter"
	oauth1 "github.com/dghubble/oauth1"
)

type TwitterUser struct {
	Name   string
	Bio    string
	Avatar string
}

func NewTwitterCache() *TwitterCache {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)
	twitterClient := twitter.NewClient(httpClient)
	cache := TwitterCache{
		map[string]*TwitterUser{},
		twitterClient,
	}
	return &cache
}

type TwitterCache struct {
	userCache     map[string]*TwitterUser
	twitterClient *twitter.Client
}

func (twitterCache *TwitterCache) fetchUser(name string) {
	fmt.Printf("Fetching user from twitter!\n")
	user_resp, _, err := twitterCache.twitterClient.Users.Show(&twitter.UserShowParams{
		ScreenName: name,
	})
	if err != nil {
		fmt.Printf("Error requesting twitter: %s\n", err)
	}
	user := &TwitterUser{
		user_resp.Name,
		user_resp.Description,
		user_resp.ProfileImageURL,
	}
	fmt.Printf("Sucessfully fetched user: %s\n", user)
	twitterCache.userCache[name] = user
}

func (twitterCache *TwitterCache) getUser(name string) *TwitterUser {
	user := twitterCache.userCache[name]
	if user != nil {
		return user
	}
	twitterCache.fetchUser(name)
	user = twitterCache.userCache[name]
	return user
}
