# go-hog
Go Client for mailhog:
https://github.com/mailhog/MailHog

MailHog is a convenient test utility for email integration 

This client provides a simple way to interact with Mailhog's API. To use it, ensure Mailhog is running. You can quickly 
set it up using the provided docker-compose file. Run in the project root:

`docker compose up`

Then test it is working by running; 

`go test`

Go install: 

`go get github.com/latebitflip-io/go-hog`