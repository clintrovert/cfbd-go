package cfbd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

// baseURL is the CFBD REST API base URL.
const baseURL = "https://api.collegefootballdata.com"

// Client is a minimal REST client that unmarshals responses into Protobuf
// messages using protojson (so swagger JSON like camelCase works with
// snake_case proto fields).
//
// Authentication: CFBD uses an API key as a Bearer token in the Authorization
// header.
// Example:
//
//	Authorization: Bearer <your_api_key>
//
// Notes:
// - This client is intentionally thin and predictable.
// - All methods accept a context.Context.
// - For endpoints that return JSON arrays, the client unmarshals each element
// into a message.
// - Unknown JSON fields are discarded by default to tolerate API evolution.
//
// Reference for authentication header examples:
// https://blog.collegefootballdata.com/using-api-keys-with-the-cfbd-api/
//
// If you want retries/backoff, wrap the http.RoundTripper or add middleware.
//
//go:generate echo "This pkg uses generated protobuf bindings from cfbd/proto/cfbd.proto"
type Client struct {
	baseURL   *url.URL
	apiKey    string
	client    *http.Client
	userAgent string
	unmarshal protojson.UnmarshalOptions
}

// NewClient creates a Client with sane defaults.
func NewClient(apiKey string) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse base url; %w", err)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key was not provided")
	}

	return &Client{
		baseURL: base,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "cfbd-go/1.0",
		unmarshal: protojson.UnmarshalOptions{
			DiscardUnknown: true,
			AllowPartial:   true,
		},
	}, nil
}

// -----------------------------
// games
// -----------------------------

func (c *Client) GetGames(
	ctx context.Context,
	p GetGamesRequest,
) ([]*Game, error) {
	b, err := c.doGetRequest(ctx, "/games", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /games; %w", err)
	}

	var games []*Game
	if err = c.unmarshalList(b, &games, &Game{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal games; %w", err)
	}

	return games, nil
}

// -----------------------------
// games (additional endpoints)
// -----------------------------

// GetGameTeamStats calls GET /games/teams (team box score stats).
func (c *Client) GetGameTeamStats(ctx context.Context, p GameTeamStatsRequest) ([]*GameTeamStats, error) {
	b, err := c.doGetRequest(ctx, "/games/teams", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /games/teams; %w", err)
	}
	var games []*GameTeamStats
	if err = c.unmarshalList(b, &games, &GameTeamStats{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game team stats; %w", err)
	}
	return games, nil
}

// GetGamePlayerStats calls GET /games/players (player box score stats).
func (c *Client) GetGamePlayerStats(ctx context.Context, p GamePlayerStatsRequest) ([]*GamePlayerStats, error) {
	b, err := c.doGetRequest(ctx, "/games/players", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /games/players; %w", err)
	}
	var games []*GamePlayerStats
	if err = c.unmarshalList(b, &games, &GamePlayerStats{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game player stats; %w", err)
	}
	return games, nil
}

// GetGameMedia calls GET /games/media.
func (c *Client) GetGameMedia(ctx context.Context, p GameMediaRequest) ([]*GameMedia, error) {
	b, err := c.doGetRequest(ctx, "/games/media", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /games/media; %w", err)
	}
	var games []*GameMedia
	if err = c.unmarshalList(b, &games, &GameMedia{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game media; %w", err)
	}
	return games, nil
}

// GetGameWeather calls GET /games/weather.
func (c *Client) GetGameWeather(
	ctx context.Context,
	p GameWeatherRequest,
) ([]*GameWeather, error) {
	b, err := c.doGetRequest(ctx, "/games/weather", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /games/weather; %w", err)
	}
	var games []*GameWeather
	if err = c.unmarshalList(b, &games, &GameWeather{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game weather; %w", err)
	}
	return games, nil
}

// GetAdvancedBoxScore calls GET /game/box/advanced.
func (c *Client) GetAdvancedBoxScore(
	ctx context.Context,
	gameID int32,
) (*AdvancedBoxScore, error) {
	v := url.Values{}
	v.Set("gameId", strconv.FormatInt(int64(gameID), 10))
	b, err := c.doGetRequest(ctx, "/game/box/advanced", v)
	if err != nil {
		return nil, fmt.Errorf("failed to request /game/box/advanced; %w", err)
	}

	var val AdvancedBoxScore

	if err = c.unmarshalInto(b, &val); err != nil {
		return nil, fmt.Errorf("failed to unmarshal advanced box score; %w", err)
	}

	return &val, nil
}

// GetCalendar calls GET /calendar.
func (c *Client) GetCalendar(ctx context.Context, year int32) ([]*CalendarWeek, error) {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(year), 10))
	b, err := c.doGetRequest(ctx, "/calendar", v)
	if err != nil {
		return nil, fmt.Errorf("failed to request /calendar; %w", err)
	}
	var weeks []*CalendarWeek
	if err = c.unmarshalList(b, &weeks, &CalendarWeek{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal calendar weeks; %w", err)
	}
	return weeks, nil
}

// GetTeamRecords calls GET /records.
func (c *Client) GetTeamRecords(ctx context.Context, p RecordsRequest) ([]*TeamRecords, error) {
	b, err := c.doGetRequest(ctx, "/records", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /records; %w", err)
	}
	var records []*TeamRecords
	if err = c.unmarshalList(b, &records, &TeamRecords{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team records; %w", err)
	}
	return records, nil
}

// GetLiveScoreboard calls GET /scoreboard.
func (c *Client) GetLiveScoreboard(ctx context.Context, p LiveScoreboardRequest) ([]*ScoreboardGame, error) {
	b, err := c.doGetRequest(ctx, "/scoreboard", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /scoreboard; %w", err)
	}
	var games []*ScoreboardGame
	if err = c.unmarshalList(b, &games, &ScoreboardGame{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scoreboard games; %w", err)
	}
	return games, nil
}

// -----------------------------
// drives
// -----------------------------

func (c *Client) GetDrives(ctx context.Context, p DrivesRequest) ([]*Drive, error) {
	b, err := c.doGetRequest(ctx, "/drives", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /drives; %w", err)
	}
	var drives []*Drive
	if err = c.unmarshalList(b, &drives, &Drive{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal drives; %w", err)
	}
	return drives, nil
}

// -----------------------------
// plays
// -----------------------------

func (c *Client) GetPlays(ctx context.Context, p PlaysRequest) ([]*Play, error) {
	b, err := c.doGetRequest(ctx, "/plays", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /plays; %w", err)
	}
	var plays []*Play
	if err = c.unmarshalList(b, &plays, &Play{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plays; %w", err)
	}
	return plays, nil
}

func (c *Client) GetPlayTypes(ctx context.Context) ([]*PlayType, error) {
	b, err := c.doGetRequest(ctx, "/plays/types", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /plays/types; %w", err)
	}
	var playTypes []*PlayType
	if err = c.unmarshalList(b, &playTypes, &PlayType{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal play types; %w", err)
	}
	return playTypes, nil
}

func (c *Client) GetPlayStats(ctx context.Context, p PlayStatsRequest) ([]*PlayStat, error) {
	b, err := c.doGetRequest(ctx, "/plays/stats", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /plays/stats; %w", err)
	}
	var stats []*PlayStat
	if err = c.unmarshalList(b, &stats, &PlayStat{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal play stats; %w", err)
	}
	return stats, nil
}

func (c *Client) GetPlayStatTypes(ctx context.Context) ([]*PlayStatType, error) {
	b, err := c.doGetRequest(ctx, "/plays/stats/types", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /plays/stats/types; %w", err)
	}
	var statTypes []*PlayStatType
	if err = c.unmarshalList(b, &statTypes, &PlayStatType{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal play stat types; %w", err)
	}
	return statTypes, nil
}

// live play-by-play
func (c *Client) GetLivePlays(ctx context.Context, gameID int32) (*LiveGame, error) {
	q := url.Values{}
	q.Set("gameId", strconv.FormatInt(int64(gameID), 10))
	b, err := c.doGetRequest(ctx, "/live/plays", q)
	if err != nil {
		return nil, fmt.Errorf("failed to request /live/plays; %w", err)
	}
	var game LiveGame
	if err = c.unmarshalInto(b, &game); err != nil {
		return nil, fmt.Errorf("failed to unmarshal live game; %w", err)
	}
	return &game, nil
}

// -----------------------------
// teams
// -----------------------------

func (c *Client) GetTeams(ctx context.Context, p TeamsRequest) ([]*Team, error) {
	b, err := c.doGetRequest(ctx, "/teams", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /teams; %w", err)
	}
	var teams []*Team
	if err = c.unmarshalList(b, &teams, &Team{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal teams; %w", err)
	}
	return teams, nil
}

func (c *Client) GetTeamsFBS(ctx context.Context, p TeamsFbsRequest) ([]*Team, error) {
	b, err := c.doGetRequest(ctx, "/teams/fbs", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /teams/fbs; %w", err)
	}
	var teams []*Team
	if err = c.unmarshalList(b, &teams, &Team{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal teams; %w", err)
	}
	return teams, nil
}

func (c *Client) GetTeamMatchup(ctx context.Context, p TeamMatchupRequest) (*Matchup, error) {
	b, err := c.doGetRequest(ctx, "/teams/matchup", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /teams/matchup; %w", err)
	}
	var matchup Matchup
	if err = c.unmarshalInto(b, &matchup); err != nil {
		return nil, fmt.Errorf("failed to unmarshal matchup; %w", err)
	}
	return &matchup, nil
}

func (c *Client) GetTeamATS(ctx context.Context, p TeamATSRequest) ([]*TeamATS, error) {
	b, err := c.doGetRequest(ctx, "/teams/ats", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /teams/ats; %w", err)
	}
	var teams []*TeamATS
	if err = c.unmarshalList(b, &teams, &TeamATS{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team ATS; %w", err)
	}
	return teams, nil
}

// -----------------------------
// roster
// -----------------------------

func (c *Client) GetRoster(ctx context.Context, p RosterRequest) ([]*RosterPlayer, error) {
	b, err := c.doGetRequest(ctx, "/roster", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /roster; %w", err)
	}
	var players []*RosterPlayer
	if err = c.unmarshalList(b, &players, &RosterPlayer{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal roster players; %w", err)
	}
	return players, nil
}

// -----------------------------
// talent
// -----------------------------

func (c *Client) GetTeamTalent(ctx context.Context, p TalentRequest) ([]*TeamTalent, error) {
	b, err := c.doGetRequest(ctx, "/talent", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /talent; %w", err)
	}
	var talents []*TeamTalent
	if err = c.unmarshalList(b, &talents, &TeamTalent{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team talent; %w", err)
	}
	return talents, nil
}

// -----------------------------
// conferences
// -----------------------------

func (c *Client) GetConferences(ctx context.Context) ([]*Conference, error) {
	b, err := c.doGetRequest(ctx, "/conferences", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /conferences; %w", err)
	}
	var conferences []*Conference
	if err = c.unmarshalList(b, &conferences, &Conference{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conferences; %w", err)
	}
	return conferences, nil
}

// -----------------------------
// venues
// -----------------------------

func (c *Client) GetVenues(ctx context.Context) ([]*Venue, error) {
	b, err := c.doGetRequest(ctx, "/venues", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /venues; %w", err)
	}
	var venues []*Venue
	if err = c.unmarshalList(b, &venues, &Venue{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal venues; %w", err)
	}
	return venues, nil
}

// -----------------------------
// coaches
// -----------------------------

func (c *Client) GetCoaches(ctx context.Context, p CoachesRequest) ([]*Coach, error) {
	b, err := c.doGetRequest(ctx, "/coaches", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /coaches; %w", err)
	}
	var coaches []*Coach
	if err = c.unmarshalList(b, &coaches, &Coach{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal coaches; %w", err)
	}
	return coaches, nil
}

// -----------------------------
// players
// -----------------------------

func (c *Client) SearchPlayers(ctx context.Context, p PlayerSearchRequest) ([]*PlayerSearchResult, error) {
	b, err := c.doGetRequest(ctx, "/player/search", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /player/search; %w", err)
	}
	var players []*PlayerSearchResult
	if err = c.unmarshalList(b, &players, &PlayerSearchResult{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player search results; %w", err)
	}
	return players, nil
}

func (c *Client) GetPlayerUsage(ctx context.Context, p PlayerUsageRequest) ([]*PlayerUsage, error) {
	b, err := c.doGetRequest(ctx, "/player/usage", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /player/usage; %w", err)
	}
	var usage []*PlayerUsage
	if err = c.unmarshalList(b, &usage, &PlayerUsage{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player usage; %w", err)
	}
	return usage, nil
}

func (c *Client) GetReturningProduction(ctx context.Context, p ReturningProductionRequest) ([]*ReturningProduction, error) {
	b, err := c.doGetRequest(ctx, "/player/returning", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /player/returning; %w", err)
	}
	var production []*ReturningProduction
	if err = c.unmarshalList(b, &production, &ReturningProduction{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal returning production; %w", err)
	}
	return production, nil
}

func (c *Client) GetTransferPortal(ctx context.Context, p PlayerPortalRequest) ([]*PlayerTransfer, error) {
	b, err := c.doGetRequest(ctx, "/player/portal", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /player/portal; %w", err)
	}
	var transfers []*PlayerTransfer
	if err = c.unmarshalList(b, &transfers, &PlayerTransfer{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player transfers; %w", err)
	}
	return transfers, nil
}

// -----------------------------
// rankings
// -----------------------------

func (c *Client) GetRankings(ctx context.Context, p RankingsRequest) ([]*PollWeek, error) {
	b, err := c.doGetRequest(ctx, "/rankings", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /rankings; %w", err)
	}
	var rankings []*PollWeek
	if err = c.unmarshalList(b, &rankings, &PollWeek{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rankings; %w", err)
	}
	return rankings, nil
}

// -----------------------------
// betting lines
// -----------------------------

func (c *Client) GetLines(ctx context.Context, p LinesRequest) ([]*BettingGame, error) {
	b, err := c.doGetRequest(ctx, "/lines", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /lines; %w", err)
	}
	var games []*BettingGame
	if err = c.unmarshalList(b, &games, &BettingGame{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal betting games; %w", err)
	}
	return games, nil
}

// -----------------------------
// recruiting
// -----------------------------

func (c *Client) GetRecruitingPlayers(ctx context.Context, p RecruitingPlayersRequest) ([]*Recruit, error) {
	b, err := c.doGetRequest(ctx, "/recruiting/players", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /recruiting/players; %w", err)
	}
	var recruits []*Recruit
	if err = c.unmarshalList(b, &recruits, &Recruit{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recruits; %w", err)
	}
	return recruits, nil
}

func (c *Client) GetRecruitingTeams(ctx context.Context, p RecruitingTeamsRequest) ([]*TeamRecruitingRanking, error) {
	b, err := c.doGetRequest(ctx, "/recruiting/teams", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /recruiting/teams; %w", err)
	}
	var rankings []*TeamRecruitingRanking
	if err = c.unmarshalList(b, &rankings, &TeamRecruitingRanking{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team recruiting rankings; %w", err)
	}
	return rankings, nil
}

func (c *Client) GetRecruitingGroups(ctx context.Context, p RecruitingGroupsRequest) ([]*AggregatedTeamRecruiting, error) {
	b, err := c.doGetRequest(ctx, "/recruiting/groups", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /recruiting/groups; %w", err)
	}
	var groups []*AggregatedTeamRecruiting
	if err = c.unmarshalList(b, &groups, &AggregatedTeamRecruiting{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal aggregated team recruiting; %w", err)
	}
	return groups, nil
}

// -----------------------------
// ratings
// -----------------------------

func (c *Client) GetRatingsSP(ctx context.Context, p RatingsSpRequest) ([]*TeamSP, error) {
	b, err := c.doGetRequest(ctx, "/ratings/sp", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ratings/sp; %w", err)
	}
	var ratings []*TeamSP
	if err = c.unmarshalList(b, &ratings, &TeamSP{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team SP ratings; %w", err)
	}
	return ratings, nil
}

func (c *Client) GetRatingsSPConferences(ctx context.Context, p RatingsSpConferencesRequest) ([]*ConferenceSP, error) {
	b, err := c.doGetRequest(ctx, "/ratings/sp/conferences", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ratings/sp/conferences; %w", err)
	}
	var conferences []*ConferenceSP
	if err = c.unmarshalList(b, &conferences, &ConferenceSP{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conference SP ratings; %w", err)
	}
	return conferences, nil
}

func (c *Client) GetRatingsSRS(ctx context.Context, p RatingsSrsRequest) ([]*TeamSRS, error) {
	b, err := c.doGetRequest(ctx, "/ratings/srs", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ratings/srs; %w", err)
	}
	var ratings []*TeamSRS
	if err = c.unmarshalList(b, &ratings, &TeamSRS{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team SRS ratings; %w", err)
	}
	return ratings, nil
}

func (c *Client) GetRatingsElo(ctx context.Context, p RatingsEloRequest) ([]*TeamElo, error) {
	b, err := c.doGetRequest(ctx, "/ratings/elo", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ratings/elo; %w", err)
	}
	var ratings []*TeamElo
	if err = c.unmarshalList(b, &ratings, &TeamElo{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team Elo ratings; %w", err)
	}
	return ratings, nil
}

func (c *Client) GetRatingsFPI(ctx context.Context, p RatingsFpiRequest) ([]*TeamFPI, error) {
	b, err := c.doGetRequest(ctx, "/ratings/fpi", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ratings/fpi; %w", err)
	}
	var ratings []*TeamFPI
	if err = c.unmarshalList(b, &ratings, &TeamFPI{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team FPI ratings; %w", err)
	}
	return ratings, nil
}

// -----------------------------
// metrics
// -----------------------------

func (c *Client) GetPredictedPoints(ctx context.Context, p PredictedPointsRequest) ([]*PredictedPointsValue, error) {
	b, err := c.doGetRequest(ctx, "/ppa/predicted", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ppa/predicted; %w", err)
	}
	var values []*PredictedPointsValue
	if err = c.unmarshalList(b, &values, &PredictedPointsValue{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal predicted points values; %w", err)
	}
	return values, nil
}

func (c *Client) GetPpaTeams(ctx context.Context, p PpaTeamsRequest) ([]*TeamSeasonPredictedPointsAdded, error) {
	b, err := c.doGetRequest(ctx, "/ppa/teams", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ppa/teams; %w", err)
	}
	var teams []*TeamSeasonPredictedPointsAdded
	if err = c.unmarshalList(b, &teams, &TeamSeasonPredictedPointsAdded{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team season PPA; %w", err)
	}
	return teams, nil
}

func (c *Client) GetPpaGames(ctx context.Context, p PpaGamesRequest) ([]*TeamGamePredictedPointsAdded, error) {
	b, err := c.doGetRequest(ctx, "/ppa/games", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ppa/games; %w", err)
	}
	var games []*TeamGamePredictedPointsAdded
	if err = c.unmarshalList(b, &games, &TeamGamePredictedPointsAdded{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team game PPA; %w", err)
	}
	return games, nil
}

func (c *Client) GetPlayerPpaGames(ctx context.Context, p PlayerPpaGamesRequest) ([]*PlayerGamePredictedPointsAdded, error) {
	b, err := c.doGetRequest(ctx, "/ppa/players/games", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ppa/players/games; %w", err)
	}
	var games []*PlayerGamePredictedPointsAdded
	if err = c.unmarshalList(b, &games, &PlayerGamePredictedPointsAdded{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player game PPA; %w", err)
	}
	return games, nil
}

func (c *Client) GetPlayerPpaSeason(ctx context.Context, p PlayerPpaSeasonRequest) ([]*PlayerSeasonPredictedPointsAdded, error) {
	b, err := c.doGetRequest(ctx, "/ppa/players/season", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /ppa/players/season; %w", err)
	}
	var players []*PlayerSeasonPredictedPointsAdded
	if err = c.unmarshalList(b, &players, &PlayerSeasonPredictedPointsAdded{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player season PPA; %w", err)
	}
	return players, nil
}

// Win probabilities by game.
func (c *Client) GetWinProbability(ctx context.Context, gameID int32) ([]*PlayWinProbability, error) {
	v := url.Values{}
	v.Set("gameId", strconv.FormatInt(int64(gameID), 10))
	b, err := c.doGetRequest(ctx, "/metrics/wp", v)
	if err != nil {
		return nil, fmt.Errorf("failed to request /metrics/wp; %w", err)
	}
	var probs []*PlayWinProbability
	if err = c.unmarshalList(b, &probs, &PlayWinProbability{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal win probabilities; %w", err)
	}
	return probs, nil
}

func (c *Client) GetPregameWinProbability(ctx context.Context, p PregameWpRequest) ([]*PregameWinProbability, error) {
	b, err := c.doGetRequest(ctx, "/metrics/wp/pregame", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /metrics/wp/pregame; %w", err)
	}
	var probs []*PregameWinProbability
	if err = c.unmarshalList(b, &probs, &PregameWinProbability{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pregame win probabilities; %w", err)
	}
	return probs, nil
}

// Field goal expected points values.
func (c *Client) GetFieldGoalEP(ctx context.Context) ([]*FieldGoalEP, error) {
	b, err := c.doGetRequest(ctx, "/metrics/fg/ep", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /metrics/fg/ep; %w", err)
	}
	var fgs []*FieldGoalEP
	if err = c.unmarshalList(b, &fgs, &FieldGoalEP{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal field goal EP; %w", err)
	}
	return fgs, nil
}

// -----------------------------
// stats
// -----------------------------

func (c *Client) GetPlayerSeasonStats(ctx context.Context, p PlayerSeasonStatsRequest) ([]*PlayerStat, error) {
	b, err := c.doGetRequest(ctx, "/stats/player/season", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /stats/player/season; %w", err)
	}
	var stats []*PlayerStat
	if err = c.unmarshalList(b, &stats, &PlayerStat{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player season stats; %w", err)
	}
	return stats, nil
}

func (c *Client) GetTeamSeasonStats(ctx context.Context, p TeamSeasonStatsRequest) ([]*TeamStat, error) {
	b, err := c.doGetRequest(ctx, "/stats/season", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /stats/season; %w", err)
	}
	var stats []*TeamStat
	if err = c.unmarshalList(b, &stats, &TeamStat{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team season stats; %w", err)
	}
	return stats, nil
}

func (c *Client) GetStatsCategories(ctx context.Context) ([]string, error) {
	b, err := c.doGetRequest(ctx, "/stats/categories", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /stats/categories; %w", err)
	}
	var out []string
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stats categories; %w", err)
	}
	return out, nil
}

func (c *Client) GetAdvancedSeasonStats(ctx context.Context, p AdvancedSeasonStatsRequest) ([]*AdvancedSeasonStat, error) {
	b, err := c.doGetRequest(ctx, "/stats/season/advanced", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /stats/season/advanced; %w", err)
	}
	var stats []*AdvancedSeasonStat
	if err = c.unmarshalList(b, &stats, &AdvancedSeasonStat{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal advanced season stats; %w", err)
	}
	return stats, nil
}

func (c *Client) GetAdvancedGameStats(ctx context.Context, p AdvancedGameStatsRequest) ([]*AdvancedGameStat, error) {
	b, err := c.doGetRequest(ctx, "/stats/game/advanced", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /stats/game/advanced; %w", err)
	}
	var stats []*AdvancedGameStat
	if err = c.unmarshalList(b, &stats, &AdvancedGameStat{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal advanced game stats; %w", err)
	}
	return stats, nil
}

func (c *Client) GetHavocGameStats(ctx context.Context, p HavocGameStatsRequest) ([]*GameHavocStats, error) {
	b, err := c.doGetRequest(ctx, "/stats/game/havoc", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /stats/game/havoc; %w", err)
	}
	var stats []*GameHavocStats
	if err = c.unmarshalList(b, &stats, &GameHavocStats{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game havoc stats; %w", err)
	}
	return stats, nil
}

// -----------------------------
// draft
// -----------------------------

func (c *Client) GetDraftTeams(ctx context.Context) ([]*DraftTeam, error) {
	b, err := c.doGetRequest(ctx, "/draft/teams", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /draft/teams; %w", err)
	}
	var teams []*DraftTeam
	if err = c.unmarshalList(b, &teams, &DraftTeam{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal draft teams; %w", err)
	}
	return teams, nil
}

func (c *Client) GetDraftPositions(ctx context.Context) ([]*DraftPosition, error) {
	b, err := c.doGetRequest(ctx, "/draft/positions", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /draft/positions; %w", err)
	}
	var positions []*DraftPosition
	if err = c.unmarshalList(b, &positions, &DraftPosition{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal draft positions; %w", err)
	}
	return positions, nil
}

func (c *Client) GetDraftPicks(ctx context.Context, p DraftPicksRequest) ([]*DraftPick, error) {
	b, err := c.doGetRequest(ctx, "/draft/picks", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /draft/picks; %w", err)
	}
	var picks []*DraftPick
	if err = c.unmarshalList(b, &picks, &DraftPick{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal draft picks; %w", err)
	}
	return picks, nil
}

// -----------------------------
// adjusted metrics (wepa)
// -----------------------------

func (c *Client) GetWepaTeamSeason(ctx context.Context, p WepaTeamSeasonRequest) ([]*AdjustedTeamMetrics, error) {
	b, err := c.doGetRequest(ctx, "/wepa/team/season", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /wepa/team/season; %w", err)
	}
	var teams []*AdjustedTeamMetrics
	if err = c.unmarshalList(b, &teams, &AdjustedTeamMetrics{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal adjusted team metrics; %w", err)
	}
	return teams, nil
}

func (c *Client) GetWepaPlayersPassing(ctx context.Context, p WepaPlayersRequest) ([]*PlayerWeightedEPA, error) {
	b, err := c.doGetRequest(ctx, "/wepa/players/passing", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /wepa/players/passing; %w", err)
	}
	var players []*PlayerWeightedEPA
	if err = c.unmarshalList(b, &players, &PlayerWeightedEPA{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player weighted EPA (passing); %w", err)
	}
	return players, nil
}

func (c *Client) GetWepaPlayersRushing(ctx context.Context, p WepaPlayersRequest) ([]*PlayerWeightedEPA, error) {
	b, err := c.doGetRequest(ctx, "/wepa/players/rushing", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /wepa/players/rushing; %w", err)
	}
	var players []*PlayerWeightedEPA
	if err = c.unmarshalList(b, &players, &PlayerWeightedEPA{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player weighted EPA (rushing); %w", err)
	}
	return players, nil
}

func (c *Client) GetWepaPlayersKicking(ctx context.Context, p WepaKickersRequest) ([]*KickerPAAR, error) {
	b, err := c.doGetRequest(ctx, "/wepa/players/kicking", p.values())
	if err != nil {
		return nil, fmt.Errorf("failed to request /wepa/players/kicking; %w", err)
	}
	var kickers []*KickerPAAR
	if err = c.unmarshalList(b, &kickers, &KickerPAAR{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kicker PAAR; %w", err)
	}
	return kickers, nil
}

// GetInfo todo: describe.
func (c *Client) GetInfo(ctx context.Context) (*UserInfo, error) {
	var userInfo UserInfo
	b, err := c.doGetRequest(ctx, "/info", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /info endpoint; %w", err)
	}

	// /info may return null if not authenticated.
	if isJSONNull(b) {
		return nil, nil
	}

	if err = c.unmarshalInto(b, &userInfo); err != nil {
		return nil, fmt.Errorf("unable to retrieve user information; %w", err)
	}

	return &userInfo, nil
}
