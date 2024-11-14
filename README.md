# GATOR (Blog AggreGATOR)

### Depends on go and a Postgres DB
Clone the Repo and go install

Create a file called ".gatorconfig.json" in the users home directory
This file should contain the following;

```
{
    "db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
    "current_user_name":"jameson"
}
```

Use `gator` to run any of the following commands;

* login - change the user to your username
* register - add a new user
* reset - DO NOT USE, deletes all user and blog data from the DB
* users - lists currently registered users
* agg - this command should be left running in a background terminal, scrapes the included feeds for new articles
* addfeed - add a new feed to the aggregator
* feeds - list the currently tracked feeds for all users
* follow - follow a new feed
* following - list feeds the current user is following
* unfollow - un follow a feed for the current user
* browse - get recent new articles from followed feeds

