# Buying Frenzy

Table of Content:
- [Buying Frenzy](#buying-frenzy)
  - [Prerequisite](#prerequisite)
  - [How to Start](#how-to-start)
    - [ETL](#etl)
    - [Webserver](#webserver)
  - [API Documentation](#api-documentation)
    - [Restaurant List](#restaurant-list)
    - [Process Purchase History](#process-purchase-history)


## Prerequisite

- `docker` or `podman`
- `psql`
- `go 1.18`
- `make`
- `curl`

## How to Start

All services/apps is run under docker. First we need to run `make docker-run-devdb` (it will be run postgres database docker), this step is mandatory


### ETL

Run `make docker-run-migrate` for migrating data to database

> First we extract and transform json file into several csv file, then `COPY` (import) all csv file into database

### Webserver

Run `make docker-run-webserver` for starting web server API

## API Documentation

This documentaion for webserver API

### Restaurant List

This endpoint is used for get list of restaurant based on applied filter

```
GET /restaurants?[filter]
```

You can combine every filter listed below separated by `&`

| Filter             | Type   | Explaination                                                                                      | Example                  |
| ------------------ | ------ | ------------------------------------------------------------------------------------------------- | ------------------------ |
| `ndishes`          | number | How many dish menu in restaurant                                                                  | `ndishes=10`             |
| `dish_price_start` | float  | How much dish price low in restaurant (note: this filter value must be below `dish_price_end`)    | `dish_price_start=80.2`  |
| `dish_price_end`   | float  | How much dish price high in restaurant (note: this filter value must be above `dish_price_start`) | `dish_price_start=100.2` |
| `name`             | string | Which resturant having name `name`                                                                | `name=cowfish`           |
| `day`              | string | Which restaurant open on day `day`                                                                | `day=monday`             |
| `start_time`       | string | When restaurant will be open (note: in `hh:mm` 24H format)                                        | `start_time=13:00`       |
| `end_time`         | string | When restaurant will be close (note: in `hh:mm` 24H format)                                       | `start_time=22:00`       |

### Process Purchase History

This endpoint is used for processing purchased history, it will calculate user purchased history, deduct balance from user, and add new balance to restaurant

```
GET /process/{user_id}
```


