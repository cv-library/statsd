# statsd

A [StatsD](https://github.com/etsy/statsd) client for Go.

[![GoDoc](https://godoc.org/github.com/cv-library/statsd?status.png)](https://godoc.org/github.com/cv-library/statsd)

## Example

``` go
import (
    "github.com/cv-library/statsd"
)

func main() {
    statsd.Address = "localhost:8125"

    timer := statsd.Timer()

    // Do stuff

    // Note: Lack of support for sampling yet.
    timer.send("metric.name", "metric.name2")
}
```

## License

Released under the [MIT license](http://www.opensource.org/licenses/mit-license.php). See `LICENSE.md` file for details.
