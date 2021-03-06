package kitkafka

import (
	"fmt"

	"github.com/go-kit/kit/log"
)

// KafkaLogAdapter is an log adapter bridging kitlog and kafka.
type KafkaLogAdapter struct {
	Logging log.Logger
}

// Printf implements kafka log interface.
func (k KafkaLogAdapter) Printf(s string, i ...interface{}) {
	k.Logging.Log("msg", fmt.Sprintf(s, i...))
}
