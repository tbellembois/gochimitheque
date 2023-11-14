# Chimithèque

Chimithèque is an open source *chemical product, biological reagent and lab consumables* management application started by the ENS-Lyon (France) and co-developped with the Université Clermont-Auvergne (France). It is written in *Golang*.

The project was started in 2015 and has moved to Github in 2017.

## Team

- *projet leaders*: Delphine Pitrat - ENS-Lyon (delphine[dot]pitrat[at]ens-lyon[dot]fr) - Thomas Bellembois - UCA
- *technical referent - chemistry*: Delphine Pitrat  

## Screenshot

![screenshot](screenshot.png)

## Web browser compatibility

Chimithèque may NOT work with Microsoft Internet Explorer/Edge.  
It was tested successfully with Firefox and Chrome/Chromium.

# Links

- project home: <https://github.com/tbellembois/gochimitheque>
- docker image: <https://hub.docker.com/repository/docker/tbellembois/gochimitheque>

# 2.1.0 News!

Here is the list of the major changes from the `2.1.0` version:
- the only supported installation is with docker
- the authentication is based on OpenID
- the LDAP configuration has been removed from Chimithèque to be managed by the OpenID server
- the initial database is automatically populated from another Chimithèque instance

# Requirements

- a *Linux AMD64* machine
- [Docker](https://docs.docker.com/get-docker/) and [docker-compose](https://docs.docker.com/compose/install/)
- an SMTP server

# Installation

Retrieve the Chimithèque `docker-compose.yml` file:
```bash
  wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/docker-compose.yml -O docker-compose.yml
```

Change the `services > casdoor > environment > origin` URL.
Change the `services > chimitheque > environment > CHIMITHEQUE_APPURL` URL.
Change the `services > chimitheque > environment > CHIMITHEQUE_OIDCISSUER` URL.

The three URLs must be the same. This is the base URL of your Chimithèque instance.

Create the data directories for the containers:
```bash
  mkdir -p /data/docker-casdoor/casdoor-db
  mkdir -p /data/docker-casdoor/casdoor-init
  mkdir -p /data/docker-nginx/nginx-auth/certs
  mkdir -p /data/docker-nginx/nginx-templates
  mkdir -p /data/docker-chimitheque/chimitheque-db
  chmod o+rwx /data
```

Retrieve the Nginx configuration:
```bash
  wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/config/default.conf.template -O /data/docker-nginx/nginx-templates/default.conf.template
```

TODO edit

Retrieve the CasDoor configuration:
```bash
  wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/config/init_data.json -O /data/docker-casdoor/casdoor-init/init_data.json
```

TODO edit

Start up:
```bash
  docker-compose up -d
```

## Connection

Then open your web browser at `http://your_chimitheque_url/chimitheque`

Et voilà !

Now login with the email `admin@chimitheque.fr` and password `chimitheque`, and change the password immediatly.

# Administrators

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

# Upgrades

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

This solution is available when creating a new product card with a drop down menu and suggest different preregistered product categories. 
It is possible to create a new category if concerned product does not feet with already existing suggestions. This solution allows in main menu, with advanced research, to show only products called with a specific category, and thus have a global vision on a specific class of products. 

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
