# auth-aws

## Installation instructions

If you want to compile from source:

```bash
go get -u github.com/tmiller/auth-aws
```

## Usage instructions

### Running
Run the program by executing `auth-aws`

### Configuration
To pass inputs into the program there are two ways that are loaded in the
following order:

1. config file
2. environment variables


The config file is located at ~/.config/auth-aws/config.ini and uses the
following format:

```ini
[adfs]
user = foo
pass = bar
host = federated.host.name
```

Here are the environment variables available:

* ADFS_USER
* ADFS_PASS
* ADFS_HOST

Finally if any variables are missing from the config file or the environment
variables, then the program will ask you to supply them when necessary.
