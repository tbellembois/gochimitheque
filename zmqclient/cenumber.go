package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type IsCeNumberReq struct {
	IsCeNumber string `json:"IsCeNumber"`
}

// Response.
type IsCeNumberOk struct {
	Ok bool
}
type IsCeNumberErr struct {
	Err string
}

func IsCeNumber(req string) (bool, error) {
	var (
		s   *zmq.Socket
		err error
	)

	if s, err = Zctx.NewSocket(zmq.REQ); err != nil {
		return false, err
	}
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
	)

	if message, err = json.Marshal(IsCeNumberReq{
		IsCeNumber: req,
	}); err != nil {
		return false, err
	}

	if _, err = s.Send(string(message), 0); err != nil {
		return false, err
	}

	if msg, err := s.Recv(0); err != nil {
		return false, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp IsCeNumberOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return false, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp IsCeNumberErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return false, err
			}

			return false, errors.New(resp.Err)

		}

	}

	return false, nil

}
