# Dependencies
---
* Golang
* Postgres

# How to Install
---
* Install the latest versions of go and postgres
    * Mac OS: `brew install postgresql@16`
    * Linux / WSL (Debian): `sudo apt install postgresql postgresql-contrib`
* Open the cli of your choice and enter the command `go install github.com/DavidDoyle20/gatorcli@latest`
* Set the dburl
    * Replace the url with your connection string
    * Mac OS: `gatorcli dburl postgres://postgres:@localhost:5432/gator`
    * Linux: `gatorcli dburl postgres://postgres:postgres@localhost:5432/gator`

# Commands
---
* **Login**: Takes a username and logs into that user
* **Register**: Creates a new user with the given username and logs into it
* **Reset**: Removes all users
* **Users**: Displays all users
* **Agg**: Takes a duration and updates all of the feeds after each duration
* **Addfeed**: Takes a name and a url and creates a new feed and adds it to the current users feeds
* **Feeds**: Displays all feeds that exist
* **Follow**: Takes a feed url from the existing feeds and adds it to the current users feeds
* **Following**: Displays all feeds the current user is following
* **Unfollow**: Takes a feed url from the feeds a user follows and unfollows that user
* **Browse**: Takes an optional limit (default 2) and displays that amount of feeds that the current user is following
* **Dburl**: Takes a connection string and updates the config
