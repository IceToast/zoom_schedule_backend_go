# Goth-Fiber: Multi-Provider Authentication for Go [![GoDoc](https://godoc.org/github.com/shareed2k/goth_fiber?status.svg)](https://godoc.org/github.com/markbates/goth)

Is wrapper for [goth library](https://github.com/markbates/goth) to use with [fiber Framework](https://github.com/gofiber/fiber), provides a simple, clean, and idiomatic way to write authentication
packages for Go web applications.

Unlike other similar packages, Goth, lets you write OAuth, OAuth2, or any other
protocol providers, as long as they implement the `Provider` and `Session` interfaces.

## Installation

```text
$ go get github.com/shareed2k/goth_fiber
```

## Supported Providers

- Amazon
- Apple
- Auth0
- Azure AD
- Battle.net
- Bitbucket
- Box
- Cloud Foundry
- Dailymotion
- Deezer
- Digital Ocean
- Discord
- Dropbox
- Eve Online
- Facebook
- Fitbit
- Gitea
- GitHub
- Gitlab
- Google
- Google+ (deprecated)
- Heroku
- InfluxCloud
- Instagram
- Intercom
- Kakao
- Lastfm
- Linkedin
- LINE
- Mailru
- Meetup
- MicrosoftOnline
- Naver
- Nextcloud
- OneDrive
- OpenID Connect (auto discovery)
- Paypal
- SalesForce
- Shopify
- Slack
- Soundcloud
- Spotify
- Steam
- Strava
- Stripe
- Tumblr
- Twitch
- Twitter
- Typetalk
- Uber
- VK
- Wepay
- Xero
- Yahoo
- Yammer
- Yandex

## Examples

See the [examples](examples) folder for a working application that lets users authenticate
through Twitter, Facebook, Google Plus etc.

To run the example either clone the source from GitHub

```text
$ git clone git@github.com/shareed2k/goth_fiber.git
```

```text
$ go get github.com/shareed2k/goth_fiber
```

```text
$ cd goth_fiber/examples
$ go get -v
$ go build
$ ./examples
```

Now open up your browser and go to [http://localhost:8088](http://localhost:8088) to see the example.

To actually use the different providers, please make sure you set environment variables. Example given in the examples/main.go file

## Security Notes

By default, goth uses a `Session` from the `gofiber/session` package to store session data.

As configured, this default store (`session.Session`) will generate cookies with `session.Config`:

```go
    session.Config{
        Key:    "dinosaurus",       // default: "sessionid"
        Lookup: "header",           // default: "cookie"
        Domain: "google.com",       // default: ""
        Expires: 30 * time.Minutes, // default: 2 * time.Hour
        Secure:  true,              // default: false
    }
```

To tailor these fields for your application, you can override the `goth_fiber.Session` variable at startup.

The following snippet shows one way to do this:

```go
    // optional config
    config := session.Config{
        Key:    "dinosaurus",       // default: "sessionid"
        Lookup: "header",           // default: "cookie"
        Domain: "google.com",       // default: ""
        Expires: 30 * time.Minutes, // default: 2 * time.Hour
        Secure:  true,              // default: false
     }

    // create session handler
    sessions := session.New(config)

    goth_fiber.Session = sessions
```

## Issues

Issues always stand a significantly better chance of getting fixed if they are accompanied by a
pull request.

## Contributing

Would I love to see more providers? Certainly! Would you love to contribute one? Hopefully, yes!

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Write Tests!
4. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request
