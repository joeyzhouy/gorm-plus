
# gorm-plus
generate GORM entity from database

## Usage
```BASH
go get github.com/joeyzhouy/gorm-plus
gorm-plus generate --host 127.0.0.1 --u root --p root --d test --target /path/to/target
```

## Example

```bash
gorm-plus generate --host 127.0.0.1 --u root --p root --d test --target /path/to/target
```

```shell
usage: gorm-plus generate

create golang entity.

Flags:
  --help              Show context-sensitive help (also try --help-long and --help-man).
  --u="root"          user for connect to database, default: root
  --p="root"          password for connect to database, default: root
  --port=3306         specify port, default: 3306
  --d="dbName"        specify database name
  --json              add json annotation
  --host="127.0.0.1"  specify host
  --target=TARGET     specify target dir
  --tables=TABLES     specify tables to generate, separator ","default all tables in database
  --protocol="mysql"  specify protocol, ex: mysql

```



## Support

- MySQL
