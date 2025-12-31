package main

import (
   "context"
   "fmt"
   "os"

   "github.com/clintrovert/cfbd-go/cfbd"
)

func main() {
   ctx := context.Background()
   client, _ := cfbd.New(os.Getenv("CFBD_API_KEY"))

   // set to true if you are a patreon sub, some endpoints require it.
   isPatreonSubscriber := true

   // GAMES
   printGames(ctx, client)
   printGameTeams(ctx, client)
   printGamePlayers(ctx, client)
   printMedia(ctx, client)
   // if isPatreonSubscriber {
   //    printWeather(ctx, client)
   // }
   printRecords(ctx, client)
   printCalendar(ctx, client)
   printScoreboard(ctx, client)
   printGameAdvancedBoxScore(ctx, client)

   // DRIVES
   printDrives(ctx, client)

   // PLAYS
   printPlays(ctx, client)
   printPlayTypes(ctx, client)
   printPlayStats(ctx, client)
   printPlayStatTypes(ctx, client)
   // printLivePlays(ctx, client) // must use a live game ID

   // TEAMS
   printTeams(ctx, client)
   printFbsTeams(ctx, client)
   printTeamMatchups(ctx, client)
   printTeamATS(ctx, client)
   printRoster(ctx, client)
   printTeamTalentComposite(ctx, client)

   // CONFERENCES
   printConferences(ctx, client)

   // VENUES
   printVenues(ctx, client)

   // COACHES
   printCoaches(ctx, client)

   // PLAYERS
   printPlayerSearch(ctx, client)
   printPlayerUsage(ctx, client)
   printReturningPlayers(ctx, client)
   printTransferPortalPlayers(ctx, client)

   // RANKINGS
   printRankings(ctx, client)

   // BETTING
   printBettingLines(ctx, client)

   // RECRUITING
   printPlayerRecruitingRankings(ctx, client)
   printTeamRecruitingRankings(ctx, client)
   printTeamPositionGroupRecruitingRankings(ctx, client)

   // RATINGS
   printSpPlusRatings(ctx, client)
   printConferenceSpPlusRatings(ctx, client)
   printSRSRatings(ctx, client)
   printEloRatings(ctx, client)
   printFPIRatings(ctx, client)

   // METRICS
   printPPAByDownAndDistance(ctx, client)
   printPPAByTeam(ctx, client) // FAILED

   // STATS
   printPlayerSeasonStats(ctx, client)
   printTeamSeasonStats(ctx, client)
   printTeamStatCategories(ctx, client)
   printAdvancedSeasonStats(ctx, client)
   printAdvancedGameStats(ctx, client)
   printHavocGameStats(ctx, client)

   // DRAFT
   printDraftTeams(ctx, client)
   printDraftPositions(ctx, client)
   printDraftPicks(ctx, client)

   // ADJUSTED METRICS
   if isPatreonSubscriber {
      printTeamSeasonWEPA(ctx, client)
      printPlayerPassingWEPA(ctx, client)
      printPlayerRushingWEPA(ctx, client)
   }

   // INFO
   printInfo(ctx, client)
}

func printGames(ctx context.Context, client *cfbd.Client) {
   games, err := client.GetGames(ctx, cfbd.GetGamesRequest{Year: 2025})
   if err != nil {
      fmt.Printf("error occurred retrieving games: %s", err.Error())
   }

   fmt.Println("================= GAMES =================")
   for _, game := range games {
      fmt.Println(game.String())
   }
}

func printGameTeams(ctx context.Context, client *cfbd.Client) {
   teams, err := client.GetGameTeams(ctx, cfbd.GetGameTeamsRequest{
      Year: 2025, Conference: "sec"},
   )
   if err != nil {
      fmt.Printf("error occurred retrieving games: %s", err.Error())
   }

   fmt.Println("================= GAME TEAMS =================")
   for _, team := range teams {
      fmt.Println(team.String())
   }
}

func printGamePlayers(ctx context.Context, client *cfbd.Client) {
   players, err := client.GetGamePlayers(
      ctx, cfbd.GetGamePlayersRequest{Year: 2025, Conference: "sec"},
   )
   if err != nil {
      fmt.Printf("error occurred retrieving games: %s", err.Error())
   }

   fmt.Println("================= GAME PLAYERS =================")
   for _, player := range players {
      fmt.Println(player.String())
   }
}

func printMedia(ctx context.Context, client *cfbd.Client) {
   media, err := client.GetGameMedia(ctx, cfbd.GetGameMediaRequest{Year: 2025})
   if err != nil {
      fmt.Printf("error occurred retrieving game media: %s", err.Error())
   }

   fmt.Println("================= GAME MEDIA =================")
   for _, m := range media {
      fmt.Println(m.String())
   }
}

func printRecords(ctx context.Context, client *cfbd.Client) {
   records, err := client.GetTeamRecords(
      ctx,
      cfbd.GetRecordsRequest{Team: "Texas"},
   )
   if err != nil {
      fmt.Printf("error occurred retrieving calendar: %s", err.Error())
   }

   fmt.Println("================= RECORDS =================")
   for _, record := range records {
      fmt.Println(record.String())
   }
}

func printCalendar(ctx context.Context, client *cfbd.Client) {
   weeks, err := client.GetCalendar(ctx, 2025)
   if err != nil {
      fmt.Printf("error occurred retrieving calendar: %s", err.Error())
   }

   fmt.Println("================= CALENDAR =================")
   for _, week := range weeks {
      fmt.Println(week.String())
   }
}

func printScoreboard(ctx context.Context, client *cfbd.Client) {
   panic("implement")
}

func printGameAdvancedBoxScore(ctx context.Context, client *cfbd.Client) {
   boxScore, err := client.GetAdvancedBoxScore(ctx, 401767768)
   if err != nil {
      fmt.Printf("error occurred get advanced box score: %s", err.Error())
   }

   fmt.Println("================= ADVANCED BOX SCORE =================")
   fmt.Println(boxScore.String())
}

func printDrives(ctx context.Context, client *cfbd.Client) {
   drives, err := client.GetDrives(ctx, cfbd.GetDrivesRequest{Year: 2025})
   if err != nil {
      fmt.Printf("error occurred retrieving drives: %s", err.Error())
   }

   fmt.Println("================= DRIVES =================")
   for _, drive := range drives {
      fmt.Println(drive.String())
   }
}

func printPlays(ctx context.Context, client *cfbd.Client) {
   plays, err := client.GetPlays(ctx, cfbd.GetPlaysRequest{Year: 2025, Week: 1})
   if err != nil {
      fmt.Printf("error occurred retrieving plays: %s", err.Error())
   }

   fmt.Println("================= PLAYS =================")
   for _, play := range plays {
      fmt.Println(play.String())
   }
}

func printPlayTypes(ctx context.Context, client *cfbd.Client) {
   types, err := client.GetPlayTypes(ctx)
   if err != nil {
      fmt.Printf("error occurred retrieving play types: %s", err.Error())
   }

   fmt.Println("================= PLAY TYPES =================")
   for _, t := range types {
      fmt.Println(t.String())
   }
}

func printPlayStats(ctx context.Context, client *cfbd.Client) {
   stats, err := client.GetPlayStats(ctx, cfbd.GetPlayStatsRequest{
      Year: 2025, Week: 2,
   })
   if err != nil {
      fmt.Printf("error occurred retrieving play stats: %s", err.Error())
   }

   fmt.Println("================= PLAY STATS =================")
   for _, stat := range stats {
      fmt.Println(stat.String())
   }
}

func printPlayStatTypes(ctx context.Context, client *cfbd.Client) {
   types, err := client.GetPlayStatTypes(ctx)
   if err != nil {
      fmt.Printf("error occurred retrieving play stat types: %s", err.Error())
   }

   fmt.Println("================= PLAY STATISTIC TYPES =================")
   for _, t := range types {
      fmt.Println(t.String())
   }
}

func printLivePlays(ctx context.Context, client *cfbd.Client) {
   plays, err := client.GetLivePlays(ctx, 401778326)
   if err != nil {
      fmt.Printf("error occurred retrieving All Teams: %s", err.Error())
   }

   fmt.Println("================= LIVE PLAYS SCOREBOARD =================")
   fmt.Println(plays.String())
}

func printTeams(ctx context.Context, client *cfbd.Client) {
   teams, err := client.GetTeams(ctx, cfbd.GetTeamsRequest{})
   if err != nil {
      fmt.Printf("error occurred retrieving All Teams: %s", err.Error())
   }

   fmt.Println("================= ALL TEAMS =================")
   for _, team := range teams {
      fmt.Println(team.String())
   }
}

func printFbsTeams(ctx context.Context, client *cfbd.Client) {
   teams, err := client.GetFBSTeams(ctx, cfbd.GetFBSTeamsRequest{})
   if err != nil {
      fmt.Printf("error occurred retrieving FBS Teams: %s", err.Error())
   }

   fmt.Println("================= FBS TEAMS =================")
   for _, team := range teams {
      fmt.Println(team.String())
   }
}

func printTeamMatchups(ctx context.Context, client *cfbd.Client) {
   matchup, err := client.GetTeamMatchup(
      ctx, cfbd.GetTeamMatchupRequest{Team1: "Texas", Team2: "Oklahoma"},
   )
   if err != nil {
      fmt.Printf("error occurred retrieving Team Matchups: %s", err.Error())
   }

   fmt.Println("================= TEAMS MATCHUPS =================")
   fmt.Println(matchup.String())
}

func printTeamATS(ctx context.Context, client *cfbd.Client) {
   teams, err := client.GetTeamATS(ctx, cfbd.GetTeamATSRequest{Year: 2025})
   if err != nil {
      fmt.Printf("error occurred retrieving Team ATS: %s", err.Error())
   }

   fmt.Println("================= TEAMS AGAINST THE SPREAD =================")
   for _, team := range teams {
      fmt.Println(team.String())
   }
}

func printRoster(ctx context.Context, client *cfbd.Client) {
   rosters, err := client.GetRoster(
      ctx, cfbd.GetRosterRequest{},
   )
   if err != nil {
      fmt.Printf("error occurred retrieving rosters: %s", err.Error())
   }

   fmt.Println("================= ROSTERS =================")
   for _, roster := range rosters {
      fmt.Println(roster.String())
   }
}

func printTeamTalentComposite(ctx context.Context, client *cfbd.Client) {
   talents, err := client.GetTeamTalentComposite(
      ctx, cfbd.GetTalentCompositeRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf("error occurred retrieving talent composite: %s", err.Error())
   }

   fmt.Println("================= 247 TEAM TALENT COMPOSITE =================")
   for _, talent := range talents {
      fmt.Println(talent.String())
   }
}

func printConferences(ctx context.Context, client *cfbd.Client) {
   conferences, err := client.GetConferences(ctx)
   if err != nil {
      fmt.Printf("error occurred retrieving conferences: %s", err.Error())
   }

   fmt.Println("================= CONFERENCES =================")
   for _, conference := range conferences {
      fmt.Println(conference.String())
   }
}

func printVenues(ctx context.Context, client *cfbd.Client) {
   venues, err := client.GetVenues(ctx)
   if err != nil {
      fmt.Printf("error occurred retrieving venues: %s", err.Error())
   }

   fmt.Println("================= VENUES =================")
   for _, venue := range venues {
      fmt.Println(venue.String())
   }
}

func printCoaches(ctx context.Context, client *cfbd.Client) {
   coaches, err := client.GetCoaches(ctx, cfbd.GetCoachesRequest{})
   if err != nil {
      fmt.Printf("error occurred retrieving coaches: %s", err.Error())
   }

   fmt.Println("================= COACHES =================")
   for _, coach := range coaches {
      fmt.Println(coach.String())
   }
}

func printPlayerSearch(ctx context.Context, client *cfbd.Client) {
   players, err := client.SearchPlayers(
      ctx, cfbd.SearchPlayersRequest{SearchTerm: "Smith"},
   )
   if err != nil {
      fmt.Printf("error occurred searching players: %s", err.Error())
   }

   fmt.Println("================= PLAYER SEARCH =================")
   for _, player := range players {
      fmt.Println(player.String())
   }
}

// THIS FAILS
func printPlayerUsage(ctx context.Context, client *cfbd.Client) {
   players, err := client.GetPlayerUsage(
      ctx, cfbd.GetPlayerUsageRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf("error occurred requesting player usage: %s", err.Error())
   }

   fmt.Println("================= PLAYER USAGE =================")
   for _, player := range players {
      fmt.Println(player.String())
   }
}

func printReturningPlayers(ctx context.Context, client *cfbd.Client) {
   players, err := client.GetReturningProduction(
      ctx, cfbd.GetReturningProductionRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting returning production: %s", err.Error(),
      )
   }

   fmt.Println("================= RETURNING PRODUCTION =================")
   for _, player := range players {
      fmt.Println(player.String())
   }
}

func printTransferPortalPlayers(ctx context.Context, client *cfbd.Client) {
   players, err := client.GetTransferPortalPlayers(
      ctx, cfbd.GetTransferPortalPlayersRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf("error occurred requesting portal players: %s", err.Error())
   }

   fmt.Println("================= TRANSFER PORTAL PLAYERS =================")
   for _, player := range players {
      fmt.Println(player.String())
   }
}

func printRankings(ctx context.Context, client *cfbd.Client) {
   rankings, err := client.GetRankings(ctx, cfbd.GetRankingsRequest{Year: 2025})
   if err != nil {
      fmt.Printf("error occurred requesting rankings: %s", err.Error())
   }

   fmt.Println("================= RANKINGS =================")
   for _, ranking := range rankings {
      fmt.Println(ranking.String())
   }
}

func printBettingLines(ctx context.Context, client *cfbd.Client) {
   lines, err := client.GetBettingLines(
      ctx, cfbd.GetBettingLinesRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf("error occurred requesting betting lines: %s", err.Error())
   }

   fmt.Println("================= BETTING LINES =================")
   for _, line := range lines {
      fmt.Println(line.String())
   }
}

func printPlayerRecruitingRankings(ctx context.Context, client *cfbd.Client) {
   players, err := client.GetPlayerRecruitingRankings(
      ctx, cfbd.GetPlayersRecruitingRankingsRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting player recuriting rankings: %s",
         err.Error(),
      )
   }

   fmt.Println("================= PLAYER RECRUITING RANKINGS =================")
   for _, player := range players {
      fmt.Println(player.String())
   }
}

func printTeamRecruitingRankings(ctx context.Context, client *cfbd.Client) {
   teams, err := client.GetTeamRecruitingRankings(
      ctx, cfbd.GetTeamRecruitingRankingsRequest{},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting team recuriting rankings: %s", err.Error(),
      )
   }

   fmt.Println("================= TEAM RECRUITING RANKINGS =================")
   for _, team := range teams {
      fmt.Println(team.String())
   }
}

func printTeamPositionGroupRecruitingRankings(
   ctx context.Context, client *cfbd.Client,
) {
   groups, err := client.GetTeamPositionGroupRecruitingRankings(
      ctx, cfbd.GetTeamPositionGroupRecruitingRankingsRequest{},
   )
   if err != nil {
      fmt.Printf("error occurred requesting recuriting groups: %s", err.Error())
   }

   fmt.Println("================= TEAM POSITION GROUP " +
      "RECRUITING RANKINGS =================")
   for _, group := range groups {
      fmt.Println(group.String())
   }
}

func printSRSRatings(ctx context.Context, client *cfbd.Client) {
   ratings, err := client.GetSRSRatings(
      ctx, cfbd.GetSRSRatingsRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf("error occurred requesting SRS ratings: %s", err.Error())
   }

   fmt.Println("================= SRS RATINGS =================")
   for _, rating := range ratings {
      fmt.Printf(rating.String())
   }
}

func printFPIRatings(ctx context.Context, client *cfbd.Client) {
   ratings, err := client.GetFPIRatings(
      ctx, cfbd.GetFPIRatingsRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting FPI ratings: %s", err.Error(),
      )
   }

   fmt.Println("================= FPI RATINGS =================")
   for _, rating := range ratings {
      if rating.Fpi != nil {
         fmt.Printf("%d, %s - %f\n", rating.Year, rating.Team, rating.Fpi.Value)
      }
   }
}

func printConferenceSpPlusRatings(ctx context.Context, client *cfbd.Client) {
   conferences, err := client.GetConferenceSPPlusRatings(
      ctx, cfbd.GetConferenceSPPlusRatingsRequest{},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting conference SP+ ratings: %s", err.Error(),
      )
   }

   fmt.Println("================= SP+ CONFERENCE RATINGS =================")
   for _, conference := range conferences {
      fmt.Println(conference.String())
   }
}

func printSpPlusRatings(ctx context.Context, client *cfbd.Client) {
   ratings, err := client.GetTeamSPPlusRatings(
      ctx, cfbd.GetSPPlusRatingsRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf("error occurred requesting SP+ ratings: %s", err.Error())
   }

   fmt.Println("================= SP+ RATINGS =================")
   for _, rating := range ratings {
      fmt.Println(rating.String())
   }
}

func printEloRatings(ctx context.Context, client *cfbd.Client) {
   ratings, err := client.GetEloRatings(ctx, cfbd.GetEloRatingsRequest{})
   if err != nil {
      fmt.Printf("error occurred requesting elo ratings: %s", err.Error())
   }

   fmt.Println("================= ELO RATINGS =================")
   for _, rating := range ratings {
      fmt.Println(rating.String())
   }
}

func printPPAByDownAndDistance(ctx context.Context, client *cfbd.Client) {
   ppa, err := client.GetPredictedPoints(
      ctx, cfbd.GetPredictedPointsRequest{Distance: 10, Down: 2},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting ppa by down and distance: %s",
         err.Error(),
      )
   }

   fmt.Println("================= PPA BY DOWN AND DISTANCE =================")
   for _, p := range ppa {
      fmt.Println(p.String())
   }
}

// FAILED
func printPPAByTeam(ctx context.Context, client *cfbd.Client) {
   ppaByTeam, err := client.GetTeamsPPA(
      ctx, cfbd.GetTeamsPPARequest{Team: "Texas"},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting ppa by team: %s",
         err.Error(),
      )
   }

   ppaByYear, err := client.GetTeamsPPA(
      ctx, cfbd.GetTeamsPPARequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting ppa by year: %s",
         err.Error(),
      )
   }

   fmt.Println("================= PPA BY TEAM =================")
   for _, p := range ppaByTeam {
      fmt.Println(p.String())
   }

   fmt.Println("\n\n================= PPA BY YEAR =================")
   for _, p := range ppaByYear {
      fmt.Println(p.String())
   }
}

func printPlayerSeasonStats(ctx context.Context, client *cfbd.Client) {
   stats, err := client.GetPlayerSeasonStats(
      ctx, cfbd.GetPlayerSeasonStatsRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting player season stats: %s",
         err.Error(),
      )
   }

   fmt.Println("================= PLAYER SEASON STATS =================")
   for _, stat := range stats {
      fmt.Println(stat.String())
   }
}

func printTeamSeasonStats(ctx context.Context, client *cfbd.Client) {
   stats, err := client.GetTeamSeasonStats(
      ctx, cfbd.GetTeamSeasonStatsRequest{Year: 2024},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting team season stats: %s",
         err.Error(),
      )
   }

   fmt.Println("================= TEAM SEASON STATS =================")
   for _, stat := range stats {
      fmt.Println(stat.String())
   }
}

func printTeamStatCategories(ctx context.Context, client *cfbd.Client) {
   stats, err := client.GetStatCategories(ctx)
   if err != nil {
      fmt.Printf(
         "error occurred requesting team stat categories: %s",
         err.Error(),
      )
   }

   fmt.Println("================= TEAM STAT CATEGORIES =================")
   for _, stat := range stats {
      fmt.Println(stat)
   }
}

func printAdvancedSeasonStats(ctx context.Context, client *cfbd.Client) {
   stats, err := client.GetAdvancedSeasonStats(
      ctx, cfbd.GetAdvancedSeasonStatsRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting advanced season stats: %s",
         err.Error(),
      )
   }
   fmt.Println("================= ADVANCED SEASON STATS =================")
   for _, stat := range stats {
      fmt.Println(stat.String())
   }
}

func printAdvancedGameStats(ctx context.Context, client *cfbd.Client) {
   stats, err := client.GetAdvancedGameStats(
      ctx, cfbd.GetAdvancedGameStatsRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting advanced game stats: %s",
         err.Error(),
      )
   }
   fmt.Println("================= ADVANCED GAME STATS =================")
   for _, stat := range stats {
      fmt.Println(stat.String())
   }
}

func printHavocGameStats(ctx context.Context, client *cfbd.Client) {
   stats, err := client.GetHavocGameStats(
      ctx, cfbd.GetHavocGameStatsRequest{Year: 2025},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting havoc game stats: %s",
         err.Error(),
      )
   }
   fmt.Println("================= HAVOC GAME STATS =================")
   for _, stat := range stats {
      fmt.Println(stat.String())
   }
}

func printDraftTeams(ctx context.Context, client *cfbd.Client) {
   teams, err := client.GetDraftTeams(ctx)
   if err != nil {
      fmt.Printf(
         "error occurred requesting draft teams: %s",
         err.Error(),
      )
   }

   fmt.Println("================= DRAFT TEAMS =================")
   for _, t := range teams {
      fmt.Println(t.String())
   }
}

func printDraftPositions(ctx context.Context, client *cfbd.Client) {
   positions, err := client.GetDraftPositions(ctx)
   if err != nil {
      fmt.Printf(
         "error occurred requesting draft positions: %s",
         err.Error(),
      )
   }

   fmt.Println("================= DRAFT POSITIONS =================")
   for _, p := range positions {
      fmt.Println(p.String())
   }
}

func printDraftPicks(ctx context.Context, client *cfbd.Client) {
   picks, err := client.GetDraftPicks(ctx, cfbd.GetDraftPicksRequest{})
   if err != nil {
      fmt.Printf(
         "error occurred requesting draft picks: %s",
         err.Error(),
      )
   }

   fmt.Println("================= DRAFT PICKS =================")
   for _, p := range picks {
      fmt.Println(p.String())
   }
}

func printTeamSeasonWEPA(ctx context.Context, client *cfbd.Client) {
   wepa, err := client.GetTeamSeasonWEPA(ctx, cfbd.GetTeamSeasonWEPARequest{})
   if err != nil {
      fmt.Printf(
         "error occurred requesting team season WEPA: %s",
         err.Error(),
      )
   }

   fmt.Println("================= TEAM SEASON WEPA =================")
   for _, w := range wepa {
      fmt.Println(w.String())
   }
}

func printPlayerPassingWEPA(ctx context.Context, client *cfbd.Client) {
   wepa, err := client.GetPlayerPassingWEPA(
      ctx, cfbd.GetPlayerWEPARequest{},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting player passing WEPA: %s",
         err.Error(),
      )
   }

   fmt.Println("================= PLAYER PASSING WEPA =================")
   for _, w := range wepa {
      fmt.Println(w.String())
   }
}

func printPlayerRushingWEPA(ctx context.Context, client *cfbd.Client) {
   wepa, err := client.GetPlayerRushingWEPA(
      ctx, cfbd.GetPlayerWEPARequest{},
   )
   if err != nil {
      fmt.Printf(
         "error occurred requesting player rushing WEPA: %s",
         err.Error(),
      )
   }

   fmt.Println("================= PLAYER RUSHING WEPA =================")
   for _, w := range wepa {
      fmt.Println(w.String())
   }
}

func printInfo(ctx context.Context, client *cfbd.Client) {
   info, err := client.GetInfo(ctx)
   if err != nil {
      fmt.Printf(
         "error occurred requesting info: %s",
         err.Error(),
      )
   }

   fmt.Println("================= INFO =================")
   fmt.Println(info.String())
}
