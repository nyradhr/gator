# gator

Gator is an RSS feed aggregator built in Go as a CLI tool.

## 🔧 Requirements

- **Go 1.23.0 or later** - [Download here](https://golang.org/dl/)
- **PostgreSQL** - Any recent version should work
- Ensure `$GOPATH/bin` is in your PATH environment variable


## 📦 Installation

```sh
go install github.com/nyradhr/gator@latest
```


## ⚙️ Setup

1. Create a PostgreSQL database named gator (or your preferred name)

2. Create a config file named .gatorconfig.json in your home directory:

```json
{"db_url":"postgres://username:password@localhost:5432/database_name?sslmode=disable","current_user_name":""}
```

Replace "username", "password", and "database_name" with your PostgreSQL credentials.


## 🧪 Usage


### Register a new user
```sh
gator register <username>
```
Adds a new user and sets it as the current user.

### List all users
```sh
gator users
```

### Login 
```sh
gator login <username>
```
Login with the given username and set it as current user.

### Add and follow feeds
```sh
gator addfeed <name> <feed_url>
```
Adds a new feed to the aggregator and automatically adds it to the "followed" list for the current user.

### Follow feed
```sh
gator follow <feed_url>
```
Used for already existing feeds that the current user is not yet following.

### Unfollow feed
```
gator unfollow <feed_url>
```

### List all feeds
```sh
gator feeds
```

### List all followed feeds
```sh
gator following
```

### Start aggregating posts
```sh
gator agg <time_interval>
# Example: gator agg 30s
```
Starts a loop that will cycle the feeds and show their posts with the specified time interval.

### Browse posts
```sh
gator browse [limit]
```
Show a list of posts from followed feeds, ordered by most recent publication date. The number of posts, if not specified, defaults to 2.

### Delete all users
```sh
gator reset
```
⚠️ Warning: This command deletes ALL users and their data. Use with caution!
