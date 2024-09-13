package validate

import (
	"encoding/json"
	"io"
	"logIngestor/logIngestor/log/logtype"
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

type Decoder struct {
	validator *Validator
}

// Decode and validate the incoing request
func NewDecoder() *Decoder {
	return &Decoder{
		validator: &Validator{},
	}
}

func (d *Decoder) Decode(r io.ReadCloser, log *logtype.Log) error {
	defer r.Close()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	if err := dec.Decode(log); err != nil {
		return err
	}

	if err := d.validator.Validate(log); err != nil {
		return err
	}

	return nil
}
