package common

import (
	"context"
	"testing"
)

func TestUser(t *testing.T) {
	var data map[string]interface{}
	resp, err := api.client.R().
		SetResult(&data).
		Get("user")
	if err != nil {
		return
	}
	t.Log(resp)
	t.Log(data)
}

func TestAPI_ParseLink(t *testing.T) {
	raw := `<https://api.github.com/user/followers?per_page=10&page=1>; rel="prev", <https://api.github.com/user/followers?per_page=10&page=3>; rel="next", <https://api.github.com/user/followers?per_page=10&page=13>; rel="last", <https://api.github.com/user/followers?per_page=10&page=1>; rel="first"`
	link := api.ParseLink(raw)
	t.Log(link)
}

func TestAPI_Next(t *testing.T) {
	var data []Follower
	link, err := api.Next(context.TODO(), "https://api.github.com/user/followers?per_page=10&page=2", &data)
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
	t.Log(link)
}

func TestAPI_GetFollowers(t *testing.T) {
	followers, err := api.GetFollowers(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log(len(followers))
	t.Log(followers)
}

func TestAPI_GetFollowing(t *testing.T) {
	following, err := api.GetFollowing(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Log(len(following))
	t.Log(following)
}
