# EpicWine

The Epic Wine API.

From the Instructions, I have written a Go API with minimal dependancies.

Choices made were based around simplicity of running this demo. Which is why
there is no database that is loaded from the contents of the CSV.  Depending on
Production use cases, I would modify to use a different library than csvfile, to one
that accesses a database for greater scale and redundancy.  Depending on how metrics
are pulled and processed, could get more detailed metrics off a Prometheus exporter
with its own endpoint on `/metrics` for instance.  I've also been interested in
exploring OpenTracing libraries/collectors.

## Building

This is a Go application that reads and writes to a csv file.  Running it
is done with two commands executed with `make`.  Requirements are `make`,
`curl` and of coarse the `go` binary for building the application.

The first command, will download the CSV file, to your local filesystem.

```bash
make getwine
```

Second command, can be one of two options.

```bash
make run
```

Or

```bash
make build
```

First one, will just start the binary, and it will start listening on Port `:5000`
Second command builds a binary named `epicwine`. Which will take some flags for
configuring its runtime.

There is also an alternative to make deploying/pushing to other Systems a little
easier.  And that is a Docker container.

```bash
make docker
```

And it can be tested locally again like:

```bash
docker run -p 5000:5000 epicwine:latest
```

## Runtime

Go Process is built with a couple configurable options by default. If building
the binary, you could pass the option `-h` to see a few more of the run time options
that are adjustable.  Such as `port` and the file name and location of the `epicwine`
binary.

I was testing using a http cli tool called `httpie`.  The api will take various
arguments and parameters. Querying the API would be done with the following
commands:

```bash
http localhost:5000/status

http localhost:5000/wine limit==10 offset==1337

http localhost:5000/wine/6549

http PUT localhost:5000/wine country=USA description="This is my wine description" designation="I dont know what this means" points="93" price="12.98" province="Mordor" region_1="South" region_2="Eastern" title="foobar"

http PUT localhost:5000/wine country=USA description="This is my wine description" designation="I dont know what this means" points="93" price="12.98" province="Mordor" region_1="South" region_2="Eastern" taster_name="Herbie Hancock" taster_twitter_handle="\@herbiebanana" title="Detour" variety="IPA" winery="Uinta"
```
