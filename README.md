# gumtree-auto

Automatically log in to gumtree.com.au and repost any expiring ads. 

Reads the username and password from ENV. 

Once built and tested, just schedule the docker image to via `cron` daily etc.

# Run

```sh
$ docker run -e USERNAME=<username> -e PASSWORD=<password> ghcr.io/mcoops/gumtree-auto:latest

2024/02/24 02:36:34 Logging into gumtree
2024/02/24 02:36:43 Logged in successfully!
2024/02/24 02:36:43 Gathering ads which need reposting
2024/02/24 02:36:51 No ads required reposting
```

# Docker build

```sh
$ docker build -t gumtree-auto .

$ docker run -e USERNAME=<username> -e PASSWORD=<password> gumtree-auto
```
