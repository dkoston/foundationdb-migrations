# foundationdb-migrations

## Foreword

Not every wants to use migrations on key/value stores. For those who do, this
database migration tool handles migrations for FoundationDB

### Examples of Why

For example, you have a service which relies on some information to be in the 
database such as a list of exchanges to pull prices from and the symbols that 
are listed on those exchanges. 

- In this case, we add a migration to add those exchanges and symbols

Then, you support more exchanges and symbols:

- We add another migration to add those, and maybe disable one who's API is 
broken.

We want this to programatically happen and want the migrations to be launched
before new code that needs them is running.

On Kubernetes, we run something like [dumb-init](https://github.com/Yelp/dumb-init)
as our container command where first we run the migrations and then the API
server. With a healthcheck on the API server container, the load balancer won't
flip traffic over to the new containers until they are all live (migrations run)
.

## License

See [MIT-License.md](./MIT-License.md)

## Install

    clone this repo into $GOPATH/github.com/dkoston/foundationdb-migrations
    run: `make install`

## Configuration and command line options (all commands under Usage)

### fdb.cluster

By default `fdbm` will look in `./db/fdb.cluster` for a FoundationDB cluster 
file. 

Alternatively, you can pass in the cluster file location with 
`-f <path/to/file.name>` when calling `fdbm`.

For example:

`fdbm -f /conf/fdb.cluster up`


### FDB Version

The fdb client requires you to define the API version. By default, we use `510`.
To specify a version, add it as a comment to `fdb.cluster` as defined above or
you can pass in the fdb version with `-v <version>` when calling `fdbm`. 

For example:

`fdbm -v 510 up`


### Migration files location

By default, `fdbm` will look in `./db/migrations` for migrations. To specify a
different directory when calling `fdbm` use `-m </path/to/migrations>`.

For example:

`fdbm -m /configs/database/migrations up`


## Usage

### Creating migrations

You should create migrations using the `fdbm` command so they contain all the
necessary features and take advantage of new features. Do not copy/paste old
migration file contents.

`fdbm create <migration_name>`

After creation, you need to edit the contents to actually do something.

#### Embedded migrations

Feel free to import `github.com/dkoston/fdbm/lib/fdbm` into your applications
and use the public APIs.

### Running migrations

When running migrations, `fdbm` will look at the db/migrations directory and run
any migrations which have not been previously run.

To do so, run:

`fdbm up`

### Rolling back migrations

WARNING: not every migration can be easily rolled back. If you have changed the
db and your code, you will have to roll back both your db and the code. Also, 
the ability to roll back migrations purely depends on your ability to write a 
migration function that will return data to its previous state. Unlike SQL where
data is structured, your application may have altered the data in such a way
that it cannot be rolled back (hopefully not). USE AT YOUR OWN RISK.

`fdbm down`

The above command will rollback the last applied migration. To rollback multiple
, you will have to run the command multiple times.


### Re-running a failed migration

For convience, you can rollback and re-run the last migration with:

`fdbm redo`

### Check the status of migrations

`fdmb status`

    $ fdbm: migration status
    $    Status         Date                      File
    $   =========================================================================
    $   Success        20180529155230           db/migrations/20130106222315_and_again.go
    $   Rolled back    20180529160754           db/migrations/20180529152851_test2.go


You will see the timestamp of when applied migrations were run against the db.
Anything listed as "Pending" has not yet been run.


### Get the last migration number run against the database

`fdbm dbversion`

    $ fdbm dbversion
    $ fdbm: dbversion 003

This will print the number of the last migration run. i.e. (003_alter_the_data_from_one.go)

## Migration file format

The file should contain 2 functions, one named `XXXX_Up()` and one named 
`XXXX_Down()` where `XXXX` is the migration number.

The function named `XXXX_Up()` will be run with `fdbm up` and then function
named `XXXX_Down()` will be run with `fdbm down`.

You may have other functions in the file but those are the two that `fdbm` will
look for an run automatically with `up` and `down`.

### Example migration file

The following file will set a key/value pair to `hello: world` on `fdbm up` and 
then removethat key/value pair on `fdbm down`:

```go
package main

import (
    "github.com/apple/foundationdb/bindings/go/src/fdb"
)


func Up_20130106222315(t fdb.Transactor) error {
    _, err := t.Transact(func (tr fdb.Transaction) (ret interface{}, err error) {
        tr.Set(fdb.Key("hello"), []byte("world"))
        return
    })

    return err
}

func Down_20130106222315(t fdb.Transactor) error {
    _, err := t.Transact(func (tr fdb.Transaction) (ret interface{}, err error) {
        tr.Clear(fdb.Key("hello"))
        return
    })

    return err
}
```



# Contributors

Thank you!

* Dave Koston (dkoston)

# Previous Contributors

Thanks to those who built goose (which fdbm is based heavily on)

* Josh Bleecher Snyder (josharian)
* Abigail Walthall (ghthor)
* Daniel Heath (danielrheath)
* Chris Baynes (chris_baynes)
* Michael Gerow (gerow)
* Vytautas Šaltenis (rtfb)
* James Cooper (coopernurse)
* Gyepi Sam (gyepisam)
* Matt Sherman (clipperhouse)
* runner_mei
* John Luebs (jkl1337)
* Luke Hutton (lukehutton)
* Kevin Gorjan (kevingorjan)
* Brendan Fosberry (Fozz)
* Nate Guerin (gusennan)