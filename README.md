# Go Gator 🐊

Go Gator is a command-line RSS feed aggregator written in Go. It allows you to manage users, subscribe to RSS feeds, and browse the latest posts from your subscriptions, all from the comfort of your terminal.

## Prerequisites

Before you begin, you will need to have the following installed on your system:

- **Go**: The programming language this application is built with. You can find installation instructions at [go.dev/doc/install](https://go.dev/doc/install).
- **PostgreSQL**: A powerful, open-source relational database used to store all application data. You can find installation instructions on the [official PostgreSQL website](https://www.postgresql.org/download/).

## Installation

You can install the `gator` CLI directly using the `go install` command:

```bash
go install github.com/daemosity/go-gator@latest
```

This command will download the source code, compile it, and place the `gator` executable in your `$GOPATH/bin` directory. Ensure this directory is in your system's `PATH` to run the command from anywhere.

## Configuration

Go Gator uses a JSON file named `.gatorconfig.json` located in your user's home directory (`~/`) to store its configuration.

### 1. Create the Config File

First, you need to create the configuration file. You can do this with the following command:

```bash
touch ~/.gatorconfig.json
```

### 2. Add Configuration Details

Open the newly created `~/.gatorconfig.json` file in a text editor and add the following content.

```json
{
  "db_url": "postgresql://user:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

You must update the `db_url` with your own PostgreSQL connection string. The format is `postgresql://<user>:<password>@<host>:<port>/<dbname>?sslmode=disable`. The `current_user_name` will be populated automatically when you use the `register` or `login` command.

### 3. Set up the Database

The application uses the `goose` tool to manage database migrations.

First, install `goose`:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Next, from the root of the project directory, run the migrations to set up the database schema. You will need to have your `db_url` from the config file handy.

```bash
# It is recommended to set this as an environment variable
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgresql://user:password@localhost:5432/gator?sslmode=disable"

# Run the migrations
goose -dir sql/schema up
```

This command will create the necessary tables (`users`, `feeds`, `posts`, etc.) in your database.

## Usage

Once installed and configured, you can start using Go Gator. Here are some of the main commands:

### Register a New User

To start, you need to register a new user. This command will also log you in automatically.

```bash
gator register my_username
```

### Login as an Existing User

If you have already registered, you can log in with:

```bash
gator login my_username
```

### Add and Follow a New Feed

To add a new RSS feed to the system and automatically follow it:

```bash
gator addfeed "The Go Blog" "https://go.dev/blog/feed.xml"
```

### View Your Followed Feeds

To see a list of all the feeds you are currently following:

```bash
gator following
```

### Browse Recent Posts

To see the latest posts from your followed feeds:

```bash
# Show the 5 most recent posts
gator browse 5
```

### List All Users

To see all users registered in the system, with your current user highlighted:

```bash
gator users
```

Enjoy using Go Gator!
