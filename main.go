package main

import (
	"context"
	"errors"
	"follow/common"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/thoas/go-funk"
	"log"
	"os"
)

func main() {
	cloudfunction.Start(Sync)
}

var RequestForbidden = errors.New("request forbidden")

func Sync(ctx context.Context) (err error) {
	api := common.NewAPI(os.Getenv("USERNAME"), os.Getenv("TOKEN"))
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
	log.Println(map[string]interface{}{
		"add": add,
		"del": del,
	})
	return
}
