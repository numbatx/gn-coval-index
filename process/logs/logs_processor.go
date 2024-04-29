package logs

import (
	"github.com/numbatx/gn-coval-index"
	"github.com/numbatx/gn-coval-index/process/utility"
	"github.com/numbatx/gn-coval-index/schema"
	"github.com/numbatx/gn-core/core"
	"github.com/numbatx/gn-core/core/check"
	"github.com/numbatx/gn-core/data"
)

type logsProcessor struct {
	pubKeyConverter core.PubkeyConverter
}

// NewLogsProcessor creates a new instance of logs processor
func NewLogsProcessor(pubKeyConverter core.PubkeyConverter) (*logsProcessor, error) {
	if check.IfNil(pubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}

	return &logsProcessor{
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// ProcessLogs converts logs data to a specific structure defined by avro schema
func (lp *logsProcessor) ProcessLogs(logs []*data.LogData) []*schema.Log {
	allLogs := make([]*schema.Log, 0, len(logs))

	for _, currLog := range logs {
		processedLog := lp.processLog(currLog)
		if processedLog != nil {
			allLogs = append(allLogs, processedLog)
		}
	}

	return allLogs
}

func (lp *logsProcessor) processLog(logData *data.LogData) *schema.Log {
	if logData == nil || check.IfNil(logData.LogHandler) {
		return nil
	}

	return &schema.Log{
		ID:      []byte(logData.TxHash),
		Address: utility.EncodePubKey(lp.pubKeyConverter, logData.LogHandler.GetAddress()),
		Events:  lp.processEvents(logData.LogHandler.GetLogEvents()),
	}
}

func (lp *logsProcessor) processEvents(events []data.EventHandler) []*schema.Event {
	allEvents := make([]*schema.Event, 0, len(events))

	for _, currEvent := range events {
		processedEvent := lp.processEvent(currEvent)

		if processedEvent != nil {
			allEvents = append(allEvents, processedEvent)
		}
	}

	return allEvents
}

func (lp *logsProcessor) processEvent(event data.EventHandler) *schema.Event {
	if check.IfNil(event) {
		return nil
	}

	return &schema.Event{
		Address:    utility.EncodePubKey(lp.pubKeyConverter, event.GetAddress()),
		Identifier: event.GetIdentifier(),
		Topics:     event.GetTopics(),
		Data:       event.GetData(),
	}
}
