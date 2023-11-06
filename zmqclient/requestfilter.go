package zmqclient

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type RequestFilter struct {
	RequestFilter string `json:"RequestFilter"`
}

// Response.
type RequestFilterOk struct {
	Ok Filter
}
type RequestFilterErr struct {
	Err string
}

func Request_filter(req string) (Filter, error) {
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)
	if message, err = json.Marshal(RequestFilter{
		RequestFilter: req,
	}); err != nil {
		return Filter{}, err
	}

	s.Send(string(message), 0)

	if msg, err := s.Recv(0); err != nil {
		return Filter{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp RequestFilterOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Filter{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp RequestFilterErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Filter{}, err
			}

			return Filter{}, err

		}

	}

	return Filter{}, nil

}
