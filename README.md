- [Chimithèque](#chimithèque)
- [Web browser compatibility](#web-browser-compatibility)
- [Download](#download)
- [Requirements](#requirements)
- [Quick start](#quick-start)
- [Production installation](#production-installation)
  - [The Docker way (recommended)](#the-docker-way-recommended)
    - [Configuration](#configuration)
  - [The classic way](#the-classic-way)
    - [Installation](#installation)
    - [Configuration](#configuration-1)
  - [Binary command line and Docker parameters](#binary-command-line-and-docker-parameters)
  - [Setting up application administrators](#setting-up-application-administrators)
- [Database backup](#database-backup)
- [Chimithèque V2 initial database import](#chimithèque-v2-initial-database-import)
  - [Principle](#principle)
  - [Importing from another public instance](#importing-from-another-public-instance)
- [Chimithèque V1 database migration](#chimithèque-v1-database-migration)
  - [Read me first !](#read-me-first-)
  - [Export](#export)
    - [PostgreSQL](#postgresql)
  - [Import](#import)
- [Upgrades](#upgrades)
  - [Classic installation](#classic-installation)
  - [Docker installation](#docker-installation)
- [Support](#support)
- [V1/V2 version](#v1v2-version)
- [List of public database Chimithèque instances](#list-of-public-database-chimithèque-instances)
- [Get the latest development compiled version](#get-the-latest-development-compiled-version)
- [Compile from sources](#compile-from-sources)

# Chimithèque

Chimithèque is an open source *chemical product management* application started by the ENS-Lyon (France) and co-developped with the Université Clermont-Auvergne (France). It is written in *Golang*.

*projet leader*: Delphine Pitrat (delphine[dot]pitrat[at]ens-lyon[dot]fr)

The project has started in 2015 and has moved to Github in 2017.

Main goals:
- *simplicity*: do one think (stores products) but do it well
- *security*: provide a global listing of the chemicals products storages
- *cost/ecology*: share chemical products to avoid waste

![screenshot](screenshot.png)

# Web browser compatibility

Chimithèque does NOT work with Microsoft Internet Explorer/Edge.  
It was tested successfully with Firefox and Chrome/Chromium.

# Download

Chimithèque releases can be downloaded here: <https://github.com/tbellembois/gochimitheque/releases>.

Download the `gochimitheque` binary (in the assets section), not the source code archive.

Permanent link to the latest release: <https://github.com/tbellembois/gochimitheque/releases/latest/download/gochimitheque>

# Requirements

- a *Linux AMD64* machine with `Glibc2.28` minimum
- an SMTP server (for password recovery - optionnal for the quick start)

Chimithèque is statically compiled and then does not require other dependencies.

# Quick start

1. download the latest `gochimitheque` binary here <https://github.com/tbellembois/gochimitheque/releases/latest/download/gochimitheque>
2. uncompress is in a directory
3. run `./gochimitheque`
4. open your web browser at `http://localhost:8081/login`

Et voilà !

Now login with the email `admin@chimitheque.fr` and password `chimitheque`, and change the password immediatly.

Do *not* use this mode in production.

# Production installation

## The Docker way (recommended)

Install [Docker](https://docs.docker.com/get-docker/) and [docker-compose](https://docs.docker.com/compose/install/).

Retrieve the Chimithèque `docker-compose.yml` file:
```bash
  wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/docker-compose.yml -O docker-compose.yml
```

Create the data directories for the Nginx and Chimithèque containers:
```bash
  mkdir -p /data/docker-nginx/nginx-auth/certs
  mkdir -p /data/docker-nginx/nginx-templates
  mkdir -p /data/docker-chimitheque/chimitheque-db
```

Retrieve the Nginx configuration:
```bash
  wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/documents/system/nginx-chimitheque.conf -O /data/docker-nginx/nginx-templates/nginx-chimitheque.conf.template
```

Start up:
```bash
  docker-compose up -d
```

### Configuration

You may need to adapt the default Docker configuration to enable SSL, change the proxy URL or setup the SMTP server adress.

This must be done:
- in the `docker-compose.yml` file with the environment variables
- in the `nginx-chimitheque.conf.template` file

## The classic way

### Installation

It is strongly recommended to run Chimithèque behind an HTTP proxy server with SSL.

- create a dedicated user

```
    groupadd --system chimitheque
    useradd --system chimitheque --gid chimitheque
    mkdir /usr/local/chimitheque
```

- drop the gochimitheque binary into the `/usr/local/chimitheque` directory

- setup permissions

```
    chown -R  chimitheque:chimitheque /usr/local/chimitheque
```

### Configuration

- configure (look at the following section) and install the systemd script `documents/system//chimitheque.service` in `/etc/systemd/system` and enable/start it with `systemctl enable chimitheque.service; systemctl start chimitheque.service`

- install and adapt the Nginx configuration `documents/system/nginx-chimitheque.conf` in `/etc/nginx/server-available/nginx.conf` and link it with `ln -s /etc/nginx/server-available/nginx.conf /etc/nginx/server-enable/nginx.conf`

## Binary command line and Docker parameters

The following parameters can be passed to the Chimithèque binary.  
For a Docker installation use the environment variables commented in the `docker-compose.yml` file.
Note that some command line parameters are not mapped to Docker environment variable. This is expected.

- `-listenport`: application listening port - default = `8081`
- `-proxyurl`: proxy base URL with no trailing slash
- `-proxypath`: proxy path - default = `/`
- `-mailserveraddress`: SMTP server address - *REQUIRED*
- `-mailserverport`: SMTP server port - *REQUIRED*
- `-mailserversender`: SMTP server sender email - *REQUIRED*
- `-mailserverusetls`: use an SMTP TLS connection - default = `false`
- `-mailservertlsskipverify`: skip SSL verification - default = `false`
- `-enablepublicproductsendpoint`: enable public products endpoint - default = `false`
- `-admins`: comma separated list of administrators emails that must be present in the database
- `-logfile`: output log file - by default logs are sent to stdout
- `-debug`: debug mode, do not enable in production

One shot commands:
- `-resetadminpassword`: reset the `admin@chimitheque.fr` admin password to `chimitheque`
- `-updateqrcode`: regenerate the storages QR codes

> example:
>
> `gochimitheque -proxyurl=https://appserver.foo.fr -proxypath=/chimitheque/ -admins=john.bar@foo.fr,jean.dupont@foo.fr -mailserveraddress=smtp.foo.fr -mailserverport=25 -mailserversender=noreply@foo.fr"`
>
> will run the appplication behind a proxy at the URL `https://appserver.foo.fr/chimitheque` with 2 additionnal administrators `john.bar@foo.fr` and `jean.dupont@foo.fr`

## Setting up application administrators

A static administrator `admin@chimitheque.fr` is created during the installation. His password must be changed after the first connection.

You can add additional administrators with the `-admins` command line parameters. Note that those admins *must already be present* in the database.

> example: `-admins=john.bar@foo.com,jean.dupont@foo.com`

# Database backup

Chimithèque uses a local *sqlite* database. You are strongly encouraged to schedule regular plain text dump in a separate machine in case of disk failure.

You can backup the database with:
```bash
    sqlite3 /path/to/chimitheque/storage.db ".backup '/path/to/backup/storage.sq3'"
```
# Chimithèque V2 initial database import

## Principle

Each Chimithèque application administrator can share its products database with `-enablepublicproductsendpoint`.
Note that only products informations (product cards) will be shared.

You need at least one other public Chimithèque instance to be able to populate your new Chimithèque database.

## Importing from another public instance

```bash
    ./gochimitheque -importfrom=[publicInstance]
```

example:
```bash
    ./gochimitheque -importfrom=https://chimitheque.ens-lyon.fr
```

# Chimithèque V1 database migration

## Read me first !

Product and storage archives and history are NOT imported into Chimithèque V2.  
If you need to keep those informations keep a Chimithèque V1 instance for archive purposes.

## Export

Databases of the V1 version must be exported into `CSV` (with headers) to be imported in the V2 version.

### PostgreSQL

In PostgreSQL this can be done with the command:

```bash
    SCHEMA="public"; DB="{chimitheque-db-name}"; psql -U {chimitheque-user} -h {chimitheque-host} -p {chimitheque-port} -Atc "select tablename from pg_tables where schemaname='$SCHEMA'" $DB | while read TBL; do psql -U {chimitheque-db-name} -h {chimitheque-host} -p {chimitheque-port} -c "COPY $SCHEMA.$TBL TO STDOUT WITH CSV HEADER" $DB > $TBL.csv; done;
```

This will generate one CSV file per table.

## Import

You can then import to the V2 version with:

```bash
    /path/to/gochimitheque -proxyurl=https://appserver.foo.fr -importv1from=/path/to/csv
```

This is important to specify the correct `-proxyurl` parameter as it will be used to generate the storages qr codes.

# Upgrades

## Classic installation

Stop Chimithèque and replace the `gochimitheque` binary with the new one.

## Docker installation

```bash
    docker-compose down
    docker-compose pull
    docker-compose up -d --force-recreate
```

You may want to install [watchtower](https://github.com/containrrr/watchtower) to perfom automatic upgrades.

# Support

Please do not (never) contact the members of the Chimithèque development team directly.

Subscribe to the mailing list: <https://groupes.renater.fr/sympa/subscribe/chimitheque?previous_action=info> or open a Github issue.

# V1/V2 version

The v2 version has been rewritten in Golang.

- dramastically faster
- much easier to deploy (zero dependencies, embeded database)
- responsive design
- simplified GUI

# List of public database Chimithèque instances

- ENS de Lyon: `https://chimitheque.ens-lyon.fr`

If you want to share your product database please send an email to the mailing list or create a Github issue.

# Get the latest development compiled version

You can retrieve the latest auto-compiled version from the trunk. 
This is strongly not recommended in a production environment except if you know what you do.

Go in the [action](https://github.com/tbellembois/gochimitheque/actions) tab of the project. Click on the last (top of the list) successfull workflow. Download the `chimitheque` artifacts zip. The `gochimitheque` binary is in the zip file. 

# Compile from sources

Install [Go](https://golang.org/doc/install).

Install `go1.12.16` for GopherJS:
```bash
    go get golang.org/dl/go1.12.16
    go1.12.16 download
    export GOPHERJS_GOROOT=/root/sdk/go1.12.16
```

Go get the code:
```bash
    go get github.com/tbellembois/gochimitheque
```

You will find the source code in the `$GOPATH/src/github.com/tbellembois/gochimitheque` and the binary in `$GOPATH/bin/gochimitheque`.

You can modify/recompile the code with:
```bash
  cd $GOPATH/src/github.com/tbellembois/gochimitheque
  go generate
  go build
```