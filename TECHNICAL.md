# technical notes

## permissions

possible permissions:
`r`, `w`, `all` on each item

possible items:
- all
- products
- rproducts
- storages
- entities

- `entity_id` = -1 : all items
- `entity_id` = -2 : any/one of item -> -2 is never stored in the db but when used the item_id field is ommited requesting the permission table

possible combinations for a given person:

| item_name (ex: storage) | perm_name (ex: r) | entity_id (ex: 2) | notes |
| :-- | :--: | --: | :-- |
| `all`       |     `all`     |  ? | manager  ex: `all` permission on all items of entity with id `3` |
| `all`       |     `all`     | -1 | super admin |
| `all`       |     `all`     | -2 | manager of at least one entity |
| `all`     |   ?    |      ? | ex: `r` on `all` items of entity with `id` 3 |
| `all`     |   ?    |      -1| ex: `r` permission on all items of all entities |
| `all`     |   ?    |      -2| *not used* ex: `r` permission on all items of at least one entity|
| ?     |   `all`    |      ? | ex: `all` permission on `storage` of entity with `id` 3 |
| ?     |   `all`    |      -1| ex: `all` permission on `storage` of all entities |
| ?     |   `all`    |      -2| *not used* ex: `all` permission on `storage` of at least one entity |
| ?     |   ?    |   -1 | ex: `r` permission on all `storage` on all entities |
| ?     |   ?    |   ?  | ex: `r` permission on `storage` on entity with `id` 3 |
| ?     |   ?    |   -2  | *not used* ex: `r` permission on `storage` on at least one entity |

permissions on other items (suppliers, classes of compounds...) are linked to the items above

entities / people management:
- only super admins can create new entities and modify entities
- entity managers can create/update/delete people in their entities

## authorization

A custom middleware. Look at `func (env *Env) AuthorizeMiddleware(h http.Handler) http.Handler ` in `handlers\auth.go`.

## static content

JS install/upgrade:

*these 3 modules are version dependent*
npm install jquery
npm install @popperjs/core
npm install bootstrap

npm install bootstrap-colorpicker
npm install bootstrap-table
npm install jquery-validation
npm install select2
npm install --save-dev @fortawesome/fontawesome-free
npm install @mdi/font
npm install animate.css --save
npm install print-js --save
npm install pako
npm install --save qr-scanner

rsync -av ./node_modules/bootstrap/dist/js/bootstrap.min.js ./static/js/
rsync -av ./node_modules/bootstrap/dist/js/bootstrap.min.js.map ./static/js/
rsync -av ./node_modules/bootstrap-colorpicker/dist/js/bootstrap-colorpicker.min.js ./static/js/
rsync -av ./node_modules/bootstrap-colorpicker/dist/js/bootstrap-colorpicker.min.js.map ./static/js/
rsync -av ./node_modules/bootstrap-table/dist/bootstrap-table.min.js  ./static/js/
rsync -av ./node_modules/jquery/dist/jquery.min.js ./static/js/
rsync -av ./node_modules/jquery-validation/dist/jquery.validate.min.js ./static/js/
rsync -av ./node_modules/jquery-validation/dist/additional-methods.min.js  ./static/js/jquery.validate.additional-methods.min.js
rsync -av ./node_modules/@popperjs/core/dist/umd/popper.min.js ./static/js/
rsync -av ./node_modules/@popperjs/core/dist/umd/popper.min.js.map ./static/js/
rsync -av ./node_modules/select2/dist/js/select2.full.min.js ./static/js/
rsync -av ./node_modules/print-js/dist/print.js ./static/js/
rsync -av ./node_modules/pako/dist/pako.min.js ./static/js/
rsync -av ./node_modules/qr-scanner/qr-scanner.umd.min.js* ./static/js/
rsync -av ./node_modules/qr-scanner/qr-scanner-worker.min.js* ./static/js/

rsync -av ./node_modules/bootstrap/dist/css/bootstrap.min.css ./static/css/
rsync -av ./node_modules/bootstrap/dist/css/bootstrap.min.css.map ./static/css/
rsync -av ./node_modules/@fortawesome/fontawesome-free/css/all.min.css ./static/css/fontawesome.all.min.css
rsync -av ./node_modules/@fortawesome/fontawesome-free/webfonts/* ./static/webfonts/
rsync -av ./node_modules/bootstrap-colorpicker/dist/css/bootstrap-colorpicker.min.css ./static/css/
rsync -av ./node_modules/bootstrap-colorpicker/dist/css/bootstrap-colorpicker.min.css.map ./static/css/
rsync -av ./node_modules/bootstrap-table/dist/bootstrap-table.min.css ./static/css/
rsync -av ./node_modules/@mdi/font/css/materialdesignicons.min.css ./static/css/
rsync -av ./node_modules/@mdi/font/css/materialdesignicons.min.css.map ./static/css/
rsync -av ./node_modules/@mdi/font/fonts/* ./static/fonts/
rsync -av ./node_modules/select2/dist/css/select2.min.css ./static/css/
rsync -av ./node_modules/animate.css/animate.min.css ./static/css/
rsync -av ./node_modules/print-js/dist/print.css  ./static/css/
rsync -av ./node_modules/print-js/dist/print.map ./static/css/

## testers

delphine.pitrat@ens-lyon.fr, laurelise.chapellet@ens-lyon.fr, guillaume.george@ens-lyon.fr, sylvain.david@ens-lyon.fr, laure.guy@ens-lyon.fr, yann.bretonniere@ens-lyon.fr, loic.richard@irstea.fr, christophe.le-bourlot@insa-lyon.fr, julien.devemy@uca.fr, alix.tordo@ens-lyon.fr, clement.courtin@ens-lyon.fr

## API

Get admin access token and save access and refresh token in variables:
```
http --body --form POST http://keycloak:8080/keycloak/realms/chimitheque/protocol/openid-connect/token client_id=chimitheque client_secret=mysupersecret grant_type=password username=admin@chimitheque.fr password=chimitheque | jq '.["access_token","refresh_token"]' > /tmp/keycloak_token
access_token=$(head -1 /tmp/keycloak_token)
refresh_token=$(tail -1 /tmp/keycloak_token)
rm /tmp/keycloak_token
```

Create and entity:
```
http --print HBhb POST http://localhost:8081/entities entity_name=entity1 "Cookie:access_token=$access_token;refresh_token=$refresh_token"
```

Create a store location:
```
http --print HBhb POST http://localhost:8081/store_locations entity[entity_id]:=1 entity[entity_name]=entity1 store_location_name=storelocation1B store_location_can_store:=true "Cookie:access_token=$access_token;refresh_token=$refresh_token"
```

Create a person:
Set an entity manager:
