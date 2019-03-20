# Chimithèque

---

This is a work in progress !
Release planned before summer 2019.

---

Chimithèque is an open source chemical product management application initially initiated by the ENS-Lyon (France) and co-developped with the Université Clermont-Auvergne (France).

The project has started in 2015 and has moved to Github in 2017. The old subversion repository can be found here: https://sourcesup.renater.fr/scm/viewvc.php?root=chimitheque

    Why Chimithèque?

    We needed a global method to manage chemical products of the different departments and laboratories of the ENS to:

    - improve the security with a precise global listing of the chemicals products stored in the entire school
    - reduce waste by encouraging chemical products managers to search in Chimithèque if a product can be borrowed from another department before ordering a new one

![screenshot](screenshot.png)

# Quick start

You need a Linux AMD64 machine. No dependencies are required.

1. download the latest release from <https://github.com/tbellembois/gochimitheque/releases>
2. uncompress is in a directory
3. run `./gochimitheque`
4. open your web browser at `http://localhost:8081/login`

Et voilà !

Now login with the email `admin@chimitheque.com` and password `chimitheque`, and change the password immediatly.

# Production installation

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

- install the systemd script `doc/chimitheque.service` in `/etc/systemd/system` and enable/start it with `systemctl enable chimitheque.service; systemctl start chimitheque.service`

- install and adapt the apache2 configuration `doc/apache2-chimitheque.conf` in `/etc/apache2/site-available` and enable it with `a2esite apache2-chimitheque.conf`

# Database backup

Chimithèque uses a local sqlite database. You are encouraged to schedule regular plain text dump in a separate machine in case of disk failure.

# v1/v2 version

The v2 version has been rewritten in Golang.

- dramastically faster
- much easier to deploy (zero dependencies, embeded database)
- responsive design
- simplified GUI