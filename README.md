# webtransport-example

#### Update dependencies:

    `go mod tidy`

#### Run server:

    `make server`

#### Run client: 

    `make client`

#### The repository already generated empty certificates for the server. But we're skipping the security check for ease of demonstration. If you wish, you can generate your own certificates with the command:
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout mykey.key -out mycert.crt

Regarding desired requirements:

Make the communication channel secure or suggest what security measures you would
implement given more time.

> Answer: Enabling encryption and generating and verifying certificates allows us to improve security.

● Provide a plan/design for an auto-recovery mechanism for both sides (in case of a
temporary connection failure). Feel free to implement that if you have enough time.

> Answer: Adding graceful shutdowns and a stream detection retry system will allow us to let both the server and the client know when the connection is lost and wait for either side to restart.

● Provide integration tests

> Answer: This example shows only a demonstration of the client and server written in the Golang language. Running them by itself is analogous to an end-to-end test. There are too few functions for full-fledged tabular unit tests and demonstration of TDD.

● Can you think of a way for the client to auto-discover the server without the need to point
it to the exact server endpoint?

> Answer: I am sure that a more detailed investigation of the web transport mechanism will make it possible to make auto detection or will allow it to be easily implemented.