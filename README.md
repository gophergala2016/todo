
# todo


## Usage

Docker is used to run CouchDB for data sync and Redis for sessions. To start them via `docker-compose` use

```bash
$ make up
```

To shut Docker down use

```bash
$ make stop
```bash

I'm using a Mac and my default `docker-machine` is running on `192.168.99.100`. You'll have to change the IP at some places to make the app work.

1. Get your own ip

If you're on a Mac run

```bash
$ docker-machine ip default
```

On Linux you can simply use `localhost`.

2. Change golang config

3. Change ios config

Start the web app via

```bash
$ make run
```

To use the iOS app you need a Mac and XCode. Start the project from `./todo` and run the simulator.

## Go's strengths showcase

I always wanted to reuse my existing code on different platforms. I've tried a lot (PhoneGap, Xamarin, React Native, ...) to make it work and it was never really satisfying. Then I saw `gomobile` and knew I had to give it a chance.

The `item` package, which defines my todo item data model, is used on the server and the same package is used on iOSas a native Objective-C framework. To create the framework simply run `make mobile`. You need to have [gomobile](https://github.com/golang/go/wiki/Mobile) installed.

Apart from `gomobile` I'd like to show some other cool Go features 

- http error handling
- interface satisfaction
- template funcs
- error chaining
- nested templates
- template functions

