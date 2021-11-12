package main

import (
	"context"
	"errors"
	"follow/common"
	"github.com/thoas/go-funk"
	"log"
	"os"
	"time"
)

func main() {
	result, err := sync()
	if err != nil {
		panic(err)
	}
	log.Println(result)
}

var RequestForbidden = errors.New("request forbidden")

func sync() (result interface{}, err error) {
	api := common.NewAPI(os.Getenv("USER"), os.Getenv("TOKEN"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var followers []common.Follower
	if followers, err = api.GetFollowers(ctx); err != nil {
		return
	}
	var following []common.Follower
	if following, err = api.GetFollowing(ctx); err != nil {
		return
	}
	add, del := funk.Difference(followers, following)
	var ok bool
	for _, x := range add.([]common.Follower) {
		ok, err = api.AddUser(ctx, x)
		if err != nil {
			return
		}
		if !ok {
			err = RequestForbidden
			return
		}
	}
	for _, x := range del.([]common.Follower) {
		ok, err = api.DelUser(ctx, x)
		if err != nil {
			return
		}
		if !ok {
			err = RequestForbidden
			return
		}
	}
	return map[string]interface{}{
		"add": add,
		"del": del,
	}, nil
}
