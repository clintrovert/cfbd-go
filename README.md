# cfbd-go
Minimal Golang client for the collegefootballdata.com APIs.

```golang
api, err := cfbd.New("CFBD_API_KEY")
if err != nil {
  ...
}
```

#### [GetGames]() (`/games`)
```golang
req := cfbd.GetGamesRequest{
   Year: 2025,
}

games, err := api.GetGames(ctx, req)
if err != nil {
   panic("handle this more gracefully")
}

fmt.Println("=========== Games ===========\n")
for i, game := range games {
   fmt.Printf("Game %s -----------\n", i)
   
}
```


#### [GetGameTeams]() (`/games/teams`)
```golang
req := cfbd.GetGameTeamsRequest{
   Year: 2025,
}

games, err := api.GetGameTeams(ctx, req)
if err != nil {
   panic("handle this")
}

fmt.Println("=========== Games ===========\n")
for _, game := range games {
   fmt.Printf("--- Game %s ---", game.Id)
   
}