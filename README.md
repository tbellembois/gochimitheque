# Chimithèque

---

The Beta version is currently been tested at the ENS-Lyon. 
Stable release planned in September 2019.

---

Chimithèque is an open source chemical product management application initially initiated by the ENS-Lyon (France) and co-developped with the Université Clermont-Auvergne (France).

The project has started in 2015 and has moved to Github in 2017.

Main goals:

- improve the security with a precise global listing of the chemicals products stored in the entire school
- reduce waste by encouraging chemical products managers to search in Chimithèque if a product can be borrowed from another department before ordering a new one

![screenshot](screenshot.png)

# Quick start (to test the application)

You need a Linux AMD64 machine. No dependencies are required.

1. download the latest release from <https://github.com/tbellembois/gochimitheque/releases>
2. uncompress is in a directory
3. run `./gochimitheque`
4. open your web browser at `http://localhost:8081/login`

Et voilà !

Now login with the email `admin@chimitheque.com` and password `chimitheque`, and change the password immediatly.

# Production installation

## Requirements

- linux (Chimithèque can be cross compiled to run on Windows but binaries are not provided and this documentation does not cover this situation)
- an SMTP server (for password recovery)

## Installation

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

- configure (look at the following section) and install the systemd script `doc/chimitheque.service` in `/etc/systemd/system` and enable/start it with `systemctl enable chimitheque.service; systemctl start chimitheque.service`

- install and adapt the apache2 configuration `doc/apache2-chimitheque.conf` in `/etc/apache2/site-available` and enable it with `a2esite apache2-chimitheque.conf`

# Binary command line parameters

You need configure the systemd script with the following parameters:

- `-port`: application listening port - default = `8081`
- `-proxyurl`: application base URL with no trailing slash - default = `http://localhost:8081`
- `-proxypath`: application path - default = `/`
- `-mailserveraddress`: SMTP server address - *REQUIRED*
- `-mailserverport`: SMTP server port - *REQUIRED*
- `-mailserversender`: SMTP server sender email - *REQUIRED*
- `-mailserverusetls`: use an SMTP TLS connection - default = `false`
- `-mailservertlsskipverify`: skip SSL verification - default = `false`
- `-admins`: comma separated list of administrators emails
- `-logfile`: output log file - by default logs are sent to stdout
- `-debug`: debug mode, do not enable in production

# Application administrators

A static administrator `admin@chimitheque.fr` is created during the installation. His password must be changed after the first connection.

You can add add additional administrators with the `admins` command line parameters.

# Database backup

Chimithèque uses a local sqlite database. You are strongly encouraged to schedule regular plain text dump in a separate machine in case of disk failure.

You can backup the database with:
```bash
    sqlite3 /path/to/chimitheque/storage.db ".backup '/path/to/backup/storage.sq3'"
```

# Chimithèque v1 database migration

// TODO

# v1/v2 version

The v2 version has been rewritten in Golang.

- dramastically faster
- much easier to deploy (zero dependencies, embeded database)
- responsive design
- simplified GUI