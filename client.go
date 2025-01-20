package raildata

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"

	"github.com/jtarrio/raildata/api"
	rderrors "github.com/jtarrio/raildata/errors"
)

// NewClient creates a client for the RailData API.
//
// The RailData API uses a token to access all the operations, so you need to pass one or more options
// to provide either the token or the credentials to use to generate this token.
//
// Note that you can only generate 5 tokens per day, so we highly recommend to save the current token
// so you can pass it in the next call to NewClient.
//
// Example:
//
//	// Read the token from a file
//	token, err := io.ReadFile("/path/to/token-file")
//	if err != nil { return err }
//	client, err := raildata.NewClient(
//		// Provide the token to the client so it doesn't need to get a new one
//		raildata.WithToken(string(token)),
//		// Provide the username and password to the client so it can get a new token if the old one expires
//		raildata.WithCredentials(username, password),
//		// If the token changes, save it to the file so we can use it in the future
//		raildata.WithTokenUpdateListener(func(newToken string, oldToken string) {
//			_ = io.WriteFile("/path/to/token-file", []byte(newToken))
//		}),
//	)
func NewClient(options ...Option) (Client, error) {
	s := &raildataClient{
		apiBase: getEndpoint(false),
		client:  http.DefaultClient,
	}
	for _, opt := range options {
		opt(s)
	}
	return s, nil
}

type Option func(*raildataClient)

// WithCredentials sets the username and password so the client can get a new token if the old one is invalid or expires.
func WithCredentials(username string, password string) Option {
	return func(s *raildataClient) {
		s.credentials = &credentials{
			username: username,
			password: password,
		}
	}
}

// WithToken sets the initial token to use.
func WithToken(token string) Option {
	return func(s *raildataClient) {
		s.token = token
	}
}

// WithTokenUpdateListener registers a function that is called whenever the client creates a new token.
func WithTokenUpdateListener(listener TokenUpdateListener) Option {
	return func(s *raildataClient) {
		s.tokenUpdateListeners = append(s.tokenUpdateListeners, listener)
	}
}

// TokenUpdateListener is the type of a function that is called whenever the client creates a new token.
// This function is called with the new and old token so you can ensure that the correct token is always saved
// even if there are many simultaneous token changes.
type TokenUpdateListener func(newToken string, previousToken string)

// WithTestEndpoint sets the API endpoint to the test endpoint (if true) or the production endpoint (if false).
func WithTestEndpoint(testEndpoint bool) Option {
	return WithApiBase(getEndpoint(testEndpoint))
}

// WithApiBase sets the API endpoint's base URL to the specified value.
func WithApiBase(apiBase url.URL) Option {
	return func(s *raildataClient) {
		s.apiBase = apiBase
	}
}

// WithHttpClient sets the HTTP client to use.
func WithHttpClient(client *http.Client) Option {
	return func(s *raildataClient) {
		s.client = client
	}
}

type credentials struct {
	username string
	password string
}

type raildataClient struct {
	credentials          *credentials
	apiBase              url.URL
	client               *http.Client
	token                string
	tokenMutex           sync.Mutex
	tokenUpdateListeners []TokenUpdateListener
}

func (s *raildataClient) RateLimitedMethods() RateLimitedMethods {
	return s
}

func (s *raildataClient) GetToken() string {
	s.tokenMutex.Lock()
	defer s.tokenMutex.Unlock()
	return s.token
}

func (s *raildataClient) IsValidToken(ctx context.Context) (*IsValidTokenResponse, error) {
	output, err := request(api.IsValidToken, s, ctx, &api.TokenRequest{})
	if err != nil {
		return nil, err
	}
	return ParseValidTokenResponse(output)
}

func (s *raildataClient) GetStationList(ctx context.Context) (*GetStationListResponse, error) {
	output, err := request(api.GetStationList, s, ctx, &api.TokenRequest{})
	if err != nil {
		return nil, err
	}
	return ParseGetStationsList(*output)
}

func (s *raildataClient) GetStationMsg(ctx context.Context, req *GetStationMsgRequest) (*GetStationMsgResponse, error) {
	input := &api.GetStationMsgRequest{}
	if req.LineCode != nil {
		input.Line = string(*req.LineCode)
	}
	if req.StationCode != nil {
		input.Station = string(*req.StationCode)
	}
	output, err := request(api.GetStationMSG, s, ctx, input)
	if err != nil {
		return nil, err
	}
	return ParseStationMsgsList(*output), nil
}

func (s *raildataClient) GetStationSchedule(ctx context.Context, req *GetStationScheduleRequest) (*GetStationScheduleResponse, error) {
	input := &api.GetStationScheduleRequest{
		Station: string(req.StationCode),
	}
	if req.NjtOnly {
		input.NjtOnly = "true"
	} else {
		input.NjtOnly = "false"
	}
	output, err := request(api.GetStationSchedule, s, ctx, input)
	if err != nil {
		return nil, err
	}
	return ParseDailyStationInfoList(*output)
}

func (s *raildataClient) GetTrainSchedule(ctx context.Context, req *GetTrainScheduleRequest) (*GetTrainScheduleResponse, error) {
	input := &api.GetTrainScheduleRequest{
		Station: string(req.StationCode),
	}
	output, err := request(api.GetTrainSchedule, s, ctx, input)
	if err != nil {
		return nil, err
	}
	return ParseStationInfo(output), nil
}

func (s *raildataClient) GetTrainSchedule19Records(ctx context.Context, req *GetTrainSchedule19RecordsRequest) (*GetTrainScheduleResponse, error) {
	input := &api.GetTrainSchedule19RecRequest{
		Station: string(req.StationCode),
	}
	if req.LineCode != nil {
		input.Line = string(*req.LineCode)
	}
	output, err := request(api.GetTrainSchedule19Rec, s, ctx, input)
	if err != nil {
		return nil, err
	}
	return ParseStationInfo(output), nil
}

func (s *raildataClient) GetTrainStopList(ctx context.Context, req *GetTrainStopListRequest) (*GetTrainStopListResponse, error) {
	input := &api.GetTrainStopListRequest{
		Train: req.TrainId,
	}
	output, err := request(api.GetTrainStopList, s, ctx, input)
	if err != nil {
		return nil, err
	}
	return ParseStops(output), nil
}

func (s *raildataClient) GetVehicleData(ctx context.Context) (*GetVehicleDataResponse, error) {
	output, err := request(api.GetVehicleData, s, ctx, &api.TokenRequest{})
	if err != nil {
		return nil, err
	}
	return ParseVehicleDataInfoList(*output), err
}

func getEndpoint(testEndpoint bool) url.URL {
	var u *url.URL
	var err error
	if testEndpoint {
		u, err = url.Parse("https://testraildata.njtransit.com/api/TrainData")
	} else {
		u, err = url.Parse("https://raildata.njtransit.com/api/TrainData")
	}
	if err != nil {
		panic(err)
	}
	return *u
}

func (s *raildataClient) refreshToken(ctx context.Context, oldToken string) error {
	if s.credentials == nil {
		return rderrors.MissingCredentialsError
	}

	s.tokenMutex.Lock()
	defer s.tokenMutex.Unlock()
	if s.token != oldToken {
		return nil
	}

	input := &api.GetTokenRequest{
		Username: s.credentials.username,
		Password: s.credentials.password,
	}
	output, err := api.GetToken.Request(ctx, s.client, s.apiBase, input)
	if err != nil {
		return err
	}

	if output.Authenticated != "True" {
		return rderrors.BadCredentialsError
	}
	s.token = output.UserToken
	for _, listener := range s.tokenUpdateListeners {
		go listener(output.UserToken, oldToken)
	}
	return nil
}

func request[I any, O any](method api.MethodDefinition[I, O], s *raildataClient, ctx context.Context, input *I) (*O, error) {
	token := s.GetToken()
	method.SetToken(input, token)
	out, err := method.Request(ctx, s.client, s.apiBase, input)
	if !errors.Is(err, rderrors.InvalidTokenError) {
		return out, err
	}

	err = s.refreshToken(ctx, token)
	if err != nil {
		return nil, err
	}
	token = s.GetToken()
	method.SetToken(input, token)
	return method.Request(ctx, s.client, s.apiBase, input)
}
