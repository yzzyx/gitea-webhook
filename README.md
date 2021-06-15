gitea-webhook
=============

Library for implementing a webhook receiver compatible with gitea

Usage
-----

```
func main () {
    // The secret key is the same as is set up in the gitea webhook configuration
    secretKey := "123456"

    api := API{
        URL: "https://try.gitea.io",
        Username: "xxx",
        Password: "yyy",
        Token: "xxxyyy", // if set, this will be used instead of username/password
    }

    // onSuccess will be called if a request to /webhook has been successfully validated
    onSuccess := func(typ EventType, ev Event, w http.ResponseWriter, r *http.Request) {
        fmt.Printf("Received event %s\n", typ)
        fmt.Printf("Event data: %+v\n", ev)

        // update commit status
        if typ == EventTypePush {
            status := CreateStatusOption{
                Description: "Great success!",
                State:       CommitStatusSuccess,
                TargetURL:   "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
            }
            err := api.UpdateCommitState(ev.Repository.FullName, ev.After, CreateStatusOption{}) 
            if err != nil {
                log.Printf("cannot update commit state: %s", err)
            }
        }
    }


    // Expect requests to be made to "/webhook"
    http.HandleFunc("/webhook", Handler(secretKey, onSuccess))

    http.ListenAndServe(":8080", nil)
}
```
