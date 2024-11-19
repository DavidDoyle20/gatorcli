# Dependencies
---
* Golang
* Postgres

# How to Install
---
* Install the latest versions of go and postgres
    * Mac OS: `brew install postgresql@16`
    * Linux / WSL (Debian): `sudo apt install postgresql postgresql-contrib`
* If you are on windows I recommend using wsl and following the linux installation steps
* Open the cli of your choice and enter the command `go install github.com/DavidDoyle20/gatorcli@latest`
* Create the postgres db
    1. Install postgres if you havent done so already
    2. Ensure the installation worked entering the command `psql --version`
    3. (Linux only) Set postgres passowrd `sudo passwd postgres`
    4. Start the postgres server in the background 
        * Mac: `brew services start postgresql`
        * Linux: `sudo service postgresql start`
    5. Enter the psql shell
        * Mac: `psql postgres`
        * Linux: `sudo -u postgres psql`
    6. You should see a prompt that looks like `postgres=#`
    7. Create a new database. I called mine gator `CREATE DATABASE gator;`
    8. Connect to the new database `\c gator`
    9. You should see a prompt that looks like `gator=#`
    10. (Linux only) Set the password `ALTER USER postgres PASSWORD 'postgres';`
    11. Query the database to see if everything is working `SELECT version();`
    12. If everything is working you can type `exit` to leave the psql shell
    13. Get your connection string, here are examples
        * Format: protocol://username:password@host:port/database
		* Mac OS (no password or username): postgres://david:@localhost:5432/gator
		* Linux & Windows: postgres://postgres:password@localhost:5432/gator
    14. Test that your connection string works by running the command `psql postgres://postgres:postgres@localhost:5432/gator`
        * Replace the url with your connection string
    15. Finally run the command `gatorcli dburl postgres://postgres:postgres@localhost:5432/gator`
        * Replace the url with your connection string
