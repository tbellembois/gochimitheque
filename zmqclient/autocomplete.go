package zmqclient

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type AutocompleteReq struct {
	Autocomplete string `json:"Autocomplete"`
}

// Response.
type AutocompleteOk struct {
	Ok Autocomplete
}
type AutocompleteErr struct {
	Err string
}

func AutocompleteProductName(req string) (Autocomplete, error) {
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)
	if message, err = json.Marshal(AutocompleteReq{
		Autocomplete: req,
	}); err != nil {
		return Autocomplete{}, err
	}

	s.Send(string(message), 0)

	if msg, err := s.Recv(0); err != nil {
		return Autocomplete{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp AutocompleteOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Autocomplete{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp AutocompleteErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Autocomplete{}, err
			}

			return Autocomplete{}, err

		}

	}

	return Autocomplete{}, nil

}
