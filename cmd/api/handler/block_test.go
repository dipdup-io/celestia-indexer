package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dipdup-io/celestia-indexer/cmd/api/handler/responses"
	"github.com/dipdup-io/celestia-indexer/internal/storage"
	"github.com/dipdup-io/celestia-indexer/internal/storage/mock"
	"github.com/dipdup-io/celestia-indexer/internal/storage/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testBlock = storage.Block{
		Id:           1,
		Hash:         []byte{0, 1, 2, 3, 4, 5, 6, 7},
		Height:       100,
		VersionBlock: "11",
		VersionApp:   "1",
		Time:         testTime,
		MessageTypes: types.NewMsgTypeBitMask(types.MsgTypeSend),
	}

	testTime = time.Date(2023, 8, 1, 1, 1, 0, 0, time.UTC)
)

// BlockTestSuite -
type BlockTestSuite struct {
	suite.Suite
	blocks  *mock.MockIBlock
	events  *mock.MockIEvent
	echo    *echo.Echo
	handler *BlockHandler
	ctrl    *gomock.Controller
}

// SetupSuite -
func (s *BlockTestSuite) SetupSuite() {
	s.echo = echo.New()
	s.echo.Validator = NewCelestiaApiValidator()
	s.ctrl = gomock.NewController(s.T())
	s.blocks = mock.NewMockIBlock(s.ctrl)
	s.events = mock.NewMockIEvent(s.ctrl)
	s.handler = NewBlockHandler(s.blocks, s.events)
}

// TearDownSuite -
func (s *BlockTestSuite) TearDownSuite() {
	s.ctrl.Finish()
	s.Require().NoError(s.echo.Shutdown(context.Background()))
}

func TestSuiteBlock_Run(t *testing.T) {
	suite.Run(t, new(BlockTestSuite))
}

func (s *BlockTestSuite) TestGet() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.blocks.EXPECT().
		ByHeight(gomock.Any(), uint64(100)).
		Return(testBlock, nil)

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var block responses.Block
	err := json.NewDecoder(rec.Body).Decode(&block)
	s.Require().NoError(err)
	s.Require().EqualValues(1, block.Id)
	s.Require().EqualValues(100, block.Height)
	s.Require().Equal("1", block.VersionApp)
	s.Require().Equal("11", block.VersionBlock)
	s.Require().Equal("0001020304050607", block.Hash)
	s.Require().Equal(testTime, block.Time)
	s.Require().Equal([]string{string(types.MsgTypeSend)}, block.MessageTypes)
}

func (s *BlockTestSuite) TestGetInvalidBlockHeight() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height")
	c.SetParamNames("height")
	c.SetParamValues("invalid")

	s.Require().NoError(s.handler.Get(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)

	var e Error
	err := json.NewDecoder(rec.Body).Decode(&e)
	s.Require().NoError(err)
	s.Contains(e.Message, "parsing")
}

func (s *BlockTestSuite) TestList() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block")

	s.blocks.EXPECT().
		List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*storage.Block{
			&testBlock,
		}, nil)

	s.Require().NoError(s.handler.List(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var blocks []responses.Block
	err := json.NewDecoder(rec.Body).Decode(&blocks)
	s.Require().NoError(err)
	s.Require().Len(blocks, 1)
	s.Require().EqualValues(1, blocks[0].Id)
	s.Require().EqualValues(100, blocks[0].Height)
	s.Require().Equal("1", blocks[0].VersionApp)
	s.Require().Equal("11", blocks[0].VersionBlock)
	s.Require().Equal("0001020304050607", blocks[0].Hash)
	s.Require().Equal(testTime, blocks[0].Time)
	s.Require().Equal([]string{string(types.MsgTypeSend)}, blocks[0].MessageTypes)
}

func (s *BlockTestSuite) TestGetEvents() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.echo.NewContext(req, rec)
	c.SetPath("/block/:height/events")
	c.SetParamNames("height")
	c.SetParamValues("100")

	s.events.EXPECT().
		ByBlock(gomock.Any(), uint64(100)).
		Return([]storage.Event{
			{
				Id:       1,
				Height:   100,
				Time:     testTime,
				Position: 2,
				Type:     types.EventTypeBurn,
				TxId:     nil,
				Data: map[string]any{
					"test": "value",
				},
			},
		}, nil)

	s.Require().NoError(s.handler.GetEvents(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var events []responses.Event
	err := json.NewDecoder(rec.Body).Decode(&events)
	s.Require().NoError(err)
	s.Require().Len(events, 1)
	s.Require().EqualValues(1, events[0].Id)
	s.Require().EqualValues(100, events[0].Height)
	s.Require().EqualValues(2, events[0].Position)
	s.Require().Equal(testTime, events[0].Time)
	s.Require().Equal(string(types.EventTypeBurn), events[0].Type)
}
