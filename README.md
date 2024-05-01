# Url Shortner Service

A URL shortener is a tool that takes a long URL and converts it into a shorter, more manageable version. It's an essential tool for anyone looking to share links efficiently and effectively. For example you take a long url like this 
https://www.marriott.com/molestie/sed/justo/pellentesque/viverra/pede.json?justo=eget&in=congue&blandit=eget&ultrices=sempe and shorten it like this https://www.link.ly/h33cqkX.

This project allow for individual and business manage their shortlinks as well detailed anyalytics by choosing the plan deemed fit for them

## Components
There are two services for this project which are 

* Rest Api Service
* Billing Service
* Key Generation Service

The rest api service exposes all endpoints for app functionality
The billing service is a cron job that generate revenues and management of pay schedules (creattion & deletion)
The Key generation service is a cron job that run to generate unused shorlink that can be assigned later for bulk links

Explain each service

## Features Implementation
* Shortlink Management
* Multiple Teams for business to manage users with role
* Detailed Link Analytics
* Bulk links Import
* Multiple custom domains
* Tags management for link
* Link Data Export
* Link cloaking


## Dependencies
* SQLite3
* Redis
* Docker (Optional)

## Usage
Make sure you have the Redis service and sqlite3 service running on the machine as specified as the dependencies for this project. Change the configuration as it suit you in the .env.example

To run locally, cd from the root to the backend directory and run 

```
go run main.go serveapi
```

To run on Docker, cd to the root and run the docker compose file 

```
docker-compose -f docker-compose.yml up -d
```

## Upcoming Features
Frontend implementation is coming soon. If you are interested in building it , please send an Email


