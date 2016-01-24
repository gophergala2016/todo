
# todo

Showcase for syncing data effortlessly between Go backend and native iOS application. Some of the backend code is reused on the mobile application via `gomobile`.

## Usage

Docker is used to run CouchDB for data sync and Redis for sessions. To start them via `docker-compose` use

```bash
$ make up
```

To shut Docker down use

```bash
$ make stop
```

I'm using a Mac and my default `docker-machine` is running on `192.168.99.100`. You'll have to change the IP at some places to make the app work.

1. Get your own ip

If you're on a Mac run

```bash
$ docker-machine ip default
```

On Linux you can simply use `localhost`.

2. Change golang config

Change the IP in `config.go`.

```go
const ip = "192.168.99.100"
```

3. Change osx config

Change the IP in `Config.swift`.

```swift
public struct Config {

    // on osx
    // docker-machine default ip
    public static let URL = "192.168.99.100"
    
    // on linux
    // public static let URL = "localhost"

}
```

Start the web app via

```bash
$ make run
```

To use the iOS app you need a Mac and XCode. Start the project from `./todo` and run the simulator.

## Go's strengths showcase

I always wanted to reuse my existing code on different platforms. I've tried a lot (PhoneGap, Xamarin, React Native, ...) to make it work and it was never really satisfying. Then I saw `gomobile` and knew I had to give it a chance.

The `item` package, which defines my todo item data model, is used on the server and the same package is used on iOS as a native Objective-C framework. To create the framework simply run `make mobile`. You need to have [gomobile](https://github.com/golang/go/wiki/Mobile) installed.

Apart from `gomobile` I'd like to show some other cool Go features

- http error handling
- interface satisfaction
- template funcs
- error chaining
- nested templates
- template functions
- use unix timestamps to handle dates on different platforms and convert them to type `time`

## Caution

Please don't use this example in production. I had to make some tweaks (like allowing sync over http rather than https) to make the demo work within 48 hours.

## License

MIT
