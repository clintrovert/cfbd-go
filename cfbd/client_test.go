package cfbd

import (
   "context"
   "net/url"
   "os"
   "reflect"
   "testing"

   "github.com/golang/mock/gomock"
   "github.com/stretchr/testify/assert"
   "google.golang.org/protobuf/encoding/protojson"
   "google.golang.org/protobuf/types/known/structpb"
)

const (
   defaultTimeFormat = "2006-01-02T15:04:05.000Z"
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

func TestGetGames_ValidRequest_ShouldSucceed(t *testing.T) {
   tester := newTestClient(t)
   filename := "./internal/test/responses/games.json"
   bytes, _ := os.ReadFile(filename)

   tester.requestExecutor.EXPECT().
      Execute(gomock.Any(), gomock.Any(), gomock.Any()).
      Return(bytes, nil).
      Times(1)

   response, err := tester.client.GetGames(
      context.Background(), GetGamesRequest{
         Week: 1, Year: 2025, Team: "Texas", SeasonType: "regular",
      },
   )

   game := response[0]

   assert.NoError(t, err)
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
   tester := newTestClient(t)
   filename := "./internal/test/responses/games_teams.json"
   bytes, _ := os.ReadFile(filename)

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
   assert.Equal(t, len(team.Stats), 33)
   assert.Equal(t, team.Stats[0].Stat, "26")
   assert.Equal(t, team.Stats[0].Category, "firstDowns")
}

func TestGetGamesPlayers_ValidRequest_ShouldSucceed(t *testing.T) {
   tester := newTestClient(t)
   filename := "./internal/test/responses/games_players.json"
   bytes, _ := os.ReadFile(filename)

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
   assert.Equal(t, len(team.Categories), 9)

   category := team.Categories[0]
   assert.Equal(t, category.Name, "passing")
   assert.Equal(t, len(category.Types), 6)

   qbr := category.Types[5]
   assert.Equal(t, qbr.Name, "QBR")
   assert.Equal(t, len(qbr.Athletes), 3)

   qb := qbr.Athletes[0]
   assert.Equal(t, qb.Id, "4870906")
   assert.Equal(t, qb.Name, "Arch Manning")
   assert.Equal(t, qb.Stat, "81.6")
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
