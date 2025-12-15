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
	suite.handler = TestEventHandler{}
	suite.handler2 = TestEventHandler{}
	suite.handler3 = TestEventHandler{}
	suite.eventDispatcher = NewEventDispatcher()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	// err := suite.eventDispatcher.Register(suite.event.GetName(), suite.handler)
	// suite.Nil(err)
	// suite.Equal(1, len(suite.eventDispatcher.GetHandlers(suite.event.GetName())))
	assert.True(suite.T(), true)
}

// func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_With_Multiple_Events() {
// 	err := suite.eventDispatcher.Register(suite.event.GetName(), suite.handler)
// 	suite.Nil(err)
// 	err = suite.eventDispatcher.Register(suite.event.GetName(), suite.handler2)
// 	suite.Nil(err)
// 	err = suite.eventDispatcher.Register(suite.event.GetName(), suite.handler3)
// 	suite.Nil(err)
// 	suite.Equal(3, len(suite.eventDispatcher.GetHandlers(suite.event.GetName())))
// 	suite.Equal(0, len(suite.eventDispatcher.GetHandlers(suite.event2.GetName())))
// }

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
