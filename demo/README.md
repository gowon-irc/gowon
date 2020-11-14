# Bot demo

These files can be used to setup a test instance for demo or testing purposes.

To run an ircd and bot instance in docker cd to this directory and run:

    docker-compose up -d

Connect your irc client to localhost:6667 (no tls) and join the #gowon channel.

An example [tiny](https://github.com/osa1/tiny) config has been included. To use it run:

    tiny --config tiny.yml
