# nuri-cms

**nuricms is a api first content management system written in go.**

---


# Docker 

You can find every dockerfile in folder `/docker`. 

## Build the container

If you want to build one of the dockerfiles you need to enter: 

`docker build -t nuricms -f path/to/dockerfile .`

For example: 
`docker build -t nuricms -f docker/nuricms-sqlite/Dockerfile .`

## Run the container

`docker run -p 8080:8080 -it -e JWT_SECRET=my-verysuper-secret-secret-32byteslong nuricms`