# QTTF-Backend

## configuration files

Paste the configuration files under the path `./`

Configuration files:

* get `token.json` from google authorization

* get `credential.json` from [google cloud console](https://console.cloud.google.com/welcome/)

* `config.json`
```
{
    "database": {
        "host": "localhost",
        "port": 5432,
        "user": "user_name",
        "password": "passwd",
        "dbname": "dbname",
        "sslmode": "disable"
    },
    "token_path": "./token.json",
    "router": {
        "port": ":3000",
        "read_timeout": 5,
        "write_timeout": 10
    },
    "spreadsheet": {
        "spreadsheet_id": "spreadsheet_id",
        "sheet_name": "sheet_name"
    },
    "credential_path": "./credential.json"
}
```

## Getting started

To make it easy for you to get started with GitLab, here's a list of recommended next steps.

Already a pro? Just edit this README.md and make it your own. Want to make it easy? [Use the template at the bottom](#editing-this-readme)!

