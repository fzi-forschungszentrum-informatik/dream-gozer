# DREAM - GoZer
[![js-standard-style](https://img.shields.io/badge/BSD%20license-3--Clause-green)](https://opensource.org/licenses/BSD-3-Clause)
[![js-standard-style](https://img.shields.io/badge/open-access-yellow)](https://open-access.network/)
[![DOI](https://zenodo.org/badge/DOI/10.5281/zenodo.1468906.svg)](https://doi.org/10.5281/zenodo.1468906)
![js-standard-style](https://img.shields.io/badge/research-prototype-red.svg)

GoZer is a prototypical database backend and part of the research project Data REseArch Mining (DREAM). It is used to evaluate new ideas on how Open Access publications can be made available on mobile devices. It provides access to an edited subset of open access publications from the [EconStor database](https://www.econstor.eu/) via a REST-like Web API. The backend is complemented by the corresponding android application [Ploc](https://github.com/fzi-forschungszentrum-informatik/dream-ploc).

Main features:

* Keyword-based interest profile
* Personalized publication feed
* Search for experts/authors in your field of expertise
* Transparent, lightweight feedback based on Ethereum
* Subset of the Econstore database with edited metadata

<p float="left">
  <img src="/images/ploc-screenshots-lowres.png" width="800" />
</p>

## Getting Started

The following step-by-step guideline describes how to build and deploy the GoZer via Docker under Linux Ubuntu 18.04 LTS Server.
We expect you have successfully installed [Linux Ubuntu 18.04 LTS Server](http://releases.ubuntu.com/18.04/) on your computer.
Then please update your Ubuntu installation first:

```
sudo apt update
sudo apt upgrade
reboot
```

### Prerequisites (Deployment)

To build and deploy GoZer you need:

* docker and docker.compose (tested with Docker version 18.09.9 and docker.compose 1.23.2)
* git with support for large file storage (LFS) (tested with git version 2.17.1)

Install Docker and Docker-Compose via Snap and provide user permissions:

```
sudo snap install docker
sudo addgroup --system docker
sudo adduser $USER docker
newgrp docker
sudo snap disable docker
sudo snap enable docker
```

Install and initialize git with support for large file storage (LFS):

```
sudo apt install -y git git-lfs
```

### Prerequisites (Local Build)

To build GoZer yourself you need the following tools:

* git with support for large file storage (LFS) (tested with git version 2.17.1)
* GNU make (tested with version 4.1)
* Go compiler (tested with version 1.10)
* SQLite 3 (tested with version 3.22)

First we install git, make, Go and SQLite3:

```
sudo apt install -y git git-lfs make golang sqlite3 
```

## Deployment via Docker-Compose

Clone the GoZer project to your server machine:

```
git lfs clone https://github.com/fzi-forschungszentrum-informatik/dream-gozer.git
cd dream-gozer
```

Edit file `Dockerfile.gozer` and specify your server's IP or FQDN under `ARG GOZER_URL_VAR=<YOUR SERVER IP>`.
The address is used to build the android client with your specific server as a target.

```
vim Dockerfile.gozer
```

After that you can build and start the backend as background process via docker-compose:

```
docker.compose up -d
```

Stop the backend via docker-compose:

```
docker.compose down
```

### Deploy to mobile client

You need an Android smartphone with at least Android 7.0 (API level 24).
The latest version of Android that we tested successfully is Android 9 (API level 28).
Enable [developer mode](https://developer.android.com/studio/debug/dev-options) on your mobile device.
The Android smartphone must be able to contact the your server's IP.

Beside deploying GoZer we also provide a [Ploc](https://github.com/fzi-forschungszentrum-informatik/dream-ploc) client specific to the servers address by an apk file.
This apk file can be downloaded to your Android device via Web directly.
Open the `http://<YOUR SERVER IP>:8080/download/ploc` from your mobile phones Web browser and download the apk file.
Go to your download directory and start the installation process by clicking on the apk file.
As this is a evaluation prototype you need to accept the installation from self-signed apps / unknown sources.

## Local Build

The following step-by-step guideline describes how to build and start the GoZer backend yourself (without Docker) under Linux Ubuntu 18.04 LTS Server.

Clone the GoZer project to your server machine:

```
mkdir -p ~/go/src/github.com/fzi-forschungszentrum-informatik && cd !$
git lfs clone https://github.com/fzi-forschungszentrum-informatik/dream-gozer.git
cd dream-gozer
```

Install required Go packages:

```
make install_packages
```

Extract demo database with publication metadata from EconStore:

```
tar -xzf storage.tgz
```

Make sure to build [Ploc](https://github.com/fzi-forschungszentrum-informatik/dream-ploc) client with your corresponding server IP (or FQDN) and copy that apk to the GoZer root directory.
Edit file `gozer.conf` and configure GoZer.
Define the name of the apk file via `ploc_apk = "ploc.apk"` that later can be downloaded on your mobile phone.
Disable ledger by setting `enable = false` (requires running Ganache installations).

```
vim gozer.conf
```

Run unit tests and check that there are no errors:

```
make test
```

Build GoZer backend:

```
make gozer
```

Run GoZer service:

```
./gozer -f gozer.conf
```

## Development

GoZer was developed with the [Go programming language](https://golang.org/) with version 1.10 in mind.
As database we use [SQLite3](https://sqlite.org/), because its single file can be easily deployed and because its features were sufficient for our evaluation scenarios.
In the following we give an overview to the project's folder structure and GoZer's architecture.

### Folder Structure

* / - Dockerfiles, build script, configuration example, main.go, example database, license information and this readme
* /config - implementations related to the configuration file
* /model - defines the data types (aka data model) used in GoZer
* /model/ploc - defines the message types used to communicate with the mobile client
* /storage - query functions to the local database (SQLite3)
* /storage/ledger - query functions to store feedback in a [Solidity](https://solidity.readthedocs.io/en/v0.5.3/) contract
* /webapi - implementation of a REST-like API for client-server communication

### Package Dependencies

GoZer uses the following packages:

* [github.com/google/uuid](https://github.com/google/uuid) - to generate globaly unique IDs for user handling
* [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - to interface the SQLite3 database file
* [github.com/gorilla/mux](https://github.com/gorilla/mux) - to demux incoming HTTP request to their corresponding handler functions
* [golang.org/x/crypto/bcrypt](https://golang.org/x/crypto/bcrypt) - to compute and store password hashes
* [github.com/BurntSushi/toml](https://github.com/BurntSushi/toml) - to load and unmarshal a TOML-based configuration file
* [github.com/ethereum/go-ethereum](https://github.com/ethereum/go-ethereum) - to interact with a Solidity-based smart contract

Thanks to all the developers that made these packages available!

### Database Scheme

![Database Scheme](/images/database-scheme.png)

## License

This open source software is licensed under the 3-Clause BSD License - see the [LICENSE.md](LICENSE.md) file for details.
The publication metadata in our [demo database](storage.tgz) is taken with permission from the [EconStor](https://www.econstor.eu/) and it is free of
known copyright and database restrictions and in the public domain throughout the world ([CC0 - No Rights Reserved](https://creativecommons.org/share-your-work/public-domain/cc0/)).

## Acknowledgements

The DREAM project and its source code are the result of a joint project of the [Humboldt Institute for Internet and Society](https://www.hiig.de/) (HIIG) and the [FZI Research Center for Information Technology](https://www.fzi.de/). The project was funded by the [Federal Ministry of Education and Research Germany](https://www.bmbf.de/) (BMBF).

We also like to thank:

* [Bielefeld Academic Search Engine](https://base-search.net/) and [EconStor](https://www.econstor.eu/) for giving access to their metadata collections
* [Ganache](https://www.trufflesuite.com/ganache) for the possibility to easily deploy an Ethereum testbed
* [Miniflux](https://github.com/miniflux/miniflux) for its clean and inspiring open source code

... and all the open source developers for their time and their great work!

## Future Work

The prototype was used to evaluate different concepts.
It is not intended for use in production environments.
For production environment one may:

* Add HTTPs support
* Switch to PostgreSQL
* Add automatic harvesting and preparation of publication metadata
* Improve test coverage
* Migrate to Go modules
* Make more database operations atomic
* Use blockchain simulator for tests instead of running ganache node
