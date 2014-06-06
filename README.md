
NichApp
=======
Powered by [Orchestrate](https://orchestrate.io), written with [gorc](https://github.com/orchestrate-io/gorc).
***
## Dependencies
### [github.com/hailiang/htmlquery](http://github.com/hailiang/htmlquery)
Scraper tool used to obtain titles of Nicholas Cage films listed on IMDB.
### [github.com/orchestrate-io/gorc](http://github.com/orchestrate-io/gorc)
Orchestrate client for golang.
***
## Contents
### populate/
**populate/get-nc-films.go** 
scrapes the [IMDB page for Nicholas Cage](https://www.imdb.com/name/nm0000115/), 
and writes each film title and year to a local file `nfs.txt` with the format *title,year*.
Each entry is separated by a newline.

**populate/q-omdb.go** 
reads the films from `nfs.txt` and uses them to query the [OMDb API](https://www.omdbapi.com).
The query results are printed to the console. 
This is a way to test data quality before initializing the database.
The code from this file is recycled in the initial population of the Orchestrate database.

**populate/pop-orc-init.go**
Queries OMDb API, 
but this time the resulting data are conditionally PUT to Orchestrate.
*Note:* 
There was an issue with PUTing the film "Face/Off,"
the forward slash of which was ultimately replaced with
a hyphen in the `nfs.txt` file.

**populate/pull-orc-data.go**
A good means of testing the initial population of the database.
GET data from orchestrate, and write the resulting JSON to local files.

**populate/pop-orc-graph.go**
Finds relations amongst films and PUTs them to Orchestrate.
Relations of interest are documented below.
***

***
## How the Data Are Organized
Visit the [Orchestrate API Docs](https://orchestrate.io/docs/api/?go) if you're lost.
### Key/Value
#### Keys are movie titles, and values return JSON with the following fields:

| FIELD      | (TYPE)        | DESCRIPTION      |
|------------|:-------------:|-----------------:|
| Year       | (int)         | Release year     |
| Director   | ([]string)    | Director(s)      |
| Genre      | ([]string)    | Genre(s)         |
| Rated      | (string)      | MPAA Rating      |
| ImdbRating | (float64)     | IMDB User Rating |
| Poster     | (string)      | URL for Poster   |
| ImdbID     | (string)      | IMDB Unique ID   |

Missing data are represented as "N/A", with the exception of ImdbRating, 
which is 0 if empty.

### Graph
#### Graph Relations are as follows:

| RELATION       | DESCRIPTION                                 |
|----------------|---------------------------------------------|
| s_Gen_1        | Shares 1 genre                              |
| s_Gen_2        | Shares 2 genres                             |
| s_Gen_3        | Shares 3 genres                             |
| s_Imdb_strong  | ImdbRating within .1 of each other          |
| s_Imdb_medium  | ImdbRating between .1 and .2 of each other  |
| s_Imdb_weak    | ImdbRating between .2 and .3 of each other  |