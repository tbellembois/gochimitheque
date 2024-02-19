package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type PubchemGetProductByNameReq struct {
	PubchemGetProductByName string `json:"PubchemGetProductByName"`
}

// Response.
type PubchemGetProductByNameOk struct {
	Ok PubchemProduct
}
type PubchemGetProductByNameErr struct {
	Err string
}

func PubchemGetProductByName(req string) (PubchemProduct, error) {
	var (
		s   *zmq.Socket
		err error
	)

	if s, err = Zctx.NewSocket(zmq.REQ); err != nil {
		return PubchemProduct{}, err
	}
	defer s.Close()

	if err = s.Connect("tcp://localhost:5556"); err != nil {
		return PubchemProduct{}, err
	}

	var (
		message []byte
	)

	if message, err = json.Marshal(PubchemGetProductByNameReq{
		PubchemGetProductByName: req,
	}); err != nil {
		return PubchemProduct{}, err
	}

	if _, err = s.Send(string(message), 0); err != nil {
		return PubchemProduct{}, err
	}

	if msg, err := s.Recv(0); err != nil {
		return PubchemProduct{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp PubchemGetProductByNameOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return PubchemProduct{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp PubchemGetProductByNameErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return PubchemProduct{}, err
			}

			return PubchemProduct{}, errors.New(resp.Err)

		}

	}

	return PubchemProduct{}, nil

}
