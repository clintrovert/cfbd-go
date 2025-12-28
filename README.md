# cfbd-go
Go client for collegefootballdata.com APIs

```golang
ctx := context.Background()

client, err := cfbd.NewClient("YOUR_API_KEY")
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

for _, conf := range conferences {
  fmt.Println(conference);
}
```