# CRAG Policy Exceptions Tool

This is a tool intended to handle requests the Research Computing team get for higher storage quotas, long or dedicated queue access, or shared storage pools.

It used to be in the shared `go-clustertools` repository, but it was split out to make it easier to develop and test.

## Usage and Config

When correctly deployed, you should have:

 - a single statically-compiled binary called `exceptions`
 - a JSON configuration file, that is assumed to be in `~/.exceptions_db.conf` (an alternative can be used with the `--config=` option)

### Setting Up the Tool for Yourself

The exceptions CLI tool should be deployed on all our clusters in the `/shared/ucl/apps/cluster-bin` directory, which is included in the `userscripts` module. You just need to copy or create the `~/.exceptions_db.conf` file. Examples for the admin and normal user are in `ccspapp`'s home directory on Myriad:

```
mysql-exceptions-db-admin.conf  
mysql-exceptions-db-user.conf
```

You shouldn't need the `-admin` credentials unless you're planning to destroy the database.

### JSON Config File Format

Format is as follows for MySQL:

```json
{
   "db_type": "mysql",
   "db_connection_string": "my_username:my_password@tcp(mysql.rc.ucl.ac.uk:3306)/my_database_name"
}
```

Or for a local SQLite DB:

```json
{
  "db_type": "sqlite3",
  "db_connection_string": "/some_directory/some_file"
}
```

The DB type can currently only be `mysql` or `sqlite3` (if you want it to work), but if it became necessary a relatively simple addition could add other databases.

These are parameters passed directly to the GORM library's `gorm.Open` function, so you might want to check the documentation there for more comprehensive information: <http://gorm.io/docs/connecting_to_the_database.html>


## From-Scratch Setup

### Building the Tool

If you already have a Golang development environment, you may want to clone the repository into the appropriate location, and then run the `build.sh` script in the root of the repository.

Otherwise, you probably want to set the `INSTALL_PATH` environment variable and run the `install.sh` script in the root of the repository, which will create a temporary Go environment to build the tool. The default install path is the one for our clusters: `/shared/ucl/apps/cluster-bin`.

### The Database

To create the MySQL setup, you will need:

 - an empty database (ours is called `service_exceptions`)
 - a user with create table permissions
 - you might also want a lower-access user for everyday stuff

In our setup the users are called `servex_admin` and `servex_user`.

Create a config file appropriately, then run `exceptions createdb`. This will create the necessary tables.

Then run `exceptions examples` and/or `exceptions --help` for further help.

## Making a Back-Up

The entire DB contents can be exported and imported as JSON:

```
exceptions dumpjson | gzip >db_jsondump.$(date +%Y-%m-%d).json.gz
gzip -cd db_jsondump.2019-07-15.json.gz | exceptions importjson 
```

If you import to a non-empty database, this will... overlay... and is not recommended.

You can also import incomplete objects: if the ID is present it will overlay, if not, it will create new objects.
This is also not really recommended except possibly for piping output of other tools, possibly.
