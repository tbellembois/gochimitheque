package zmqclient

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type PubchemGetProductByNameReq struct {
	PubchemGetProductByName string `json:"PubchemGetProductByName"`
}

// Response.
type PubchemGetProductByNameOk struct {
	Ok Product
}
type PubchemGetProductByNameErr struct {
	Err string
}

func PubchemGetProductByName(req string) (Product, error) {
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)
	if message, err = json.Marshal(PubchemGetProductByNameReq{
		PubchemGetProductByName: req,
	}); err != nil {
		return Product{}, err
	}

	s.Send(string(message), 0)

	if msg, err := s.Recv(0); err != nil {
		return Product{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp PubchemGetProductByNameOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Product{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp PubchemGetProductByNameErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Product{}, err
			}

			return Product{}, err

		}

	}

	return Product{}, nil

}
