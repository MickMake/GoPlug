package Return

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
)

func (e Error) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, e.prefix, e.when, e.err, e.warning)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (e *Error) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &e.prefix, &e.when, &e.err, &e.warning)
	return err
}

func (e *Error) GobNewEncoder(network *io.Writer) {
	// Create an encoder and send a value.
	enc := gob.NewEncoder(*network)
	err := enc.Encode(Error{
		prefix:  e.prefix,
		when:    e.when,
		err:     e.err,
		warning: e.warning,
	})
	if err != nil {
		log.Fatal("encode:", err)
	}
}

func (e *Error) GobNewDecoder(network *io.Reader) {
	// Create a decoder and receive a value.
	dec := gob.NewDecoder(*network)
	var v Error
	err := dec.Decode(&v)
	if err != nil {
		log.Fatal("decode:", err)
	}
	*e = v
}
