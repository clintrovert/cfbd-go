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
[`GET /games`](https://apinext.collegefootballdata.com/#/games/GetGames)

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
[`GET /games/teams`](https://apinext.collegefootballdata.com/#/games/GetGameTeams)

Get team box score statistics for games.

```go
teamStats, err := client.GetGameTeams(ctx, cfbd.GetGameTeamsRequest{
    Year: 2024,
    Week: 1,
})
```

#### GetGamePlayers
[`GET /games/players`](https://apinext.collegefootballdata.com/#/games/GetGamePlayers)

Retrieve player box score statistics for games.

```go
playerStats, err := client.GetGamePlayers(ctx, cfbd.GetGamePlayersRequest{
    Year: 2024,
    Week: 1,
    Team: "Alabama",
})
```

#### GetGameMedia
[`GET /games/media`](https://apinext.collegefootballdata.com/#/games/GetGameMedia)

Get media information for games.

```go
media, err := client.GetGameMedia(ctx, cfbd.GetGameMediaRequest{
    Year: 2024,
    Week: 1,
})
```

#### GetGameWeather
[`GET /games/weather`](https://apinext.collegefootballdata.com/#/games/GetGameWeather)

Get weather information for games (Patreon required).

```go
weather, err := client.GetGameWeather(ctx, cfbd.GetGameWeatherRequest{
    Year: 2024,
    Week: 1,
})
```

#### GetAdvancedBoxScore
[`GET /game/box/advanced`](https://apinext.collegefootballdata.com/#/game/GetAdvancedBoxScore)

Get advanced box score statistics for a specific game.

```go
boxScore, err := client.GetAdvancedBoxScore(ctx, gameID)
```

#### GetCalendar
[`GET /calendar`](https://apinext.collegefootballdata.com/#/calendar/GetCalendar)

Retrieve calendar weeks for a year.

```go
weeks, err := client.GetCalendar(ctx, 2024)
```

#### GetScoreboard
[`GET /scoreboard`](https://apinext.collegefootballdata.com/#/scoreboard/GetScoreboard)

Get live scoreboard data.

```go
scoreboard, err := client.GetScoreboard(ctx, cfbd.GetScoreboardRequest{
    Conference: "SEC",
})
```

### Teams

#### GetTeams
[`GET /teams`](https://apinext.collegefootballdata.com/#/teams/GetTeams)

Retrieve team information.

```go
teams, err := client.GetTeams(ctx, cfbd.GetTeamsRequest{
    Conference: "SEC",
    Year: 2024,
})
```

#### GetFBSTeams
[`GET /teams/fbs`](https://apinext.collegefootballdata.com/#/teams/GetFBSTeams)

Get FBS (Football Bowl Subdivision) teams.

```go
fbsTeams, err := client.GetFBSTeams(ctx, cfbd.GetFBSTeamsRequest{
    Year: 2024,
})
```

#### GetTeamRecords
[`GET /records`](https://apinext.collegefootballdata.com/#/records/GetTeamRecords)

Get team records.

```go
records, err := client.GetTeamRecords(ctx, cfbd.GetTeamRecordsRequest{
    Team: "Texas",
    Year: 2024,
})
```

#### GetTeamMatchup
[`GET /teams/matchup`](https://apinext.collegefootballdata.com/#/teams/GetTeamMatchup)

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
[`GET /teams/ats`](https://apinext.collegefootballdata.com/#/teams/GetTeamATS)

Get team against-the-spread records.

```go
ats, err := client.GetTeamATS(ctx, cfbd.GetTeamATSRequest{
    Year: 2024,
    Conference: "SEC",
})
```

#### GetRoster
[`GET /roster`](https://apinext.collegefootballdata.com/#/roster/GetRoster)

Get team roster information.

```go
roster, err := client.GetRoster(ctx, cfbd.GetRosterRequest{
    Team: "Alabama",
    Year: 2024,
})
```

#### GetTeamTalentComposite
[`GET /talent`](https://apinext.collegefootballdata.com/#/talent/GetTeamTalentComposite)

Get 247 team talent composite ratings.

```go
talent, err := client.GetTeamTalentComposite(ctx, cfbd.GetTalentCompositeRequest{
    Year: 2024,
})
```

### Players

#### SearchPlayers
[`GET /player/search`](https://apinext.collegefootballdata.com/#/player/SearchPlayers)

Search for players by name or other criteria.

```go
players, err := client.SearchPlayers(ctx, cfbd.SearchPlayersRequest{
    SearchTerm: "Smith",
    Year: 2024,
    Team: "Alabama",
})
```

#### GetPlayerUsage
[`GET /player/usage`](https://apinext.collegefootballdata.com/#/player/GetPlayerUsage)

Get player usage statistics.

```go
usage, err := client.GetPlayerUsage(ctx, cfbd.GetPlayerUsageRequest{
    Year: 2024,
    Team: "Texas",
    Position: "QB",
})
```

#### GetReturningProduction
[`GET /player/returning`](https://apinext.collegefootballdata.com/#/player/GetReturningProduction)

Get returning production statistics for players.

```go
production, err := client.GetReturningProduction(ctx, cfbd.GetReturningProductionRequest{
    Year: 2024,
    Team: "Alabama",
})
```

#### GetTransferPortalPlayers
[`GET /player/portal`](https://apinext.collegefootballdata.com/#/player/GetTransferPortalPlayers)

Get transfer portal player information.

```go
transfers, err := client.GetTransferPortalPlayers(ctx, cfbd.GetTransferPortalPlayersRequest{
    Year: 2024,
})
```

### Plays and Drives

#### GetPlays
[`GET /plays`](https://apinext.collegefootballdata.com/#/plays/GetPlays)

Get play-by-play data for games.

```go
plays, err := client.GetPlays(ctx, cfbd.GetPlaysRequest{
    Year: 2024,
    Week: 1,
    Team: "Texas",
})
```

#### GetPlayTypes
[`GET /plays/types`](https://apinext.collegefootballdata.com/#/plays/GetPlayTypes)

Get all available play types.

```go
playTypes, err := client.GetPlayTypes(ctx)
```

#### GetPlayStats
[`GET /plays/stats`](https://apinext.collegefootballdata.com/#/plays/GetPlayStats)

Get play statistics.

```go
stats, err := client.GetPlayStats(ctx, cfbd.GetPlayStatsRequest{
    Year: 2024,
    Week: 1,
    GameID: 401767768,
})
```

#### GetPlayStatTypes
[`GET /plays/stats/types`](https://apinext.collegefootballdata.com/#/plays/GetPlayStatTypes)

Get all available play statistic types.

```go
statTypes, err := client.GetPlayStatTypes(ctx)
```

#### GetLivePlays
[`GET /live/plays`](https://apinext.collegefootballdata.com/#/live/GetLivePlays)

Get live play-by-play data for a game (requires live game ID).

```go
liveGame, err := client.GetLivePlays(ctx, gameID)
```

#### GetDrives
[`GET /drives`](https://apinext.collegefootballdata.com/#/drives/GetDrives)

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
[`GET /stats/player/season`](https://apinext.collegefootballdata.com/#/stats/GetPlayerSeasonStats)

Get player season statistics.

```go
stats, err := client.GetPlayerSeasonStats(ctx, cfbd.GetPlayerSeasonStatsRequest{
    Year: 2024,
    Team: "Texas",
    Category: "passing",
})
```

#### GetTeamSeasonStats
[`GET /stats/season`](https://apinext.collegefootballdata.com/#/stats/GetTeamSeasonStats)

Get team season statistics.

```go
stats, err := client.GetTeamSeasonStats(ctx, cfbd.GetTeamSeasonStatsRequest{
    Year: 2024,
    Team: "Alabama",
})
```

#### GetAdvancedSeasonStats
[`GET /stats/season/advanced`](https://apinext.collegefootballdata.com/#/stats/GetAdvancedSeasonStats)

Get advanced season statistics.

```go
stats, err := client.GetAdvancedSeasonStats(ctx, cfbd.GetAdvancedSeasonStatsRequest{
    Year: 2024,
    Team: "Texas",
    ExcludeGarbageTime: &excludeGarbage,
})
```

#### GetAdvancedGameStats
[`GET /stats/game/advanced`](https://apinext.collegefootballdata.com/#/stats/GetAdvancedGameStats)

Get advanced game statistics.

```go
stats, err := client.GetAdvancedGameStats(ctx, cfbd.GetAdvancedGameStatsRequest{
    Year: 2024,
    Week: 1,
    Team: "Alabama",
})
```

#### GetHavocGameStats
[`GET /stats/game/havoc`](https://apinext.collegefootballdata.com/#/stats/GetHavocGameStats)

Get havoc game statistics.

```go
stats, err := client.GetHavocGameStats(ctx, cfbd.GetHavocGameStatsRequest{
    Year: 2024,
    Week: 1,
    Team: "Texas",
})
```

#### GetStatCategories
[`GET /stats/categories`](https://apinext.collegefootballdata.com/#/stats/GetStatCategories)

Get all available statistics categories.

```go
categories, err := client.GetStatCategories(ctx)
```

### Ratings

#### GetTeamSPPlusRatings
[`GET /ratings/sp`](https://apinext.collegefootballdata.com/#/ratings/GetTeamSPPlusRatings)

Get SP+ (S&P+) ratings for teams.

```go
ratings, err := client.GetTeamSPPlusRatings(ctx, cfbd.GetSPPlusRatingsRequest{
    Year: 2024,
    Team: "Alabama",
})
```

#### GetConferenceSPPlusRatings
[`GET /ratings/sp/conferences`](https://apinext.collegefootballdata.com/#/ratings/GetConferenceSPPlusRatings)

Get SP+ ratings for conferences.

```go
ratings, err := client.GetConferenceSPPlusRatings(ctx, cfbd.GetConferenceSPPlusRatingsRequest{
    Year: 2024,
    Conference: "SEC",
})
```

#### GetSRSRatings
[`GET /ratings/srs`](https://apinext.collegefootballdata.com/#/ratings/GetSRSRatings)

Get SRS (Simple Rating System) ratings.

```go
ratings, err := client.GetSRSRatings(ctx, cfbd.GetSRSRatingsRequest{
    Year: 2024,
    Team: "Texas",
})
```

#### GetEloRatings
[`GET /ratings/elo`](https://apinext.collegefootballdata.com/#/ratings/GetEloRatings)

Get Elo ratings for teams.

```go
ratings, err := client.GetEloRatings(ctx, cfbd.GetEloRatingsRequest{
    Year: 2024,
    Week: 1,
    Team: "Alabama",
})
```

#### GetFPIRatings
[`GET /ratings/fpi`](https://apinext.collegefootballdata.com/#/ratings/GetFPIRatings)

Get FPI (Football Power Index) ratings.

```go
ratings, err := client.GetFPIRatings(ctx, cfbd.GetFPIRatingsRequest{
    Year: 2024,
    Team: "Texas",
})
```

### Rankings

#### GetRankings
[`GET /rankings`](https://apinext.collegefootballdata.com/#/rankings/GetRankings)

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
[`GET /lines`](https://apinext.collegefootballdata.com/#/lines/GetBettingLines)

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
[`GET /recruiting/players`](https://apinext.collegefootballdata.com/#/recruiting/GetPlayerRecruitingRankings)

Get player recruiting rankings.

```go
recruits, err := client.GetPlayerRecruitingRankings(ctx, cfbd.GetPlayersRecruitingRankingsRequest{
    Year: 2024,
    Team: "Alabama",
    Position: "QB",
})
```

#### GetTeamRecruitingRankings
[`GET /recruiting/teams`](https://apinext.collegefootballdata.com/#/recruiting/GetTeamRecruitingRankings)

Get team recruiting rankings.

```go
rankings, err := client.GetTeamRecruitingRankings(ctx, cfbd.GetTeamRecruitingRankingsRequest{
    Year: 2024,
    Team: "Texas",
})
```

#### GetTeamPositionGroupRecruitingRankings
[`GET /recruiting/groups`](https://apinext.collegefootballdata.com/#/recruiting/GetTeamPositionGroupRecruitingRankings)

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
[`GET /ppa/predicted`](https://apinext.collegefootballdata.com/#/ppa/GetPredictedPoints)

Get predicted points values by down and distance.

```go
points, err := client.GetPredictedPoints(ctx, cfbd.GetPredictedPointsRequest{
    Down: 2,
    Distance: 10,
})
```

#### GetTeamsPPA
[`GET /ppa/teams`](https://apinext.collegefootballdata.com/#/ppa/GetTeamsPPA)

Get team season PPA (Predicted Points Added) statistics.

```go
ppa, err := client.GetTeamsPPA(ctx, cfbd.GetTeamsPPARequest{
    Year: 2024,
    Team: "Texas",
})
```

#### GetGamesPPA
[`GET /ppa/games`](https://apinext.collegefootballdata.com/#/ppa/GetGamesPPA)

Get team game PPA statistics.

```go
ppa, err := client.GetGamesPPA(ctx, cfbd.GetPpaGamesRequest{
    Year: 2024,
    Week: 1,
    Team: "Alabama",
})
```

#### GetPlayersPPA
[`GET /ppa/players/games`](https://apinext.collegefootballdata.com/#/ppa/GetPlayersPPA)

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
[`GET /ppa/players/season`](https://apinext.collegefootballdata.com/#/ppa/GetPlayerSeasonPPA)

Get player season PPA statistics.

```go
ppa, err := client.GetPlayerSeasonPPA(ctx, cfbd.GetPlayerSeasonPPARequest{
    Year: 2024,
    Team: "Alabama",
    Position: "RB",
})
```

#### GetWinProbability
[`GET /metrics/wp`](https://apinext.collegefootballdata.com/#/metrics/GetWinProbability)

Get win probability data for each play in a game.

```go
probabilities, err := client.GetWinProbability(ctx, gameID)
```

#### GetPregameWinProbability
[`GET /metrics/wp/pregame`](https://apinext.collegefootballdata.com/#/metrics/GetPregameWinProbability)

Get pregame win probability data.

```go
probabilities, err := client.GetPregameWinProbability(ctx, cfbd.GetPregameWpRequest{
    Year: 2024,
    Week: 1,
    Team: "Texas",
})
```

#### GetFieldGoalExpectedPoints
[`GET /metrics/fg/ep`](https://apinext.collegefootballdata.com/#/metrics/GetFieldGoalExpectedPoints)

Get field goal expected points values.

```go
ep, err := client.GetFieldGoalExpectedPoints(ctx)
```

### Adjusted Metrics (Patreon Required)

#### GetTeamSeasonWEPA
[`GET /wepa/team/season`](https://apinext.collegefootballdata.com/#/wepa/GetTeamSeasonWEPA)

Get team season WEPA (Weighted Expected Points Added) metrics.

```go
wepa, err := client.GetTeamSeasonWEPA(ctx, cfbd.GetTeamSeasonWEPARequest{
    Year: 2024,
    Team: "Alabama",
})
```

#### GetPlayerPassingWEPA
[`GET /wepa/players/passing`](https://apinext.collegefootballdata.com/#/wepa/GetPlayerPassingWEPA)

Get player passing WEPA metrics.

```go
wepa, err := client.GetPlayerPassingWEPA(ctx, cfbd.GetPlayerWEPARequest{
    Year: 2024,
    Team: "Texas",
    Position: "QB",
})
```

#### GetPlayerRushingWEPA
[`GET /wepa/players/rushing`](https://apinext.collegefootballdata.com/#/wepa/GetPlayerRushingWEPA)

Get player rushing WEPA metrics.

```go
wepa, err := client.GetPlayerRushingWEPA(ctx, cfbd.GetPlayerWEPARequest{
    Year: 2024,
    Team: "Alabama",
    Position: "RB",
})
```

#### GetPlayerKickingWEPA
[`GET /wepa/players/kicking`](https://apinext.collegefootballdata.com/#/wepa/GetPlayerKickingWEPA)

Get kicker PAAR (Points Above Average Replacement) metrics.

```go
paar, err := client.GetPlayerKickingWEPA(ctx, cfbd.GetWepaPlayersKickingRequest{
    Year: 2024,
    Team: "Texas",
})
```

### Draft

#### GetDraftTeams
[`GET /draft/teams`](https://apinext.collegefootballdata.com/#/draft/GetDraftTeams)

Get all NFL draft teams.

```go
teams, err := client.GetDraftTeams(ctx)
```

#### GetDraftPositions
[`GET /draft/positions`](https://apinext.collegefootballdata.com/#/draft/GetDraftPositions)

Get all NFL draft positions.

```go
positions, err := client.GetDraftPositions(ctx)
```

#### GetDraftPicks
[`GET /draft/picks`](https://apinext.collegefootballdata.com/#/draft/GetDraftPicks)

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
[`GET /conferences`](https://apinext.collegefootballdata.com/#/conferences/GetConferences)

Get all available conferences.

```go
conferences, err := client.GetConferences(ctx)
```

#### GetVenues
[`GET /venues`](https://apinext.collegefootballdata.com/#/venues/GetVenues)

Get all available venues.

```go
venues, err := client.GetVenues(ctx)
```

#### GetCoaches
[`GET /coaches`](https://apinext.collegefootballdata.com/#/coaches/GetCoaches)

Get coach information.

```go
coaches, err := client.GetCoaches(ctx, cfbd.GetCoachesRequest{
    Team: "Alabama",
    Year: 2024,
})
```

### User Information

#### GetInfo
[`GET /info`](https://apinext.collegefootballdata.com/#/info/GetInfo)

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

### Branch Naming and Versioning

This project uses semantic versioning (semver) with automated version calculation based on branch names. When creating a pull request, **prefix your branch name** with one of the following to indicate the type of change:

- **`major/`** - For breaking changes that require a major version bump
  - Example: `major/breaking-api-changes`
  - Result: `v2.0.0` (if previous was `v1.x.x`)

- **`minor/`** - For new features that are backward compatible
  - Example: `minor/add-new-endpoint`
  - Result: `v1.2.0` (if previous was `v1.1.x`)

- **`patch/`** - For bug fixes and other backward-compatible changes
  - Example: `patch/fix-bug`
  - Result: `v1.1.1` (if previous was `v1.1.0`)

**Note**: If no prefix is provided, the default is a patch version bump.

### Workflow

1. Create a branch with the appropriate prefix (`major/`, `minor/`, or `patch/`)
2. Make your changes
3. Ensure tests pass (`go test ./...`)
4. Submit a pull request
5. The CI workflow will:
   - Run all tests
   - Calculate and comment the suggested next release version on your PR
6. When your PR is merged, a new release will be automatically created with the calculated version

For more details, see the [GitHub workflows](.github/workflows/) for test and release automation.
