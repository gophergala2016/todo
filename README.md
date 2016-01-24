
# todo

Todo is a todo application which automatically syncs between the cloud and your mobile device. You don't even have to be online to create a new item or edit existing ones. Simply sync when you have internet connection.

Todo consists of two parts. The first one is the web app written in Go. The second one is the native iOS app written in Swift. They share some common code via `gomobile`. That's the whole point of this project. I wanted to learn about `gomobile`.

![screencast](https://raw.githubusercontent.com/gophergala2016/todo/master/screencast.gif)

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

## Conclusion

`gomobile` is awesome. You never have to worry about your data models getting out of sync. Simply add a new field in your Go `struct` and recompile the framework for your mobile devices.

However so far `gomobile` only supports basic type like `string` and `int`. For my app I needed a timestamp. Using `time.Time` wasn't possible and therefore I had to fall back to unix timestamps. All struct fields are getters/setters in Objective-C. In my app I have to convert an object to/from a json like format. As I don't have access to the fields directly and only have getter methods, generating json programatically is more work. With proper fields and reflection this process would be much easier.

`gomobile` is still very young and I hope the Go team will continue their awesome work.

## License

MIT
