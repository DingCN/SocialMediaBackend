package web

import (
	"time"

	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
	"github.com/DingCN/SocialMediaBackend/pkg/twitterTimestamp"
)

// TweetTimeConvert - convert Tweet type to TweetTmpl type to suit
// correct timestamp display format
func TweetTimeConvert(tweet *protocol.Tweet) TweetTmpl {
	newTime := twitterTimestamp.Timestamp(tweet.Timestamp)
	newTimeString := newTime.Add(time.Hour * (-5)).Format("2006-01-02 15:04:05")
	newTweet := TweetTmpl{
		UserName:  tweet.UserName,
		Timestamp: newTimeString,
		Body:      tweet.Body,
	}
	return newTweet
}

// TweetListToTweetTmpl - convert a slice of Tweet into a slice of TweetTmpl type
func TweetListToTweetTmpl(tweetList []*protocol.Tweet) []TweetTmpl {
	var newTweetList []TweetTmpl

	for _, t := range tweetList {
		newTweetList = append(newTweetList, TweetTimeConvert(t))
	}
	return newTweetList
}
