# [WIP] wikint: Interest Comparison
This program collects data for interest comparisons and presents it in a neat api.

## Concepts
### profile database
To build the database, it crawls facebook user profiles and posts.


### interest database
To build the interest database, it scrapes wikipedia for nouns


### comparison api
The comparison api combines the collected data, so you can compare which interents corolate to other interests.


## Usage
install:
```sh
go install github.com/verassitnh/wikint
```
run the data pipelines:
```
wikint dp
```
then start up the api in another terminal:
```
wikint api
```


### api endpoints
/interests - all interests with interest rating

/interests/:interest - get most interesting topics for users interested in :interest

