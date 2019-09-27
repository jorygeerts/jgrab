Grab some data from a JSON string or stream.

Can be useful to parse an output stream of an application into something that is more readable.

## Installation

```terminal
go install github.com/jorygeerts/jgrab/cli
```

## Usage

If your application outputs a stream of JSON blobs and all you care about are the `severity` and `message` nodes,
you can filter them out like this: 

```terminal
your-stream-generating-application | jgrab --node severity --node message
```

Alternatively, if you do not pass any `--node`, all nodes will be shown (but in a more readable format).
