<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Chimithèque](#chimithèque)
  - [Team](#team)
  - [Screenshot](#screenshot)
- [Web browser compatibility](#web-browser-compatibility)
- [Links](#links)
- [Requirements](#requirements)
- [Quick start](#quick-start)
  - [With Docker (recommended)](#with-docker-recommended)
  - [Without Docker](#without-docker)
  - [Connection](#connection)
- [Production installation](#production-installation)
  - [The Docker way (recommended)](#the-docker-way-recommended)
    - [Configuration](#configuration)
  - [The classic way](#the-classic-way)
    - [Installation](#installation)
    - [Configuration](#configuration-1)
  - [Binary command line and Docker parameters](#binary-command-line-and-docker-parameters)
- [Database backup](#database-backup)
- [Initial product database import](#initial-product-database-import)
  - [Principle](#principle)
  - [Importing from another public instance](#importing-from-another-public-instance)
- [Upgrades](#upgrades)
  - [Classic installation](#classic-installation)
  - [Docker installation](#docker-installation)
- [Support](#support)
- [Use of categories and tags](#use-of-categories-and-tags)
- [Use of barecode and QRCode](#use-of-barecode-and-qrcode)
- [List of public database Chimithèque instances](#list-of-public-database-chimithèque-instances)

<!-- markdown-toc end -->

# Chimithèque

Chimithèque is an open source *chemical product, biological reagent and lab consumables* management application started by the ENS-Lyon (France) and co-developped with the Université Clermont-Auvergne (France). It is written in *Golang*.

The project was started in 2015 and has moved to Github in 2017.

## Team

- *projet leaders*: Delphine Pitrat - ENS-Lyon (delphine[dot]pitrat[at]ens-lyon[dot]fr) - Thomas Bellembois - UCA
- *technical referent - chemistry*: Delphine Pitrat  
- *technical referent - biology*: Antoine Goisnard Phd - University Clermont-Auvergne / IMOST lab. (antoine[dot]goinard[at]uca[dot]fr)  
- Marie Depresle - University Clermont-Auvergne / Biorcell3D: *biology specialist*
- Manon Roux - University Clermont-Auvergne / Biorcell3D: *chemistry specialist*

## Screenshot

![screenshot](screenshot.png)

# Web browser compatibility

Chimithèque may NOT work with Microsoft Internet Explorer/Edge.  
It was tested successfully with Firefox and Chrome/Chromium.

# Links

- chimithèque binary: <https://github.com/tbellembois/gochimitheque/releases>

- docker image: <https://hub.docker.com/repository/docker/tbellembois/gochimitheque>

# Requirements

- a *Linux AMD64* machine
- an SMTP server (for password recovery - optionnal for the quick start and if using LDAP)

Chimithèque is statically compiled and then does not require other dependencies.

# Quick start

## With Docker (recommended)

Look at the most recent tagged version at <https://hub.docker.com/r/tbellembois/gochimitheque>.

Do *NOT* use the `latest` tag as it is the development version.

Do *NOT* use this method in production.

For example fo the `2.0.8` version:
```bash
  docker run --name chimitheque -v /tmp:/data -p 127.0.0.1:8081:8081 -e CHIMITHEQUE_DOCKERPORT=8081 tbellembois/gochimitheque:2.0.8
```

## Without Docker

1. download the latest `gochimitheque` binary here <https://github.com/tbellembois/gochimitheque/releases/latest/download/gochimitheque>
2. uncompress it in a directory
3. run `./gochimitheque`

## Connection

Then open your web browser at `http://127.0.0.1:8081/chimitheque`

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

Edit the file, it is self documented. 

Create the data directories for the Nginx and Chimithèque containers:
```bash
  mkdir -p /data/docker-nginx/nginx-auth/certs
  mkdir -p /data/docker-nginx/nginx-templates
  mkdir -p /data/docker-chimitheque/chimitheque-db
  chmod o+rwx /data/docker-chimitheque
```

Retrieve the Nginx configuration:
```bash
  wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/documents/system/nginx-chimitheque.conf -O /data/docker-nginx/nginx-templates/default.conf.template
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

Run `gochimitheque --help` for a list of command line parameters. 
 
For a Docker installation use the corresponding environment variables commented in the `docker-compose.yml` file. They are self explainatory.
Note that some command line parameters are not mapped to Docker environment variable. This is expected.

> example:
>
> `gochimitheque -proxyurl=https://appserver.foo.fr -proxypath=/chimitheque/ -admins=john.bar@foo.fr,jean.dupont@foo.fr -mailserveraddress=smtp.foo.fr -mailserverport=25 -mailserversender=noreply@foo.fr"`
>
> will run the appplication behind a proxy at the URL `https://appserver.foo.fr/chimitheque` with 2 additionnal administrators `john.bar@foo.fr` and `jean.dupont@foo.fr`

Note about admins:

A static administrator `admin@chimitheque.fr` is created during the installation. His password must be changed after the first connection.

You can add a comma separated list of admins emails. Accounts must have been created in the application before. You should limit the number of admins and set entity managers instead.

> example: `-admins=john.bar@foo.com,jean.dupont@foo.com`

# Database backup

Chimithèque uses a local *sqlite* database. You are strongly encouraged to schedule regular plain text dump in a separate machine in case of disk failure.

You can backup the database with:
```bash
  sqlite3 /path/to/chimitheque/storage.db ".backup '/path/to/backup/storage.sq3'"
```

Restore it with:
```bash
  cp /path/to/backup/storage.sq3 /path/to/chimitheque/storage.db
```

# Initial product database import

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

# Use of categories and tags

For chemical and biological reagents, there is now the possibility to class products in different categories in order to make easier product research.

This solution is available when creating a new product sheet with a scrolling menu and suggest different preregistered product categories. It is possible to create a new category if concerned product does not feet with already existing suggestions. This solution allows in main menu, with advanced research, to show only products called with a specific category, and thus have a global vision on a specific class of products. 

This solution is completed with the possibility to apply tags on chemical or biological reagents, also available in product sheet section. This allows to associate a product with various fields, methods, protocols, projects, or application domains. 
Like previously, preregistered tags are proposed in a scrolling menu with the possibility to create new tags. For example, a stem cell culture medium can be associated with Stem Cells, Cell Culture or Culture Medium tags. This function may reveal particularly useful to rapidly show products associated with a specific activity, projects or method in the advanced research of Chimithéque main menu. Moreover, it is a way to personalize and adapt product research according to a lab or a structure specific needs or habits. 

# Use of barecode and QRCode

A new option is now available for creating an association between a product and a specific label: the QRCode. 
It is different from the bare-code, because it is readable by every device which have a camera and permits to access directly to the page with the product's storage. 
By default, when a product is stocked, the software create a random bare-code and a new QRcode. 
However, if a product need to be sampled, you can check the option "identical bare-code" when the number of samples is required, and it will generate the same bare-code for each new sample. 
The major advantage is that you can scan any QRcode of these strictly identical products and it will display the page of the storage with all the samples. 
Then, any of these samples could be borrowed or archived, for example. 
For instance, for conservation conditions, it could be recommended to limit freeze-thaw cycles. 
To avoid that, the product could be sampled in different dishes with the same volume or mass. 
To store them on Chimitheque, the "identical bare-code" option will permit to create QRcodes linked with all the samples, so that any of them could be destocked when one of them is used. 

# List of public database Chimithèque instances

- ENS de Lyon: `https://chimitheque.ens-lyon.fr`

If you want to share your product database please send an email to the mailing list or create a Github issue.
