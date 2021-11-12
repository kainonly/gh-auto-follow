package common

import (
	"context"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type API struct {
	client *resty.Client
}

// NewAPI Via OAuth and personal access tokens
// We recommend you use OAuth tokens to authenticate to the GitHub API.
// OAuth tokens include personal access tokens and enable the user to revoke access at any time.
func NewAPI(user string, token string) *API {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	client := resty.New().
		SetBaseURL("https://api.github.com").
		SetBasicAuth(user, token).
		SetHeader("accept", "application/vnd.github.v3+json")
	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal
	return &API{client: client}
}

type Link struct {
	Prev  string
	Next  string
	Last  string
	First string
}

func (x *API) ParseLink(raw string) *Link {
	link := new(Link)
	rels := strings.Split(raw, ",")
	for _, s := range rels {
		x := strings.Split(s, ";")
		v := strings.TrimSpace(x[0])[1 : len(x[0])-1]
		switch strings.TrimSpace(x[1]) {
		case `rel="prev"`:
			link.Prev = v
			break
		case `rel="next"`:
			link.Next = v
			break
		case `rel="last"`:
			link.Last = v
			break
		case `rel="first"`:
			link.First = v
			break
		}
	}
	return link
}

func (x *API) Next(ctx context.Context, next string, i interface{}) (link *Link, err error) {
	var resp *resty.Response
	if resp, err = x.client.R().
		SetContext(ctx).
		SetResult(i).
		Get(next); err != nil {
		return
	}
	return x.ParseLink(resp.Header().Get("Link")), nil
}

type Follower struct {
	Login string `json:"login"`
}

// GetFollowers list followers of the authenticated user
func (x *API) GetFollowers(ctx context.Context) (followers []Follower, err error) {
	var resp *resty.Response
	if resp, err = x.client.R().
		SetContext(ctx).
		SetQueryParam("per_page", "100").
		SetResult(&followers).
		Get("user/followers"); err != nil {
		return
	}
	link := x.ParseLink(resp.Header().Get("Link"))
	for link.Next != "" {
		var data []Follower
		if link, err = x.Next(ctx, link.Next, &data); err != nil {
			return
		}
		followers = append(followers, data...)
	}
	return
}

// GetFollowing list the people the authenticated user follows
func (x *API) GetFollowing(ctx context.Context) (following []Follower, err error) {
	var resp *resty.Response
	if resp, err = x.client.R().
		SetContext(ctx).
		SetQueryParam("per_page", "100").
		SetResult(&following).
		Get("user/following"); err != nil {
		return
	}
	link := x.ParseLink(resp.Header().Get("Link"))
	for link.Next != "" {
		var data []Follower
		if link, err = x.Next(ctx, link.Next, &data); err != nil {
			return
		}
		following = append(following, data...)
	}
	return
}

// AddUser follow a user
func (x *API) AddUser(ctx context.Context, user Follower) (ok bool, err error) {
	var resp *resty.Response
	if resp, err = x.client.R().
		SetContext(ctx).
		SetPathParam("username", user.Login).
		Put("user/following/{username}"); err != nil {
		return
	}
	return resp.StatusCode() == 204, nil
}

// DelUser unfollow a user
func (x *API) DelUser(ctx context.Context, user Follower) (ok bool, err error) {
	var resp *resty.Response
	if resp, err = x.client.R().
		SetContext(ctx).
		SetPathParam("username", user.Login).
		Delete("user/following/{username}"); err != nil {
		return
	}
	return resp.StatusCode() == 204, nil
}
