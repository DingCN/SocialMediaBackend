package web

import (
	"github.com/DingCN/SocialMediaBackend/pkg/protocol"
	"github.com/DingCN/SocialMediaBackend/pkg/twitterTimestamp"
)

func TweetTimeConvert(tweet *protocol.Tweet) TweetTmpl {
	newTime := twitterTimestamp.Timestamp(tweet.Timestamp)
	newTimeString := newTime.Format("2006-01-02 15:04:05")
	newTweet := TweetTmpl{
		UserName:  tweet.UserName,
		Timestamp: newTimeString,
		Body:      tweet.Body,
	}
	return newTweet
}

func TweetListToTweetTmpl(tweetList []*protocol.Tweet) []TweetTmpl {
	var newTweetList []TweetTmpl

	for _, t := range tweetList {
		newTweetList = append(newTweetList, TweetTimeConvert(t))
	}
	return newTweetList
}
