package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"github.com/elastic/beats/v7/libbeat/publisher"
)

func init() {
	outputs.RegisterType("http", makeHTTP)
}

type httpOutput struct {
	log            *logp.Logger
	url            string
	beat           beat.Info
	observer       outputs.Observer
	codec          codec.Codec
	client         *http.Client
	removeMetaData bool
}

// makeHTTP instantiates a new file output instance.
func makeHTTP(
	_ outputs.IndexManager,
	beat beat.Info,
	observer outputs.Observer,
	cfg *common.Config,
) (outputs.Group, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return outputs.Fail(err)
	}
	ho := &httpOutput{
		log:            logp.NewLogger("http"),
		beat:           beat,
		observer:       observer,
		url:            config.URL,
		removeMetaData: config.RemoveMetaData,
	}
	// disable bulk support in publisher pipeline
	if err := cfg.SetInt("bulk_max_size", -1, -1); err != nil {
		ho.log.Error("Disable bulk error: ", err)
	}

	if err := ho.init(beat, config); err != nil {
		return outputs.Fail(err)
	}

	return outputs.Success(-1, 0, ho)
}

func (out *httpOutput) init(beat beat.Info, c config) error {
	var err error

	out.codec, err = codec.CreateEncoder(beat, c.Codec)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	out.client = &http.Client{Transport: tr}

	out.log.Infof("Initialized http output. "+
		"url=%v codec=%v",
		c.URL, c.Codec)

	return nil
}

// Implement Client
func (out *httpOutput) Close() error {
	return nil
}

func (out *httpOutput) Publish(_ context.Context, batch publisher.Batch) error {
	defer batch.ACK()

	st := out.observer
	events := batch.Events()
	st.NewBatch(len(events))

	dropped := 0
	for i := range events {
		event := &events[i]
		if out.removeMetaData {
			if err := event.Content.Fields.Delete("@metadata"); err != nil {
				out.log.Warn("Delete @metadata error: ", err)
			}

		}

		out.log.Debugf("Message fields: %s", event.Content.Fields.StringToPrint())

		serializedEvent, err := out.codec.Encode(out.beat.Beat, &event.Content)
		if err != nil {
			if event.Guaranteed() {
				out.log.Errorf("Failed to serialize the event: %+v", err)
			} else {
				out.log.Warnf("Failed to serialize the event: %+v", err)
			}
			out.log.Debugf("Failed event: %v", event)

			dropped++
			continue
		}

		if err = out.send(serializedEvent); err != nil {
			st.WriteError(err)

			if event.Guaranteed() {
				out.log.Errorf("Writing event to http failed with: %+v", err)
			} else {
				out.log.Warnf("Writing event to http failed with: %+v", err)
			}

			dropped++
			continue
		}

		st.WriteBytes(len(serializedEvent))
	}

	st.Dropped(dropped)
	st.Acked(len(events) - dropped)

	return nil
}

func (out *httpOutput) String() string {
	return "http(" + out.url + ")"
}

func (out *httpOutput) send(data []byte) error {

	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", out.url, buf)
	if err != nil {
		return err
	}
	resp, err := out.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("bad response code: %d", resp.StatusCode))
	}

	return nil
}
