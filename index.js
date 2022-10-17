const { Octokit } = require('octokit');
const { username, token } = process.env;

const octokit = new Octokit({ auth: token });

async function bootstrap() {
  const listFollowersForUser = octokit.paginate.iterator(octokit.rest.users.listFollowersForUser, {
    username,
    per_page: 100
  });
  const followers = new Set();
  for await (const { data: values } of listFollowersForUser) {
    for (const value of values) {
      followers.add(value.login);
    }
  }

  const listFollowingForUser = octokit.paginate.iterator(octokit.rest.users.listFollowingForUser, {
    username,
    per_page: 100
  });
  const followings = new Set();
  for await (const { data: values } of listFollowingForUser) {
    for (const value of values) {
      followings.add(value.login);
    }
  }

  const add = new Set([...followers].filter(x => !followings.has(x)));
  const remove = new Set([...followings].filter(x => !followers.has(x)));

  for (const value of add.values()) {
    await octokit.rest.users.follow({
      username: value
    });
  }
  for (const value of remove.values()) {
    await octokit.rest.users.unfollow({
      username: value
    });
  }

  console.log('add:', add, 'remove:', remove);
}

bootstrap();
