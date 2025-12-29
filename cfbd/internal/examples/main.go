package main

import (
   "context"
   "fmt"
   "os"

   "github.com/clintrovert/cfbd-go/cfbd"
)

func main() {
   key := os.Getenv("CFBD_API_KEY")
   fmt.Println("Key: " + key)

   client, _ := cfbd.NewClient(key)
   gameId := int32(56)
   games, err := client.GetGames(
      context.Background(),
      cfbd.GetGamesRequest{GameID: &gameId},
   )
   if err != nil {
      fmt.Println(err.Error())
   }

   for _, game := range games {
      fmt.Println(game.Id)
      fmt.Println(game.SeasonType)
      fmt.Println(game.Week)
      fmt.Println(game.HomeTeam)
      fmt.Println(game.AwayTeam)
   }
}
