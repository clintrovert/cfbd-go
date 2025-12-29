package cfbd

import (
   "bytes"
   "context"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "reflect"
   "strconv"
   "time"

   "google.golang.org/protobuf/encoding/protojson"
   "google.golang.org/protobuf/proto"
)

const (
   baseURL           = "https://api.collegefootballdata.com"
   defaultTimeoutSec = 30
   userAgent         = "cfbd-go/1.0"
)

// ErrMissingAPIKey is returned if the API key provided was empty.
var ErrMissingAPIKey = errors.New("API key was not provided")

// requestExecutor wraps the http client via an interface for ease in mock
// testing.
type requestExecutor interface {
   execute(
      ctx context.Context,
      path string,
      params url.Values,
   ) ([]byte, error)
}

// Client todo:describe.
//
// Authentication: CFBD uses an API key as a Bearer token in the Authorization
// header.
// Example:
//    Authorization: Bearer <your_api_key>
//
// Notes:
// - All methods accept a cancellable context.Context.
// - For endpoints that return JSON arrays, the client unmarshals each element
// into a message.
// - Unknown JSON fields are discarded by default to tolerate future API
// releases.
//
// Reference for authentication header examples:
// https://blog.collegefootballdata.com/using-api-keys-with-the-cfbd-api/
type Client struct {
   apiKey       string
   unmarshaller protojson.UnmarshalOptions
   executor     requestExecutor
}

// NewClient creates a Client with sane defaults.
func NewClient(apiKey string) (*Client, error) {
   base, err := url.Parse(baseURL)
   if err != nil {
      return nil, fmt.Errorf("could not parse base url; %w", err)
   }

   if apiKey == "" {
      return nil, ErrMissingAPIKey
   }

   return &Client{
      apiKey: apiKey,
      executor: &restClient{
         apiKey:    apiKey,
         baseURL:   base,
         userAgent: userAgent,
         client: &http.Client{
            Timeout: defaultTimeoutSec * time.Second,
         },
      },
      unmarshaller: protojson.UnmarshalOptions{
         DiscardUnknown: true,
         AllowPartial:   true,
      },
   }, nil
}

// GetGames retrieves a list of games based on the provided request
// parameters.
func (c *Client) GetGames(
   ctx context.Context,
   request GetGamesRequest,
) ([]*Game, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/games", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /games; %w", err)
   }

   var games []*Game
   if err = c.unmarshalList(response, &games, &Game{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal games; %w", err)
   }

   return games, nil
}

// GetGameTeams retrieves team box score statistics for games based on
// the provided request parameters.
func (c *Client) GetGameTeams(
   ctx context.Context,
   request GameTeamsRequest,
) ([]*GameTeamStats, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/games/teams", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /games/teams; %w", err)
   }

   var games []*GameTeamStats
   if err = c.unmarshalList(response, &games, &GameTeamStats{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal game team stats; %w", err)
   }

   return games, nil
}

// GetGamePlayers retrieves player box score statistics for games based
// on the provided request parameters.
func (c *Client) GetGamePlayers(
   ctx context.Context,
   request GamePlayersRequest,
) ([]*GamePlayerStats, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/games/players", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /games/players; %w", err)
   }

   var games []*GamePlayerStats
   if err = c.unmarshalList(response, &games, &GamePlayerStats{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal game player stats; %w", err)
   }

   return games, nil
}

// GetGameMedia retrieves media information for games based on the provided
// request parameters.
func (c *Client) GetGameMedia(
   ctx context.Context,
   request GameMediaRequest,
) ([]*GameMedia, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/games/media", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /games/media; %w", err)
   }

   var games []*GameMedia
   if err = c.unmarshalList(response, &games, &GameMedia{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal game media; %w", err)
   }

   return games, nil
}

// GetGameWeather retrieves weather information for games based on the
// provided request parameters.
func (c *Client) GetGameWeather(
   ctx context.Context,
   request GameWeatherRequest,
) ([]*GameWeather, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/games/weather", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /games/weather; %w", err)
   }

   var games []*GameWeather
   if err = c.unmarshalList(response, &games, &GameWeather{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal game weather; %w", err)
   }

   return games, nil
}

// GetAdvancedBoxScore retrieves advanced box score statistics for the
// specified game.
func (c *Client) GetAdvancedBoxScore(
   ctx context.Context,
   gameID int32,
) (*AdvancedBoxScore, error) {
   v := url.Values{}
   v.Set("gameId", strconv.FormatInt(int64(gameID), 10))
   response, err := c.executor.execute(ctx, "/game/box/advanced", v)
   if err != nil {
      return nil, fmt.Errorf("failed to request /game/box/advanced; %w", err)
   }

   var val AdvancedBoxScore
   if err = c.unmarshal(response, &val); err != nil {
      return nil, fmt.Errorf("failed to unmarshal advanced box score; %w", err)
   }

   return &val, nil
}

// GetCalendar retrieves calendar weeks for the specified year.
func (c *Client) GetCalendar(
   ctx context.Context,
   year int32,
) ([]*CalendarWeek, error) {
   if year < 1 {
      return nil, fmt.Errorf("year is required; %w", ErrMissingRequiredParams)
   }

   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(year), 10))
   response, err := c.executor.execute(ctx, "/calendar", v)
   if err != nil {
      return nil, fmt.Errorf("failed to request /calendar; %w", err)
   }

   var weeks []*CalendarWeek
   if err = c.unmarshalList(response, &weeks, &CalendarWeek{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal calendar weeks; %w", err)
   }

   return weeks, nil
}

// GetTeamRecords retrieves team records based on the provided request
// parameters.
func (c *Client) GetTeamRecords(
   ctx context.Context,
   request RecordsRequest,
) ([]*TeamRecords, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/records", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /records; %w", err)
   }

   var records []*TeamRecords
   if err = c.unmarshalList(response, &records, &TeamRecords{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team records; %w", err)
   }

   return records, nil
}

// GetLiveScoreboard retrieves live scoreboard data based on the provided
// request parameters.
func (c *Client) GetLiveScoreboard(
   ctx context.Context,
   request LiveScoreboardRequest,
) ([]*ScoreboardGame, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/scoreboard", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /scoreboard; %w", err)
   }

   var games []*ScoreboardGame
   if err = c.unmarshalList(response, &games, &ScoreboardGame{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal scoreboard games; %w", err)
   }

   return games, nil
}

// GetDrives retrieves drive information for games based on the provided
// request parameters.
func (c *Client) GetDrives(
   ctx context.Context,
   request DrivesRequest,
) ([]*Drive, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/drives", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /drives; %w", err)
   }

   var drives []*Drive
   if err = c.unmarshalList(response, &drives, &Drive{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal drives; %w", err)
   }

   return drives, nil
}

// GetPlays retrieves play-by-play data for games based on the provided
// request parameters.
func (c *Client) GetPlays(
   ctx context.Context,
   request PlaysRequest,
) ([]*Play, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/plays", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /plays; %w", err)
   }

   var plays []*Play
   if err = c.unmarshalList(response, &plays, &Play{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal plays; %w", err)
   }

   return plays, nil
}

// GetPlayTypes retrieves all available play types.
func (c *Client) GetPlayTypes(ctx context.Context) ([]*PlayType, error) {
   response, err := c.executor.execute(ctx, "/plays/types", url.Values{})
   if err != nil {
      return nil, fmt.Errorf("failed to request /plays/types; %w", err)
   }

   var playTypes []*PlayType
   if err = c.unmarshalList(response, &playTypes, &PlayType{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal play types; %w", err)
   }

   return playTypes, nil
}

// GetPlayStats retrieves play statistics based on the provided request
// parameters.
func (c *Client) GetPlayStats(
   ctx context.Context,
   request PlayStatsRequest,
) ([]*PlayStat, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/plays/stats", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /plays/stats; %w", err)
   }

   var stats []*PlayStat
   if err = c.unmarshalList(response, &stats, &PlayStat{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal play stats; %w", err)
   }

   return stats, nil
}

// GetPlayStatTypes retrieves all available play statistic types.
func (c *Client) GetPlayStatTypes(
   ctx context.Context,
) ([]*PlayStatType, error) {
   response, err := c.executor.execute(ctx, "/plays/stats/types", url.Values{})
   if err != nil {
      return nil, fmt.Errorf("failed to request /plays/stats/types; %w", err)
   }

   var statTypes []*PlayStatType
   if err = c.unmarshalList(response, &statTypes, &PlayStatType{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal play stat types; %w", err)
   }

   return statTypes, nil
}

// GetLivePlays retrieves live play-by-play data for the specified game.
func (c *Client) GetLivePlays(
   ctx context.Context,
   gameID int32,
) (*LiveGame, error) {
   params := url.Values{}
   params.Set("gameId", strconv.FormatInt(int64(gameID), 10))

   response, err := c.executor.execute(ctx, "/live/plays", params)
   if err != nil {
      return nil, fmt.Errorf("failed to request /live/plays; %w", err)
   }

   var game LiveGame
   if err = c.unmarshal(response, &game); err != nil {
      return nil, fmt.Errorf("failed to unmarshal live game; %w", err)
   }

   return &game, nil
}

// GetTeams retrieves team information based on the provided request
// parameters.
func (c *Client) GetTeams(
   ctx context.Context,
   request TeamsRequest,
) ([]*Team, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/teams", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /teams; %w", err)
   }

   var teams []*Team
   if err = c.unmarshalList(response, &teams, &Team{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal teams; %w", err)
   }

   return teams, nil
}

// GetTeamsFBS retrieves FBS (Football Bowl Subdivision) team information
// based on the provided request parameters.
func (c *Client) GetTeamsFBS(
   ctx context.Context,
   request TeamsFbsRequest,
) ([]*Team, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/teams/fbs", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /teams/fbs; %w", err)
   }

   var teams []*Team
   if err = c.unmarshalList(response, &teams, &Team{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal teams; %w", err)
   }

   return teams, nil
}

// GetTeamMatchup retrieves historical matchup data between two teams based
// on the provided request parameters.
func (c *Client) GetTeamMatchup(
   ctx context.Context,
   request TeamMatchupRequest,
) (*Matchup, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/teams/matchup", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /teams/matchup; %w", err)
   }

   var matchup Matchup
   if err = c.unmarshal(response, &matchup); err != nil {
      return nil, fmt.Errorf("failed to unmarshal matchup; %w", err)
   }

   return &matchup, nil
}

// GetTeamATS retrieves team against-the-spread (ATS) records based on the
// provided request parameters.
func (c *Client) GetTeamATS(
   ctx context.Context,
   request TeamATSRequest,
) ([]*TeamATS, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/teams/ats", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /teams/ats; %w", err)
   }

   var teams []*TeamATS
   if err = c.unmarshalList(response, &teams, &TeamATS{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team ATS; %w", err)
   }

   return teams, nil
}

// GetRoster retrieves roster information for a team based on the provided
// request parameters.
func (c *Client) GetRoster(
   ctx context.Context,
   request RosterRequest,
) ([]*RosterPlayer, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/roster", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /roster; %w", err)
   }

   var players []*RosterPlayer
   if err = c.unmarshalList(response, &players, &RosterPlayer{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal roster players; %w", err)
   }

   return players, nil
}

// GetTeamTalent retrieves team talent ratings based on the provided request
// parameters.
func (c *Client) GetTeamTalent(
   ctx context.Context,
   request TalentRequest,
) ([]*TeamTalent, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/talent", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /talent; %w", err)
   }

   var talents []*TeamTalent
   if err = c.unmarshalList(response, &talents, &TeamTalent{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team talent; %w", err)
   }

   return talents, nil
}

// GetConferences retrieves all available conferences.
func (c *Client) GetConferences(ctx context.Context) ([]*Conference, error) {
   response, err := c.executor.execute(ctx, "/conferences", url.Values{})
   if err != nil {
      return nil, fmt.Errorf("failed to request /conferences; %w", err)
   }

   var conferences []*Conference
   if err = c.unmarshalList(response, &conferences, &Conference{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal conferences; %w", err)
   }

   return conferences, nil
}

// GetVenues retrieves all available venues.
func (c *Client) GetVenues(ctx context.Context) ([]*Venue, error) {
   response, err := c.executor.execute(ctx, "/venues", url.Values{})
   if err != nil {
      return nil, fmt.Errorf("failed to request /venues; %w", err)
   }

   var venues []*Venue
   if err = c.unmarshalList(response, &venues, &Venue{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal venues; %w", err)
   }

   return venues, nil
}

// GetCoaches retrieves coach information based on the provided request
// parameters.
func (c *Client) GetCoaches(
   ctx context.Context,
   request CoachesRequest,
) ([]*Coach, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/coaches", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /coaches; %w", err)
   }

   var coaches []*Coach
   if err = c.unmarshalList(response, &coaches, &Coach{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal coaches; %w", err)
   }

   return coaches, nil
}

// SearchPlayers searches for players based on the provided request
// parameters.
func (c *Client) SearchPlayers(
   ctx context.Context,
   request PlayerSearchRequest,
) ([]*PlayerSearchResult, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/player/search", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /player/search; %w", err)
   }

   var players []*PlayerSearchResult
   if err = c.unmarshalList(
      response, &players, &PlayerSearchResult{},
   ); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal player search results; %w", err,
      )
   }

   return players, nil
}

// GetPlayerUsage retrieves player usage statistics based on the provided
// request parameters.
func (c *Client) GetPlayerUsage(
   ctx context.Context,
   request PlayerUsageRequest,
) ([]*PlayerUsage, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/player/usage", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /player/usage; %w", err)
   }

   var usage []*PlayerUsage
   if err = c.unmarshalList(response, &usage, &PlayerUsage{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal player usage; %w", err)
   }

   return usage, nil
}

// GetReturningProduction retrieves returning production statistics for
// players based on the provided request parameters.
func (c *Client) GetReturningProduction(
   ctx context.Context,
   request ReturningProductionRequest,
) ([]*ReturningProduction, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/player/returning", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /player/returning; %w", err)
   }

   var production []*ReturningProduction
   if err = c.unmarshalList(
      response, &production, &ReturningProduction{},
   ); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal returning production; %w", err,
      )
   }

   return production, nil
}

// GetTransferPortal retrieves player transfer portal information based on
// the provided request parameters.
func (c *Client) GetTransferPortal(
   ctx context.Context,
   request PlayerPortalRequest,
) ([]*PlayerTransfer, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/player/portal", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /player/portal; %w", err)
   }

   var transfers []*PlayerTransfer
   if err = c.unmarshalList(
      response, &transfers, &PlayerTransfer{},
   ); err != nil {
      return nil, fmt.Errorf("failed to unmarshal player transfers; %w", err)
   }

   return transfers, nil
}

// GetRankings retrieves college football rankings (polls) based on the
// provided request parameters.
func (c *Client) GetRankings(
   ctx context.Context,
   request RankingsRequest,
) ([]*PollWeek, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/rankings", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /rankings; %w", err)
   }

   var rankings []*PollWeek
   if err = c.unmarshalList(response, &rankings, &PollWeek{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal rankings; %w", err)
   }

   return rankings, nil
}

// GetLines retrieves betting lines for games based on the provided request
// parameters.
func (c *Client) GetLines(
   ctx context.Context,
   request LinesRequest,
) ([]*BettingGame, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/lines", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /lines; %w", err)
   }

   var games []*BettingGame
   if err = c.unmarshalList(response, &games, &BettingGame{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal betting games; %w", err)
   }

   return games, nil
}

// GetRecruitingPlayers retrieves recruiting information for players based
// on the provided request parameters.
func (c *Client) GetRecruitingPlayers(
   ctx context.Context,
   request RecruitingPlayersRequest,
) ([]*Recruit, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/recruiting/players", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /recruiting/players; %w", err)
   }

   var recruits []*Recruit
   if err = c.unmarshalList(response, &recruits, &Recruit{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal recruits; %w", err)
   }

   return recruits, nil
}

// GetRecruitingTeams retrieves team recruiting rankings based on the
// provided request parameters.
func (c *Client) GetRecruitingTeams(
   ctx context.Context,
   request RecruitingTeamsRequest,
) ([]*TeamRecruitingRanking, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/recruiting/teams", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /recruiting/teams; %w", err)
   }

   var rankings []*TeamRecruitingRanking
   if err = c.unmarshalList(
      response, &rankings, &TeamRecruitingRanking{},
   ); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal team recruiting rankings; %w", err,
      )
   }

   return rankings, nil
}

// GetRecruitingGroups retrieves aggregated team recruiting information
// based on the provided request parameters.
func (c *Client) GetRecruitingGroups(
   ctx context.Context,
   request RecruitingGroupsRequest,
) ([]*AggregatedTeamRecruiting, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/recruiting/groups", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /recruiting/groups; %w", err)
   }

   var groups []*AggregatedTeamRecruiting
   if err = c.unmarshalList(
      response, &groups, &AggregatedTeamRecruiting{},
   ); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal aggregated team recruiting; %w", err,
      )
   }

   return groups, nil
}

// GetRatingsSP retrieves SP+ (S&P+) ratings for teams based on the
// provided request parameters.
func (c *Client) GetRatingsSP(
   ctx context.Context,
   request RatingsSpRequest,
) ([]*TeamSP, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/ratings/sp", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /ratings/sp; %w", err)
   }

   var ratings []*TeamSP
   if err = c.unmarshalList(response, &ratings, &TeamSP{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team SP ratings; %w", err)
   }

   return ratings, nil
}

// GetRatingsSPConferences retrieves SP+ (S&P+) ratings for conferences
// based on the provided request parameters.
func (c *Client) GetRatingsSPConferences(
   ctx context.Context,
   request RatingsSpConferencesRequest,
) ([]*ConferenceSP, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/ratings/sp/conferences", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf(
         "failed to request /ratings/sp/conferences; %w", err,
      )
   }

   var conferences []*ConferenceSP
   if err = c.unmarshalList(
      response, &conferences, &ConferenceSP{},
   ); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal conference SP ratings; %w", err,
      )
   }

   return conferences, nil
}

// GetRatingsSRS retrieves SRS (Simple Rating System) ratings for teams
// based on the provided request parameters.
func (c *Client) GetRatingsSRS(
   ctx context.Context,
   request RatingsSrsRequest,
) ([]*TeamSRS, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/ratings/srs", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /ratings/srs; %w", err)
   }

   var ratings []*TeamSRS
   if err = c.unmarshalList(response, &ratings, &TeamSRS{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team SRS ratings; %w", err)
   }

   return ratings, nil
}

// GetRatingsElo retrieves Elo ratings for teams based on the provided
// request parameters.
func (c *Client) GetRatingsElo(
   ctx context.Context,
   request RatingsEloRequest,
) ([]*TeamElo, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/ratings/elo", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /ratings/elo; %w", err)
   }

   var ratings []*TeamElo
   if err = c.unmarshalList(response, &ratings, &TeamElo{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team Elo ratings; %w", err)
   }

   return ratings, nil
}

// GetRatingsFPI retrieves FPI (Football Power Index) ratings for teams
// based on the provided request parameters.
func (c *Client) GetRatingsFPI(
   ctx context.Context,
   request RatingsFpiRequest,
) ([]*TeamFPI, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/ratings/fpi", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /ratings/fpi; %w", err)
   }

   var ratings []*TeamFPI
   if err = c.unmarshalList(response, &ratings, &TeamFPI{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team FPI ratings; %w", err)
   }

   return ratings, nil
}

// GetPredictedPoints retrieves predicted points values based on the
// provided request parameters.
func (c *Client) GetPredictedPoints(
   ctx context.Context,
   request PredictedPointsRequest,
) ([]*PredictedPointsValue, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/ppa/predicted", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /ppa/predicted; %w", err)
   }

   var values []*PredictedPointsValue
   if err = c.unmarshalList(
      response, &values, &PredictedPointsValue{},
   ); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal predicted points values; %w", err,
      )
   }

   return values, nil
}

// GetPpaTeams retrieves team season PPA (Predicted Points Added) statistics
// based on the provided request parameters.
func (c *Client) GetPpaTeams(
   ctx context.Context,
   request PpaTeamsRequest,
) ([]*TeamSeasonPredictedPointsAdded, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/ppa/teams", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /ppa/teams; %w", err)
   }

   var teams []*TeamSeasonPredictedPointsAdded
   if err = c.unmarshalList(
      response, &teams, &TeamSeasonPredictedPointsAdded{},
   ); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team season PPA; %w", err)
   }

   return teams, nil
}

// GetPpaGames retrieves team game PPA (Predicted Points Added) statistics
// based on the provided request parameters.
func (c *Client) GetPpaGames(
   ctx context.Context,
   request PpaGamesRequest,
) ([]*TeamGamePredictedPointsAdded, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/ppa/games", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /ppa/games; %w", err)
   }

   var games []*TeamGamePredictedPointsAdded
   if err = c.unmarshalList(
      response, &games, &TeamGamePredictedPointsAdded{},
   ); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team game PPA; %w", err)
   }

   return games, nil
}

// GetPlayerPpaGames retrieves player game PPA (Predicted Points Added)
// statistics based on the provided request parameters.
func (c *Client) GetPlayerPpaGames(
   ctx context.Context,
   request PlayerPpaGamesRequest,
) ([]*PlayerGamePredictedPointsAdded, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/ppa/players/games", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /ppa/players/games; %w", err)
   }

   var games []*PlayerGamePredictedPointsAdded
   if err = c.unmarshalList(
      response, &games, &PlayerGamePredictedPointsAdded{},
   ); err != nil {
      return nil, fmt.Errorf("failed to unmarshal player game PPA; %w", err)
   }

   return games, nil
}

// GetPlayerPpaSeason retrieves player season PPA (Predicted Points Added)
// statistics based on the provided request parameters.
func (c *Client) GetPlayerPpaSeason(
   ctx context.Context,
   request PlayerPpaSeasonRequest,
) ([]*PlayerSeasonPredictedPointsAdded, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/ppa/players/season", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /ppa/players/season; %w", err)
   }

   var players []*PlayerSeasonPredictedPointsAdded
   if err = c.unmarshalList(
      response, &players, &PlayerSeasonPredictedPointsAdded{},
   ); err != nil {
      return nil, fmt.Errorf("failed to unmarshal player season PPA; %w", err)
   }

   return players, nil
}

// GetWinProbability retrieves win probability data for each play in the
// specified game.
func (c *Client) GetWinProbability(
   ctx context.Context,
   gameID int32,
) ([]*PlayWinProbability, error) {
   params := url.Values{}
   params.Set("gameId", strconv.FormatInt(int64(gameID), 10))

   response, err := c.executor.execute(ctx, "/metrics/wp", params)
   if err != nil {
      return nil, fmt.Errorf("failed to request /metrics/wp; %w", err)
   }

   var probs []*PlayWinProbability
   if err = c.unmarshalList(
      response, &probs, &PlayWinProbability{},
   ); err != nil {
      return nil, fmt.Errorf("failed to unmarshal win probabilities; %w", err)
   }

   return probs, nil
}

// GetPregameWinProbability retrieves pregame win probability data based on
// the provided request parameters.
func (c *Client) GetPregameWinProbability(
   ctx context.Context,
   request PregameWpRequest,
) ([]*PregameWinProbability, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/metrics/wp/pregame", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /metrics/wp/pregame; %w", err)
   }

   var probs []*PregameWinProbability
   if err = c.unmarshalList(
      response, &probs, &PregameWinProbability{},
   ); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal pregame win probabilities; %w", err,
      )
   }

   return probs, nil
}

// GetFieldGoalEP retrieves field goal expected points values.
func (c *Client) GetFieldGoalEP(ctx context.Context) ([]*FieldGoalEP, error) {
   response, err := c.executor.execute(ctx, "/metrics/fg/ep", url.Values{})
   if err != nil {
      return nil, fmt.Errorf("failed to request /metrics/fg/ep; %w", err)
   }

   var fgs []*FieldGoalEP
   if err = c.unmarshalList(response, &fgs, &FieldGoalEP{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal field goal EP; %w", err)
   }

   return fgs, nil
}

// GetPlayerSeasonStats retrieves player season statistics based on the
// provided request parameters.
func (c *Client) GetPlayerSeasonStats(
   ctx context.Context,
   request PlayerSeasonStatsRequest,
) ([]*PlayerStat, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/stats/player/season", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /stats/player/season; %w", err)
   }

   var stats []*PlayerStat
   if err = c.unmarshalList(response, &stats, &PlayerStat{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal player season stats; %w", err)
   }

   return stats, nil
}

// GetTeamSeasonStats retrieves team season statistics based on the provided
// request parameters.
func (c *Client) GetTeamSeasonStats(
   ctx context.Context,
   request TeamSeasonStatsRequest,
) ([]*TeamStat, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/stats/season", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /stats/season; %w", err)
   }

   var stats []*TeamStat
   if err = c.unmarshalList(response, &stats, &TeamStat{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal team season stats; %w", err)
   }

   return stats, nil
}

// GetStatsCategories retrieves all available statistics categories.
func (c *Client) GetStatsCategories(ctx context.Context) ([]string, error) {
   response, err := c.executor.execute(ctx, "/stats/categories", url.Values{})
   if err != nil {
      return nil, fmt.Errorf("failed to request /stats/categories; %w", err)
   }

   var out []string
   if err := json.Unmarshal(response, &out); err != nil {
      return nil, fmt.Errorf("failed to unmarshal stats categories; %w", err)
   }

   return out, nil
}

// GetAdvancedSeasonStats retrieves advanced season statistics based on the
// provided request parameters.
func (c *Client) GetAdvancedSeasonStats(
   ctx context.Context,
   request AdvancedSeasonStatsRequest,
) ([]*AdvancedSeasonStat, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/stats/season/advanced", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf(
         "failed to request /stats/season/advanced; %w", err,
      )
   }

   var stats []*AdvancedSeasonStat
   if err = c.unmarshalList(
      response, &stats, &AdvancedSeasonStat{},
   ); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal advanced season stats; %w", err,
      )
   }

   return stats, nil
}

// GetAdvancedGameStats retrieves advanced game statistics based on the
// provided request parameters.
func (c *Client) GetAdvancedGameStats(
   ctx context.Context,
   req AdvancedGameStatsRequest,
) ([]*AdvancedGameStat, error) {
   if err := req.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   resp, err := c.executor.execute(ctx, "/stats/game/advanced", req.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /stats/game/advanced; %w", err)
   }

   var stats []*AdvancedGameStat
   if err = c.unmarshalList(resp, &stats, &AdvancedGameStat{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal advanced game stats; %w", err)
   }

   return stats, nil
}

// GetHavocGameStats retrieves havoc game statistics based on the provided
// request parameters.
func (c *Client) GetHavocGameStats(
   ctx context.Context,
   request HavocGameStatsRequest,
) ([]*GameHavocStats, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/stats/game/havoc", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /stats/game/havoc; %w", err)
   }

   var stats []*GameHavocStats
   if err = c.unmarshalList(response, &stats, &GameHavocStats{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal game havoc stats; %w", err)
   }

   return stats, nil
}

// GetDraftTeams retrieves all NFL draft teams.
func (c *Client) GetDraftTeams(ctx context.Context) ([]*DraftTeam, error) {
   response, err := c.executor.execute(ctx, "/draft/teams", url.Values{})
   if err != nil {
      return nil, fmt.Errorf("failed to request /draft/teams; %w", err)
   }

   var teams []*DraftTeam
   if err = c.unmarshalList(response, &teams, &DraftTeam{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal draft teams; %w", err)
   }

   return teams, nil
}

// GetDraftPositions retrieves all NFL draft positions.
func (c *Client) GetDraftPositions(
   ctx context.Context,
) ([]*DraftPosition, error) {
   response, err := c.executor.execute(ctx, "/draft/positions", url.Values{})
   if err != nil {
      return nil, fmt.Errorf("failed to request /draft/positions; %w", err)
   }

   var positions []*DraftPosition
   if err = c.unmarshalList(
      response, &positions, &DraftPosition{},
   ); err != nil {
      return nil, fmt.Errorf("failed to unmarshal draft positions; %w", err)
   }

   return positions, nil
}

// GetDraftPicks retrieves NFL draft picks based on the provided request
// parameters.
func (c *Client) GetDraftPicks(
   ctx context.Context,
   request DraftPicksRequest,
) ([]*DraftPick, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(ctx, "/draft/picks", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /draft/picks; %w", err)
   }

   var picks []*DraftPick
   if err = c.unmarshalList(response, &picks, &DraftPick{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal draft picks; %w", err)
   }

   return picks, nil
}

// GetWepaTeamSeason retrieves team season WEPA (Weighted Expected Points
// Added) metrics based on the provided request parameters.
func (c *Client) GetWepaTeamSeason(
   ctx context.Context,
   request WepaTeamSeasonRequest,
) ([]*AdjustedTeamMetrics, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   resp, err := c.executor.execute(ctx, "/wepa/team/season", request.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /wepa/team/season; %w", err)
   }

   var teams []*AdjustedTeamMetrics
   if err = c.unmarshalList(resp, &teams, &AdjustedTeamMetrics{}); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal adjusted team metrics; %w", err,
      )
   }

   return teams, nil
}

// GetWepaPlayersPassing retrieves player passing WEPA (Weighted Expected
// Points Added) metrics based on the provided request parameters.
func (c *Client) GetWepaPlayersPassing(
   ctx context.Context,
   request WepaPlayersRequest,
) ([]*PlayerWeightedEPA, error) {
   if err := request.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   resp, err := c.executor.execute(
      ctx, "/wepa/players/passing", request.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /wepa/players/passing; %w", err)
   }

   var players []*PlayerWeightedEPA
   if err = c.unmarshalList(resp, &players, &PlayerWeightedEPA{}); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal player weighted EPA (passing); %w", err,
      )
   }

   return players, nil
}

// GetWepaPlayersRushing retrieves player rushing WEPA (Weighted Expected
// Points Added) metrics based on the provided request parameters.
func (c *Client) GetWepaPlayersRushing(
   ctx context.Context,
   req WepaPlayersRequest,
) ([]*PlayerWeightedEPA, error) {
   if err := req.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   resp, err := c.executor.execute(ctx, "/wepa/players/rushing", req.values())
   if err != nil {
      return nil, fmt.Errorf("failed to request /wepa/players/rushing; %w", err)
   }

   var players []*PlayerWeightedEPA
   if err = c.unmarshalList(resp, &players, &PlayerWeightedEPA{}); err != nil {
      return nil, fmt.Errorf(
         "failed to unmarshal player weighted EPA (rushing); %w", err,
      )
   }

   return players, nil
}

// GetWepaPlayersKicking retrieves kicker PAAR (Points Above Average
// Replacement) metrics based on the provided request parameters.
func (c *Client) GetWepaPlayersKicking(
   ctx context.Context,
   req WepaKickersRequest,
) ([]*KickerPAAR, error) {
   if err := req.validate(); err != nil {
      return nil, fmt.Errorf("request could not be validated; %w", err)
   }

   response, err := c.executor.execute(
      ctx, "/wepa/players/kicking", req.values(),
   )
   if err != nil {
      return nil, fmt.Errorf("failed to request /wepa/players/kicking; %w", err)
   }

   var kickers []*KickerPAAR
   if err = c.unmarshalList(response, &kickers, &KickerPAAR{}); err != nil {
      return nil, fmt.Errorf("failed to unmarshal kicker PAAR; %w", err)
   }

   return kickers, nil
}

// GetInfo retrieves information about the authenticated user's API key.
// Returns nil if the user is not authenticated.
func (c *Client) GetInfo(ctx context.Context) (*UserInfo, error) {
   response, err := c.executor.execute(ctx, "/info", url.Values{})
   if err != nil {
      return nil, fmt.Errorf("failed to request /info endpoint; %w", err)
   }

   // /info may return null if not authenticated.
   if isJSONNull(response) {
      return nil, nil
   }

   var userInfo UserInfo
   if err = c.unmarshal(response, &userInfo); err != nil {
      return nil, fmt.Errorf("unable to retrieve user information; %w", err)
   }

   return &userInfo, nil
}

func (c *Client) unmarshal(b []byte, out proto.Message) error {
   if out == nil {
      return fmt.Errorf("out cannot be nil")
   }
   if len(bytes.TrimSpace(b)) == 0 || isJSONNull(b) {
      return nil
   }

   if err := c.unmarshaller.Unmarshal(b, out); err != nil {
      return fmt.Errorf("")
   }

   return nil
}

func (c *Client) unmarshalList(
   b []byte, out any, prototype proto.Message,
) error {
   if len(bytes.TrimSpace(b)) == 0 || isJSONNull(b) {
      return nil
   }
   if prototype == nil {
      return fmt.Errorf("prototype cannot be nil (e.g. &pb.Drive{})")
   }

   rv := reflect.ValueOf(out)
   if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Slice {
      return fmt.Errorf("out must be pointer to slice, got %T", out)
   }

   var raws []json.RawMessage
   if err := json.Unmarshal(b, &raws); err != nil {
      return err
   }

   slice := rv.Elem()
   for _, raw := range raws {
      if isJSONNull(raw) {
         continue
      }

      msg := proto.Clone(prototype)
      if err := c.unmarshaller.Unmarshal(raw, msg); err != nil {
         return err
      }

      // Ensure msg type matches slice element type
      msgV := reflect.ValueOf(msg)
      if !msgV.Type().AssignableTo(slice.Type().Elem()) {
         return fmt.Errorf(
            "prototype type %T not assignable to slice element type %s",
            msg, slice.Type().Elem(),
         )
      }

      slice = reflect.Append(slice, msgV)
   }

   rv.Elem().Set(slice)
   return nil
}
