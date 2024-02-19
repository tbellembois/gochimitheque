package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type PubchemAutocompleteReq struct {
	PubchemAutocomplete string `json:"PubchemAutocomplete"`
}

// Response.
type PubchemAutocompleteOk struct {
	Ok PubchemAutocomplete
}
type PubchemAutocompleteErr struct {
	Err string
}

func PubchemAutocompleteProductName(req string) (PubchemAutocomplete, error) {
	var (
		s   *zmq.Socket
		err error
	)

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	if err = s.Connect("tcp://localhost:5556"); err != nil {
		return PubchemAutocomplete{}, err
	}

	var (
		message []byte
	)

	if message, err = json.Marshal(PubchemAutocompleteReq{
		PubchemAutocomplete: req,
	}); err != nil {
		return PubchemAutocomplete{}, err
	}

	if _, err = s.Send(string(message), 0); err != nil {
		return PubchemAutocomplete{}, err
	}

	if msg, err := s.Recv(0); err != nil {
		return PubchemAutocomplete{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp PubchemAutocompleteOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return PubchemAutocomplete{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp PubchemAutocompleteErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return PubchemAutocomplete{}, err
			}

			return PubchemAutocomplete{}, errors.New(resp.Err)

		}

	}

	return PubchemAutocomplete{}, nil

}
