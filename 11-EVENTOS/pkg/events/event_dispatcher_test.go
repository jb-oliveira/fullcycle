package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	PayLoad interface{}
}

func (t TestEvent) GetName() string {
	return t.Name
}
func (t TestEvent) GetPayload() interface{} {
	return t.PayLoad
}
func (t TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	Id string
}

func (t TestEventHandler) Handle(event EventInterface, params ...interface{}) {
	//fmt.Println("TestEventHandler")
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.event = TestEvent{
		Name:    "TestEvent",
		PayLoad: "TestPayload",
	}
	suite.event2 = TestEvent{
		Name:    "TestEvent2",
		PayLoad: "TestPayload2",
	}
	suite.handler = TestEventHandler{
		Id: "TestEventHandler",
	}
	suite.handler2 = TestEventHandler{
		Id: "TestEventHandler2",
	}
	suite.handler3 = TestEventHandler{
		Id: "TestEventHandler3",
	}
	suite.eventDispatcher = NewEventDispatcher()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.handlers[suite.event.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][1])

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

}
func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Equal(ErrHandlerAlreadyRegistered, err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	suite.eventDispatcher.Clear()

	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
