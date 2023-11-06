package zmqclient

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type IsCasNumber struct {
	IsCasNumber string `json:"IsCasNumber"`
}

// Response.
type IsCasNumberOk struct {
	Ok bool
}
type IsCasNumberErr struct {
	Err string
}

func Is_cas_number(req string) (bool, error) {
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)
	if message, err = json.Marshal(IsCasNumber{
		IsCasNumber: req,
	}); err != nil {
		return false, err
	}

	s.Send(string(message), 0)

	if msg, err := s.Recv(0); err != nil {
		return false, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp IsCasNumberOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return false, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp IsCasNumberErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return false, err
			}

			return false, err

		}

	}

	return false, nil

}
