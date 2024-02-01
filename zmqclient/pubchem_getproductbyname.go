package zmqclient

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type GetProductByNameReq struct {
	GetProductByName string `json:"GetProductByName"`
}

// Response.
type GetProductByNameOk struct {
	Ok Product
}
type GetProductByNameErr struct {
	Err string
}

func GetProductByName(req string) (Product, error) {
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)
	if message, err = json.Marshal(GetProductByNameReq{
		GetProductByName: req,
	}); err != nil {
		return Product{}, err
	}

	s.Send(string(message), 0)

	if msg, err := s.Recv(0); err != nil {
		return Product{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp GetProductByNameOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Product{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp GetProductByNameErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return Product{}, err
			}

			return Product{}, err

		}

	}

	return Product{}, nil

}
