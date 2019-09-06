<img src='https://github.com/okkur/reposeed/blob/master/media/logo.svg' width='500'/>

Hosted version for [Reposeed](https://github.com/okkur/reposeed) repository base file generation

 [![state](https://img.shields.io/badge/state-beta-blue.svg)]() [![release](https://img.shields.io/github/release/okkur/reposeed.svg)](https://github.com/okkur/reposeed/releases) [![license](https://img.shields.io/github/license/okkur/reposeed.svg)](LICENSE)

**NOTE: This is a work-in-progress, we do not consider it production ready. Use at your own risk.**

# Reposeed Server
Hosted version to simplify usage of Reposeed from a web interface.

## Using Reposeed Server
```
go get github.com/okkur/reposeed-server
```

After installation create an .env file
```
echo 'PORT=":8080"' >> .env
echo 'STORAGE="./storage/zips/"' >> .env # Path to your storage folder for zip files
```

To start the server run ```reposeed-server```.

For Reposeed server to run correctly the .env file is necessary.
Take a look at our full [documentation](/docs).

## Support
For detailed information on support options see our [support guide](/SUPPORT.md).

## Helping out
Best place to start is our [contribution guide](/CONTRIBUTING.md).

----

*Code is licensed under the [Apache License, Version 2.0](/LICENSE).*  
*Documentation/examples are licensed under [Creative Commons BY-SA 4.0](/docs/LICENSE).*  
*Illustrations, trademarks and third-party resources are owned by their respective party and are subject to different licensing.*

*The Reposeed logo was created by [Florin Luca](https://99designs.com/profiles/florinluca)*

---

Copyright 2017 - The RepoSeed-server authors
