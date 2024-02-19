package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type IsCasNumberReq struct {
	IsCasNumber string `json:"IsCasNumber"`
}

// Response.
type IsCasNumberOk struct {
	Ok bool
}
type IsCasNumberErr struct {
	Err string
}

func IsCasNumber(req string) (bool, error) {
	var (
		s   *zmq.Socket
		err error
	)

	if s, err = Zctx.NewSocket(zmq.REQ); err != nil {
		return false, err
	}
	defer s.Close()

	if err = s.Connect("tcp://localhost:5556"); err != nil {
		return false, err
	}

	var (
		message []byte
	)

	if message, err = json.Marshal(IsCasNumberReq{
		IsCasNumber: req,
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

			return false, errors.New(resp.Err)

		}

	}

	return false, nil

}
