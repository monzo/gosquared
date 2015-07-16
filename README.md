# A golang client for gosquared

```go
package main

import (
        "github.com/mondough/gosquared"
)

type ExampleEventData struct {
        NumTransactions int
        Type string
        Tags []string
}

func main() {
        client := gosquared.NewClient("apiKey", "siteToken")
        data := &ExampleEventData{
                NumTransactions: 102,
                Type: "foo",
                Tags: []string{ "bar", "foobar"},
        }
        client.Event("event name", data, "personId (can be blank)")
}

```
