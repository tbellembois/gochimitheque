package zmqclient

import (
	"encoding/json"
	"errors"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type CasbinMatchEntityHasMembersReq struct {
	CasbinMatchEntityHasMembers int64 `json:"CasbinMatchEntityHasMembers"`
}

// string, int

// Response.
type CasbinMatchEntityHasMembersOk struct {
	Ok json.RawMessage
}
type CasbinMatchEntityHasMembersErr struct {
	Err string
}

func CasbinMatchEntityHasMembers(entity_id int64) (json.RawMessage, error) {
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

	if message, err = json.Marshal(CasbinMatchEntityHasMembersReq{
		CasbinMatchEntityHasMembers: entity_id,
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

			var resp CasbinMatchEntityHasMembersOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp CasbinMatchEntityHasMembersErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return json.RawMessage{}, err
			}

			return json.RawMessage{}, errors.New(resp.Err)

		}

	}

	return json.RawMessage{}, nil
}
