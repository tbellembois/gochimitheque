- [API examples](#api-examples)
  - [Authentication](#authentication)
    - [Get a JWT token](#get-a-jwt-token)
  - [Storages](#storages)
    - [Get a storage by ID](#get-a-storage-by-id)
  - [Units](#units)
    - [Get units](#get-units)

# API examples

## Authentication

### Get a JWT token

- request

```bash
curl "http://localhost:8081/get-token" \
  -X POST \
  -d "{\n  \"person_email\": \"jean@dupont@uca.fr\",\n  \"person_password\": \"mypassword\"\n}" \
  -H "content-type: application/json" 
```

- response 

`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGNoaW1pdGhlcXVlLmZyIiwiZXhwIjoxNjMxODI1MTM2LCJpZCI6MH0.jw7Pthq__x3G1CqLiZPsGEzqKOKNNkm1CMP8tAkzpgs`

## Storages

### Get a storage by ID

- request

```bash
curl "http://localhost:8081/storages/89" \
  -H "cookie: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGNoaW1pdGhlcXVlLmZyIiwiZXhwIjoxNjMyMTUyODI0LCJpZCI6MH0.949dV9LwDCyF64zzLHXziHaF0pS4ztKkMZWiyPTkicc" 
```

- response

```json
{
  "storage_id": {
    "Int64": 89,
    "Valid": true
  },
  "storage_creationdate": "2021-01-25T15:27:59.361361961Z",
  "storage_modificationdate": "2021-01-25T15:27:59.361362124Z",
 ...
  "storage_quantity": {
    "Float64": 24,
    "Valid": true
  },
  "storage_nbitem": 0,
 ...
  "storage_number_of_unit": {
    "Int64": 0,
    "Valid": false
  },
  "storage_number_of_bag": {
    "Int64": 0,
    "Valid": false
  },
  "storage_number_of_carton": {
    "Int64": 0,
    "Valid": false
  },
 ...
  "storelocation": {
    "storelocation_id": {
      "Int64": 7,
      "Valid": true
    },
    "storelocation_name": {
      "String": "Frigo",
      "Valid": true
    },
    "storelocation_canstore": {
      "Bool": false,
      "Valid": false
    },
...
  },
  "unit_quantity": {
    "unit_id": {
      "Int64": 2,
      "Valid": true
    },
    "unit_label": {
      "String": "mL",
      "Valid": true
    },
    "unit_type": {
      "String": "",
      "Valid": false
    },
    "unit": null
  },
...
  "supplier": {
    "c": 0,
    "supplier_id": {
      "Int64": 42,
      "Valid": true
    },
    "supplier_label": {
      "String": "ThermoFisher",
      "Valid": true
    }
  },
}
```

## Units

### Get units

- request

```bash
curl "http://localhost:8081/storages/units" \
  -H "cookie: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGNoaW1pdGhlcXVlLmZyIiwiZXhwIjoxNjMyMTUyODI0LCJpZCI6MH0.949dV9LwDCyF64zzLHXziHaF0pS4ztKkMZWiyPTkicc" 
```

- response 

```json
{
   "rows":[
      {
         "unit_id":{
            "Int64":14,
            "Valid":true
         },
         "unit_label":{
            "String":"nM",
            "Valid":true
         },
         "unit_type":{
            "String":"concentration",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":15,
            "Valid":true
         },
         "unit_label":{
            "String":"µM",
            "Valid":true
         },
         "unit_type":{
            "String":"concentration",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":16,
            "Valid":true
         },
         "unit_label":{
            "String":"mM",
            "Valid":true
         },
         "unit_type":{
            "String":"concentration",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":17,
            "Valid":true
         },
         "unit_label":{
            "String":"ng/L",
            "Valid":true
         },
         "unit_type":{
            "String":"concentration",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":18,
            "Valid":true
         },
         "unit_label":{
            "String":"µg/L",
            "Valid":true
         },
         "unit_type":{
            "String":"concentration",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":19,
            "Valid":true
         },
         "unit_label":{
            "String":"mg/L",
            "Valid":true
         },
         "unit_type":{
            "String":"concentration",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":20,
            "Valid":true
         },
         "unit_label":{
            "String":"g/L",
            "Valid":true
         },
         "unit_type":{
            "String":"concentration",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":21,
            "Valid":true
         },
         "unit_label":{
            "String":"%",
            "Valid":true
         },
         "unit_type":{
            "String":"concentration",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":22,
            "Valid":true
         },
         "unit_label":{
            "String":"X",
            "Valid":true
         },
         "unit_type":{
            "String":"concentration",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":1,
            "Valid":true
         },
         "unit_label":{
            "String":"L",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":2,
            "Valid":true
         },
         "unit_label":{
            "String":"mL",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":3,
            "Valid":true
         },
         "unit_label":{
            "String":"µL",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":4,
            "Valid":true
         },
         "unit_label":{
            "String":"kg",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":5,
            "Valid":true
         },
         "unit_label":{
            "String":"g",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":6,
            "Valid":true
         },
         "unit_label":{
            "String":"mg",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":7,
            "Valid":true
         },
         "unit_label":{
            "String":"µg",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":8,
            "Valid":true
         },
         "unit_label":{
            "String":"m",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":9,
            "Valid":true
         },
         "unit_label":{
            "String":"dm",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":10,
            "Valid":true
         },
         "unit_label":{
            "String":"cm",
            "Valid":true
         },
         "unit_type":{
            "String":"quantity",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":11,
            "Valid":true
         },
         "unit_label":{
            "String":"°K",
            "Valid":true
         },
         "unit_type":{
            "String":"temperature",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":12,
            "Valid":true
         },
         "unit_label":{
            "String":"°F",
            "Valid":true
         },
         "unit_type":{
            "String":"temperature",
            "Valid":true
         },
         "unit":null
      },
      {
         "unit_id":{
            "Int64":13,
            "Valid":true
         },
         "unit_label":{
            "String":"°C",
            "Valid":true
         },
         "unit_type":{
            "String":"temperature",
            "Valid":true
         },
         "unit":null
      }
   ],
   "total":22
}
```