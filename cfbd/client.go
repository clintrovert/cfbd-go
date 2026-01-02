// Package cfbd provides a minimal, type-safe Golang client for the College Football Data API.
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
	"strings"
	"time"

	"github.com/clintrovert/cfbd-go/cfbd/internal/httpget"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	baseURL           = "https://api.collegefootballdata.com"
	defaultTimeoutSec = 30
	userAgent         = "cfbd-go/1.0.0"

	// Common query param keys, consolidating here so they aren't misspelled.
	idKey                 = "id"
	yearKey               = "year"
	weekKey               = "week"
	teamKey               = "team"
	homeKey               = "home"
	awayKey               = "away"
	gameIDKey             = "gameId"
	minYearKey            = "minYear"
	maxYearKey            = "maxYear"
	offenseKey            = "offense"
	defenseKey            = "defense"
	endWeekKey            = "endWeek"
	endYearKey            = "endYear"
	lastNameKey           = "lastName"
	categoryKey           = "category"
	opponentKey           = "opponent"
	positionKey           = "position"
	playerIDKey           = "playerId"
	thresholdKey          = "threshold"
	startYearKey          = "startYear"
	firstNameKey          = "firstName"
	startWeekKey          = "startWeek"
	conferenceKey         = "conference"
	seasonTypeKey         = "seasonType"
	classificationKey     = "classification"
	offenseConferenceKey  = "offenseConference"
	defenseConferenceKey  = "defenseConference"
	excludeGarbageTimeKey = "excludeGarbageTime"
)

var (
	// ErrMissingAPIKey is returned if the API key provided was empty.
	ErrMissingAPIKey = errors.New("API key was not provided")
	// ErrMissingRequiredParams is returned when a request is missing
	// required parameters.
	ErrMissingRequiredParams = errors.New("request missing required params")
)

// httpGetExecutor wraps the http client with an interface for ease in mock
// testing.
type httpGetExecutor interface {
	Execute(
		ctx context.Context,
		path string,
		params url.Values,
	) ([]byte, error)
}

// Client is a REST client for the College Football Data (CFBD) API. It
// provides methods to retrieve college football statistics, game data, team
// information, and more.
//
// Authentication: CFBD uses an API key as a Bearer token in the Authorization
// header. Example:
//
//	Authorization: Bearer <your_api_key>
//
// Reference for authentication header examples:
// https://blog.collegefootballdata.com/using-api-keys-with-the-cfbd-api/
//
// Notes:
// - All methods accept a cancellable context.Context.
// - Unknown JSON fields are discarded by default to tolerate future API
// releases.
type Client struct {
	apiKey       string
	unmarshaller protojson.UnmarshalOptions
	httpGet      httpGetExecutor
}

// New creates a new Client.
func New(apiKey string) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse base url; %w", err)
	}

	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	return &Client{
		apiKey: apiKey,
		httpGet: &httpget.Client{
			APIKey:    apiKey,
			BaseURL:   base,
			UserAgent: userAgent,
			HTTPClient: &http.Client{
				Timeout: defaultTimeoutSec * time.Second,
			},
		},
		unmarshaller: protojson.UnmarshalOptions{
			DiscardUnknown: true,
			AllowPartial:   true,
		},
	}, nil
}

// ================================ GET /games ================================

// GetGamesRequest is the request configuration for the resource
// located at GET /games.
//
//	Year or Game ID is required.
type GetGamesRequest struct {
	// Year is required if GameID is not set.
	Year int32
	// GameID is required if Year is not set.
	GameID int32
	// SeasonType is optional.
	SeasonType string
	// Week is optional.
	Week int32
	// Team is optional.
	Team string
	// Home is optional.
	Home string
	// Away is optional.
	Away string
	// Conference is optional.
	Conference string
	// Classification is optional.
	Classification string
}

// GetGames retrieves a list of games based on the provided request
// parameters.
//
// Calls GET /games.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for games
func (c *Client) GetGames(
	ctx context.Context,
	request GetGamesRequest,
) ([]*Game, error) {
	if request.GameID < 1 && request.Year < 1 {
		return nil, fmt.Errorf(
			"year or ID must be set; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, idKey, request.GameID)
	setInt32(values, yearKey, request.Year)
	setInt32(values, weekKey, request.Week)
	setString(values, teamKey, request.Team)
	setString(values, homeKey, request.Home)
	setString(values, awayKey, request.Away)
	setString(values, seasonTypeKey, request.SeasonType)
	setString(values, conferenceKey, request.Conference)
	setString(values, classificationKey, request.Classification)

	response, err := c.httpGet.Execute(ctx, "/games", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /games; %w", err)
	}

	var games []*Game
	if err = c.unmarshalList(response, &games, &Game{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal games; %w", err)
	}

	return games, nil
}

// ============================= GET /games/teams ==============================

// GetGameTeamsRequest is the request configuration for the resource
// located at GET /games/teams.
//
//	Year is required if GameID is not set.
//	GameID is required if Year is not set.
type GetGameTeamsRequest struct {
	// Year is required if GameID is not set.
	Year int32
	// GameID is required if Year is not set.
	GameID int32
	// Week is optional.
	Week int32
	// SeasonType is optional.
	SeasonType string
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
	// Classification is optional.
	Classification string
}

// GetGameTeams retrieves team box score statistics for games based on
// the provided request parameters.
//
// Calls GET /games/teams.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for game team statistics
func (c *Client) GetGameTeams(
	ctx context.Context,
	request GetGameTeamsRequest,
) ([]*GameTeamStats, error) {
	if request.GameID < 1 && request.Year < 1 {
		return nil, fmt.Errorf(
			"year or ID must be set; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setInt32(values, weekKey, request.Week)
	setString(values, seasonTypeKey, request.SeasonType)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)
	setString(values, classificationKey, request.Classification)
	setInt32(values, gameIDKey, request.GameID)

	response, err := c.httpGet.Execute(ctx, "/games/teams", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /games/teams; %w", err)
	}

	var games []*GameTeamStats
	if err = c.unmarshalList(response, &games, &GameTeamStats{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game team stats; %w", err)
	}

	return games, nil
}

// ============================ GET /games/players ============================

// GetGamePlayersRequest is the request configuration for the resource
// located at GET /games/players.
//
//	Year is required if GameID is not set.
//	GameID is required if Year is not set.
type GetGamePlayersRequest struct {
	// Year is required if GameID is not set.
	Year int32
	// GameID is required if Year is not set.
	GameID int32
	// Week is optional.
	Week int32
	// SeasonType is optional.
	SeasonType string
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
	// Category is optional.
	Category string
}

// GetGamePlayers retrieves player box score statistics for games based
// on the provided request parameters.
//
// Calls GET /games/players.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for game player statistics
func (c *Client) GetGamePlayers(
	ctx context.Context,
	request GetGamePlayersRequest,
) ([]*GamePlayerStats, error) {
	if request.GameID < 1 && request.Year < 1 {
		return nil, fmt.Errorf(
			"year or ID must be set; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setInt32(values, weekKey, request.Week)
	setString(values, seasonTypeKey, request.SeasonType)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)
	setInt32(values, gameIDKey, request.GameID)
	setString(values, categoryKey, request.Category)

	response, err := c.httpGet.Execute(ctx, "/games/players", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /games/players; %w", err)
	}

	var games []*GamePlayerStats
	if err = c.unmarshalList(response, &games, &GamePlayerStats{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game player stats; %w", err)
	}

	return games, nil
}

// ============================= GET /games/media =============================

// GetGameMediaRequest is the request configuration for the resource
// located at GET /games/media.
//
//	Year is required.
type GetGameMediaRequest struct {
	// Year is required.
	Year int32
	// Week is optional.
	Week int32
	// SeasonType is optional.
	SeasonType string
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
	// MediaType is optional.
	MediaType string
}

// GetGameMedia retrieves media information for games based on the provided
// request parameters.
//
// Calls GET /games/media.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for game media
func (c *Client) GetGameMedia(
	ctx context.Context,
	request GetGameMediaRequest,
) ([]*GameMedia, error) {
	if request.Year < 1 {
		return nil, fmt.Errorf("year must be set; %w", ErrMissingRequiredParams)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setInt32(values, weekKey, request.Week)
	setString(values, seasonTypeKey, request.SeasonType)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)
	setString(values, "mediaType", request.MediaType)

	response, err := c.httpGet.Execute(ctx, "/games/media", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /games/media; %w", err)
	}

	var games []*GameMedia
	if err = c.unmarshalList(response, &games, &GameMedia{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game media; %w", err)
	}

	return games, nil
}

// ============================ GET /games/weather ============================

// GetGameWeatherRequest is the request configuration for the resource
// located at GET /games/weather.
//
//	Year is required if GameID is not set.
//	GameID is required if Year is not set.
type GetGameWeatherRequest struct {
	// Year is required if GameID is not set.
	Year int32
	// GameID is required if Year is not set.
	GameID int32
	// Week is optional.
	Week int32
	// SeasonType is optional.
	SeasonType string
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
}

// GetGameWeather retrieves weather information for games based on the
// provided request parameters.
//
// Calls GET /games/weather.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for game weather data
func (c *Client) GetGameWeather(
	ctx context.Context,
	request GetGameWeatherRequest,
) ([]*GameWeather, error) {
	if request.GameID < 1 && request.Year < 1 {
		return nil, fmt.Errorf(
			"year or ID must be set; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setInt32(values, weekKey, request.Week)
	setString(values, seasonTypeKey, request.SeasonType)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)
	setInt32(values, gameIDKey, request.GameID)

	response, err := c.httpGet.Execute(ctx, "/games/weather", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /games/weather; %w", err)
	}

	var games []*GameWeather
	if err = c.unmarshalList(response, &games, &GameWeather{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game weather; %w", err)
	}

	return games, nil
}

// ========================= GET /game/box/advanced ==========================

// GetAdvancedBoxScoreRequest is the request configuration for the resource
// located at GET /game/box/advanced.
type GetAdvancedBoxScoreRequest struct {
	// GameID is required.
	GameID int32
}

// GetAdvancedBoxScore retrieves advanced box score statistics for the
// specified game.
//
// Calls GET /game/box/advanced.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains the game ID for the advanced box score
func (c *Client) GetAdvancedBoxScore(
	ctx context.Context,
	request GetAdvancedBoxScoreRequest,
) (*AdvancedBoxScore, error) {
	if request.GameID < 1 {
		return nil, fmt.Errorf(
			"game ID is required; %w", ErrMissingRequiredParams,
		)
	}

	v := url.Values{}
	setInt32(v, idKey, request.GameID)
	response, err := c.httpGet.Execute(ctx, "/game/box/advanced", v)
	if err != nil {
		return nil, fmt.Errorf("failed to request /game/box/advanced; %w", err)
	}

	var val AdvancedBoxScore
	if err = c.unmarshal(response, &val); err != nil {
		return nil, fmt.Errorf("failed to unmarshal advanced box score; %w", err)
	}

	return &val, nil
}

// ================================ GET /calendar ===============================

// GetCalendarRequest is the request configuration for the resource
// located at GET /calendar.
type GetCalendarRequest struct {
	// Year is required.
	Year int32
}

// GetCalendar retrieves calendar weeks for the specified year.
//
// Calls GET /calendar.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains the calendar year to retrieve weeks for
func (c *Client) GetCalendar(
	ctx context.Context,
	request GetCalendarRequest,
) ([]*CalendarWeek, error) {
	if request.Year < 1 {
		return nil, fmt.Errorf("year is required; %w", ErrMissingRequiredParams)
	}

	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	response, err := c.httpGet.Execute(ctx, "/calendar", v)
	if err != nil {
		return nil, fmt.Errorf("failed to request /calendar; %w", err)
	}

	var weeks []*CalendarWeek
	if err = c.unmarshalList(response, &weeks, &CalendarWeek{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal calendar weeks; %w", err)
	}

	return weeks, nil
}

// ================================ GET /records ===============================

// GetTeamRecordsRequest is the request configuration for the resource
// located at GET /records.
//
//	Year is required if Team is not set.
//	Team is required if Year is not set.
type GetTeamRecordsRequest struct {
	// Year is required if Team is not set.
	Year int32
	// Team is required if Year is not set.
	Team string
	// Conference is optional.
	Conference string
}

// GetTeamRecords retrieves team records based on the provided request
// parameters.
//
// Calls GET /records.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for team records
func (c *Client) GetTeamRecords(
	ctx context.Context,
	request GetTeamRecordsRequest,
) ([]*TeamRecords, error) {
	if request.Year < 1 && strings.TrimSpace(request.Team) == "" {
		return nil, fmt.Errorf(
			"year or team must be set; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)

	response, err := c.httpGet.Execute(ctx, "/records", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /records; %w", err)
	}

	var records []*TeamRecords
	if err = c.unmarshalList(response, &records, &TeamRecords{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team records; %w", err)
	}

	return records, nil
}

// ============================== GET /scoreboard ==============================

// GetScoreboardRequest is the request configuration for the resource
// located at GET /scoreboard.
type GetScoreboardRequest struct {
	// Classification is optional.
	Classification string
	// Conference is optional.
	Conference string
}

// GetScoreboard retrieves live scoreboard data based on the provided
// request parameters.
//
// Calls GET /scoreboard.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for scoreboard games
func (c *Client) GetScoreboard(
	ctx context.Context,
	request GetScoreboardRequest,
) ([]*Scoreboard, error) {
	values := url.Values{}
	setString(values, "classification", request.Classification)
	setString(values, conferenceKey, request.Conference)

	response, err := c.httpGet.Execute(ctx, "/scoreboard", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /scoreboard; %w", err)
	}

	var games []*Scoreboard
	if err = c.unmarshalList(response, &games, &Scoreboard{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scoreboard games; %w", err)
	}

	return games, nil
}

// ================================ GET /drives ================================

// GetDrivesRequest is the request configuration for the resource
// located at GET /drives.
//
//	Year is required.
type GetDrivesRequest struct {
	// Year is required.
	Year int32
	// SeasonType is optional.
	SeasonType string
	// Week is optional.
	Week int32
	// Team is optional.
	Team string
	// Offense is optional.
	Offense string
	// Defense is optional.
	Defense string
	// Conference is optional.
	Conference string
	// OffenseConference is optional.
	OffenseConference string
	// DefenseConference is optional.
	DefenseConference string
	// Classification is optional.
	Classification string
}

// GetDrives retrieves drive information for games based on the provided
// request parameters.
//
// Calls GET /drives.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for drives
func (c *Client) GetDrives(
	ctx context.Context,
	request GetDrivesRequest,
) ([]*Drive, error) {
	if request.Year < 1 {
		return nil, fmt.Errorf("year must be set; %w", ErrMissingRequiredParams)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, seasonTypeKey, request.SeasonType)
	setInt32(values, weekKey, request.Week)
	setString(values, teamKey, request.Team)
	setString(values, offenseKey, request.Offense)
	setString(values, defenseKey, request.Defense)
	setString(values, conferenceKey, request.Conference)
	setString(values, offenseConferenceKey, request.OffenseConference)
	setString(values, defenseConferenceKey, request.DefenseConference)
	setString(values, classificationKey, request.Classification)

	response, err := c.httpGet.Execute(ctx, "/drives", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /drives; %w", err)
	}

	var drives []*Drive
	if err = c.unmarshalList(response, &drives, &Drive{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal drives; %w", err)
	}

	return drives, nil
}

// ================================ GET /plays =================================

// GetPlaysRequest is the request configuration for the resource
// located at GET /plays.
//
//	Year is required.
//	Week is required.
type GetPlaysRequest struct {
	// Year is required.
	Year int32
	// Week is required.
	Week int32
	// Team is optional.
	Team string
	// Offense is optional.
	Offense string
	// Defense is optional.
	Defense string
	// OffenseConference is optional.
	OffenseConference string
	// DefenseConference is optional.
	DefenseConference string
	// Conference is optional.
	Conference string
	// PlayType is optional.
	PlayType string
	// SeasonType is optional.
	SeasonType string
	// Classification is optional.
	Classification string
}

// GetPlays retrieves play-by-play data for games based on the provided
// request parameters.
//
// Calls GET /plays.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for plays
func (c *Client) GetPlays(
	ctx context.Context,
	request GetPlaysRequest,
) ([]*Play, error) {
	if request.Year < 1 {
		return nil, fmt.Errorf("year must be set; %w", ErrMissingRequiredParams)
	}

	if request.Week < 1 {
		return nil, fmt.Errorf("week must be set; %w", ErrMissingRequiredParams)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setInt32(values, weekKey, request.Week)
	setString(values, teamKey, request.Team)
	setString(values, offenseKey, request.Offense)
	setString(values, defenseKey, request.Defense)
	setString(values, offenseConferenceKey, request.OffenseConference)
	setString(values, defenseConferenceKey, request.DefenseConference)
	setString(values, conferenceKey, request.Conference)
	setString(values, "playType", request.PlayType)
	setString(values, seasonTypeKey, request.SeasonType)
	setString(values, classificationKey, request.Classification)

	response, err := c.httpGet.Execute(ctx, "/plays", values)
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
//
// Calls GET /plays/types.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
func (c *Client) GetPlayTypes(ctx context.Context) ([]*PlayType, error) {
	response, err := c.httpGet.Execute(ctx, "/plays/types", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /plays/types; %w", err)
	}

	var playTypes []*PlayType
	if err = c.unmarshalList(response, &playTypes, &PlayType{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal play types; %w", err)
	}

	return playTypes, nil
}

// ============================= GET /plays/stats ============================

// GetPlayStatsRequest is the request configuration for the resource
// located at GET /plays/stats.
type GetPlayStatsRequest struct {
	// Year is optional.
	Year int32
	// Week is optional.
	Week int32
	// Team is optional.
	Team string
	// GameID is optional.
	GameID int32
	// AthleteID is optional.
	AthleteID int32
	// StatTypeID is optional.
	StatTypeID int32
	// SeasonType is optional.
	SeasonType string
	// Conference is optional.
	Conference string
}

// GetPlayStats retrieves play statistics.
//
// Calls GET /plays/stats.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for play statistics
func (c *Client) GetPlayStats(
	ctx context.Context,
	request GetPlayStatsRequest,
) ([]*PlayStat, error) {
	// All fields are optional, no validation needed.
	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setInt32(values, weekKey, request.Week)
	setString(values, teamKey, request.Team)
	setInt32(values, gameIDKey, request.GameID)
	setInt32(values, "athleteId", request.AthleteID)
	setInt32(values, "statTypeId", request.StatTypeID)
	setString(values, seasonTypeKey, request.SeasonType)
	setString(values, conferenceKey, request.Conference)

	response, err := c.httpGet.Execute(ctx, "/plays/stats", values)
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
//
// Calls GET /plays/stats/types.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
func (c *Client) GetPlayStatTypes(
	ctx context.Context,
) ([]*PlayStatType, error) {
	response, err := c.httpGet.Execute(ctx, "/plays/stats/types", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /plays/stats/types; %w", err)
	}

	var statTypes []*PlayStatType
	if err = c.unmarshalList(response, &statTypes, &PlayStatType{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal play stat types; %w", err)
	}

	return statTypes, nil
}

// =============================== GET /live/plays ==============================

// GetLivePlaysRequest is the request configuration for the resource
// located at GET /live/plays.
type GetLivePlaysRequest struct {
	// GameID is required.
	GameID int32
}

// GetLivePlays retrieves live play-by-play data for the specified game.
//
// Calls GET /live/plays.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains the game ID for live play data
func (c *Client) GetLivePlays(
	ctx context.Context,
	request GetLivePlaysRequest,
) (*LiveGame, error) {
	if request.GameID < 1 {
		return nil, fmt.Errorf(
			"game ID is required; %w", ErrMissingRequiredParams,
		)
	}

	params := url.Values{}
	setInt32(params, gameIDKey, request.GameID)

	response, err := c.httpGet.Execute(ctx, "/live/plays", params)
	if err != nil {
		return nil, fmt.Errorf("failed to request /live/plays; %w", err)
	}

	var game LiveGame
	if err = c.unmarshal(response, &game); err != nil {
		return nil, fmt.Errorf("failed to unmarshal live game; %w", err)
	}

	return &game, nil
}

// ================================ GET /teams =================================

// GetTeamsRequest is the request configuration for the resource
// located at GET /teams.
type GetTeamsRequest struct {
	// Conference is optional.
	Conference string
	// Year is optional.
	Year int32
}

// GetTeams retrieves team information based on the provided request
// parameters.
//
// Calls GET /teams.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for teams
func (c *Client) GetTeams(
	ctx context.Context,
	request GetTeamsRequest,
) ([]*Team, error) {
	values := url.Values{}
	setString(values, conferenceKey, request.Conference)
	setInt32(values, yearKey, request.Year)

	response, err := c.httpGet.Execute(ctx, "/teams", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /teams; %w", err)
	}

	var teams []*Team
	if err = c.unmarshalList(response, &teams, &Team{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal teams; %w", err)
	}

	return teams, nil
}

// ============================== GET /teams/fbs ===============================

// GetFBSTeamsRequest is the request configuration for the resource
// located at GET /teams/fbs.
type GetFBSTeamsRequest struct {
	// Year is optional.
	Year int32
}

// GetFBSTeams retrieves FBS (Football Bowl Subdivision) team information
// based on the provided request parameters.
//
// Calls GET /teams/fbs.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for FBS teams
func (c *Client) GetFBSTeams(
	ctx context.Context,
	request GetFBSTeamsRequest,
) ([]*Team, error) {
	values := url.Values{}
	setInt32(values, yearKey, request.Year)

	response, err := c.httpGet.Execute(ctx, "/teams/fbs", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /teams/fbs; %w", err)
	}

	var teams []*Team
	if err = c.unmarshalList(response, &teams, &Team{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal teams; %w", err)
	}

	return teams, nil
}

// ============================ GET /teams/matchup ============================

// GetTeamMatchupRequest is the request configuration for the resource
// located at GET /teams/matchup.
type GetTeamMatchupRequest struct {
	// Team1 is optional.
	Team1 string
	// Team2 is optional.
	Team2 string
	// MinYear is optional.
	MinYear int32
	// MaxYear is optional.
	MaxYear int32
}

// GetTeamMatchup retrieves historical matchup data between two teams based
// on the provided request parameters.
//
// Calls GET /teams/matchup.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains the parameters for the team matchup query
func (c *Client) GetTeamMatchup(
	ctx context.Context,
	request GetTeamMatchupRequest,
) (*Matchup, error) {
	if strings.TrimSpace(request.Team1) == "" ||
		strings.TrimSpace(request.Team2) == "" {
		return nil, fmt.Errorf(
			"team1 and team2 are required; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setString(values, "team1", request.Team1)
	setString(values, "team2", request.Team2)
	setInt32(values, minYearKey, request.MinYear)
	setInt32(values, maxYearKey, request.MaxYear)

	response, err := c.httpGet.Execute(ctx, "/teams/matchup", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /teams/matchup; %w", err)
	}

	var matchup Matchup
	if err = c.unmarshal(response, &matchup); err != nil {
		return nil, fmt.Errorf("failed to unmarshal matchup; %w", err)
	}

	return &matchup, nil
}

// ============================== GET /teams/ats ===============================

// GetTeamATSRequest is the request configuration for the resource
// located at GET /teams/ats.
type GetTeamATSRequest struct {
	// Year is required.
	Year int32
	// Conference is optional.
	Conference string
	// Team is optional.
	Team string
}

// GetTeamATS retrieves team against-the-spread (ATS) records based on the
// provided request parameters.
//
// Calls GET /teams/ats.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for team ATS records
func (c *Client) GetTeamATS(
	ctx context.Context,
	request GetTeamATSRequest,
) ([]*TeamATS, error) {
	if request.Year < 0 {
		return nil, fmt.Errorf("year is required; %w", ErrMissingRequiredParams)
	}
	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, conferenceKey, request.Conference)
	setString(values, teamKey, request.Team)

	response, err := c.httpGet.Execute(ctx, "/teams/ats", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /teams/ats; %w", err)
	}

	var teams []*TeamATS
	if err = c.unmarshalList(response, &teams, &TeamATS{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team ATS; %w", err)
	}

	return teams, nil
}

// ================================ GET /roster ================================

// GetRosterRequest is the request configuration for the resource
// located at GET /roster.
type GetRosterRequest struct {
	// Team is optional.
	Team string
	// Year is optional.
	Year int32
	// Classification is optional.
	Classification string
}

// GetRoster retrieves roster information for a team based on the provided
// request parameters.
//
// Calls GET /roster.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for roster players
func (c *Client) GetRoster(
	ctx context.Context,
	request GetRosterRequest,
) ([]*RosterPlayer, error) {
	values := url.Values{}
	setString(values, teamKey, request.Team)
	setInt32(values, yearKey, request.Year)
	setString(values, classificationKey, request.Classification)

	response, err := c.httpGet.Execute(ctx, "/roster", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /roster; %w", err)
	}

	var players []*RosterPlayer
	if err = c.unmarshalList(response, &players, &RosterPlayer{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal roster players; %w", err)
	}

	return players, nil
}

// ================================ GET /talent ================================

// GetTalentCompositeRequest is the request configuration for the resource
// located at GET /talent.
type GetTalentCompositeRequest struct {
	// Year is required.
	Year int32
}

// GetTeamTalentComposite retrieves 247 team talent composite ratings
//
// Calls GET /talent.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for team talent ratings
func (c *Client) GetTeamTalentComposite(
	ctx context.Context,
	request GetTalentCompositeRequest,
) ([]*TeamTalent, error) {
	if request.Year < 1 {
		return nil, fmt.Errorf("year is required; %w", ErrMissingRequiredParams)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)

	response, err := c.httpGet.Execute(ctx, "/talent", values)
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
//
// Calls GET /conferences.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
func (c *Client) GetConferences(ctx context.Context) ([]*Conference, error) {
	response, err := c.httpGet.Execute(ctx, "/conferences", url.Values{})
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
//
// Calls GET /venues.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
func (c *Client) GetVenues(ctx context.Context) ([]*Venue, error) {
	response, err := c.httpGet.Execute(ctx, "/venues", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /venues; %w", err)
	}

	var venues []*Venue
	if err = c.unmarshalList(response, &venues, &Venue{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal venues; %w", err)
	}

	return venues, nil
}

// =============================== GET /coaches ===============================

// GetCoachesRequest is the request configuration for the resource
// located at GET /coaches.
type GetCoachesRequest struct {
	// FirstName is optional.
	FirstName string
	// LastName is optional.
	LastName string
	// Team is optional.
	Team string
	// Year is optional.
	Year int32
	// MinYear is optional.
	MinYear int32
	// MaxYear is optional.
	MaxYear int32
}

// GetCoaches retrieves coach information based on the provided request
// parameters.
//
// Calls GET /coaches.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for coaches
func (c *Client) GetCoaches(
	ctx context.Context,
	request GetCoachesRequest,
) ([]*Coach, error) {
	values := url.Values{}
	setString(values, firstNameKey, request.FirstName)
	setString(values, lastNameKey, request.LastName)
	setString(values, teamKey, request.Team)
	setInt32(values, yearKey, request.Year)
	setInt32(values, minYearKey, request.MinYear)
	setInt32(values, maxYearKey, request.MaxYear)

	response, err := c.httpGet.Execute(ctx, "/coaches", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /coaches; %w", err)
	}

	var coaches []*Coach
	if err = c.unmarshalList(response, &coaches, &Coach{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal coaches; %w", err)
	}

	return coaches, nil
}

// ============================ GET /player/search ============================

// SearchPlayersRequest is the request configuration for the resource
// located at GET /player/search.
type SearchPlayersRequest struct {
	// SearchTerm is optional.
	SearchTerm string
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// Position is optional.
	Position string
}

// SearchPlayers searches for players based on the provided request
// parameters.
//
// Calls GET /player/search.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains the parameters for the player search
func (c *Client) SearchPlayers(
	ctx context.Context,
	request SearchPlayersRequest,
) ([]*PlayerSearchResult, error) {
	if request.SearchTerm == "" {
		return nil, fmt.Errorf(
			"search term is required; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setString(values, "searchTerm", request.SearchTerm)
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, positionKey, request.Position)

	response, err := c.httpGet.Execute(ctx, "/player/search", values)
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

// ============================= GET /player/usage ============================

// GetPlayerUsageRequest is the request configuration for the resource
// located at GET /player/usage.
type GetPlayerUsageRequest struct {
	// Year is optional.
	Year int32
	// Conference is optional.
	Conference string
	// Position is optional.
	Position string
	// Team is optional.
	Team string
	// PlayerID is optional.
	PlayerID int32
	// ExcludeGarbageTime is optional.
	ExcludeGarbageTime *bool
}

// GetPlayerUsage retrieves player usage statistics based on the provided
// request parameters.
//
// Calls GET /player/usage.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for player usage statistics
func (c *Client) GetPlayerUsage(
	ctx context.Context,
	request GetPlayerUsageRequest,
) ([]*PlayerUsage, error) {
	if request.Year < 1 {
		return nil, fmt.Errorf("year must be set; %w", ErrMissingRequiredParams)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, conferenceKey, request.Conference)
	setString(values, positionKey, request.Position)
	setString(values, teamKey, request.Team)
	setInt32(values, playerIDKey, request.PlayerID)
	setBool(values, excludeGarbageTimeKey, request.ExcludeGarbageTime)

	response, err := c.httpGet.Execute(ctx, "/player/usage", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /player/usage; %w", err)
	}

	var usage []*PlayerUsage
	if err = c.unmarshalList(response, &usage, &PlayerUsage{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player usage; %w", err)
	}

	return usage, nil
}

// =========================== GET /player/returning ==========================

// GetReturningProductionRequest is the request configuration for the resource
// located at GET /player/returning.
type GetReturningProductionRequest struct {
	// Year is required of Team is not set.
	Year int32
	// Team is required if Year is not set.
	Team string
	// Conference is optional.
	Conference string
}

// GetReturningProduction retrieves returning production statistics for
// players based on the provided request parameters.
//
// Calls GET /player/returning.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for returning production data
func (c *Client) GetReturningProduction(
	ctx context.Context,
	request GetReturningProductionRequest,
) ([]*ReturningProduction, error) {
	if request.Year < 1 && strings.TrimSpace(request.Team) == "" {
		return nil, fmt.Errorf(
			"year or team must be set; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)

	response, err := c.httpGet.Execute(ctx, "/player/returning", values)
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

// ============================ GET /player/portal ============================

// GetTransferPortalPlayersRequest is the request configuration for the resource
// located at GET /player/portal.
type GetTransferPortalPlayersRequest struct {
	// Year is required.
	Year int32
}

// GetTransferPortalPlayers retrieves player transfer portal information based on
// the provided request parameters.
//
// Calls GET /player/portal.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for transfer portal data
func (c *Client) GetTransferPortalPlayers(
	ctx context.Context,
	request GetTransferPortalPlayersRequest,
) ([]*PlayerTransfer, error) {
	if request.Year < 1 {
		return nil, fmt.Errorf("year must be set; %w", ErrMissingRequiredParams)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)

	response, err := c.httpGet.Execute(ctx, "/player/portal", values)
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

// ============================== GET /rankings ===============================

// GetRankingsRequest is the request configuration for the resource
// located at GET /rankings.
type GetRankingsRequest struct {
	// Year is required.
	Year int32
	// SeasonType is optional.
	SeasonType string
	// Week is optional.
	Week float64
}

// GetRankings retrieves college football rankings (polls) based on the
// provided request parameters.
//
// Calls GET /rankings.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for rankings
func (c *Client) GetRankings(
	ctx context.Context,
	request GetRankingsRequest,
) ([]*PollWeek, error) {
	if request.Year < 1 {
		return nil, fmt.Errorf("year must be set; %w", ErrMissingRequiredParams)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, seasonTypeKey, request.SeasonType)
	setFloat64(values, weekKey, request.Week)

	response, err := c.httpGet.Execute(ctx, "/rankings", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /rankings; %w", err)
	}

	var rankings []*PollWeek
	if err = c.unmarshalList(response, &rankings, &PollWeek{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rankings; %w", err)
	}

	return rankings, nil
}

// ================================ GET /lines =================================

// GetBettingLinesRequest is the request configuration for the resource
// located at GET /lines.
type GetBettingLinesRequest struct {
	// GameID is optional.
	GameID int32
	// Year is optional.
	Year int32
	// SeasonType is optional.
	SeasonType string
	// Week is optional.
	Week int32
	// Team is optional.
	Team string
	// Home is optional.
	Home string
	// Away is optional.
	Away string
	// Conference is optional.
	Conference string
	// Provider is optional.
	Provider string
}

// GetBettingLines retrieves betting lines for games based on the provided
// request parameters.
//
// Calls GET /lines.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for betting lines
func (c *Client) GetBettingLines(
	ctx context.Context,
	request GetBettingLinesRequest,
) ([]*BettingGame, error) {
	if request.Year < 1 && request.GameID < 1 {
		return nil, fmt.Errorf(
			"year or game ID must be set; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, gameIDKey, request.GameID)
	setInt32(values, yearKey, request.Year)
	setString(values, seasonTypeKey, request.SeasonType)
	setInt32(values, weekKey, request.Week)
	setString(values, teamKey, request.Team)
	setString(values, homeKey, request.Home)
	setString(values, awayKey, request.Away)
	setString(values, conferenceKey, request.Conference)
	setString(values, "provider", request.Provider)

	response, err := c.httpGet.Execute(ctx, "/lines", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /lines; %w", err)
	}

	var games []*BettingGame
	if err = c.unmarshalList(response, &games, &BettingGame{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal betting games; %w", err)
	}

	return games, nil
}

// =========================== GET /recruiting/players =========================

// GetPlayersRecruitingRankingsRequest is the request configuration for the resource
// located at GET /recruiting/players.
type GetPlayersRecruitingRankingsRequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// Position is optional.
	Position string
	// State is optional.
	State string
	// Classification is optional.
	Classification string
}

// GetPlayerRecruitingRankings retrieves recruiting information for players based
// on the provided request parameters.
//
// Calls GET /recruiting/players.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for recruiting players
func (c *Client) GetPlayerRecruitingRankings(
	ctx context.Context,
	request GetPlayersRecruitingRankingsRequest,
) ([]*Recruit, error) {
	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, positionKey, request.Position)
	setString(values, "state", request.State)
	setString(values, classificationKey, request.Classification)

	response, err := c.httpGet.Execute(ctx, "/recruiting/players", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /recruiting/players; %w", err)
	}

	var recruits []*Recruit
	if err = c.unmarshalList(response, &recruits, &Recruit{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recruits; %w", err)
	}

	return recruits, nil
}

// =========================== GET /recruiting/teams ===========================

// GetTeamRecruitingRankingsRequest is the request configuration for the
// resource located at GET /recruiting/teams.
type GetTeamRecruitingRankingsRequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
}

// GetTeamRecruitingRankings retrieves team recruiting rankings based on the
// provided request parameters.
//
// Calls GET /recruiting/teams.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for team recruiting rankings
func (c *Client) GetTeamRecruitingRankings(
	ctx context.Context,
	request GetTeamRecruitingRankingsRequest,
) ([]*TeamRecruitingRanking, error) {
	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)

	response, err := c.httpGet.Execute(ctx, "/recruiting/teams", values)
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

// =========================== GET /recruiting/groups ==========================

// GetTeamPositionGroupRecruitingRankingsRequest is the request configuration for the resource
// located at GET /recruiting/groups.
type GetTeamPositionGroupRecruitingRankingsRequest struct {
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
	// RecruitType is optional.
	RecruitType string
	// StartYear is optional.
	StartYear int32
	// EndYear is optional.
	EndYear int32
}

// GetTeamPositionGroupRecruitingRankings retrieves aggregated team recruiting information
// based on the provided request parameters.
//
// Calls GET /recruiting/groups.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for recruiting groups
func (c *Client) GetTeamPositionGroupRecruitingRankings(
	ctx context.Context,
	request GetTeamPositionGroupRecruitingRankingsRequest,
) ([]*AggregatedTeamRecruiting, error) {
	values := url.Values{}
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)
	setString(values, "recruitType", request.RecruitType)
	setInt32(values, startYearKey, request.StartYear)
	setInt32(values, endYearKey, request.EndYear)

	response, err := c.httpGet.Execute(ctx, "/recruiting/groups", values)
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

// ============================== GET /ratings/sp ==============================

// GetSPPlusRatingsRequest is the request configuration for the resource
// located at GET /ratings/sp.
type GetSPPlusRatingsRequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
}

// GetTeamSPPlusRatings retrieves SP+ (S&P+) ratings for teams based on the
// provided request parameters.
//
// Calls GET /ratings/sp.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for SP+ ratings
func (c *Client) GetTeamSPPlusRatings(
	ctx context.Context,
	request GetSPPlusRatingsRequest,
) ([]*TeamSP, error) {
	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setString(v, teamKey, request.Team)

	response, err := c.httpGet.Execute(ctx, "/ratings/sp", v)
	if err != nil {
		return nil, fmt.Errorf("failed to request /ratings/sp; %w", err)
	}

	var ratings []*TeamSP
	if err = c.unmarshalList(response, &ratings, &TeamSP{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team SP ratings; %w", err)
	}

	return ratings, nil
}

// ======================= GET /ratings/sp/conferences ========================

// GetConferenceSPPlusRatingsRequest is the request configuration for the
// resource located at GET /ratings/sp/conferences.
type GetConferenceSPPlusRatingsRequest struct {
	// Year is optional.
	Year int32
	// Conference is optional.
	Conference string
}

// GetConferenceSPPlusRatings retrieves SP+ (S&P+) ratings for conferences
// based on the provided request parameters.
//
// Calls GET /ratings/sp/conferences.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for conference SP+ ratings
func (c *Client) GetConferenceSPPlusRatings(
	ctx context.Context,
	request GetConferenceSPPlusRatingsRequest,
) ([]*ConferenceSP, error) {
	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setString(v, conferenceKey, request.Conference)

	response, err := c.httpGet.Execute(
		ctx, "/ratings/sp/conferences", v,
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

// ============================= GET /ratings/srs ==============================

// GetSRSRatingsRequest is the request configuration for the resource
// located at GET /ratings/srs.
type GetSRSRatingsRequest struct {
	// Year is required if Team is not specified.
	Year int32
	// Team is required if Year is not specified.
	Team string
	// Conference is optional.
	Conference string
}

// GetSRSRatings retrieves SRS (Simple Rating System) ratings for teams
// based on the provided request parameters.
//
// Calls GET /ratings/srs.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for SRS ratings
func (c *Client) GetSRSRatings(
	ctx context.Context,
	request GetSRSRatingsRequest,
) ([]*TeamSRS, error) {
	if request.Year < 1 && strings.TrimSpace(request.Team) == "" {
		return nil, fmt.Errorf(
			"year or team required; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)

	response, err := c.httpGet.Execute(ctx, "/ratings/srs", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /ratings/srs; %w", err)
	}

	var ratings []*TeamSRS
	if err = c.unmarshalList(response, &ratings, &TeamSRS{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team SRS ratings; %w", err)
	}

	return ratings, nil
}

// ============================= GET /ratings/elo ==============================

// GetEloRatingsRequest is the request configuration for the resource
// located at GET /ratings/elo.
type GetEloRatingsRequest struct {
	// Year is optional.
	Year int32
	// Week is optional.
	Week int32
	// SeasonType is optional.
	SeasonType string
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
}

// GetEloRatings retrieves Elo ratings for teams based on the provided
// request parameters.
//
// Calls GET /ratings/elo.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for Elo ratings
func (c *Client) GetEloRatings(
	ctx context.Context,
	request GetEloRatingsRequest,
) ([]*TeamElo, error) {
	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setInt32(values, weekKey, request.Week)
	setString(values, seasonTypeKey, request.SeasonType)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)

	response, err := c.httpGet.Execute(ctx, "/ratings/elo", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /ratings/elo; %w", err)
	}

	var ratings []*TeamElo
	if err = c.unmarshalList(response, &ratings, &TeamElo{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team Elo ratings; %w", err)
	}

	return ratings, nil
}

// ============================= GET /ratings/fpi ==============================

// GetFPIRatingsRequest is the request configuration for the resource
// located at GET /ratings/fpi.
type GetFPIRatingsRequest struct {
	// Year is required if Team is not specified.
	Year int32
	// Team is required if Year is not specified.
	Team string
	// Conference is optional.
	Conference string
}

// GetFPIRatings retrieves FPI (Football Power Index) ratings for teams
// based on the provided request parameters.
//
// Calls GET /ratings/fpi.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for FPI ratings
func (c *Client) GetFPIRatings(
	ctx context.Context,
	request GetFPIRatingsRequest,
) ([]*TeamFPI, error) {
	if request.Year < 1 && strings.TrimSpace(request.Team) == "" {
		return nil, fmt.Errorf(
			"year or team required; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)

	response, err := c.httpGet.Execute(ctx, "/ratings/fpi", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /ratings/fpi; %w", err)
	}

	var ratings []*TeamFPI
	if err = c.unmarshalList(response, &ratings, &TeamFPI{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team FPI ratings; %w", err)
	}

	return ratings, nil
}

// =========================== GET /ppa/predicted =============================

// GetPredictedPointsRequest is the request configuration for the resource
// located at GET /ppa/predicted.
type GetPredictedPointsRequest struct {
	// Down is required.
	Down int32
	// Distance is required.
	Distance int32
}

// GetPredictedPoints retrieves predicted points values based on the
// provided request parameters.
//
// Calls GET /ppa/predicted.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for predicted points
func (c *Client) GetPredictedPoints(
	ctx context.Context,
	request GetPredictedPointsRequest,
) ([]*PredictedPointsValue, error) {
	if request.Distance < 0 {
		return nil, fmt.Errorf(
			"distance is required; %w", ErrMissingRequiredParams,
		)
	}

	if request.Down < 1 {
		return nil, fmt.Errorf(
			"down is required; %w", ErrMissingRequiredParams,
		)
	}

	v := url.Values{}
	setInt32(v, "down", request.Down)
	setInt32(v, "distance", request.Distance)

	response, err := c.httpGet.Execute(ctx, "/ppa/predicted", v)
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

// ============================== GET /ppa/teams ===============================

// GetTeamsPPARequest is the request configuration for the resource
// located at GET /ppa/teams.
type GetTeamsPPARequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
	// ExcludeGarbageTime is optional.
	ExcludeGarbageTime *bool
}

// GetTeamsPPA retrieves team season PPA (Predicted Points Added) statistics
// based on the provided request parameters.
//
// Calls GET /ppa/teams.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for team season PPA
func (c *Client) GetTeamsPPA(
	ctx context.Context,
	request GetTeamsPPARequest,
) ([]*TeamSeasonPredictedPointsAdded, error) {
	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setString(v, teamKey, request.Team)
	setString(v, conferenceKey, request.Conference)
	setBool(v, excludeGarbageTimeKey, request.ExcludeGarbageTime)

	response, err := c.httpGet.Execute(ctx, "/ppa/teams", v)
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

// ============================= GET /ppa/games ================================

// GetPpaGamesRequest is the request configuration for the resource
// located at GET /ppa/games.
type GetPpaGamesRequest struct {
	// Year is optional.
	Year int32
	// Week is optional.
	Week int32
	// SeasonType is optional.
	SeasonType string
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
	// ExcludeGarbageTime is optional.
	ExcludeGarbageTime *bool
}

// GetGamesPPA retrieves team game PPA (Predicted Points Added) statistics
// based on the provided request parameters.
//
// Calls GET /ppa/games.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for team game PPA
func (c *Client) GetGamesPPA(
	ctx context.Context,
	request GetPpaGamesRequest,
) ([]*TeamGamePredictedPointsAdded, error) {
	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setInt32(v, weekKey, request.Week)
	setString(v, seasonTypeKey, request.SeasonType)
	setString(v, teamKey, request.Team)
	setString(v, conferenceKey, request.Conference)
	setBool(v, excludeGarbageTimeKey, request.ExcludeGarbageTime)

	response, err := c.httpGet.Execute(ctx, "/ppa/games", v)
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

// ========================== GET /ppa/players/games ===========================

// GetPlayerPpaGamesRequest is the request configuration for the resource
// located at GET /ppa/players/games.
type GetPlayerPpaGamesRequest struct {
	// Year is optional.
	Year int32
	// Week is optional.
	Week int32
	// SeasonType is optional.
	SeasonType string
	// Team is optional.
	Team string
	// Position is optional.
	Position string
	// PlayerID is optional.
	PlayerID string
	// Threshold is optional.
	Threshold float64
	// ExcludeGarbageTime is optional.
	ExcludeGarbageTime *bool
}

// GetPlayersPPA retrieves player game PPA (Predicted Points Added)
// statistics based on the provided request parameters.
//
// Calls GET /ppa/players/games.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for player game PPA
func (c *Client) GetPlayersPPA(
	ctx context.Context,
	request GetPlayerPpaGamesRequest,
) ([]*PlayerGamePredictedPointsAdded, error) {
	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setInt32(v, weekKey, request.Week)
	setString(v, seasonTypeKey, request.SeasonType)
	setString(v, teamKey, request.Team)
	setString(v, positionKey, request.Position)
	setString(v, playerIDKey, request.PlayerID)
	setFloat64(v, thresholdKey, request.Threshold)
	setBool(v, excludeGarbageTimeKey, request.ExcludeGarbageTime)

	response, err := c.httpGet.Execute(
		ctx, "/ppa/players/games", v,
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

// ========================= GET /ppa/players/season ==========================

// GetPlayerSeasonPPARequest is the request configuration for the resource
// located at GET /ppa/players/season.
type GetPlayerSeasonPPARequest struct {
	// Year is optional.
	Year int32
	// Conference is optional.
	Conference string
	// Team is optional.
	Team string
	// Position is optional.
	Position string
	// PlayerID is optional.
	PlayerID string
	// Threshold is optional.
	Threshold float64
	// ExcludeGarbageTime is optional.
	ExcludeGarbageTime *bool
}

// GetPlayerSeasonPPA retrieves player season PPA (Predicted Points Added)
// statistics based on the provided request parameters.
//
// Calls GET /ppa/players/season.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for player season PPA
func (c *Client) GetPlayerSeasonPPA(
	ctx context.Context,
	request GetPlayerSeasonPPARequest,
) ([]*PlayerSeasonPredictedPointsAdded, error) {
	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setString(v, conferenceKey, request.Conference)
	setString(v, teamKey, request.Team)
	setString(v, positionKey, request.Position)
	setString(v, playerIDKey, request.PlayerID)
	setFloat64(v, thresholdKey, request.Threshold)
	setBool(v, excludeGarbageTimeKey, request.ExcludeGarbageTime)

	response, err := c.httpGet.Execute(
		ctx, "/ppa/players/season", v,
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

// =============================== GET /metrics/wp ==============================

// GetWinProbabilityRequest is the request configuration for the resource
// located at GET /metrics/wp.
type GetWinProbabilityRequest struct {
	// GameID is required.
	GameID int32
}

// GetWinProbability retrieves win probability data for each play in the
// specified game.
//
// Calls GET /metrics/wp.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains the game ID for win probability data
func (c *Client) GetWinProbability(
	ctx context.Context,
	request GetWinProbabilityRequest,
) ([]*PlayWinProbability, error) {
	if request.GameID < 1 {
		return nil, fmt.Errorf(
			"game ID is required; %w", ErrMissingRequiredParams,
		)
	}

	params := url.Values{}
	setInt32(params, gameIDKey, request.GameID)

	response, err := c.httpGet.Execute(ctx, "/metrics/wp", params)
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

// ========================= GET /metrics/wp/pregame ===========================

// GetPregameWpRequest is the request configuration for the resource
// located at GET /metrics/wp/pregame.
type GetPregameWpRequest struct {
	// Year is optional.
	Year int32
	// Week is optional.
	Week int32
	// SeasonType is optional.
	SeasonType string
	// Team is optional.
	Team string
}

// GetPregameWinProbability retrieves pregame win probability data based on
// the provided request parameters.
//
// Calls GET /metrics/wp/pregame.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for pregame win probabilities
func (c *Client) GetPregameWinProbability(
	ctx context.Context,
	request GetPregameWpRequest,
) ([]*PregameWinProbability, error) {
	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setInt32(v, weekKey, request.Week)
	setString(v, teamKey, request.Team)
	setString(v, seasonTypeKey, request.SeasonType)

	response, err := c.httpGet.Execute(
		ctx, "/metrics/wp/pregame", v,
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

// GetFieldGoalExpectedPoints retrieves field goal expected points values.
//
// Calls GET /metrics/fg/ep.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
func (c *Client) GetFieldGoalExpectedPoints(
	ctx context.Context,
) ([]*FieldGoalEP, error) {
	response, err := c.httpGet.Execute(ctx, "/metrics/fg/ep", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /metrics/fg/ep; %w", err)
	}

	var ep []*FieldGoalEP
	if err = c.unmarshalList(response, &ep, &FieldGoalEP{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal field goal EP; %w", err)
	}

	return ep, nil
}

// ========================== GET /stats/player/season =========================

// GetPlayerSeasonStatsRequest is the request configuration for the resource
// located at GET /stats/player/season.
type GetPlayerSeasonStatsRequest struct {
	// Year is optional.
	Year int32
	// Conference is optional.
	Conference string
	// Team is optional.
	Team string
	// StartWeek is optional.
	StartWeek int32
	// EndWeek is optional.
	EndWeek int32
	// SeasonType is optional.
	SeasonType string
	// Category is optional.
	Category string
}

// GetPlayerSeasonStats retrieves player season statistics based on the
// provided request parameters.
//
// Calls GET /stats/player/season.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for player season statistics
func (c *Client) GetPlayerSeasonStats(
	ctx context.Context,
	request GetPlayerSeasonStatsRequest,
) ([]*PlayerStat, error) {
	if request.Year < 1 {
		return nil, fmt.Errorf("year is required; %w", ErrMissingRequiredParams)
	}

	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setString(v, conferenceKey, request.Conference)
	setString(v, teamKey, request.Team)
	setInt32(v, startWeekKey, request.StartWeek)
	setInt32(v, endWeekKey, request.EndWeek)
	setString(v, seasonTypeKey, request.SeasonType)
	setString(v, categoryKey, request.Category)

	response, err := c.httpGet.Execute(
		ctx, "/stats/player/season", v,
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

// ============================== GET /stats/season ============================

// GetTeamSeasonStatsRequest is the request configuration for the resource
// located at GET /stats/season.
type GetTeamSeasonStatsRequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
	// StartWeek is optional.
	StartWeek int32
	// EndWeek is optional.
	EndWeek int32
}

// GetTeamSeasonStats retrieves team season statistics based on the provided
// request parameters.
//
// Calls GET /stats/season.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for team season statistics
func (c *Client) GetTeamSeasonStats(
	ctx context.Context,
	request GetTeamSeasonStatsRequest,
) ([]*TeamStat, error) {
	if request.Year < 1 && request.Team == "" {
		return nil, fmt.Errorf(
			"year or team is required; %w", ErrMissingRequiredParams,
		)
	}

	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)
	setInt32(values, startWeekKey, request.StartWeek)
	setInt32(values, endWeekKey, request.EndWeek)

	response, err := c.httpGet.Execute(ctx, "/stats/season", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /stats/season; %w", err)
	}

	var stats []*TeamStat
	if err = c.unmarshalList(response, &stats, &TeamStat{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team season stats; %w", err)
	}

	return stats, nil
}

// GetStatCategories retrieves all available statistics categories.
//
// Calls GET /stats/categories.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
func (c *Client) GetStatCategories(ctx context.Context) ([]string, error) {
	response, err := c.httpGet.Execute(ctx, "/stats/categories", url.Values{})
	if err != nil {
		return nil, fmt.Errorf("failed to request /stats/categories; %w", err)
	}

	var out []string
	if err := json.Unmarshal(response, &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stats categories; %w", err)
	}

	return out, nil
}

// ========================= GET /stats/season/advanced =======================

// GetAdvancedSeasonStatsRequest is the request configuration for the resource
// located at GET /stats/season/advanced.
type GetAdvancedSeasonStatsRequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// ExcludeGarbageTime is optional.
	ExcludeGarbageTime *bool
	// StartWeek is optional.
	StartWeek int32
	// EndWeek is optional.
	EndWeek int32
}

// GetAdvancedSeasonStats retrieves advanced season statistics based on the
// provided request parameters.
//
// Calls GET /stats/season/advanced.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for advanced season statistics
func (c *Client) GetAdvancedSeasonStats(
	ctx context.Context,
	request GetAdvancedSeasonStatsRequest,
) ([]*AdvancedSeasonStat, error) {
	if request.Year < 1 && request.Team == "" {
		return nil, fmt.Errorf(
			"year or team is required; %w", ErrMissingRequiredParams,
		)
	}

	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setString(v, teamKey, request.Team)
	setBool(v, excludeGarbageTimeKey, request.ExcludeGarbageTime)
	setInt32(v, startWeekKey, request.StartWeek)
	setInt32(v, endWeekKey, request.EndWeek)

	response, err := c.httpGet.Execute(
		ctx, "/stats/season/advanced", v,
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

// ========================= GET /stats/game/advanced ==========================

// GetAdvancedGameStatsRequest is the request configuration for the resource
// located at GET /stats/game/advanced.
type GetAdvancedGameStatsRequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// Week is optional.
	Week float64
	// Opponent is optional.
	Opponent string
	// ExcludeGarbageTime is optional.
	ExcludeGarbageTime *bool
	// SeasonType is optional.
	SeasonType string
}

// GetAdvancedGameStats retrieves advanced game statistics based on the
// provided request parameters.
//
// Calls GET /stats/game/advanced.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
//	req  contains filtering options for advanced game statistics
func (c *Client) GetAdvancedGameStats(
	ctx context.Context,
	req GetAdvancedGameStatsRequest,
) ([]*AdvancedGameStat, error) {
	if req.Year < 1 && req.Team == "" {
		return nil, fmt.Errorf(
			"year or team is required; %w", ErrMissingRequiredParams,
		)
	}

	v := url.Values{}
	setInt32(v, yearKey, req.Year)
	setString(v, teamKey, req.Team)
	setFloat64(v, weekKey, req.Week)
	setString(v, opponentKey, req.Opponent)
	setBool(v, excludeGarbageTimeKey, req.ExcludeGarbageTime)
	setString(v, seasonTypeKey, req.SeasonType)

	resp, err := c.httpGet.Execute(ctx, "/stats/game/advanced", v)
	if err != nil {
		return nil, fmt.Errorf("failed to request /stats/game/advanced; %w", err)
	}

	var stats []*AdvancedGameStat
	if err = c.unmarshalList(resp, &stats, &AdvancedGameStat{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal advanced game stats; %w", err)
	}

	return stats, nil
}

// ========================== GET /stats/game/havoc ============================

// GetHavocGameStatsRequest is the request configuration for the resource
// located at GET /stats/game/havoc.
type GetHavocGameStatsRequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// Week is optional.
	Week float64
	// Opponent is optional.
	Opponent string
	// SeasonType is optional.
	SeasonType string
}

// GetHavocGameStats retrieves havoc game statistics based on the provided
// request parameters.
//
// Calls GET /stats/game/havoc.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for game havoc statistics
func (c *Client) GetHavocGameStats(
	ctx context.Context,
	request GetHavocGameStatsRequest,
) ([]*GameHavocStats, error) {
	if request.Year < 1 && request.Team == "" {
		return nil, fmt.Errorf(
			"year or team is required; %w", ErrMissingRequiredParams,
		)
	}

	v := url.Values{}
	setInt32(v, yearKey, request.Year)
	setString(v, teamKey, request.Team)
	setFloat64(v, weekKey, request.Week)
	setString(v, opponentKey, request.Opponent)
	setString(v, seasonTypeKey, request.SeasonType)

	response, err := c.httpGet.Execute(
		ctx, "/stats/game/havoc", v,
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
//
// Calls GET /draft/teams.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
func (c *Client) GetDraftTeams(ctx context.Context) ([]*DraftTeam, error) {
	response, err := c.httpGet.Execute(ctx, "/draft/teams", url.Values{})
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
//
// Calls GET /draft/positions.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
func (c *Client) GetDraftPositions(
	ctx context.Context,
) ([]*DraftPosition, error) {
	response, err := c.httpGet.Execute(ctx, "/draft/positions", url.Values{})
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

// ============================== GET /draft/picks =============================

// GetDraftPicksRequest is the request configuration for the resource
// located at GET /draft/picks.
type GetDraftPicksRequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// School is optional.
	School string
	// Conference is optional.
	Conference string
	// Position is optional.
	Position string
}

// GetDraftPicks retrieves NFL draft picks based on the provided request
// parameters.
//
// Calls GET /draft/picks.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for draft picks
func (c *Client) GetDraftPicks(
	ctx context.Context,
	request GetDraftPicksRequest,
) ([]*DraftPick, error) {
	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, "school", request.School)
	setString(values, conferenceKey, request.Conference)
	setString(values, positionKey, request.Position)
	response, err := c.httpGet.Execute(ctx, "/draft/picks", values)
	if err != nil {
		return nil, fmt.Errorf("failed to request /draft/picks; %w", err)
	}

	var picks []*DraftPick
	if err = c.unmarshalList(response, &picks, &DraftPick{}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal draft picks; %w", err)
	}

	return picks, nil
}

// =========================== GET /wepa/team/season ==========================

// GetTeamSeasonWEPARequest is the request configuration for the resource
// located at GET /wepa/team/season.
type GetTeamSeasonWEPARequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
}

// GetTeamSeasonWEPA retrieves team season WEPA (Weighted Expected Points
// Added) metrics based on the provided request parameters.
//
// Calls GET /wepa/team/season.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for team season WEPA metrics
func (c *Client) GetTeamSeasonWEPA(
	ctx context.Context,
	request GetTeamSeasonWEPARequest,
) ([]*AdjustedTeamMetrics, error) {
	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)

	resp, err := c.httpGet.Execute(ctx, "/wepa/team/season", values)
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

// ======================== GET /wepa/players/passing ========================

// GetPlayerWEPARequest is the request configuration for the resource
// located at GET /wepa/players/passing.
type GetPlayerWEPARequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
	// Position is optional.
	Position string
}

// GetPlayerPassingWEPA retrieves player passing WEPA (Weighted Expected
// Points Added) metrics based on the provided request parameters.
//
// Calls GET /wepa/players/passing.
//
// The behavior depends on the provided parameters:
//
//	ctx      controls request cancellation
//	request  contains filtering options for player passing WEPA metrics
func (c *Client) GetPlayerPassingWEPA(
	ctx context.Context,
	request GetPlayerWEPARequest,
) ([]*PlayerWeightedEPA, error) {
	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)
	setString(values, positionKey, request.Position)

	resp, err := c.httpGet.Execute(
		ctx, "/wepa/players/passing", values,
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

// GetPlayerRushingWEPA retrieves player rushing WEPA (Weighted Expected
// Points Added) metrics based on the provided request parameters.
//
// Calls GET /wepa/players/rushing.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
//	req  contains filtering options for player rushing WEPA metrics
func (c *Client) GetPlayerRushingWEPA(
	ctx context.Context,
	request GetPlayerWEPARequest,
) ([]*PlayerWeightedEPA, error) {
	values := url.Values{}
	setInt32(values, yearKey, request.Year)
	setString(values, teamKey, request.Team)
	setString(values, conferenceKey, request.Conference)
	setString(values, positionKey, request.Position)

	resp, err := c.httpGet.Execute(ctx, "/wepa/players/rushing", values)
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

// ======================== GET /wepa/players/kicking =========================

// GetWepaPlayersKickingRequest is the request configuration for the resource
// located at GET /wepa/players/kicking.
type GetWepaPlayersKickingRequest struct {
	// Year is optional.
	Year int32
	// Team is optional.
	Team string
	// Conference is optional.
	Conference string
}

// GetPlayerKickingWEPA retrieves kicker PAAR (Points Above Average
// Replacement) metrics based on the provided request parameters.
//
// Calls GET /wepa/players/kicking.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
//	req  contains filtering options for kicker PAAR metrics
func (c *Client) GetPlayerKickingWEPA(
	ctx context.Context,
	req GetWepaPlayersKickingRequest,
) ([]*KickerPAAR, error) {
	v := url.Values{}
	setInt32(v, yearKey, req.Year)
	setString(v, teamKey, req.Team)
	setString(v, conferenceKey, req.Conference)

	response, err := c.httpGet.Execute(
		ctx, "/wepa/players/kicking", v,
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
//
// Calls GET /info.
//
// The behavior depends on the provided parameters:
//
//	ctx  controls request cancellation
func (c *Client) GetInfo(ctx context.Context) (*UserInfo, error) {
	response, err := c.httpGet.Execute(ctx, "/info", url.Values{})
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
		return fmt.Errorf("prototype cannot be nil (e.g. &Drive{})")
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

func isJSONNull(b []byte) bool {
	return bytes.Equal(bytes.TrimSpace(b), []byte("null"))
}

func setString(v url.Values, key string, val string) {
	if strings.TrimSpace(val) == "" {
		return
	}

	// Explicitly URL encode the string value before setting it as a query parameter.
	// This ensures proper encoding of special characters in all string parameters.
	trimmed := strings.TrimSpace(val)
	encoded := url.QueryEscape(trimmed)
	v.Set(key, encoded)
}

func setInt32(v url.Values, key string, val int32) {
	if val == 0 {
		return
	}

	v.Set(key, strconv.FormatInt(int64(val), 10))
}

func setFloat64(v url.Values, key string, val float64) {
	if val == float64(0) {
		return
	}

	v.Set(key, strconv.FormatFloat(val, 'f', -1, 64))
}

// setBool sets a boolean value in url.Values if it's not nil.
// The key parameter is kept for consistency with other set functions.
//
//nolint:unparam // key is kept for API consistency even though currently only one key is used
func setBool(v url.Values, key string, val *bool) {
	if val == nil {
		return
	}

	v.Set(key, strconv.FormatBool(*val))
}
