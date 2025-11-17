# Chimithèque

Chimithèque is an open source *chemical product, biological reagent and lab consumables* management application started by the ENS-Lyon (France) and co-developped with the Université Clermont-Auvergne (France). It is written in *Golang* and *Rust*.

The project was started in 2015 and has moved to Github in 2017.

## Team

- *projet leaders*: Delphine Pitrat - ENS-Lyon (delphine[dot]pitrat[at]ens-lyon[dot]fr) - Thomas Bellembois - UCA (thomas[dot]bellembois[at]uca[dot]fr)
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

Here is the list of the major technical changes from the `2.1.0` version:
- the only supported installation is with Docker (this may change in the future)
- the authentication is based on OpenID managed by the Keycloak application.
- the LDAP configuration has been removed from Chimithèque to be managed by the OpenID server

# Requirements

- a *Linux AMD64* machine (glibc 2.29 min)
- [Docker](https://docs.docker.com/get-docker/) and [docker-compose](https://docs.docker.com/compose/install/)
- an SMTP server
- an HTTPS certificate
- the sqlite command line tool if upgrading an existing installation

# Upgrading from 2.0.8

Important: if you upgrade to a `2.1.*` version coming from a `2.0.*` version you *must* first perform the upgrades to the `2.0.8` version.

1. Backup your *entire* installation folder.

## Export users

1. Retrieve the latest release of the `chimitheque_people_keycloak_exporter` binary from <https://github.com/tbellembois/chimitheque_people_keycloak_exporter/releases>.

2. Copy the binary in your *current* Chimithèque installation (where the `storage.db` file is).

3. Run the binary:
```
chmod +x chimitheque_people_keycloak_exporter
./chimitheque_people_keycloak_exporter
```

The exporter will create a `keycloak.json` file. Keep it for later use.

Note that the exporter will panic if your database contains duplicate emails.

## Migrate database

1. From your *current* Chimithèque installation (where the `storage.db` file is), retrieve the `sql` files

```bash
wget https://raw.githubusercontent.com/tbellembois/chimitheque_db/refs/heads/main/src/resources/shema.sql
wget https://raw.githubusercontent.com/tbellembois/chimitheque_db/refs/heads/main/src/resources/migration.sql
```

2. Run the migration script

```bash
sqlite3 chimitheque.sqlite < shema.sql && sqlite3 chimitheque.sqlite < migration.sql
```

This will create a new `chimitheque.sqlite` file. Keep it for later.

# Installation

1. Retrieve the Chimithèque `docker-compose.yml` and `compose-prod.env` files:
```bash
cd /root
wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/docker-compose.yml -O docker-compose.yml
wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/compose-prod.env -O .env
```

2. Edit the `.env` file, it is self documented.

3. Create the `data` directory (and sub directories) for the container data.
```bash
mkdir -p /root/docker/keycloak
mkdir /data
mkdir -p /data/docker-keycloak/templates/
mkdir -p /data/docker-nginx/nginx-templates/
mkdir -p /data/docker-nginx/nginx-conf/
mkdir -p /data/docker-nginx/nginx-auth/certs/
mkdir -p /data/docker-chimitheque/chimitheque-db/
```
> If you want to choose another directory you will have to replace the `/data` strings in the `docker-compose.yml` file (`volumes` sections). In this documentation we assume that the default directory is kept.

4. Retrieve the containers configuration files:
```bash
wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/docker/keycloak/Dockerfile -O /root/docker/keycloak/Dockerfile
wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/docker/keycloak/chimitheque-realm-template.json -O /data/docker-keycloak/templates/chimitheque-realm-template.json
wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/docker/keycloak/chimitheque-users-0.json -O /data/docker-keycloak/templates/chimitheque-users-0.json
wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/docker/nginx/default.conf.template -O /data/docker-nginx/nginx-templates/default.conf.template
wget https://raw.githubusercontent.com/tbellembois/gochimitheque/master/docker/nginx/nginx.conf -O /data/docker-nginx/nginx-conf/nginx.conf
```

5. Copy your https certificate `crt` and `key` files in `/data/docker-nginx/nginx-auth/certs/`. Your certificate *must* contain the certification chain.

6. If you upgrade from a previous version copy the `chimitheque.sqlite` file in `/data/docker-chimitheque/chimitheque-db/`.

7. Configure Nginx, edit the `/data/docker-nginx/nginx-templates/default.conf.template` file. The sections to edit are spotted with the `# CONFIGURE:` string.

8. Start up and Wait a moment (it can take several minutes for the containers to start):
```bash
docker compose up -d
```

# Configuration

## Admin user creation

1. Connect to the OIDC server at <https://your_chimitheque_url/keycloak> with the username `admin@chimitheque.fr` and the value of your `KEYCLOAK_ADMIN_PASSWORD` for password.

2. On the top left corner drop-down list choose the `chimitheque` realm. Then click on `Users` on the left column.

3. Click the `Create new user` button and enter the following informations:
```
Email verified: yes
Email: admin@chimitheque.fr
```
And click the `Create` button.

4. Then click on the `Credentials` tab and the `Set password` button. Enter the value of your `KEYCLOAK_ADMIN_PASSWORD`, uncheck `Temporary` and click on `Save`.

## Importing previous users (migration from 2.0.8 only)

If you migrate from a `2.0.8` version you should have a `keycloak.json` file from the `Upgrading from 2.0.8` section.

1. Click on `Realm settings` on the left colums, then on the top right drop-down list `Action` choose `Partial import`.

2. Browse and upload your `keycloak.json` file, click on the `Choose the resources your want to import` checkbox and click on the `Import` button.

## Setup the smtp configuration

1. Click on `Realm settings` on the left colums, then the `Email` tab.

2. Fill in the required information.

## Additionnal configuration

You mail want to enable/disable user registration as well as activate the LDAP connectivity. Please refer to the [Keycloak documentation](https://www.keycloak.org/docs/latest/server_admin/#_ldap) for more information.

# Connection

Then open your web browser at `https://your_chimitheque_url`

Et voilà !

Now login with the email `admin@chimitheque.fr` and the value of your `KEYCLOAK_ADMIN_PASSWORD` password.

# Administrators

A static administrator `admin@chimitheque.fr` with id `1` is created during the installation. It is hardcoded and must not be deleted.
You can add a comma separated list of admins emails. You should limit the number of admins and set entity managers instead. Always keep the `admin@chimitheque.fr` account.
Non existing accounts will be created.

> example: `-admins=admin@chimitheque.fr,john.bar@foo.com,jean.dupont@foo.com`

# Users management

Users permissions are still managed in the Chimithèque application (by admins and entity managers). But user creation and deletion are now managed by the embeded Keycloak application.
There are two ways to manage users:
1. enable user registration in Keycloak (easiest way)
People will have the possibility to create their own account but will NOT be able to connect to Chimithèque until they are affected to an entity.
2. disable user registration in Keycloak (harder way)
You will have to create users manually in Keycloak. Currently only the account `admin@chimitheque.fr` can access Keycloak, not other admins nor Chimithèque managers. This will be fixed in a next release.

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
