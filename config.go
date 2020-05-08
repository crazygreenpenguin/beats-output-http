package http

import (
	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"net/url"
)

type config struct {
	URL            string       `config:"url"`
	Codec          codec.Config `config:"codec"`
	RemoveMetaData bool         `config:"remove_metadata"`
}

var (
	defaultConfig = config{
		URL:            "http://127.0.0.1:8090/message",
		RemoveMetaData: false,
	}
)

func (c *config) Validate() error {
	_, err := url.ParseRequestURI(c.URL)

	return err
}
