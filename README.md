# cfbd-go
Minimal Golang client for the collegefootballdata.com APIs.

```golang
client, err := cfbd.NewClient("CFBD_API_KEY")
if err != nil {
  ...
}

request := &cfbd.GetTeamsRequest{
  ...
}

conferences, err := client.GetTeams(ctx, request)
if err != nil {
  ...
}

for _, conference := range conferences {
  fmt.Println(conference);
}
```