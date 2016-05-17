Networker
=========

Networker is a simple networking client.

## Installing

To install networker package:

    $ go get github.com/intelex-llc/networker

## How to use

```go
import (
    "github.com/intelex-llc/networker"
)
```

#### Simple GET request

```go
body, resp, err := networker.Get("https://someurl.com", nil).Do()

```

#### Simple POST request

You can use String, Map or interface{} for the request body.

```go

type Data struct {
    Foo string
}

data := Data{
    "bar",
}

body, resp, err := networker.Post("https://someurl.com", nil, networker.JSON, data).Do()

```

#### Another way

```go

req := networker.New(GET).
                Url("https://someurl.com").
                Query(map[string]string{"foo": "bar"}).
                Header("foo", "bar")

body, resp, err := req.Do()

```

#### Add custom header

```go

body, resp, err := networker.Get("https://someurl.com", nil).
                            Header("Authorization", "Bearer Zm9zY236Ym9z23Y28=").
                            Do()

```

#### Set cookies

```go

body, resp, err := networker.Get("https://someurl.com", nil).
                            Cookie("foo", "bar").
                            Do()

```

#### Use basic auth

```go

body, resp, err := networker.Get("https://someurl.com", nil).
                            BaseAuth("username", "password").
                            Do()

```

## License

Networker is distributed under the
[MIT License](https://opensource.org/licenses/MIT),
see LICENSE.txt for more information.
