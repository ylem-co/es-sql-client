ElasticSearch SQL Client
=======

Example usage

```go
logger := log.StandardLogger()
lvl, _ := log.ParseLevel("trace")
logger.SetLevel(lvl)

client := CreateWithBaseUrl(
    context.Background(),
    "http://localhost:9200",
    nil,
)

_, err := client.Version(nil)
if err != nil {
    panic(err)
}

client.SetLogger(logger)

result, err := client.SqlQuery(`SELECT * FROM kibana_sample_data_ecommerce WHERE "category" = 'Men''s Shoes' LIMIT 2`)
if err != nil {
    panic(err)
}

b, err := json.Marshal(result.Rows)

fmt.Printf("result: %s\n", string(b))
```
