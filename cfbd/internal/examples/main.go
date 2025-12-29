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
