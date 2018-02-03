# go-warp10 package

This package contains tools which aims to help you build applications over Warp10 platform.

## GTS Helper

- Parse sensision format into GTS structs
- Get Sensision (Warp10 input format) value of GTS struct

Example:

```go
gts, err := ParseGTSArrayFromString("1234// my.metric{} 10\n1234// my.metric{} 10")

```

## Queries

Build basic query

```go
n := time.Now()
d := time.Hour

c := w.NewClient("https://warp10.gra1.metrics.ovh.net")
c.ReadToken = "READ_TOKEN"

q := c.NewQuery()
q = q.Fetch("", "os.cpu", map[string]string{}, n, d)
q = q.Bucketize(w.BSum, n, 5*time.Second, 0)
q = q.Reduce(w.RSum, []string{})

result, err := q.Exec()

```