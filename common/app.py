import os
import requests


class App:
    def __init__(self):
        self.base_uri = 'https://api.github.com'
        self.auth = (os.environ['username'], os.environ['token'])

    def bootstrap(self) -> dict:
        result = {
            'increase': [],
            'remove': []
        }
        followers = self.__get_followers()
        followings = self.__get_followings()
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
            response = self.__del_user(unfollower['login'])
            if response.status_code == 204:
                result['remove'].append(unfollower['login'])

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
            response = self.__add_user(unfollowing['login'])
            if response.status_code == 204:
                result['increase'].append(unfollowing['login'])
        return result

    def __fetch(self, method: str, path: str, params: dict = None) -> requests.Response:
        return requests.request(
            method=method,
            url=self.base_uri + path,
            auth=self.auth,
            headers={
                'accept': 'application/vnd.github.v3+json'
            },
            params=params
        )

    def __get_followers(self) -> []:
        default = 1
        followers = []
        while True:
            response = self.__fetch('GET', '/user/followers', {
                'page': default,
                'per_page': 100
            })
            lists = response.json()
            if len(lists) != 0:
                followers.extend(lists)
                default += 1
            else:
                break
        return followers

    def __get_followings(self) -> []:
        default = 1
        followings = []
        while True:
            response = self.__fetch('GET', '/user/following', {
                'page': default,
                'per_page': 100
            })
            lists = response.json()
            if len(lists) != 0:
                followings.extend(lists)
                default += 1
            else:
                break
        return followings

    def __del_user(self, username: str) -> requests.Response:
        return self.__fetch('DELETE', f'/user/following/{username}')

    def __add_user(self, username: str) -> requests.Response:
        return self.__fetch('PUT', f'/user/following/{username}')
