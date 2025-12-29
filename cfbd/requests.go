package cfbd

import (
   "errors"
   "fmt"
   "net/url"
   "strconv"
   "strings"
)

var ErrMissingRequiredParams = errors.New("request missing required params")

// Request for /games

// GetGamesRequest matches the typical /games filters.
type GetGamesRequest struct {
   Year       int32
   SeasonType string
   Week       int32
   Team       string
   Home       string
   Away       string
   Conference string
   Division   string
   GameID     int32
}

func (p GetGamesRequest) validate() error {
   if p.GameID > 0 {
      return nil
   }

   if p.Year < 1 {
      return fmt.Errorf("year or ID must be set; %w", ErrMissingRequiredParams)
   }

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

// Request for /games/teams

type GetGameTeamsRequest struct {
   Year           int32
   Week           int32
   SeasonType     string
   Team           string
   Conference     string
   Classification string
   GameID         int32
}

func (p GetGameTeamsRequest) validate() error {
   if p.GameID > 0 {
      return nil
   }

   if p.Year < 1 {
      return fmt.Errorf("year or ID must be set; %w", ErrMissingRequiredParams)
   }

   return nil
}

func (p GetGameTeamsRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(p.Year), 10))
   setInt32(v, "week", p.Week)
   setString(v, "seasonType", p.SeasonType)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   setString(v, "classification", p.Classification)
   setInt32(v, "gameId", p.GameID)
   return v
}

// Request for /games/players

// GetGamePlayersRequest todo:describe.
type GetGamePlayersRequest struct {
   // Year is a required field is GameID is not set.
   Year int32
   // Week is an optional field.
   Week int32
   // SeasonType todo:describe.
   SeasonType string
   // Team todo:describe.
   Team string
   // Conference todo:describe.
   Conference string
   // GameID is a required field if Year is not set.
   GameID int32
   // Category todo:describe.
   Category string
}

func (p GetGamePlayersRequest) validate() error {
   if p.GameID > 0 {
      return nil
   }

   if p.Year < 1 {
      return fmt.Errorf("year or ID must be set; %w", ErrMissingRequiredParams)
   }

   return nil
}

func (p GetGamePlayersRequest) values() url.Values {
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

// Request for /games/media

type GetGameMediaRequest struct {
   Year       int32
   Week       int32
   SeasonType string
   Team       string
   Conference string
   MediaType  string
}

func (p GetGameMediaRequest) validate() error {
   if p.Year < 1 {
      return fmt.Errorf("year must be set; %w", ErrMissingRequiredParams)
   }

   return nil
}

func (p GetGameMediaRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(p.Year), 10))
   setInt32(v, "week", p.Week)
   setString(v, "seasonType", p.SeasonType)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   setString(v, "mediaType", p.MediaType)
   return v
}

// Request for /games/weather

type GetGameWeatherRequest struct {
   Year       int32
   Week       int32
   SeasonType string
   Team       string
   Conference string
   GameID     int32
}

func (req GetGameWeatherRequest) validate() error {
   if req.GameID > 0 {
      return nil
   }

   if req.Year < 1 {
      return fmt.Errorf("year or ID must be set; %w", ErrMissingRequiredParams)
   }

   return nil
}

func (req GetGameWeatherRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(req.Year), 10))
   setInt32(v, "week", req.Week)
   setString(v, "seasonType", req.SeasonType)
   setString(v, "team", req.Team)
   setString(v, "conference", req.Conference)
   setInt32(v, "gameId", req.GameID)
   return v
}

// Request for /games/records

// GetRecordsRequest todo:describe.
type GetRecordsRequest struct {
   Year       int32
   Team       string
   Conference string
}

func (p GetRecordsRequest) validate() error {
   if p.Year > 0 {
      return nil
   }

   if strings.TrimSpace(p.Team) == "" {
      return fmt.Errorf(
         "year or team must be set; %w", ErrMissingRequiredParams,
      )
   }

   return nil
}

func (p GetRecordsRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(p.Year), 10))
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   return v
}

// GetScoreboardRequest todo:describe.
type GetScoreboardRequest struct {
   Division   string // fbs,fcs,ii,iii
   Conference string
}

func (p GetScoreboardRequest) validate() error {
   return nil
}

func (p GetScoreboardRequest) values() url.Values {
   v := url.Values{}
   setString(v, "division", p.Division)
   setString(v, "conference", p.Conference)
   return v
}

// GetDrivesRequest todo:describe.
type GetDrivesRequest struct {
   Year int32 // required

   SeasonType        string
   Week              int32
   Team              string
   Offense           string
   Defense           string
   Conference        string
   OffenseConference string
   DefenseConference string
   Classification    string // fbs,fcs,ii,iii
}

func (p GetDrivesRequest) validate() error {
   return nil
}

func (p GetDrivesRequest) values() url.Values {
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

// GetPlaysRequest todo:describe.
type GetPlaysRequest struct {
   Year int32 // required
   Week int32 // required

   Team              string
   Offense           string
   Defense           string
   OffenseConference string
   DefenseConference string
   Conference        string
   PlayType          string
   SeasonType        string
   Classification    string
}

func (p GetPlaysRequest) validate() error {
   return nil
}

func (p GetPlaysRequest) values() url.Values {
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

// GetPlayStatsRequest todo:describe.
type GetPlayStatsRequest struct {
   Year       int32
   Week       int32
   Team       string
   GameID     int32
   AthleteID  int32
   StatTypeID int32
   SeasonType string
   Conference string
}

func (p GetPlayStatsRequest) validate() error {
   return nil
}

func (p GetPlayStatsRequest) values() url.Values {
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

// GetTeamsRequest todo:describe.
type GetTeamsRequest struct {
   Conference string
   Year       int32
}

func (p GetTeamsRequest) validate() error {
   return nil
}

func (p GetTeamsRequest) values() url.Values {
   v := url.Values{}
   setString(v, "conference", p.Conference)
   setInt32(v, "year", p.Year)
   return v
}

// GetTeamsFbsRequest todo:describe.
type GetTeamsFbsRequest struct {
   Year int32
}

func (p GetTeamsFbsRequest) validate() error {
   return nil
}

func (p GetTeamsFbsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   return v
}

// GetTeamMatchupRequest todo:describe.
type GetTeamMatchupRequest struct {
   Team1   string
   Team2   string
   MinYear int32
   MaxYear int32
}

func (p GetTeamMatchupRequest) validate() error {
   return nil
}

func (p GetTeamMatchupRequest) values() url.Values {
   v := url.Values{}
   v.Set("team1", p.Team1)
   v.Set("team2", p.Team2)
   setInt32(v, "minYear", p.MinYear)
   setInt32(v, "maxYear", p.MaxYear)
   return v
}

// GetTeamATSRequest todo:describe.
type GetTeamATSRequest struct {
   Year       int32
   Conference string
   Team       string
}

func (p GetTeamATSRequest) validate() error {
   return nil
}

func (p GetTeamATSRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(p.Year), 10))
   setString(v, "conference", p.Conference)
   setString(v, "team", p.Team)
   return v
}

// GetRosterRequest todo:describe.
type GetRosterRequest struct {
   Team           string
   Year           int32
   Classification string
}

func (p GetRosterRequest) validate() error {
   return nil
}

func (p GetRosterRequest) values() url.Values {
   v := url.Values{}
   setString(v, "team", p.Team)
   setInt32(v, "year", p.Year)
   setString(v, "classification", p.Classification)
   return v
}

// GetTalentRequest todo:describe.
type GetTalentRequest struct {
   Year int32
}

func (p GetTalentRequest) validate() error {
   return nil
}

func (p GetTalentRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(p.Year), 10))
   return v
}

// GetCoachesRequest todo:describe.
type GetCoachesRequest struct {
   FirstName string
   LastName  string
   Team      string
   Year      int32
   MinYear   int32
   MaxYear   int32
}

func (p GetCoachesRequest) validate() error {
   return nil
}

func (p GetCoachesRequest) values() url.Values {
   v := url.Values{}
   setString(v, "firstName", p.FirstName)
   setString(v, "lastName", p.LastName)
   setString(v, "team", p.Team)
   setInt32(v, "year", p.Year)
   setInt32(v, "minYear", p.MinYear)
   setInt32(v, "maxYear", p.MaxYear)
   return v
}

// GetPlayerSearchRequest todo:describe.
type GetPlayerSearchRequest struct {
   SearchTerm string
   Year       int32
   Team       string
   Position   string
}

func (p GetPlayerSearchRequest) validate() error {
   return nil
}

func (p GetPlayerSearchRequest) values() url.Values {
   v := url.Values{}
   v.Set("searchTerm", p.SearchTerm)
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "position", p.Position)
   return v
}

// GetPlayerUsageRequest todo:describe.
type GetPlayerUsageRequest struct {
   Year               int32
   Conference         string
   Position           string
   Team               string
   PlayerID           int32
   ExcludeGarbageTime *bool
}

func (p GetPlayerUsageRequest) validate() error {
   return nil
}

func (p GetPlayerUsageRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(p.Year), 10))
   setString(v, "conference", p.Conference)
   setString(v, "position", p.Position)
   setString(v, "team", p.Team)
   setInt32(v, "playerId", p.PlayerID)
   setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
   return v
}

// GetReturningProductionRequest todo:describe.
type GetReturningProductionRequest struct {
   Year       int32
   Team       string
   Conference string
}

func (p GetReturningProductionRequest) validate() error {
   return nil
}

func (p GetReturningProductionRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   return v
}

// GetPlayerPortalRequest todo:describe.
type GetPlayerPortalRequest struct {
   Year int32
}

func (p GetPlayerPortalRequest) validate() error {
   return nil
}

func (p GetPlayerPortalRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(p.Year), 10))
   return v
}

// GetRankingsRequest todo:describe.
type GetRankingsRequest struct {
   Year       int32
   SeasonType string
   Week       float64
}

func (p GetRankingsRequest) validate() error {
   return nil
}

func (p GetRankingsRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(p.Year), 10))
   setString(v, "seasonType", p.SeasonType)
   setFloat64(v, "week", p.Week)
   return v
}

// GetBettingLinesRequest todo:describe.
type GetBettingLinesRequest struct {
   GameID     int32
   Year       int32
   SeasonType string
   Week       int32
   Team       string
   Home       string
   Away       string
   Conference string
   Provider   string
}

func (p GetBettingLinesRequest) validate() error {
   return nil
}

func (p GetBettingLinesRequest) values() url.Values {
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

// GetRecruitingPlayersRequest todo: describe.
type GetRecruitingPlayersRequest struct {
   Year           int32
   Team           string
   Position       string
   State          string
   Classification string
}

func (p GetRecruitingPlayersRequest) validate() error {
   return nil
}

func (p GetRecruitingPlayersRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "position", p.Position)
   setString(v, "state", p.State)
   setString(v, "classification", p.Classification)
   return v
}

// GetTeamRecruitingRankingsRequest todo: describe.
type GetTeamRecruitingRankingsRequest struct {
   Year int32
   Team string
}

func (p GetTeamRecruitingRankingsRequest) validate() error {
   return nil
}

func (p GetTeamRecruitingRankingsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   return v
}

type GetRecruitingGroupsRequest struct {
   Team        string
   Conference  string
   RecruitType string
   StartYear   int32
   EndYear     int32
}

func (p GetRecruitingGroupsRequest) validate() error {
   return nil
}

func (p GetRecruitingGroupsRequest) values() url.Values {
   v := url.Values{}
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   setString(v, "recruitType", p.RecruitType)
   setInt32(v, "startYear", p.StartYear)
   setInt32(v, "endYear", p.EndYear)
   return v
}

type GetSPPlusRatingsRequest struct {
   Year int32
   Team string
}

func (p GetSPPlusRatingsRequest) validate() error {
   return nil
}

func (p GetSPPlusRatingsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   return v
}

type GetConferenceSPPlusRatingsRequest struct {
   Year       int32
   Conference string
}

func (p GetConferenceSPPlusRatingsRequest) validate() error {
   return nil
}

func (p GetConferenceSPPlusRatingsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "conference", p.Conference)
   return v
}

type GetSRSRatingsRequest struct {
   Year       int32
   Team       string
   Conference string
}

func (p GetSRSRatingsRequest) validate() error {
   return nil
}

func (p GetSRSRatingsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   return v
}

type GetEloRatingsRequest struct {
   Year       int32
   Week       int32
   SeasonType string
   Team       string
   Conference string
}

func (p GetEloRatingsRequest) validate() error {
   return nil
}

func (p GetEloRatingsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setInt32(v, "week", p.Week)
   setString(v, "seasonType", p.SeasonType)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   return v
}

type GetFPIRatingsRequest struct {
   Year       int32
   Team       string
   Conference string
}

func (p GetFPIRatingsRequest) validate() error {
   return nil
}

func (p GetFPIRatingsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   return v
}

// GetPredictedPointsRequest points values by down and distance.
type GetPredictedPointsRequest struct {
   Down     int32
   Distance int32
}

func (p GetPredictedPointsRequest) validate() error {
   return nil
}

func (p GetPredictedPointsRequest) values() url.Values {
   v := url.Values{}
   v.Set("down", strconv.FormatInt(int64(p.Down), 10))
   v.Set("distance", strconv.FormatInt(int64(p.Distance), 10))
   return v
}

// GetTeamsPPARequest season PPA (predicted points added) metrics.
type GetTeamsPPARequest struct {
   Year               int32
   Team               string
   Conference         string
   ExcludeGarbageTime *bool
}

func (p GetTeamsPPARequest) validate() error {
   return nil
}

func (p GetTeamsPPARequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
   return v
}

// GetPpaGamesRequest PPA by game.
type GetPpaGamesRequest struct {
   Year               int32
   Week               int32
   SeasonType         string
   Team               string
   Conference         string
   ExcludeGarbageTime *bool
}

func (p GetPpaGamesRequest) validate() error {
   return nil
}

func (p GetPpaGamesRequest) values() url.Values {
   v := url.Values{}
   v.Set("year", strconv.FormatInt(int64(p.Year), 10))
   setInt32(v, "week", p.Week)
   setString(v, "seasonType", p.SeasonType)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
   return v
}

// GetPlayerPpaGamesRequest todo:describe.
type GetPlayerPpaGamesRequest struct {
   Year               int32
   Week               int32
   SeasonType         string
   Team               string
   Position           string
   PlayerID           string
   Threshold          float64
   ExcludeGarbageTime *bool
}

func (p GetPlayerPpaGamesRequest) validate() error {
   return nil
}

func (p GetPlayerPpaGamesRequest) values() url.Values {
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

// GetPlayerSeasonPPARequest todo:describe.
type GetPlayerSeasonPPARequest struct {
   Year               int32
   Conference         string
   Team               string
   Position           string
   PlayerID           string
   Threshold          float64
   ExcludeGarbageTime *bool
}

func (p GetPlayerSeasonPPARequest) validate() error {
   return nil
}

func (p GetPlayerSeasonPPARequest) values() url.Values {
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

// GetPregameWpRequest todo:describe.
type GetPregameWpRequest struct {
   Year       int32
   Week       int32
   SeasonType string
   Team       string
}

func (p GetPregameWpRequest) validate() error {
   return nil
}

func (p GetPregameWpRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setInt32(v, "week", p.Week)
   setString(v, "seasonType", p.SeasonType)
   setString(v, "team", p.Team)
   return v
}

// GetPlayerSeasonStatsRequest is the request configuration for the resource
// located at /stats/player/season.
type GetPlayerSeasonStatsRequest struct {
   Year       int32
   Conference string
   Team       string
   StartWeek  int32
   EndWeek    int32
   SeasonType string
   Category   string
}

func (p GetPlayerSeasonStatsRequest) validate() error {
   return nil
}

func (p GetPlayerSeasonStatsRequest) values() url.Values {
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

// GetTeamSeasonStatsRequest is the request configuration for the resource
// located at /stats/season.
type GetTeamSeasonStatsRequest struct {
   Year       int32
   Team       string
   Conference string
   StartWeek  int32
   EndWeek    int32
}

func (p GetTeamSeasonStatsRequest) validate() error {
   return nil
}

func (p GetTeamSeasonStatsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   setInt32(v, "startWeek", p.StartWeek)
   setInt32(v, "endWeek", p.EndWeek)
   return v
}

// GetAdvancedSeasonStatsRequest is the request configuration for the resource
// located at /stats/season/advanced.
type GetAdvancedSeasonStatsRequest struct {
   Year               int32
   Team               string
   ExcludeGarbageTime *bool
   StartWeek          int32
   EndWeek            int32
}

func (p GetAdvancedSeasonStatsRequest) validate() error {
   return nil
}

func (p GetAdvancedSeasonStatsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
   setInt32(v, "startWeek", p.StartWeek)
   setInt32(v, "endWeek", p.EndWeek)
   return v
}

// GetAdvancedGameStatsRequest is the request configuration for the resource
// located at /stats/game/advanced.
type GetAdvancedGameStatsRequest struct {
   Year               int32
   Team               string
   Week               float64
   Opponent           string
   ExcludeGarbageTime *bool
   SeasonType         string
}

func (p GetAdvancedGameStatsRequest) validate() error {
   return nil
}

func (p GetAdvancedGameStatsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setFloat64(v, "week", p.Week)
   setString(v, "opponent", p.Opponent)
   setBool(v, "excludeGarbageTime", p.ExcludeGarbageTime)
   setString(v, "seasonType", p.SeasonType)
   return v
}

// GetGameHavocStatsRequest is the request configuration for the resource
// located at /stats/game/havoc.
type GetGameHavocStatsRequest struct {
   Year       int32
   Team       string
   Week       float64
   Opponent   string
   SeasonType string
}

func (p GetGameHavocStatsRequest) validate() error {
   return nil
}

func (p GetGameHavocStatsRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setFloat64(v, "week", p.Week)
   setString(v, "opponent", p.Opponent)
   setString(v, "seasonType", p.SeasonType)
   return v
}

// GetDraftPicksRequest is the request configuration for the resource
// located at /draft/picks.
type GetDraftPicksRequest struct {
   Year       int32
   Team       string
   School     string
   Conference string
   Position   string
}

func (p GetDraftPicksRequest) validate() error {
   return nil
}

func (p GetDraftPicksRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "school", p.School)
   setString(v, "conference", p.Conference)
   setString(v, "position", p.Position)
   return v
}

// GetTeamSeasonWEPARequest is the request configuration for the resource
// located at /wepa/team/season.
type GetTeamSeasonWEPARequest struct {
   Year       int32
   Team       string
   Conference string
}

func (p GetTeamSeasonWEPARequest) validate() error {
   return nil
}

func (p GetTeamSeasonWEPARequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   return v
}

// GetWepaPlayersPassingRequest is the request configuration for the resource
// located at /wepa/players/passing.
type GetWepaPlayersPassingRequest struct {
   Year       int32
   Team       string
   Conference string
   Position   string
}

func (p GetWepaPlayersPassingRequest) validate() error {
   return nil
}

func (p GetWepaPlayersPassingRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   setString(v, "position", p.Position)
   return v
}

// GetWepaPlayersKickingRequest is the request configuration for the resource
// located at /wepa/players/kicking.
type GetWepaPlayersKickingRequest struct {
   Year       int32
   Team       string
   Conference string
}

func (p GetWepaPlayersKickingRequest) validate() error {
   return nil
}

func (p GetWepaPlayersKickingRequest) values() url.Values {
   v := url.Values{}
   setInt32(v, "year", p.Year)
   setString(v, "team", p.Team)
   setString(v, "conference", p.Conference)
   return v
}

func setString(v url.Values, key string, val string) {
   if strings.TrimSpace(val) == "" {
      return
   }

   v.Set(key, strings.TrimSpace(val))
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

func setBool(v url.Values, key string, val *bool) {
   if val == nil {
      return
   }

   v.Set(key, strconv.FormatBool(*val))
}
