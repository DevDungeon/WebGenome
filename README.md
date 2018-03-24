# Web Genome

## Overview

A breadth first web crawler that stores HTTP headers in a MongoDB database with
a web front end all written in Go. http://www.webgno.me is a https://www.devdungeon.com project.

## Website

* [http://www.webgeno.me](http://www.webgeno.me)
* [DevDungeon.com Web Genome project page](http://www.devdungeon.com/content/web-genome)


## Setup

### Create a Linux user

	sudo useradd webgenome

### Set up the Go environment

Set up your GOPATH to be /home/webgenome/go

	# In ~/.bashrc
	export GOPATH=/home/webgenome/gospace

### Get the packages

    # All at once with
    go get github.cm/DevDungeon/WebGenome...
	
	# Or individually
    go get github.com/DevDungeon/WebGenome
    go get github.com/DevDungeon/WebGenome/core
    go get github.com/DevDungeon/WebGenome/website
    go get github.com/DevDungeon/WebGenome/worker_http
    
### Setting up database

Create a MongoDB database and seed it with a domain. Add an index on the name field to really speed things up:

	sudo apt install mongodb
	mongo
	> use webgenome
	> db.domains.insert({'name':'www.devdungeon.com'})
	> db.domains.createIndex({name:1})

#### Sample database queries
	
	db.getCollectionNames()
	db.showCollections()
	db.domains.getIndexes()
	db.domains.stats()
	db.domains.count()
	db.domains.find({name:'www.devdungeon.com'})
	db.domains.count({lastchecked:{$exists:true}, skipped: null})
	db.domains.find({headers: {$elemMatch: {value: {$regex: 'Cookie'}}}}).pretty()
	db.domains.find({headers: {$elemMatch: {key: {$regex: 'Drupal'}}}}).pretty()

### Run website using systemd

The systemd directory contains a sample service file that can be used to run the website as a service.
	
	sudo cp /home/webgenome/go/src/github.com/DevDungeon/WebGenome/systemctl/webgenome.service /etc/systemd/system/
	sudo chown root:root /etc/systemd/system/webgenome.service
	
	sudo vim /etc/systemd/system/webgenome.service # Double check settings
	
	systemctl webgenome enable
	systemctl webgenome start

### Nginx reverse proxy

The web server will listen on port 3000 by default.
Access it directly or set up a reverse proxy with nginx like this:

	# /etc/nginx/conf.d/webgenome.conf
	server {  # Redirect non-www to www
		listen 80;
		server_name webgeno.me;
		return 301 $scheme://www.webgeno.me$request_uri;
	}
	server {
		listen 80;
		server_name www.webgeno.me;
		location / {
			proxy_set_header X-Real-IP $remote_addr;
			proxy_pass http://localhost:3000;
		}
	}

### Running worker_http

Here is an example usage of running the crawler:

	worker_http --host=localhost --database=webgenome --collection=domains --max-threads=4 --http-timeout=30 --batch-size=100 --verbose

## Updating

	# Update the source and executables
    go get -u github.com/DevDungeon/WebGenome...
	
	# Restart the service
	systemctl restart webgenome

## Source Code

* [WebGenome (GitHub.com)](https://www.github.com/DevDungeon/WebGenome)

## Contact

support@webgeno.me

## License

GNU GPL v2. See LICENSE.txt.

## Notes

You can kill the http_worker at any time and restart it without causing any problems.
If you run multiple instances of the worker at the same time it will end up checking
a lot of the domains multiple times. If you want to crawl more just increase the
number of threads to the worker. I was able to run it with 256 threads on a small
Linode computer.

The website has a hard-coded static directory currently and should be run with
the current working directory of website/. There are multiple database connections
also hard-coded in the website.go file. Yeah, yeah... it needs to be refactored
and dried up.

The worker is not run as a service because it may fill up your disk space and it
should be run in verbose mode at the beginning so you can tune and make sure
it's not hammering nested subdomains on a single site.

## Changelog

v1.1 - 2018/03/23 - Clean up files, sync disjointed repo, and relaunch
v1.0 - 2016/11/18 - Initial stable release

## Screenshots

![Screenshot of worker](screenshots/worker_http.png)

![Screenshot of website](screenshots/website.png)
