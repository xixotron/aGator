# agreGator from Boot.Dev project

This is an RSS agregator from [boot.dev](https://boot.dev)'s course guided project/course
[Build a Blog Aggregator in Go](https://www.boot.dev/courses/build-blog-aggregator-golang)

## Requirements

To be able to run this aggregator you'll need:

* Install Posgres, with a dabase for this agregator
* Go complier in order to build the executable

## Installation

In order to install agreGator in your system youll need to use [go install](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies)
`go install https://github.com/xixotron/aGator`
or if want to download the source first:

```bash
cd aGator
go install .
```

## Configuration

First create a database in your postgres server using a client
(ex: [psql](https://www.postgresql.org/docs/current/app-psql.html))
`CREATE DATABASE <dbname>;` we suggest "gator" as the name.

Then the necesary tables need to be created inside the new database,
for this we have provided in [sql/schema](/sql/schemal) migration files to
configure the database.

We used [Goose](https://github.com/pressly/goose) to apply these migration
files to the database.
`goose postgres <connection_string> up`

The [connection string](https://www.geeksforgeeks.org/postgresql/postgresql-connection-string/)
also needs to be saved to a configuration file named `.gatorconfig.josn` in your
user folder `~/.gatorconfig.json` which

```json
{
  "db_url": "postgres://<user>:<password>@<host>:5432/<dbname>"
}
```

it's recommended to dissable ssl when using localhost as the host `...<dbname>?sslmode=disable"}`

## Usage

To use the proram you will need to create a user:

```bash
aGator adduser <name>
```

Configure the feeds the user follows:

```bash
aGator addfeed "<feed_name>" <url>
```
feed_name should be a friendly name for the feed we want to follow.

Then after we have feeds added we should run the aggregator's feed parsing loop,
this log to the terminal as it parses feeds, `Ctrl+C` should stop the loop:

```bash
aGator agg <time_between_reqests>
```

The time should be a "duration" ex: `5m` for 5 minutes, or `1h` for 1 hour.
We strongly advise against going bellow `30s` to avoid "DDoSing" the rss servers, or being banned from them for excessive requests.

Lastly as a user who is following the feeds `aGator login <user_name>` now we can
read the last entries to the feeds:

```bash
aGator browse [limit]
```
