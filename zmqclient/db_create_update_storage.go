package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type DBCreateUpdateStorageReq struct {
	DBCreateUpdateStorage []interface{} `json:"DBCreateUpdateStorage"`
}

// Response.
type DBCreateUpdateStorageOk struct {
	Ok json.RawMessage
}
type DBCreateUpdateStorageErr struct {
	Err string
}

func DBCreateUpdateStorage(storagetRawString json.RawMessage, nb_items int, identical_barecode bool) (json.RawMessage, error) {
	var (
		s   *zmq.Socket
		err error
	)

	if s, err = Zctx.NewSocket(zmq.REQ); err != nil {
		return json.RawMessage{}, err
	}
	defer s.Close()

	if err = s.Connect("tcp://localhost:5556"); err != nil {
		return json.RawMessage{}, err
	}

	var (
		message []byte
	)

	req := make([]interface{}, 0)
	req = append(req, storagetRawString)
	req = append(req, nb_items)
	req = append(req, identical_barecode)

	if message, err = json.Marshal(DBCreateUpdateStorageReq{
		DBCreateUpdateStorage: req,
	}); err != nil {
		return json.RawMessage{}, err
	}

	if _, err = s.Send(string(message), 0); err != nil {
		return json.RawMessage{}, err
	}

	if msg, err := s.Recv(0); err != nil {
		return json.RawMessage{}, err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp DBCreateUpdateStorageOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp DBCreateUpdateStorageErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return json.RawMessage{}, errors.New(resp.Err)

		}

	}

	return json.RawMessage{}, nil
}
