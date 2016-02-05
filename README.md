Khronos
=======

Khronos is a modern replacement of cron for microservice architecture.


## Run on dev mode

To run khronos on dev mode you will need: `make`, `docker-compose` and `Docker`.

### Basic start (Single run mode)

The simplest way to run Khronos on your dev machine is doing:

    $ make dev

This command will build all the environment, compile and run the application

### Development mode

The most confortable way to run the dev environment is doing:

    $ make up

This command will build the environment and run the shell where you can exectute
any command

### Settings on dev

To run with specific settings you can use the `KHRONOS_CONFIG_FILE` env var,
for example:

    $ KHRONOS_CONFIG_FILE="`pwd`/environment/dev/settings.json" go run ./main.go

### Others

There are other comands like `make test` to run the tests, `make app_build` to
build the app binary and many more. Check the [Makefile](Makefile) for all the
commands

## Changelog

Check [Changelog](CHANGELOG.md)


## Authors

Check [AUTHORS](AUTHORS)

## License

Check [LICENSE](LICENSE)
