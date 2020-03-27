gitea-webhook
=============

Library for implementing a webhook receiver compatible with gitea

Usage
-----

```
func main () {
    // The secret key is the same as is set up in the gitea webhook configuration
    secretKey := "123456"

    // onSuccess will be called if a request to /webhook has been successfully validated
    onSuccess := func(typ EventType, ev Event, w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Received event %s\n", typ)
        fmt.Printf("Event data: %+v\n", ev)
        // do something with the information
    }


    // Expect requests to be made to "/webhook"
    http.HandleFunc("/webhook", Handler(secretKey, onSuccess))

    http.ListenAndServe(":8080", nil)
}
```
