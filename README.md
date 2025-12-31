# cfbd-go

A minimal, type-safe Golang client for the [College Football Data API](https://collegefootballdata.com/).

## Features

- **Type-safe**: All API responses are strongly typed using Protocol Buffers
- **Comprehensive**: Supports all endpoints covering games, teams, players, stats, ratings, and more
- **Minimal dependencies**: Lightweight with no unnecessary dependencies
- **Context-aware**: All methods support `context.Context` for cancellation and timeouts
- **Future-proof**: Unknown JSON fields are discarded by default to tolerate future API releases

## Installation

```bash
go get github.com/clintrovert/cfbd-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/clintrovert/cfbd-go/cfbd"
)

func main() {
    // Initialize the client with your API key
    client, err := cfbd.New(os.Getenv("CFBD_API_KEY"))
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    // Get games for a specific year
    games, err := client.GetGames(ctx, cfbd.GetGamesRequest{
        Year: 2024,
        Week: 1,
    })
    if err != nil {
        panic(err)
    }

    for _, game := range games {
        fmt.Printf("%s vs %s\n", game.AwayTeam, game.HomeTeam)
    }
}
```

## Authentication

The CFBD API requires an API key for authentication. You can obtain an API key by:

1. Signing up at [collegefootballdata.com](https://collegefootballdata.com/)
2. Getting a free API key or subscribing to Patreon for additional features

The API key is passed as a Bearer token in the Authorization header. Some endpoints require a Patreon subscription.

## API Methods

### Games

#### GetGames
Retrieve game information with filtering options.

```go
games, err := client.GetGames(ctx, cfbd.GetGamesRequest{
    Year: 2024,
    Week: 1,
    Team: "Texas",
    Conference: "Big 12",
    SeasonType: "regular",
})
```

#### GetGameTeams
Get team box score statistics for games.

```go
teamStats, err := client.GetGameTeams(ctx, cfbd.GetGameTeamsRequest{
    Year: 2024,
    Week: 1,
})
```

#### GetGamePlayers
Retrieve player box score statistics for games.

```go
playerStats, err := client.GetGamePlayers(ctx, cfbd.GetGamePlayersRequest{
    Year: 2024,
    Week: 1,
    Team: "Alabama",
})
```

#### GetGameMedia
Get media information for games.

```go
media, err := client.GetGameMedia(ctx, cfbd.GetGameMediaRequest{
    Year: 2024,
    Week: 1,
})
```

#### GetGameWeather
Get weather information for games (Patreon required).

```go
weather, err := client.GetGameWeather(ctx, cfbd.GetGameWeatherRequest{
    Year: 2024,
    Week: 1,
})
```

#### GetAdvancedBoxScore
Get advanced box score statistics for a specific game.

```go
boxScore, err := client.GetAdvancedBoxScore(ctx, gameID)
```

#### GetCalendar
Retrieve calendar weeks for a year.

```go
weeks, err := client.GetCalendar(ctx, 2024)
```

#### GetScoreboard
Get live scoreboard data.

```go
scoreboard, err := client.GetScoreboard(ctx, cfbd.GetScoreboardRequest{
    Conference: "SEC",
})
```

### Teams

#### GetTeams
Retrieve team information.

```go
teams, err := client.GetTeams(ctx, cfbd.GetTeamsRequest{
    Conference: "SEC",
    Year: 2024,
})
```

#### GetFBSTeams
Get FBS (Football Bowl Subdivision) teams.

```go
fbsTeams, err := client.GetFBSTeams(ctx, cfbd.GetFBSTeamsRequest{
    Year: 2024,
})
```

#### GetTeamRecords
Get team records.

```go
records, err := client.GetTeamRecords(ctx, cfbd.GetTeamRecordsRequest{
    Team: "Texas",
    Year: 2024,
})
```

#### GetTeamMatchup
Get historical matchup data between two teams.

```go
matchup, err := client.GetTeamMatchup(ctx, cfbd.GetTeamMatchupRequest{
    Team1: "Texas",
    Team2: "Oklahoma",
    MinYear: 2020,
    MaxYear: 2024,
})
```

#### GetTeamATS
Get team against-the-spread records.

```go
ats, err := client.GetTeamATS(ctx, cfbd.GetTeamATSRequest{
    Year: 2024,
    Conference: "SEC",
})
```

#### GetRoster
Get team roster information.

```go
roster, err := client.GetRoster(ctx, cfbd.GetRosterRequest{
    Team: "Alabama",
    Year: 2024,
})
```

#### GetTeamTalentComposite
Get 247 team talent composite ratings.

```go
talent, err := client.GetTeamTalentComposite(ctx, cfbd.GetTalentCompositeRequest{
    Year: 2024,
})
```

### Players

#### SearchPlayers
Search for players by name or other criteria.

```go
players, err := client.SearchPlayers(ctx, cfbd.SearchPlayersRequest{
    SearchTerm: "Smith",
    Year: 2024,
    Team: "Alabama",
})
```

#### GetPlayerUsage
Get player usage statistics.

```go
usage, err := client.GetPlayerUsage(ctx, cfbd.GetPlayerUsageRequest{
    Year: 2024,
    Team: "Texas",
    Position: "QB",
})
```

#### GetReturningProduction
Get returning production statistics for players.

```go
production, err := client.GetReturningProduction(ctx, cfbd.GetReturningProductionRequest{
    Year: 2024,
    Team: "Alabama",
})
```

#### GetTransferPortalPlayers
Get transfer portal player information.

```go
transfers, err := client.GetTransferPortalPlayers(ctx, cfbd.GetTransferPortalPlayersRequest{
    Year: 2024,
})
```

### Plays and Drives

#### GetPlays
Get play-by-play data for games.

```go
plays, err := client.GetPlays(ctx, cfbd.GetPlaysRequest{
    Year: 2024,
    Week: 1,
    Team: "Texas",
})
```

#### GetPlayTypes
Get all available play types.

```go
playTypes, err := client.GetPlayTypes(ctx)
```

#### GetPlayStats
Get play statistics.

```go
stats, err := client.GetPlayStats(ctx, cfbd.GetPlayStatsRequest{
    Year: 2024,
    Week: 1,
    GameID: 401767768,
})
```

#### GetPlayStatTypes
Get all available play statistic types.

```go
statTypes, err := client.GetPlayStatTypes(ctx)
```

#### GetLivePlays
Get live play-by-play data for a game (requires live game ID).

```go
liveGame, err := client.GetLivePlays(ctx, gameID)
```

#### GetDrives
Get drive information for games.

```go
drives, err := client.GetDrives(ctx, cfbd.GetDrivesRequest{
    Year: 2024,
    Week: 1,
    Team: "Alabama",
})
```

### Statistics

#### GetPlayerSeasonStats
Get player season statistics.

```go
stats, err := client.GetPlayerSeasonStats(ctx, cfbd.GetPlayerSeasonStatsRequest{
    Year: 2024,
    Team: "Texas",
    Category: "passing",
})
```

#### GetTeamSeasonStats
Get team season statistics.

```go
stats, err := client.GetTeamSeasonStats(ctx, cfbd.GetTeamSeasonStatsRequest{
    Year: 2024,
    Team: "Alabama",
})
```

#### GetAdvancedSeasonStats
Get advanced season statistics.

```go
stats, err := client.GetAdvancedSeasonStats(ctx, cfbd.GetAdvancedSeasonStatsRequest{
    Year: 2024,
    Team: "Texas",
    ExcludeGarbageTime: &excludeGarbage,
})
```

#### GetAdvancedGameStats
Get advanced game statistics.

```go
stats, err := client.GetAdvancedGameStats(ctx, cfbd.GetAdvancedGameStatsRequest{
    Year: 2024,
    Week: 1,
    Team: "Alabama",
})
```

#### GetHavocGameStats
Get havoc game statistics.

```go
stats, err := client.GetHavocGameStats(ctx, cfbd.GetHavocGameStatsRequest{
    Year: 2024,
    Week: 1,
    Team: "Texas",
})
```

#### GetStatCategories
Get all available statistics categories.

```go
categories, err := client.GetStatCategories(ctx)
```

### Ratings

#### GetTeamSPPlusRatings
Get SP+ (S&P+) ratings for teams.

```go
ratings, err := client.GetTeamSPPlusRatings(ctx, cfbd.GetSPPlusRatingsRequest{
    Year: 2024,
    Team: "Alabama",
})
```

#### GetConferenceSPPlusRatings
Get SP+ ratings for conferences.

```go
ratings, err := client.GetConferenceSPPlusRatings(ctx, cfbd.GetConferenceSPPlusRatingsRequest{
    Year: 2024,
    Conference: "SEC",
})
```

#### GetSRSRatings
Get SRS (Simple Rating System) ratings.

```go
ratings, err := client.GetSRSRatings(ctx, cfbd.GetSRSRatingsRequest{
    Year: 2024,
    Team: "Texas",
})
```

#### GetEloRatings
Get Elo ratings for teams.

```go
ratings, err := client.GetEloRatings(ctx, cfbd.GetEloRatingsRequest{
    Year: 2024,
    Week: 1,
    Team: "Alabama",
})
```

#### GetFPIRatings
Get FPI (Football Power Index) ratings.

```go
ratings, err := client.GetFPIRatings(ctx, cfbd.GetFPIRatingsRequest{
    Year: 2024,
    Team: "Texas",
})
```

### Rankings

#### GetRankings
Get college football rankings (polls).

```go
rankings, err := client.GetRankings(ctx, cfbd.GetRankingsRequest{
    Year: 2024,
    Week: 1,
    SeasonType: "regular",
})
```

### Betting

#### GetBettingLines
Get betting lines for games.

```go
lines, err := client.GetBettingLines(ctx, cfbd.GetBettingLinesRequest{
    Year: 2024,
    Week: 1,
    Provider: "fanduel",
})
```

### Recruiting

#### GetPlayerRecruitingRankings
Get player recruiting rankings.

```go
recruits, err := client.GetPlayerRecruitingRankings(ctx, cfbd.GetPlayersRecruitingRankingsRequest{
    Year: 2024,
    Team: "Alabama",
    Position: "QB",
})
```

#### GetTeamRecruitingRankings
Get team recruiting rankings.

```go
rankings, err := client.GetTeamRecruitingRankings(ctx, cfbd.GetTeamRecruitingRankingsRequest{
    Year: 2024,
    Team: "Texas",
})
```

#### GetTeamPositionGroupRecruitingRankings
Get aggregated team recruiting information by position group.

```go
groups, err := client.GetTeamPositionGroupRecruitingRankings(ctx, cfbd.GetTeamPositionGroupRecruitingRankingsRequest{
    Team: "Alabama",
    StartYear: 2020,
    EndYear: 2024,
})
```

### Metrics

#### GetPredictedPoints
Get predicted points values by down and distance.

```go
points, err := client.GetPredictedPoints(ctx, cfbd.GetPredictedPointsRequest{
    Down: 2,
    Distance: 10,
})
```

#### GetTeamsPPA
Get team season PPA (Predicted Points Added) statistics.

```go
ppa, err := client.GetTeamsPPA(ctx, cfbd.GetTeamsPPARequest{
    Year: 2024,
    Team: "Texas",
})
```

#### GetGamesPPA
Get team game PPA statistics.

```go
ppa, err := client.GetGamesPPA(ctx, cfbd.GetPpaGamesRequest{
    Year: 2024,
    Week: 1,
    Team: "Alabama",
})
```

#### GetPlayersPPA
Get player game PPA statistics.

```go
ppa, err := client.GetPlayersPPA(ctx, cfbd.GetPlayerPpaGamesRequest{
    Year: 2024,
    Week: 1,
    Team: "Texas",
    Position: "QB",
})
```

#### GetPlayerSeasonPPA
Get player season PPA statistics.

```go
ppa, err := client.GetPlayerSeasonPPA(ctx, cfbd.GetPlayerSeasonPPARequest{
    Year: 2024,
    Team: "Alabama",
    Position: "RB",
})
```

#### GetWinProbability
Get win probability data for each play in a game.

```go
probabilities, err := client.GetWinProbability(ctx, gameID)
```

#### GetPregameWinProbability
Get pregame win probability data.

```go
probabilities, err := client.GetPregameWinProbability(ctx, cfbd.GetPregameWpRequest{
    Year: 2024,
    Week: 1,
    Team: "Texas",
})
```

#### GetFieldGoalExpectedPoints
Get field goal expected points values.

```go
ep, err := client.GetFieldGoalExpectedPoints(ctx)
```

### Adjusted Metrics (Patreon Required)

#### GetTeamSeasonWEPA
Get team season WEPA (Weighted Expected Points Added) metrics.

```go
wepa, err := client.GetTeamSeasonWEPA(ctx, cfbd.GetTeamSeasonWEPARequest{
    Year: 2024,
    Team: "Alabama",
})
```

#### GetPlayerPassingWEPA
Get player passing WEPA metrics.

```go
wepa, err := client.GetPlayerPassingWEPA(ctx, cfbd.GetPlayerWEPARequest{
    Year: 2024,
    Team: "Texas",
    Position: "QB",
})
```

#### GetPlayerRushingWEPA
Get player rushing WEPA metrics.

```go
wepa, err := client.GetPlayerRushingWEPA(ctx, cfbd.GetPlayerWEPARequest{
    Year: 2024,
    Team: "Alabama",
    Position: "RB",
})
```

#### GetPlayerKickingWEPA
Get kicker PAAR (Points Above Average Replacement) metrics.

```go
paar, err := client.GetPlayerKickingWEPA(ctx, cfbd.GetWepaPlayersKickingRequest{
    Year: 2024,
    Team: "Texas",
})
```

### Draft

#### GetDraftTeams
Get all NFL draft teams.

```go
teams, err := client.GetDraftTeams(ctx)
```

#### GetDraftPositions
Get all NFL draft positions.

```go
positions, err := client.GetDraftPositions(ctx)
```

#### GetDraftPicks
Get NFL draft picks.

```go
picks, err := client.GetDraftPicks(ctx, cfbd.GetDraftPicksRequest{
    Year: 2024,
    Team: "Dallas Cowboys",
    School: "Alabama",
})
```

### Reference Data

#### GetConferences
Get all available conferences.

```go
conferences, err := client.GetConferences(ctx)
```

#### GetVenues
Get all available venues.

```go
venues, err := client.GetVenues(ctx)
```

#### GetCoaches
Get coach information.

```go
coaches, err := client.GetCoaches(ctx, cfbd.GetCoachesRequest{
    Team: "Alabama",
    Year: 2024,
})
```

### User Information

#### GetInfo
Get information about the authenticated user's API key.

```go
info, err := client.GetInfo(ctx)
if info != nil {
    fmt.Printf("User: %s\n", info.Email)
}
```

## Error Handling

The client returns errors for various conditions:

- `ErrMissingAPIKey`: Returned when the API key is empty
- `ErrMissingRequiredParams`: Returned when required parameters are missing
- Network and API errors are wrapped with context

```go
client, err := cfbd.New("")
if err == cfbd.ErrMissingAPIKey {
    // Handle missing API key
}

games, err := client.GetGames(ctx, cfbd.GetGamesRequest{})
if err != nil {
    // Handle API error
}
```

## Context Support

All methods accept a `context.Context` for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

games, err := client.GetGames(ctx, cfbd.GetGamesRequest{Year: 2024})
```

## Patreon Subscriptions

Some endpoints require a Patreon subscription:

- Weather data (`GetGameWeather`)
- Adjusted metrics (WEPA endpoints)
- Some advanced statistics

Check the [CFBD API documentation](https://collegefootballdata.com/) for the latest list of Patreon-only features.

## Examples

See the [examples directory](cfbd/internal/examples/main.go) for comprehensive usage examples covering all API endpoints.

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
