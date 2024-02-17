package eventhandler

import (
	"fmt"
	"xclientapi/server"
	"xcom/global"
)

func Init() {
	server.RabbitMq().OnQueueMsgFanout("event.bus", fmt.Sprintf("%s.%s.event", global.Project, global.Id), on_event)
}

func on_event(data []byte) error {
	return nil
}
