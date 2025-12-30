package main

import (
   "context"
   "fmt"
   "os"

   "github.com/clintrovert/cfbd-go/cfbd"
)

func main() {
   client, _ := cfbd.NewClient(os.Getenv("CFBD_API_KEY"))

   // gameId := int32(2025)
   games, err := client.GetTeams(context.Background(), cfbd.TeamsRequest{})
   if err != nil {
      fmt.Println(err.Error())
   }

   for _, game := range games {
      if game.School == "Texas" {
         fmt.Println(game.String())
      }

      if game.Classification != nil {
         fmt.Println(game.Classification.GetValue())
      }
   }
}
