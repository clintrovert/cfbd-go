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
   confs, err := client.GetConferences(context.Background())
   if err != nil {
      fmt.Println(err.Error())
   }

   for _, conf := range confs {
      fmt.Println(conf.Name)
   }
}
