# dc-mongo-remote-test
Try out running mongo as a remote service

This project was to test out how to connect an app to a remote service of mongo.  also brought in some viper for setting mannagment.

When set up this way viper will pull in environmental variables overriding the settings.toml file.  NEAT!

Based on this tutorial about maintaining a clean development enviroment with docker: https://outcrawl.com/clean-development-environment-docker , https://github.com/tinrab/go-mongo-docker-example

TODO:
play with the nginx because the static site did not deploy correctly