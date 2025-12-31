package cfbd

import (
	"context"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	defaultTimeFormat      = "2006-01-02T15:04:05.000Z"
	testResponsePathPrefix = "./internal/test/responses/"
	testYear               = 2025
	testWeek               = 2
	testTeam               = "Texas"
)

type testClient struct {
	client          *Client
	requestExecutor *mockHttpGetExecutor
}

func newTestClient(t *testing.T) *testClient {
	ctrl := gomock.NewController(t)
	exec := newMockHttpGetExecutor(ctrl)

	return &testClient{
		client: &Client{
			apiKey: "",
			unmarshaller: protojson.UnmarshalOptions{
				AllowPartial:   true,
				DiscardUnknown: true,
			},
			httpGet: exec,
		},
		requestExecutor: exec,
	}
}

// setupTestWithFile creates a test client and reads the test response file.
// The filename should be relative to the test response directory
// (e.g., "games.json").
func setupTestWithFile(t *testing.T, filename string) (*testClient, []byte) {
	tester := newTestClient(t)
	fullPath := testResponsePathPrefix + filename
	bytes, err := os.ReadFile(fullPath)
	require.NoError(t, err, "failed to read test file: %s", fullPath)
	return tester, bytes
}

func TestGetGames_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "games.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetGames(
		context.Background(), GetGamesRequest{
			Week: 1, Year: testYear, Team: testTeam, SeasonType: "regular",
		},
	)

	require.NoError(t, err)
	require.NotEmpty(t, response)
	game := response[0]
	assert.Equal(t, game.Id, int32(401752677))
	assert.Equal(t, game.Season, int32(2025))
	assert.Equal(t, game.SeasonType, "regular")
	assert.Equal(t, game.StartTime_TBD, false)
	assert.Equal(t, game.Completed, true)
	assert.Equal(t, game.NeutralSite, false)
	assert.Equal(t, game.ConferenceGame, false)
	assert.Equal(t, game.Week, int32(1))
	assert.Equal(t, game.VenueId.Value, int32(3861))
	assert.Equal(t, game.Venue.Value, "Ohio Stadium")
	assert.Equal(t, game.HomeId.Value, int32(194))
	assert.Equal(t, game.HomeTeam, "Ohio State")
	assert.Equal(t, game.HomeClassification.Value, "fbs")
	assert.Equal(t, game.HomeConference.Value, "Big Ten")
	assert.Equal(t, convertToInt32Slice(game.HomeLineScores.Values), []int32{
		int32(0), int32(7), int32(0), int32(7),
	})
	assert.Equal(t, game.HomePostgameWinProbability.Value, 0.750937283039093)
	assert.Equal(t, game.HomePregameElo.Value, int32(1974))
	assert.Equal(t, game.HomePostgameElo.Value, int32(1977))
	assert.Equal(t, game.AwayId.Value, int32(251))
	assert.Equal(t, game.AwayTeam, "Texas")
	assert.Equal(t, game.AwayClassification.Value, "fbs")
	assert.Equal(t, game.AwayConference.Value, "SEC")
	assert.Equal(t, game.AwayPoints.Value, int32(7))
	assert.Equal(t, convertToInt32Slice(game.AwayLineScores.Values), []int32{
		int32(0), int32(0), int32(0), int32(7),
	})
	assert.Equal(t, game.AwayPostgameWinProbability.Value, 0.24906271696090698)
	assert.Equal(t, game.AwayPregameElo.Value, int32(1861))
	assert.Equal(t, game.AwayPostgameElo.Value, int32(1858))
	assert.Equal(t, game.Highlights.Value, "")
	assert.Nil(t, game.Attendance)
	assert.Nil(t, game.Notes)
	assert.Equal(t,
		response[0].StartDate.AsTime().Format(defaultTimeFormat),
		"2025-08-30T16:00:00.000Z",
	)
}

func TestGetGamesTeams_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "games_teams.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetGameTeams(
		context.Background(), GetGameTeamsRequest{GameID: 401752723},
	)

	game := response[0]
	team := game.Teams[1]

	assert.NoError(t, err)
	assert.Equal(t, game.Id, int32(401752723))
	assert.Equal(t, team.TeamId, int32(251))
	assert.Equal(t, team.Team, "Texas")
	assert.Equal(t, team.Conference.Value, "SEC")
	assert.Equal(t, team.HomeAway, "home")
	assert.Equal(t, team.Points.Value, int32(55))
	assert.Len(t, team.Stats, 33)
	assert.Equal(t, team.Stats[0].Stat, "26")
	assert.Equal(t, team.Stats[0].Category, "firstDowns")
}

func TestGetGamesPlayers_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "games_players.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetGamePlayers(
		context.Background(), GetGamePlayersRequest{GameID: 401752723},
	)

	game := response[0]
	team := game.Teams[0]

	assert.NoError(t, err)
	assert.Equal(t, game.Id, int32(401752723))
	assert.Equal(t, team.Team, "Texas")
	assert.Equal(t, team.Conference.Value, "SEC")
	assert.Equal(t, team.HomeAway, "home")
	assert.Equal(t, team.Points.Value, int32(55))
	assert.Len(t, team.Categories, 9)

	category := team.Categories[0]
	assert.Equal(t, category.Name, "passing")
	assert.Len(t, category.Types, 6)

	qbr := category.Types[5]
	assert.Equal(t, qbr.Name, "QBR")
	assert.Len(t, qbr.Athletes, 3)

	qb := qbr.Athletes[0]
	assert.Equal(t, qb.Id, "4870906")
	assert.Equal(t, qb.Name, "Arch Manning")
	assert.Equal(t, qb.Stat, "81.6")
}

func TestGetGamesMedia_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "games_media.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetGameMedia(
		context.Background(), GetGameMediaRequest{
			Team: testTeam, Week: testWeek, Year: testYear,
		},
	)

	media := response[0]

	assert.NoError(t, err)
	assert.Equal(t, media.Id, int32(401752693))
	assert.Equal(t, media.Season, int32(2025))
	assert.Equal(t, media.Week, int32(2))
	assert.Equal(t, media.IsStartTime_TBD, false)
	assert.Equal(t, media.HomeTeam, "Texas")
	assert.Equal(t, media.HomeConference.Value, "SEC")
	assert.Equal(t, media.AwayTeam, "San José State")
	assert.Equal(t, media.AwayConference.Value, "Mountain West")
	assert.Equal(t, media.MediaType, "tv")
	assert.Equal(t, media.Outlet, "ABC")
	assert.Equal(t,
		media.StartTime.AsTime().Format(defaultTimeFormat),
		"2025-09-06T16:00:00.000Z",
	)
}

func TestGetGamesWeather_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "games_weather.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetGameWeather(
		context.Background(), GetGameWeatherRequest{
			GameID: 401767476,
		},
	)

	media := response[0]

	assert.NoError(t, err)
	assert.Equal(t, media.Id, int32(401767476))
	assert.Equal(t, media.Season, int32(2025))
	assert.Equal(t, media.Week, int32(1))
	assert.Equal(t, media.GameIndoors, false)
	assert.Equal(t, media.HomeTeam, "Nicholls")
	assert.Equal(t, media.HomeConference.Value, "Southland")
	assert.Equal(t, media.AwayTeam, "Incarnate Word")
	assert.Equal(t, media.AwayConference.Value, "Southland")
	assert.Equal(t, media.VenueId.Value, int32(3779))
	assert.Equal(t, media.Temperature.Value, 89.6)
	assert.Equal(t, media.DewPoint.Value, 73.4)
	assert.Equal(t, media.Humidity.Value, float64(59))
	assert.Equal(t, media.Precipitation.Value, float64(0.004))
	assert.Equal(t, media.Snowfall.Value, float64(0))
	assert.Equal(t, media.WindDirection.Value, float64(340))
	assert.Equal(t, media.WindSpeed.Value, float64(8.1))
	assert.Equal(t, media.Pressure.Value, float64(1014))
	assert.Equal(t, media.WeatherConditionCode.Value, float64(7))
	assert.Equal(t, media.WeatherCondition.Value, "Light Rain")
	assert.Equal(t, media.Venue.Value, "Manning Field at John L. Guidry Stadium")
	assert.Equal(t,
		media.StartTime.AsTime().Format(defaultTimeFormat),
		"2025-08-23T17:00:00.000Z",
	)
}

func TestGetRecords_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "records.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetTeamRecords(
		context.Background(), GetTeamRecordsRequest{
			Team: testTeam, Year: testYear,
		},
	)

	team := response[0]

	assert.NoError(t, err)
	assert.Equal(t, team.Year, int32(2025))
	assert.Equal(t, team.TeamId.Value, int32(251))
	assert.Equal(t, team.Team, "Texas")
	assert.Equal(t, team.Classification.Value, "fbs")
	assert.Equal(t, team.Conference, "SEC")
	assert.Equal(t, team.Division, "")
	assert.Equal(t, team.ExpectedWins.Value, float64(7.891881301999092))
	assert.Equal(t, team.Total.Games, int32(12))
	assert.Equal(t, team.Total.Wins, int32(9))
	assert.Equal(t, team.Total.Losses, int32(3))
	assert.Equal(t, team.Total.Ties, int32(0))

	assert.Equal(t, team.ConferenceGames.Games, int32(8))
	assert.Equal(t, team.ConferenceGames.Wins, int32(6))
	assert.Equal(t, team.ConferenceGames.Losses, int32(2))
	assert.Equal(t, team.ConferenceGames.Ties, int32(0))

	assert.Equal(t, team.HomeGames.Games, int32(6))
	assert.Equal(t, team.HomeGames.Wins, int32(6))
	assert.Equal(t, team.HomeGames.Losses, int32(0))
	assert.Equal(t, team.HomeGames.Ties, int32(0))

	assert.Equal(t, team.AwayGames.Games, int32(5))
	assert.Equal(t, team.AwayGames.Wins, int32(2))
	assert.Equal(t, team.AwayGames.Losses, int32(3))
	assert.Equal(t, team.AwayGames.Ties, int32(0))

	assert.Equal(t, team.NeutralSiteGames.Games, int32(1))
	assert.Equal(t, team.NeutralSiteGames.Wins, int32(1))
	assert.Equal(t, team.NeutralSiteGames.Losses, int32(0))
	assert.Equal(t, team.NeutralSiteGames.Ties, int32(0))

	assert.Equal(t, team.RegularSeason.Games, int32(12))
	assert.Equal(t, team.RegularSeason.Wins, int32(9))
	assert.Equal(t, team.RegularSeason.Losses, int32(3))
	assert.Equal(t, team.RegularSeason.Ties, int32(0))

	assert.Equal(t, team.Postseason.Games, int32(0))
	assert.Equal(t, team.Postseason.Wins, int32(0))
	assert.Equal(t, team.Postseason.Losses, int32(0))
	assert.Equal(t, team.Postseason.Ties, int32(0))
}

func TestGetCalendar_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "calendar.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetCalendar(
		context.Background(), testYear,
	)

	require.NoError(t, err)
	assert.Len(t, response, 17)

	week := response[0]
	assert.Equal(t, week.Season, int32(2025))
	assert.Equal(t, week.Week, int32(1))
	assert.Equal(t, week.SeasonType, "regular")
	assert.Equal(t,
		week.StartDate.AsTime().Format(defaultTimeFormat),
		"2025-08-23T07:00:00.000Z",
	)
	assert.Equal(t,
		week.EndDate.AsTime().Format(defaultTimeFormat),
		"2025-09-02T06:59:00.000Z",
	)
	assert.Equal(t,
		week.FirstGameStart.AsTime().Format(defaultTimeFormat),
		"2025-08-23T07:00:00.000Z",
	)
	assert.Equal(t,
		week.LastGameStart.AsTime().Format(defaultTimeFormat),
		"2025-09-02T06:59:00.000Z",
	)
}

func TestGetScoreboard_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "scoreboard.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetScoreboard(
		context.Background(), GetScoreboardRequest{},
	)

	// week := response[0]
	assert.NoError(t, err)

	score := response[0]
	assert.Equal(t, score.Id, int32(401762521))
	assert.Equal(t,
		score.StartDate.AsTime().Format(defaultTimeFormat),
		"2025-12-13T20:00:00.000Z",
	)
	assert.Equal(t, score.StartTime_TBD, false)
	assert.Equal(t, score.Tv.Value, "CBS")
	assert.Equal(t, score.NeutralSite, true)
	assert.Equal(t, score.ConferenceGame, true)
	assert.Equal(t, score.Status, "completed")
	assert.Nil(t, score.Period)
	assert.Nil(t, score.Clock)
	assert.Equal(t, score.Situation.Value, "3rd & 13 at ARMY 42")
	assert.Equal(t, score.Possession.Value, "home")
	assert.Equal(t,
		score.LastPlay.Value,
		"(00:45) Kneel down by Navy at Army42 for loss of 2 yards",
	)
}

func TestGetAdvancedBoxScore_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "advanced_box_score.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetAdvancedBoxScore(
		context.Background(), 401752677,
	)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Test gameInfo
	assert.NotNil(t, response.GameInfo)
	assert.Equal(t, response.GameInfo.HomeTeam, "Ohio State")
	assert.Equal(t, response.GameInfo.HomePoints, int32(14))
	assert.Equal(t, response.GameInfo.HomeWinProb, 0.750937283039093)
	assert.Equal(t, response.GameInfo.AwayTeam, "Texas")
	assert.Equal(t, response.GameInfo.AwayPoints, int32(7))
	assert.Equal(t, response.GameInfo.AwayWinProb, 0.24906271696090698)
	assert.Equal(t, response.GameInfo.HomeWinner, true)
	assert.Equal(t, response.GameInfo.Excitement, 3.8398377534422923)

	// Test teams PPA
	assert.NotNil(t, response.Teams)
	assert.NotNil(t, response.Teams.Ppa)
	assert.Len(t, response.Teams.Ppa, 2)

	ohioStatePPA := response.Teams.Ppa[0]
	assert.Equal(t, ohioStatePPA.Team, "Ohio State")
	assert.Equal(t, ohioStatePPA.Plays, int32(54))
	assert.NotNil(t, ohioStatePPA.Overall)
	assert.Equal(t, ohioStatePPA.Overall.Total, 0.0357)
	assert.NotNil(t, ohioStatePPA.Overall.Quarter1)
	assert.Equal(t, ohioStatePPA.Overall.Quarter1.Value, -0.096)
	assert.NotNil(t, ohioStatePPA.Passing)
	assert.Equal(t, ohioStatePPA.Passing.Total, 0.319)
	assert.NotNil(t, ohioStatePPA.Rushing)
	assert.Equal(t, ohioStatePPA.Rushing.Total, -0.131)

	texasPPA := response.Teams.Ppa[1]
	assert.Equal(t, texasPPA.Team, "Texas")
	assert.Equal(t, texasPPA.Plays, int32(67))

	// Test cumulative PPA
	assert.NotNil(t, response.Teams.CumulativePpa)
	assert.Len(t, response.Teams.CumulativePpa, 2)

	// Test success rates
	assert.NotNil(t, response.Teams.SuccessRates)
	assert.Len(t, response.Teams.SuccessRates, 2)
	ohioStateSR := response.Teams.SuccessRates[0]
	assert.Equal(t, ohioStateSR.Team, "Ohio State")
	assert.NotNil(t, ohioStateSR.Overall)
	assert.Equal(t, ohioStateSR.Overall.Total, 0.333)

	// Test explosiveness
	assert.NotNil(t, response.Teams.Explosiveness)
	assert.Len(t, response.Teams.Explosiveness, 2)

	// Test rushing stats
	assert.NotNil(t, response.Teams.Rushing)
	assert.Len(t, response.Teams.Rushing, 2)
	ohioStateRush := response.Teams.Rushing[0]
	assert.Equal(t, ohioStateRush.Team, "Ohio State")
	assert.Equal(t, ohioStateRush.PowerSuccess, 0.5)
	assert.Equal(t, ohioStateRush.StuffRate, 0.265)
	assert.Equal(t, ohioStateRush.LineYards, 64.0)

	// Test havoc
	assert.NotNil(t, response.Teams.Havoc)
	assert.Len(t, response.Teams.Havoc, 2)

	// Test scoring opportunities
	assert.NotNil(t, response.Teams.ScoringOpportunities)
	assert.Len(t, response.Teams.ScoringOpportunities, 2)
	ohioStateSO := response.Teams.ScoringOpportunities[0]
	assert.Equal(t, ohioStateSO.Team, "Ohio State")
	assert.Equal(t, ohioStateSO.Opportunities, int32(2))
	assert.Equal(t, ohioStateSO.Points, int32(14))

	// Test field position
	assert.NotNil(t, response.Teams.FieldPosition)
	assert.Len(t, response.Teams.FieldPosition, 2)

	// Test players usage
	assert.NotNil(t, response.Players)
	assert.NotNil(t, response.Players.Usage)
	assert.Greater(t, len(response.Players.Usage), 0)

	// Test players PPA
	assert.NotNil(t, response.Players.Ppa)
	assert.Greater(t, len(response.Players.Ppa), 0)
	playerPPA := response.Players.Ppa[0]
	assert.NotNil(t, playerPPA.Player)
	assert.NotNil(t, playerPPA.Average)
	assert.NotNil(t, playerPPA.Cumulative)
}

func TestGetDrives_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "drives.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetDrives(
		context.Background(), GetDrivesRequest{
			Year: testYear,
			Week: testWeek,
			Team: testTeam,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Len(t, response, 1)

	drive := response[0]
	assert.Equal(t, drive.Id, "4017526931")
	assert.Equal(t, drive.GameId, int32(401752693))
	assert.Equal(t, drive.Offense, "Texas")
	assert.Equal(t, drive.OffenseConference.Value, "SEC")
	assert.Equal(t, drive.Defense, "San José State")
	assert.Equal(t, drive.DefenseConference.Value, "Mountain West")
	assert.Equal(t, drive.DriveNumber.Value, int32(1))
	assert.Equal(t, drive.Scoring, false)
	assert.Equal(t, drive.StartPeriod, int32(1))
	assert.Equal(t, drive.StartYardline, int32(25))
	assert.Equal(t, drive.StartYardsToGoal, int32(75))
	assert.NotNil(t, drive.StartTime)
	assert.Equal(t, drive.StartTime.Minutes.Value, int32(15))
	assert.Equal(t, drive.StartTime.Seconds.Value, int32(0))
	assert.Equal(t, drive.EndPeriod, int32(1))
	assert.Equal(t, drive.EndYardline, int32(21))
	assert.Equal(t, drive.EndYardsToGoal, int32(79))
	assert.NotNil(t, drive.EndTime)
	assert.Equal(t, drive.EndTime.Minutes.Value, int32(13))
	assert.Equal(t, drive.EndTime.Seconds.Value, int32(54))
	assert.NotNil(t, drive.Elapsed)
	assert.Equal(t, drive.Elapsed.Minutes.Value, int32(1))
	assert.Equal(t, drive.Elapsed.Seconds.Value, int32(6))
	assert.Equal(t, drive.Plays, int32(4))
	assert.Equal(t, drive.Yards, int32(-4))
	assert.Equal(t, drive.DriveResult, "PUNT")
	assert.Equal(t, drive.IsHomeOffense, true)
	assert.Equal(t, drive.StartOffenseScore, int32(0))
	assert.Equal(t, drive.StartDefenseScore, int32(0))
	assert.Equal(t, drive.EndOffenseScore, int32(0))
	assert.Equal(t, drive.EndDefenseScore, int32(0))
}

func TestGetPlays_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "plays.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetPlays(
		context.Background(), GetPlaysRequest{
			Year: testYear,
			Week: testWeek,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Len(t, response, 2)

	// Test first play
	play1 := response[0]
	assert.Equal(t, play1.GameId, int32(401752693))
	assert.Equal(t, play1.DriveId, "40175269311")
	assert.Equal(t, play1.Id, "401752693102869401")
	assert.Equal(t, play1.DriveNumber.Value, int32(11))
	assert.Equal(t, play1.PlayNumber.Value, int32(2))
	assert.Equal(t, play1.Offense, "Texas")
	assert.Equal(t, play1.OffenseConference.Value, "SEC")
	assert.Equal(t, play1.OffenseScore, int32(21))
	assert.Equal(t, play1.Defense, "San José State")
	assert.Equal(t, play1.DefenseConference.Value, "Mountain West")
	assert.Equal(t, play1.DefenseScore, int32(0))
	assert.Equal(t, play1.Home, "Texas")
	assert.Equal(t, play1.Away, "San José State")
	assert.Equal(t, play1.Period, int32(2))
	assert.NotNil(t, play1.Clock)
	assert.Equal(t, play1.Clock.Minutes.Value, int32(13))
	assert.Equal(t, play1.Clock.Seconds.Value, int32(5))
	assert.Equal(t, play1.OffenseTimeouts.Value, int32(3))
	assert.Equal(t, play1.DefenseTimeouts.Value, int32(3))
	assert.Equal(t, play1.Yardline, int32(99))
	assert.Equal(t, play1.YardsToGoal, int32(1))
	assert.Equal(t, play1.Down, int32(2))
	assert.Equal(t, play1.Distance, int32(0))
	assert.Equal(t, play1.YardsGained, int32(-15))
	assert.Equal(t, play1.Scoring, false)
	assert.Equal(t, play1.PlayType, "Penalty")
	assert.NotNil(t, play1.PlayText)
	assert.Contains(t, play1.PlayText.Value, "Baxter, CJ rush for 1 yard")
	assert.Nil(t, play1.Ppa)
	assert.NotNil(t, play1.Wallclock)

	// Test second play (scoring play)
	play2 := response[1]
	assert.Equal(t, play2.GameId, int32(401752693))
	assert.Equal(t, play2.DriveId, "40175269311")
	assert.Equal(t, play2.Id, "401752693102874301")
	assert.Equal(t, play2.DriveNumber.Value, int32(11))
	assert.Equal(t, play2.PlayNumber.Value, int32(3))
	assert.Equal(t, play2.Offense, "Texas")
	assert.Equal(t, play2.OffenseConference.Value, "SEC")
	assert.Equal(t, play2.OffenseScore, int32(28))
	assert.Equal(t, play2.Defense, "San José State")
	assert.Equal(t, play2.DefenseConference.Value, "Mountain West")
	assert.Equal(t, play2.DefenseScore, int32(0))
	assert.Equal(t, play2.Home, "Texas")
	assert.Equal(t, play2.Away, "San José State")
	assert.Equal(t, play2.Period, int32(2))
	assert.NotNil(t, play2.Clock)
	assert.Equal(t, play2.Clock.Minutes.Value, int32(12))
	assert.Equal(t, play2.Clock.Seconds.Value, int32(56))
	assert.Equal(t, play2.OffenseTimeouts.Value, int32(3))
	assert.Equal(t, play2.DefenseTimeouts.Value, int32(3))
	assert.Equal(t, play2.Yardline, int32(84))
	assert.Equal(t, play2.YardsToGoal, int32(16))
	assert.Equal(t, play2.Down, int32(2))
	assert.Equal(t, play2.Distance, int32(0))
	assert.Equal(t, play2.YardsGained, int32(16))
	assert.Equal(t, play2.Scoring, true)
	assert.Equal(t, play2.PlayType, "Passing Touchdown")
	assert.NotNil(t, play2.PlayText)
	assert.Contains(t, play2.PlayText.Value, "Arch Manning pass complete")
	assert.Nil(t, play2.Ppa)
	assert.NotNil(t, play2.Wallclock)
}

func TestGetPlayTypes_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "play_types.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetPlayTypes(context.Background())

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Greater(t, len(response), 0)

	// Helper function to find play type by ID
	findPlayType := func(id int32) *PlayType {
		for _, pt := range response {
			if pt.Id == id {
				return pt
			}
		}
		return nil
	}

	// Test Rush
	rush := findPlayType(5)
	assert.NotNil(t, rush)
	assert.Equal(t, rush.Text, "Rush")
	assert.NotNil(t, rush.Abbreviation)
	assert.Equal(t, rush.Abbreviation.Value, "RUSH")

	// Test Pass Incompletion (has null abbreviation)
	passIncomplete := findPlayType(3)
	assert.NotNil(t, passIncomplete)
	assert.Equal(t, passIncomplete.Text, "Pass Incompletion")
	assert.Nil(t, passIncomplete.Abbreviation)

	// Test Kickoff
	kickoff := findPlayType(53)
	assert.NotNil(t, kickoff)
	assert.Equal(t, kickoff.Text, "Kickoff")
	assert.NotNil(t, kickoff.Abbreviation)
	assert.Equal(t, kickoff.Abbreviation.Value, "K")

	// Test Punt
	punt := findPlayType(52)
	assert.NotNil(t, punt)
	assert.Equal(t, punt.Text, "Punt")
	assert.NotNil(t, punt.Abbreviation)
	assert.Equal(t, punt.Abbreviation.Value, "PUNT")

	// Test Penalty
	penalty := findPlayType(8)
	assert.NotNil(t, penalty)
	assert.Equal(t, penalty.Text, "Penalty")
	assert.NotNil(t, penalty.Abbreviation)
	assert.Equal(t, penalty.Abbreviation.Value, "PEN")

	// Test Passing Touchdown
	passingTD := findPlayType(67)
	assert.NotNil(t, passingTD)
	assert.Equal(t, passingTD.Text, "Passing Touchdown")
	assert.NotNil(t, passingTD.Abbreviation)
	assert.Equal(t, passingTD.Abbreviation.Value, "TD")

	// Test Rushing Touchdown
	rushingTD := findPlayType(68)
	assert.NotNil(t, rushingTD)
	assert.Equal(t, rushingTD.Text, "Rushing Touchdown")
	assert.NotNil(t, rushingTD.Abbreviation)
	assert.Equal(t, rushingTD.Abbreviation.Value, "TD")

	// Test Field Goal Good
	fieldGoal := findPlayType(59)
	assert.NotNil(t, fieldGoal)
	assert.Equal(t, fieldGoal.Text, "Field Goal Good")
	assert.NotNil(t, fieldGoal.Abbreviation)
	assert.Equal(t, fieldGoal.Abbreviation.Value, "FG")

	// Test Uncategorized (has null abbreviation)
	uncategorized := findPlayType(999)
	assert.NotNil(t, uncategorized)
	assert.Equal(t, uncategorized.Text, "Uncategorized")
	assert.Nil(t, uncategorized.Abbreviation)

	// Test Sack (has null abbreviation)
	sack := findPlayType(7)
	assert.NotNil(t, sack)
	assert.Equal(t, sack.Text, "Sack")
	assert.Nil(t, sack.Abbreviation)
}

func TestGetPlayStats_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "plays_stats.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetPlayStats(
		context.Background(), GetPlayStatsRequest{
			Year: testYear,
			Week: testWeek,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Len(t, response, 2)

	stat1 := response[0]
	assert.Equal(t, stat1.GameId, float64(401752693))
	assert.Equal(t, stat1.Season, float64(2025))
	assert.Equal(t, stat1.Week, float64(2))
	assert.Equal(t, stat1.Team, "Texas")
	assert.Equal(t, stat1.Conference, "SEC")
	assert.Equal(t, stat1.Opponent, "San José State")
	assert.Equal(t, stat1.TeamScore, float64(0))
	assert.Equal(t, stat1.OpponentScore, float64(0))
	assert.Equal(t, stat1.DriveId, "4017526931")
	assert.Equal(t, stat1.PlayId, "401752693101854901")
	assert.Equal(t, stat1.Period, float64(1))
	assert.NotNil(t, stat1.Clock)
	assert.Equal(t, stat1.Clock.Minutes.Value, float64(14))
	assert.Equal(t, stat1.Clock.Seconds.Value, float64(50))
	assert.Equal(t, stat1.YardsToGoal, float64(75))
	assert.Equal(t, stat1.Down, float64(1))
	assert.Equal(t, stat1.Distance, float64(10))
	assert.Equal(t, stat1.AthleteId, "4870906")
	assert.Equal(t, stat1.AthleteName, "Arch Manning")
	assert.Equal(t, stat1.StatType, "Completion")
	assert.Equal(t, stat1.Stat, float64(6))

	stat2 := response[1]
	assert.Equal(t, stat2.GameId, float64(401752693))
	assert.Equal(t, stat2.Season, float64(2025))
	assert.Equal(t, stat2.Week, float64(2))
	assert.Equal(t, stat2.Team, "Texas")
	assert.Equal(t, stat2.Conference, "SEC")
	assert.Equal(t, stat2.Opponent, "San José State")
	assert.Equal(t, stat2.TeamScore, float64(0))
	assert.Equal(t, stat2.OpponentScore, float64(0))
	assert.Equal(t, stat2.DriveId, "4017526931")
	assert.Equal(t, stat2.PlayId, "401752693101859201")
	assert.Equal(t, stat2.Period, float64(1))
	assert.NotNil(t, stat2.Clock)
	assert.Equal(t, stat2.Clock.Minutes.Value, float64(14))
	assert.Equal(t, stat2.Clock.Seconds.Value, float64(7))
	assert.Equal(t, stat2.YardsToGoal, float64(79))
	assert.Equal(t, stat2.Down, float64(2))
	assert.Equal(t, stat2.Distance, float64(14))
	assert.Equal(t, stat2.AthleteId, "4870906")
	assert.Equal(t, stat2.AthleteName, "Arch Manning")
	assert.Equal(t, stat2.StatType, "Incompletion")
	assert.Equal(t, stat2.Stat, float64(1))
}

func TestGetPlayStatTypes_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "play_stats_types.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetPlayStatTypes(context.Background())

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Greater(t, len(response), 0)

	// Helper function to find play stat type by ID
	findPlayStatType := func(id int32) *PlayStatType {
		for _, pst := range response {
			if pst.Id == id {
				return pst
			}
		}
		return nil
	}

	// Test Incompletion
	incompletion := findPlayStatType(1)
	require.NotNil(t, incompletion)
	assert.Equal(t, incompletion.Name, "Incompletion")

	// Test Target
	target := findPlayStatType(2)
	require.NotNil(t, target)
	assert.Equal(t, target.Name, "Target")

	// Test Completion
	completion := findPlayStatType(4)
	require.NotNil(t, completion)
	assert.Equal(t, completion.Name, "Completion")

	// Test Reception
	reception := findPlayStatType(5)
	require.NotNil(t, reception)
	assert.Equal(t, reception.Name, "Reception")

	// Test Tackle
	tackle := findPlayStatType(6)
	require.NotNil(t, tackle)
	assert.Equal(t, tackle.Name, "Tackle")

	// Test Rush
	rush := findPlayStatType(7)
	require.NotNil(t, rush)
	assert.Equal(t, rush.Name, "Rush")

	// Test Fumble
	fumble := findPlayStatType(8)
	require.NotNil(t, fumble)
	assert.Equal(t, fumble.Name, "Fumble")

	// Test Sack
	sack := findPlayStatType(12)
	require.NotNil(t, sack)
	assert.Equal(t, sack.Name, "Sack")

	// Test Kickoff
	kickoff := findPlayStatType(13)
	require.NotNil(t, kickoff)
	assert.Equal(t, kickoff.Name, "Kickoff")

	// Test Punt
	punt := findPlayStatType(16)
	require.NotNil(t, punt)
	assert.Equal(t, punt.Name, "Punt")

	// Test Interception
	interception := findPlayStatType(21)
	require.NotNil(t, interception)
	assert.Equal(t, interception.Name, "Interception")

	// Test Touchdown
	touchdown := findPlayStatType(22)
	require.NotNil(t, touchdown)
	assert.Equal(t, touchdown.Name, "Touchdown")

	// Test Field Goal Made
	fieldGoalMade := findPlayStatType(24)
	require.NotNil(t, fieldGoalMade)
	assert.Equal(t, fieldGoalMade.Name, "Field Goal Made")

	// Test Field Goal Missed
	fieldGoalMissed := findPlayStatType(25)
	require.NotNil(t, fieldGoalMissed)
	assert.Equal(t, fieldGoalMissed.Name, "Field Goal Missed")

	// Test QB Hurry
	qbHurry := findPlayStatType(26)
	require.NotNil(t, qbHurry)
	assert.Equal(t, qbHurry.Name, "QB Hurry")
}

func TestGetLivePlays_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "live_plays.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetLivePlays(context.Background(), 401778330)

	require.NoError(t, err)
	require.NotNil(t, response)

	// Test game-level fields
	assert.Equal(t, response.Id, int32(401778330))
	assert.Equal(t, response.Status, "In Progress")
	require.NotNil(t, response.Period)
	assert.Equal(t, response.Period.Value, int32(1))
	assert.Equal(t, response.Clock, "3:15")
	assert.Equal(t, response.Possession, "Michigan")
	require.NotNil(t, response.Down)
	assert.Equal(t, response.Down.Value, int32(4))
	require.NotNil(t, response.Distance)
	assert.Equal(t, response.Distance.Value, int32(6))
	require.NotNil(t, response.YardsToGoal)
	assert.Equal(t, response.YardsToGoal.Value, int32(19))

	// Test teams
	require.NotNil(t, response.Teams)
	assert.Len(t, response.Teams, 2)

	texasTeam := response.Teams[0]
	assert.Equal(t, texasTeam.TeamId, int32(251))
	assert.Equal(t, texasTeam.Team, "Texas")
	assert.Equal(t, texasTeam.HomeAway, "home")
	assert.Equal(t, texasTeam.LineScores, []int32{3})
	assert.Equal(t, texasTeam.Points, int32(3))
	assert.Equal(t, texasTeam.Drives, int32(2))
	assert.Equal(t, texasTeam.ScoringOpportunities, int32(1))
	assert.Equal(t, texasTeam.PointsPerOpportunity, 3.0)
	require.NotNil(t, texasTeam.AverageStartYardLine)
	assert.Equal(t, texasTeam.AverageStartYardLine.Value, 76.0)
	assert.Equal(t, texasTeam.Plays, int32(14))
	assert.Equal(t, texasTeam.LineYards, 13.0)
	assert.Equal(t, texasTeam.LineYardsPerRush, 4.3)
	assert.Equal(t, texasTeam.SuccessRate, 0.429)
	assert.Equal(t, texasTeam.Explosiveness, 1.128)
	require.NotNil(t, texasTeam.DeserveToWin)
	assert.Equal(t, texasTeam.DeserveToWin.Value, 0.467)

	michiganTeam := response.Teams[1]
	assert.Equal(t, michiganTeam.TeamId, int32(130))
	assert.Equal(t, michiganTeam.Team, "Michigan")
	assert.Equal(t, michiganTeam.HomeAway, "away")
	assert.Equal(t, michiganTeam.Points, int32(3))
	assert.Equal(t, michiganTeam.Drives, int32(3))

	// Test drives
	require.NotNil(t, response.Drives)
	assert.Len(t, response.Drives, 5)

	// Test first drive
	drive1 := response.Drives[0]
	assert.Equal(t, drive1.Id, "4017783302")
	assert.Equal(t, drive1.OffenseId, int32(251))
	assert.Equal(t, drive1.Offense, "Texas")
	assert.Equal(t, drive1.DefenseId, int32(130))
	assert.Equal(t, drive1.Defense, "Michigan")
	assert.Equal(t, drive1.PlayCount, int32(9))
	assert.Equal(t, drive1.Yards, int32(50))
	assert.Equal(t, drive1.StartPeriod, int32(1))
	require.NotNil(t, drive1.StartClock)
	assert.Equal(t, drive1.StartClock.Value, "15:00")
	assert.Equal(t, drive1.StartYardsToGoal, int32(75))
	require.NotNil(t, drive1.EndPeriod)
	assert.Equal(t, drive1.EndPeriod.Value, int32(1))
	require.NotNil(t, drive1.EndClock)
	assert.Equal(t, drive1.EndClock.Value, "12:05")
	require.NotNil(t, drive1.EndYardsToGoal)
	assert.Equal(t, drive1.EndYardsToGoal.Value, int32(25))
	require.NotNil(t, drive1.Duration)
	assert.Equal(t, drive1.Duration.Value, "2:55")
	assert.Equal(t, drive1.ScoringOpportunity, true)
	assert.Equal(t, drive1.Result, "Field Goal")
	assert.Equal(t, drive1.PointsGained, int32(3))

	// Test plays in first drive
	require.NotNil(t, drive1.Plays)
	assert.Len(t, drive1.Plays, 1)
	play1 := drive1.Plays[0]
	assert.Equal(t, play1.Id, "4017783303")
	assert.Equal(t, play1.HomeScore, int32(0))
	assert.Equal(t, play1.AwayScore, int32(0))
	assert.Equal(t, play1.Period, int32(1))
	assert.Equal(t, play1.Clock, "15:00")
	assert.NotNil(t, play1.WallClock)
	assert.Equal(t, play1.TeamId, int32(130))
	assert.Equal(t, play1.Team, "Michigan")
	assert.Equal(t, play1.Down, int32(1))
	assert.Equal(t, play1.Distance, int32(10))
	assert.Equal(t, play1.YardsToGoal, int32(65))
	assert.Equal(t, play1.YardsGained, int32(0))
	assert.Equal(t, play1.PlayTypeId, int32(53))
	assert.Equal(t, play1.PlayType, "Kickoff")
	assert.Nil(t, play1.Epa)
	assert.Equal(t, play1.GarbageTime, false)
	assert.Equal(t, play1.Success, false)
	assert.Equal(t, play1.RushPass, "other")
	assert.Equal(t, play1.DownType, "standard")
	assert.Contains(t, play1.PlayText, "B.Sunderland kickoff")

	// Test drive with null end fields (drive 5)
	drive5 := response.Drives[4]
	assert.Equal(t, drive5.Id, "4017783308")
	assert.Nil(t, drive5.StartClock)
	assert.Nil(t, drive5.EndPeriod)
	assert.Nil(t, drive5.EndClock)
	assert.Nil(t, drive5.EndYardsToGoal)
	assert.Equal(t, drive5.Result, "")
	assert.Equal(t, drive5.PointsGained, int32(0))

	// Test play with EPA value
	drive4 := response.Drives[3]
	require.NotNil(t, drive4.Plays)
	assert.Len(t, drive4.Plays, 1)
	playWithEPA := drive4.Plays[0]
	require.NotNil(t, playWithEPA.Epa)
	assert.Equal(t, playWithEPA.Epa.Value, 0.911)
	assert.Equal(t, playWithEPA.Success, true)
	assert.Equal(t, playWithEPA.RushPass, "rush")
}

func TestGetTeams_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "teams.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetTeams(
		context.Background(), GetTeamsRequest{
			Conference: "SEC",
		},
	)

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Greater(t, len(response), 0)

	// Test single team (Alabama - first team)
	team := response[0]
	assert.Equal(t, team.Id, int32(333))
	assert.Equal(t, team.School, "Alabama")
	require.NotNil(t, team.Mascot)
	assert.Equal(t, team.Mascot.Value, "Crimson Tide")
	require.NotNil(t, team.Abbreviation)
	assert.Equal(t, team.Abbreviation.Value, "ALA")
	require.NotNil(t, team.AlternateNames)
	assert.Greater(t, len(team.AlternateNames.Values), 0)
	require.NotNil(t, team.Conference)
	assert.Equal(t, team.Conference.Value, "SEC")
	assert.Nil(t, team.Division) // null in JSON
	require.NotNil(t, team.Classification)
	assert.Equal(t, team.Classification.Value, "fbs")
	require.NotNil(t, team.Color)
	assert.Equal(t, team.Color.Value, "#9e1632")
	require.NotNil(t, team.AlternateColor)
	assert.Equal(t, team.AlternateColor.Value, "#ffffff")
	require.NotNil(t, team.Logos)
	assert.Greater(t, len(team.Logos.Values), 0)
	assert.Equal(t, team.Twitter, "@AlabamaFTBL")

	// Test location/venue
	require.NotNil(t, team.Location)
	location := team.Location
	require.NotNil(t, location.Id)
	assert.Equal(t, location.Id.Value, int32(3657))
	require.NotNil(t, location.Name)
	assert.Equal(t, location.Name.Value, "Bryant-Denny Stadium")
	require.NotNil(t, location.City)
	assert.Equal(t, location.City.Value, "Tuscaloosa")
	require.NotNil(t, location.State)
	assert.Equal(t, location.State.Value, "AL")
	require.NotNil(t, location.Zip)
	assert.Equal(t, location.Zip.Value, "35487")
	require.NotNil(t, location.CountryCode)
	assert.Equal(t, location.CountryCode.Value, "US")
	require.NotNil(t, location.Timezone)
	assert.Equal(t, location.Timezone.Value, "America/Chicago")
	require.NotNil(t, location.Latitude)
	assert.Equal(t, location.Latitude.Value, 33.2082752)
	require.NotNil(t, location.Longitude)
	assert.Equal(t, location.Longitude.Value, -87.5503836)
	require.NotNil(t, location.Elevation)
	assert.Equal(t, location.Elevation.Value, "70.05136108")
	require.NotNil(t, location.Capacity)
	assert.Equal(t, location.Capacity.Value, int32(101821))
	require.NotNil(t, location.ConstructionYear)
	assert.Equal(t, location.ConstructionYear.Value, int32(1929))
	require.NotNil(t, location.Grass)
	assert.Equal(t, location.Grass.Value, true)
	require.NotNil(t, location.Dome)
	assert.Equal(t, location.Dome.Value, false)
}

func TestGetFBSTeams_ValidRequest_ShouldSucceed(t *testing.T) {
	tester, bytes := setupTestWithFile(t, "teams_fbs.json")

	tester.requestExecutor.EXPECT().
		Execute(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bytes, nil).
		Times(1)

	response, err := tester.client.GetFBSTeams(
		context.Background(), GetFBSTeamsRequest{
			Year: testYear,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Greater(t, len(response), 0)

	// Test single team (Air Force - first team)
	team := response[0]
	assert.Equal(t, team.Id, int32(2005))
	assert.Equal(t, team.School, "Air Force")
	require.NotNil(t, team.Mascot)
	assert.Equal(t, team.Mascot.Value, "Falcons")
	require.NotNil(t, team.Abbreviation)
	assert.Equal(t, team.Abbreviation.Value, "AF")
	require.NotNil(t, team.AlternateNames)
	assert.Greater(t, len(team.AlternateNames.Values), 0)
	require.NotNil(t, team.Conference)
	assert.Equal(t, team.Conference.Value, "Mountain West")
	assert.Nil(t, team.Division) // null in JSON
	require.NotNil(t, team.Classification)
	assert.Equal(t, team.Classification.Value, "fbs")
	require.NotNil(t, team.Color)
	assert.Equal(t, team.Color.Value, "#004a7b")
	require.NotNil(t, team.AlternateColor)
	assert.Equal(t, team.AlternateColor.Value, "#ffffff")
	require.NotNil(t, team.Logos)
	assert.Greater(t, len(team.Logos.Values), 0)
	assert.Equal(t, team.Twitter, "@AF_Football")

	// Test location/venue
	require.NotNil(t, team.Location)
	location := team.Location
	require.NotNil(t, location.Id)
	assert.Equal(t, location.Id.Value, int32(3713))
	require.NotNil(t, location.Name)
	assert.Equal(t, location.Name.Value, "Falcon Stadium")
	require.NotNil(t, location.City)
	assert.Equal(t, location.City.Value, "Colorado Springs")
	require.NotNil(t, location.State)
	assert.Equal(t, location.State.Value, "CO")
	require.NotNil(t, location.Zip)
	assert.Equal(t, location.Zip.Value, "80840")
	require.NotNil(t, location.CountryCode)
	assert.Equal(t, location.CountryCode.Value, "US")
	require.NotNil(t, location.Timezone)
	assert.Equal(t, location.Timezone.Value, "America/Denver")
	require.NotNil(t, location.Latitude)
	assert.Equal(t, location.Latitude.Value, 38.9969701)
	require.NotNil(t, location.Longitude)
	assert.Equal(t, location.Longitude.Value, -104.8436165)
	require.NotNil(t, location.Elevation)
	assert.Equal(t, location.Elevation.Value, "2024.875732")
	require.NotNil(t, location.Capacity)
	assert.Equal(t, location.Capacity.Value, int32(46692))
	require.NotNil(t, location.ConstructionYear)
	assert.Equal(t, location.ConstructionYear.Value, int32(1962))
	require.NotNil(t, location.Grass)
	assert.Equal(t, location.Grass.Value, false)
	require.NotNil(t, location.Dome)
	assert.Equal(t, location.Dome.Value, false)
}

func convertToInt32Slice(values []*structpb.Value) []int32 {
	results := make([]int32, len(values))
	for i, v := range values {
		results[i] = int32(v.GetNumberValue())
	}

	return results
}

// mock of request httpGet below

// mockHttpGetExecutor is a mock of httpGetExecutor interface.
type mockHttpGetExecutor struct {
	ctrl     *gomock.Controller
	recorder *mockHttpGetExecutorMockRecorder
}

// mockHttpGetExecutorMockRecorder is the mock recorder for mockHttpGetExecutor.
type mockHttpGetExecutorMockRecorder struct {
	mock *mockHttpGetExecutor
}

// newMockHttpGetExecutor creates a new mock instance.
func newMockHttpGetExecutor(ctrl *gomock.Controller) *mockHttpGetExecutor {
	mock := &mockHttpGetExecutor{ctrl: ctrl}
	mock.recorder = &mockHttpGetExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *mockHttpGetExecutor) EXPECT() *mockHttpGetExecutorMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *mockHttpGetExecutor) Execute(
	ctx context.Context, path string, params url.Values,
) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, path, params)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *mockHttpGetExecutorMockRecorder) Execute(
	ctx, path, params interface{},
) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock, "Execute", reflect.TypeOf((*mockHttpGetExecutor)(nil).Execute),
		ctx, path, params,
	)
}
