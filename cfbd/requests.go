package cfbd

import (
	"net/url"
	"strconv"
)

// GetGamesRequest matches the typical /games filters.
type GetGamesRequest struct {
	Year       int32
	SeasonType *string
	Week       *int32
	Team       *string
	Home       *string
	Away       *string
	Conference *string
	Division   *string
	GameID     *int32
}

func (p GetGamesRequest) validate() error {
	return nil
}

func (p GetGamesRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setString(v, "seasonType", p.SeasonType)
	setInt32(v, "week", p.Week)
	setString(v, "team", p.Team)
	setString(v, "home", p.Home)
	setString(v, "away", p.Away)
	setString(v, "conference", p.Conference)
	setString(v, "division", p.Division)
	setInt32(v, "id", p.GameID)
	return v
}

type GameTeamStatsRequest struct {
	Year int32 // required

	Week       *int32
	SeasonType *string
	Team       *string
	Conference *string
	GameID     *int32
}

func (p GameTeamStatsRequest) validate() error {
	return nil
}

func (p GameTeamStatsRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setInt32(v, "week", p.Week)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	setInt32(v, "gameId", p.GameID)
	return v
}

type GamePlayerStatsRequest struct {
	Year int32 // required

	Week       *int32
	SeasonType *string
	Team       *string
	Conference *string
	GameID     *int32
	Category   *string
}

func (p GamePlayerStatsRequest) validate() error {
	return nil
}

func (p GamePlayerStatsRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setInt32(v, "week", p.Week)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	setInt32(v, "gameId", p.GameID)
	setString(v, "category", p.Category)
	return v
}

type GameMediaRequest struct {
	Year int32 // required

	Week       *int32
	SeasonType *string
	Team       *string
	Conference *string
	MediaType  *string // tv, radio, web, ppv, mobile
}

func (p GameMediaRequest) validate() error {
	return nil
}

func (p GameMediaRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setInt32(v, "week", p.Week)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	setString(v, "mediaType", p.MediaType)
	return v
}

type GameWeatherRequest struct {
	Year int32 // required

	Week       *int32
	SeasonType *string
	Team       *string
	Conference *string
	GameID     *int32
}

func (p GameWeatherRequest) validate() error {
	return nil
}

func (p GameWeatherRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setInt32(v, "week", p.Week)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	setInt32(v, "gameId", p.GameID)
	return v
}

type RecordsRequest struct {
	Year int32 // required

	Team       *string
	Conference *string
}

func (p RecordsRequest) validate() error {
	return nil
}

func (p RecordsRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	return v
}

type LiveScoreboardRequest struct {
	Division   *string // fbs,fcs,ii,iii
	Conference *string
}

func (p LiveScoreboardRequest) validate() error {
	return nil
}

func (p LiveScoreboardRequest) values() url.Values {
	v := url.Values{}
	setString(v, "division", p.Division)
	setString(v, "conference", p.Conference)
	return v
}

type DrivesRequest struct {
	Year int32 // required

	SeasonType        *string
	Week              *int32
	Team              *string
	Offense           *string
	Defense           *string
	Conference        *string
	OffenseConference *string
	DefenseConference *string
	Classification    *string // fbs,fcs,ii,iii
}

func (p DrivesRequest) validate() error {
	return nil
}

func (p DrivesRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setString(v, "seasonType", p.SeasonType)
	setInt32(v, "week", p.Week)
	setString(v, "team", p.Team)
	setString(v, "offense", p.Offense)
	setString(v, "defense", p.Defense)
	setString(v, "conference", p.Conference)
	setString(v, "offenseConference", p.OffenseConference)
	setString(v, "defenseConference", p.DefenseConference)
	setString(v, "classification", p.Classification)
	return v
}

type PlaysRequest struct {
	Year int32 // required
	Week int32 // required

	Team              *string
	Offense           *string
	Defense           *string
	OffenseConference *string
	DefenseConference *string
	Conference        *string
	PlayType          *string
	SeasonType        *string
	Classification    *string
}

func (p PlaysRequest) validate() error {
	return nil
}

func (p PlaysRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	v.Set("week", strconv.FormatInt(int64(p.Week), 10))
	setString(v, "team", p.Team)
	setString(v, "offense", p.Offense)
	setString(v, "defense", p.Defense)
	setString(v, "offenseConference", p.OffenseConference)
	setString(v, "defenseConference", p.DefenseConference)
	setString(v, "conference", p.Conference)
	setString(v, "playType", p.PlayType)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "classification", p.Classification)
	return v
}

type PlayStatsRequest struct {
	Year       *int32
	Week       *int32
	Team       *string
	GameID     *int32
	AthleteID  *int32
	StatTypeID *int32
	SeasonType *string
	Conference *string
}

func (p PlayStatsRequest) validate() error {
	return nil
}

func (p PlayStatsRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setInt32(v, "week", p.Week)
	setString(v, "team", p.Team)
	setInt32(v, "gameId", p.GameID)
	setInt32(v, "athleteId", p.AthleteID)
	setInt32(v, "statTypeId", p.StatTypeID)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "conference", p.Conference)
	return v
}

type TeamsRequest struct {
	Conference *string
	Year       *int32
}

func (p TeamsRequest) validate() error {
	return nil
}

func (p TeamsRequest) values() url.Values {
	v := url.Values{}
	setString(v, "conference", p.Conference)
	setInt32(v, "year", p.Year)
	return v
}

type TeamsFbsRequest struct {
	Year *int32
}

func (p TeamsFbsRequest) validate() error {
	return nil
}

func (p TeamsFbsRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	return v
}

type TeamMatchupRequest struct {
	Team1   string
	Team2   string
	MinYear *int32
	MaxYear *int32
}

func (p TeamMatchupRequest) validate() error {
	return nil
}

func (p TeamMatchupRequest) values() url.Values {
	v := url.Values{}
	v.Set("team1", p.Team1)
	v.Set("team2", p.Team2)
	setInt32(v, "minYear", p.MinYear)
	setInt32(v, "maxYear", p.MaxYear)
	return v
}

type TeamATSRequest struct {
	Year       int32
	Conference *string
	Team       *string
}

func (p TeamATSRequest) validate() error {
	return nil
}

func (p TeamATSRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setString(v, "conference", p.Conference)
	setString(v, "team", p.Team)
	return v
}

type RosterRequest struct {
	Team           *string
	Year           *int32
	Classification *string
}

func (p RosterRequest) validate() error {
	return nil
}

func (p RosterRequest) values() url.Values {
	v := url.Values{}
	setString(v, "team", p.Team)
	setInt32(v, "year", p.Year)
	setString(v, "classification", p.Classification)
	return v
}

type TalentRequest struct {
	Year int32
}

func (p TalentRequest) validate() error {
	return nil
}

func (p TalentRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	return v
}

type CoachesRequest struct {
	FirstName *string
	LastName  *string
	Team      *string
	Year      *int32
	MinYear   *int32
	MaxYear   *int32
}

func (p CoachesRequest) validate() error {
	return nil
}

func (p CoachesRequest) values() url.Values {
	v := url.Values{}
	setString(v, "firstName", p.FirstName)
	setString(v, "lastName", p.LastName)
	setString(v, "team", p.Team)
	setInt32(v, "year", p.Year)
	setInt32(v, "minYear", p.MinYear)
	setInt32(v, "maxYear", p.MaxYear)
	return v
}

type PlayerSearchRequest struct {
	SearchTerm string
	Year       *int32
	Team       *string
	Position   *string
}

func (p PlayerSearchRequest) validate() error {
	return nil
}

func (p PlayerSearchRequest) values() url.Values {
	v := url.Values{}
	v.Set("searchTerm", p.SearchTerm)
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "position", p.Position)
	return v
}

type PlayerUsageRequest struct {
	Year               int32
	Conference         *string
	Position           *string
	Team               *string
	PlayerID           *int32
	ExcludeGarbageTime *bool
}

func (p PlayerUsageRequest) validate() error {
	return nil
}

func (p PlayerUsageRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setString(v, "conference", p.Conference)
	setString(v, "position", p.Position)
	setString(v, "team", p.Team)
	setInt32(v, "playerId", p.PlayerID)
	setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
	return v
}

type ReturningProductionRequest struct {
	Year       *int32
	Team       *string
	Conference *string
}

func (p ReturningProductionRequest) validate() error {
	return nil
}

func (p ReturningProductionRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	return v
}

type PlayerPortalRequest struct {
	Year int32
}

func (p PlayerPortalRequest) validate() error {
	return nil
}

func (p PlayerPortalRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	return v
}

type RankingsRequest struct {
	Year       int32
	SeasonType *string
	Week       *float64
}

func (p RankingsRequest) validate() error {
	return nil
}

func (p RankingsRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setString(v, "seasonType", p.SeasonType)
	setFloat64(v, "week", p.Week)
	return v
}

type LinesRequest struct {
	GameID     *int32
	Year       *int32
	SeasonType *string
	Week       *int32
	Team       *string
	Home       *string
	Away       *string
	Conference *string
	Provider   *string
}

func (p LinesRequest) validate() error {
	return nil
}

func (p LinesRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "gameId", p.GameID)
	setInt32(v, "year", p.Year)
	setString(v, "seasonType", p.SeasonType)
	setInt32(v, "week", p.Week)
	setString(v, "team", p.Team)
	setString(v, "home", p.Home)
	setString(v, "away", p.Away)
	setString(v, "conference", p.Conference)
	setString(v, "provider", p.Provider)
	return v
}

type RecruitingPlayersRequest struct {
	Year           *int32
	Team           *string
	Position       *string
	State          *string
	Classification *string
}

func (p RecruitingPlayersRequest) validate() error {
	return nil
}

func (p RecruitingPlayersRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "position", p.Position)
	setString(v, "state", p.State)
	setString(v, "classification", p.Classification)
	return v
}

type RecruitingTeamsRequest struct {
	Year *int32
	Team *string
}

func (p RecruitingTeamsRequest) validate() error {
	return nil
}

func (p RecruitingTeamsRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	return v
}

type RecruitingGroupsRequest struct {
	Team        *string
	Conference  *string
	RecruitType *string
	StartYear   *int32
	EndYear     *int32
}

func (p RecruitingGroupsRequest) validate() error {
	return nil
}

func (p RecruitingGroupsRequest) values() url.Values {
	v := url.Values{}
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	setString(v, "recruitType", p.RecruitType)
	setInt32(v, "startYear", p.StartYear)
	setInt32(v, "endYear", p.EndYear)
	return v
}

type RatingsSpRequest struct {
	Year *int32
	Team *string
}

func (p RatingsSpRequest) validate() error {
	return nil
}

func (p RatingsSpRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	return v
}

type RatingsSpConferencesRequest struct {
	Year       *int32
	Conference *string
}

func (p RatingsSpConferencesRequest) validate() error {
	return nil
}

func (p RatingsSpConferencesRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "conference", p.Conference)
	return v
}

type RatingsSrsRequest struct {
	Year       *int32
	Team       *string
	Conference *string
}

func (p RatingsSrsRequest) validate() error {
	return nil
}

func (p RatingsSrsRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	return v
}

type RatingsEloRequest struct {
	Year       *int32
	Week       *int32
	SeasonType *string
	Team       *string
	Conference *string
}

func (p RatingsEloRequest) validate() error {
	return nil
}

func (p RatingsEloRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setInt32(v, "week", p.Week)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	return v
}

type RatingsFpiRequest struct {
	Year       *int32
	Team       *string
	Conference *string
}

func (p RatingsFpiRequest) validate() error {
	return nil
}

func (p RatingsFpiRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	return v
}

// PredictedPointsRequest points values by down and distance.
type PredictedPointsRequest struct {
	Down     int32
	Distance int32
}

func (p PredictedPointsRequest) validate() error {
	return nil
}

func (p PredictedPointsRequest) values() url.Values {
	v := url.Values{}
	v.Set("down", strconv.FormatInt(int64(p.Down), 10))
	v.Set("distance", strconv.FormatInt(int64(p.Distance), 10))
	return v
}

// PpaTeamsRequest season PPA (predicted points added) metrics.
type PpaTeamsRequest struct {
	Year               *int32
	Team               *string
	Conference         *string
	ExcludeGarbageTime *bool
}

func (p PpaTeamsRequest) validate() error {
	return nil
}

func (p PpaTeamsRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
	return v
}

// PpaGamesRequest PPA by game.
type PpaGamesRequest struct {
	Year               int32
	Week               *int32
	SeasonType         *string
	Team               *string
	Conference         *string
	ExcludeGarbageTime *bool
}

func (p PpaGamesRequest) validate() error {
	return nil
}

func (p PpaGamesRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setInt32(v, "week", p.Week)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
	return v
}

// PlayerPpaGamesRequest todo:describe.
type PlayerPpaGamesRequest struct {
	Year               int32
	Week               *int32
	SeasonType         *string
	Team               *string
	Position           *string
	PlayerID           *string
	Threshold          *float64
	ExcludeGarbageTime *bool
}

func (p PlayerPpaGamesRequest) validate() error {
	return nil
}

func (p PlayerPpaGamesRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setInt32(v, "week", p.Week)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "team", p.Team)
	setString(v, "position", p.Position)
	setString(v, "playerId", p.PlayerID)
	setFloat64(v, "threshold", p.Threshold)
	setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
	return v
}

// PlayerPpaSeasonRequest todo:describe.
type PlayerPpaSeasonRequest struct {
	Year               *int32
	Conference         *string
	Team               *string
	Position           *string
	PlayerID           *string
	Threshold          *float64
	ExcludeGarbageTime *bool
}

func (p PlayerPpaSeasonRequest) validate() error {
	return nil
}

func (p PlayerPpaSeasonRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "conference", p.Conference)
	setString(v, "team", p.Team)
	setString(v, "position", p.Position)
	setString(v, "playerId", p.PlayerID)
	setFloat64(v, "threshold", p.Threshold)
	setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
	return v
}

// PregameWpRequest todo:describe.
type PregameWpRequest struct {
	Year       *int32
	Week       *int32
	SeasonType *string
	Team       *string
}

func (p PregameWpRequest) validate() error {
	return nil
}

func (p PregameWpRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setInt32(v, "week", p.Week)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "team", p.Team)
	return v
}

type PlayerSeasonStatsRequest struct {
	Year       int32
	Conference *string
	Team       *string
	StartWeek  *int32
	EndWeek    *int32
	SeasonType *string
	Category   *string
}

func (p PlayerSeasonStatsRequest) validate() error {
	return nil
}

func (p PlayerSeasonStatsRequest) values() url.Values {
	v := url.Values{}
	v.Set("year", strconv.FormatInt(int64(p.Year), 10))
	setString(v, "conference", p.Conference)
	setString(v, "team", p.Team)
	setInt32(v, "startWeek", p.StartWeek)
	setInt32(v, "endWeek", p.EndWeek)
	setString(v, "seasonType", p.SeasonType)
	setString(v, "category", p.Category)
	return v
}

type TeamSeasonStatsRequest struct {
	Year       *int32
	Team       *string
	Conference *string
	StartWeek  *int32
	EndWeek    *int32
}

func (p TeamSeasonStatsRequest) validate() error {
	return nil
}

func (p TeamSeasonStatsRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	setInt32(v, "startWeek", p.StartWeek)
	setInt32(v, "endWeek", p.EndWeek)
	return v
}

type AdvancedSeasonStatsRequest struct {
	Year               *int32
	Team               *string
	ExcludeGarbageTime *bool
	StartWeek          *int32
	EndWeek            *int32
}

func (p AdvancedSeasonStatsRequest) validate() error {
	return nil
}

func (p AdvancedSeasonStatsRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
	setInt32(v, "startWeek", p.StartWeek)
	setInt32(v, "endWeek", p.EndWeek)
	return v
}

type AdvancedGameStatsRequest struct {
	Year               *int32
	Team               *string
	Week               *float64
	Opponent           *string
	ExcludeGarbageTime *bool
	SeasonType         *string
}

func (p AdvancedGameStatsRequest) validate() error {
	return nil
}

func (p AdvancedGameStatsRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setFloat64(v, "week", p.Week)
	setString(v, "opponent", p.Opponent)
	setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
	setString(v, "seasonType", p.SeasonType)
	return v
}

type HavocGameStatsRequest struct {
	Year       *int32
	Team       *string
	Week       *float64
	Opponent   *string
	SeasonType *string
}

func (p HavocGameStatsRequest) validate() error {
	return nil
}

func (p HavocGameStatsRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setFloat64(v, "week", p.Week)
	setString(v, "opponent", p.Opponent)
	setString(v, "seasonType", p.SeasonType)
	return v
}

type DraftPicksRequest struct {
	Year       *int32
	Team       *string
	School     *string
	Conference *string
	Position   *string
}

func (p DraftPicksRequest) validate() error {
	return nil
}

func (p DraftPicksRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "school", p.School)
	setString(v, "conference", p.Conference)
	setString(v, "position", p.Position)
	return v
}

type WepaTeamSeasonRequest struct {
	Year       *int32
	Team       *string
	Conference *string
}

func (p WepaTeamSeasonRequest) validate() error {
	return nil
}

func (p WepaTeamSeasonRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	return v
}

type WepaPlayersRequest struct {
	Year       *int32
	Team       *string
	Conference *string
	Position   *string
}

func (p WepaPlayersRequest) validate() error {
	return nil
}

func (p WepaPlayersRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	setString(v, "position", p.Position)
	return v
}

type WepaKickersRequest struct {
	Year       *int32
	Team       *string
	Conference *string
}

func (p WepaKickersRequest) validate() error {
	return nil
}

func (p WepaKickersRequest) values() url.Values {
	v := url.Values{}
	setInt32(v, "year", p.Year)
	setString(v, "team", p.Team)
	setString(v, "conference", p.Conference)
	return v
}
