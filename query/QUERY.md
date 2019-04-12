# Queries

Build basic query

```go
n := time.Now()
d := time.Hour

c := w.NewClient("https://warp10.gra1.metrics.ovh.net")
c.ReadToken = "READ_TOKEN"

q := c.NewQuery()
q = q.Fetch("", "os.cpu", w.Labels{}, n, d)
q = q.Bucketize(w.BucketizerSum, n, 5*time.Second, 0)
q = q.Map(w.MapperRate, "", 1, 0, 0)
q = q.Reduce(w.RSum, []string{})

result, err := q.Exec()

```
