# statsd <a href=https://godoc.org/github.com/cv-library/statsd><img align=right src=https://godoc.org/github.com/cv-library/statsd?status.svg></a><a href=https://travis-ci.org/cv-library/statsd><img align=right src=https://api.travis-ci.org/cv-library/statsd.svg></a>

A [StatsD](https://github.com/etsy/statsd) client for Go.

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
