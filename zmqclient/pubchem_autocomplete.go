package zmqclient

import (
	"encoding/json"

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
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)
	if message, err = json.Marshal(PubchemAutocompleteReq{
		PubchemAutocomplete: req,
	}); err != nil {
		return PubchemAutocomplete{}, err
	}

	s.Send(string(message), 0)

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

			return PubchemAutocomplete{}, err

		}

	}

	return PubchemAutocomplete{}, nil

}
