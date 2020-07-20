import requests
import os

base_uri = 'https://api.github.com'
auth = (os.environ['USER'], os.environ['PASS'])


def fetch(path, page) -> requests.Response:
    return requests.get(
        url=base_uri + path,
        auth=auth,
        params={
            'page': page,
            'per_page': 100
        }
    )


default = 1
followers = []
while True:
    response = fetch('/user/followers', default)
    lists = response.json()
    if len(lists) != 0:
        followers.extend(lists)
        default += 1
    else:
        break

default = 1
followings = []
while True:
    response = fetch('/user/following', default)
    lists = response.json()
    if len(lists) != 0:
        followings.extend(lists)
        default += 1
    else:
        break

unfollowers = []
for following in followings:
    is_follower = False
    for follower in followers:
        if following['login'] == follower['login']:
            is_follower = True
            break
    if is_follower is False:
        unfollowers.append(following)

for unfollower in unfollowers:
    response = requests.delete(
        url=base_uri + '/user/following/' + unfollower['login'],
        auth=auth,
    )
    print(response.text)

unfollowings = []
for follower in followers:
    is_following = False
    for following in followings:
        if follower['login'] == following['login']:
            is_following = True
            break
    if is_following is False:
        unfollowings.append(follower)

for unfollowing in unfollowings:
    response = requests.put(
        url=base_uri + '/user/following/' + unfollowing['login'],
        auth=auth,
    )
    print(response.text)
