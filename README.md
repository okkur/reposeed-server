<img src='https://github.com/okkur/reposeed/blob/master/media/logo.svg' width='500'/>

A web application for [RepoSeed](https://github.com/okkur/reposeed)

 [![state](https://img.shields.io/badge/state-beta-blue.svg)]() [![release](https://img.shields.io/github/release/okkur/reposeed.svg)](https://github.com/okkur/reposeed/releases) [![license](https://img.shields.io/github/license/okkur/reposeed.svg)](LICENSE)

**NOTE: This is a work-in-progress, we do not consider it production ready. Use at your own risk.**

# RepoSeed-server
An easy way to generate documentation, license files and more boilerplate to get you started from your first commit

## Using RepoSeed-server
```
go get -v -u github.com/okkur/reposeed-server
```  
When you installed reposeed-server, create a .env file and write the following information in it
```
PORT=":8080"
STORAGE="./storage/zips/" # Path to your storage folder for zip files
```
And after that just run ```reposeed-server``` command.   
if you didn't set GOBIN path, just run the following commands
```
cd $GOPATH/src/github.com/okkur/reposeed-server
go run main.go
```
Remember that you should always have a .env file in your current directory if you want to run reposeed-server   
Take a look at our full [documentation](/docs).

## Support
For detailed information on support options see our [support guide](/SUPPORT.md).

## Helping out
Best place to start is our [contribution guide](/CONTRIBUTING.md).

----

*Code is licensed under the [Apache License, Version 2.0](/LICENSE).*  
*Documentation/examples are licensed under [Creative Commons BY-SA 4.0](/docs/LICENSE).*  
*Illustrations, trademarks and third-party resources are owned by their respective party and are subject to different licensing.*

---

Copyright 2017 - The RepoSeed-server authors
