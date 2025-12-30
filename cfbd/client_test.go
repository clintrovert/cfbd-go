package cfbd

import (
   "context"
   "net/url"
   "reflect"
   "testing"

   "github.com/golang/mock/gomock"
   "google.golang.org/protobuf/encoding/protojson"
)

type testClient struct {
   client          *Client
   requestExecutor *mockRequestExecutor
}

func newTestClient(t *testing.T) *testClient {
   ctrl := gomock.NewController(t)
   exec := newMockRequestExecutor(ctrl)

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

// func TestGetGames_ValidRequest_ShouldSucceed(t *testing.T) {
//    tester := newTestClient(t)
//
//    tester.httpGetExecutor.EXPECT().
//       Execute(gomock.Any(), gomock.Any(), gomock.Any()).
//       Return("").
//       Times(1)
//
//    actual, err := tester.client.GetGames(context.Background(), GetGamesRequest{})
//    assert.NoError(t, err)
//    assert.NotNil(t, actual)
// }

// mock of request httpGet below

// mockRequestExecutor is a mock of httpGetExecutor interface.
type mockRequestExecutor struct {
   ctrl     *gomock.Controller
   recorder *mockRequestExecutorMockRecorder
}

// mockRequestExecutorMockRecorder is the mock recorder for mockRequestExecutor.
type mockRequestExecutorMockRecorder struct {
   mock *mockRequestExecutor
}

// newMockRequestExecutor creates a new mock instance.
func newMockRequestExecutor(ctrl *gomock.Controller) *mockRequestExecutor {
   mock := &mockRequestExecutor{ctrl: ctrl}
   mock.recorder = &mockRequestExecutorMockRecorder{mock}
   return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *mockRequestExecutor) EXPECT() *mockRequestExecutorMockRecorder {
   return m.recorder
}

// execute mocks base method.
func (m *mockRequestExecutor) Execute(
   ctx context.Context,
   path string,
   params url.Values,
) ([]byte, error) {
   m.ctrl.T.Helper()
   ret := m.ctrl.Call(m, "Execute", ctx, path, params)
   ret0, _ := ret[0].([]byte)
   ret1, _ := ret[1].(error)
   return ret0, ret1
}

// execute indicates an expected call of execute.
func (mr *mockRequestExecutorMockRecorder) execute(
   ctx, path, params interface{},
) *gomock.Call {
   mr.mock.ctrl.T.Helper()
   return mr.mock.ctrl.RecordCallWithMethodType(
      mr.mock,
      "Execute",
      reflect.TypeOf((*mockRequestExecutor)(nil).Execute),
      ctx,
      path,
      params,
   )
}
