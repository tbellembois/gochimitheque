package zmqclient

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type PubchemGetCompoundByNameReq struct {
	PubchemGetCompoundByName string `json:"PubchemGetCompoundByName"`
}

// Response.
type PubchemGetCompoundByNameOk struct {
	Ok Compounds
}
type PubchemGetCompoundByNameErr struct {
	Err string
}

func PubchemGetCompoundByName(req string) (Compounds, error) {
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)
	if message, err = json.Marshal(PubchemGetCompoundByNameReq{
		PubchemGetCompoundByName: req,
	}); err != nil {
		return Compounds{}, err
	}

	s.Send(string(message), 0)

	if msg, err := s.Recv(0); err != nil {
		return Compounds{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp PubchemGetCompoundByNameOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Compounds{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp PubchemGetCompoundByNameErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Compounds{}, err
			}

			return Compounds{}, err

		}

	}

	return Compounds{}, nil

}
