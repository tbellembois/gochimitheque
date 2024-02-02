package zmqclient

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

// Request.
type SortEmpiricalFormulaReq struct {
	SortEmpiricalFormula string `json:"SortEmpiricalFormula"`
}

// Response.
type SortEmpiricalFormulaOk struct {
	Ok string
}
type SortEmpiricalFormulaErr struct {
	Err string
}

func EmpiricalFormulaFromRawString(req string) (string, error) {
	var s *zmq.Socket

	s, _ = Zctx.NewSocket(zmq.REQ)
	defer s.Close()

	s.Connect("tcp://localhost:5556")

	var (
		message []byte
		err     error
	)

	if message, err = json.Marshal(SortEmpiricalFormulaReq{
		SortEmpiricalFormula: req,
	}); err != nil {
		return "", err
	}

	s.Send(string(message), 0)

	if msg, err := s.Recv(0); err != nil {
		return "", err
	} else {

		if msg[0:5] == `{"Ok"` {

			var resp SortEmpiricalFormulaOk
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return "", err
			}

			return resp.Ok, nil

		} else if msg[0:6] == `{"Err"` {

			var resp SortEmpiricalFormulaErr
			err = json.Unmarshal([]byte(msg), &resp)

			if err != nil {
				return "", err
			}

			return "", err

		}

	}

	return "", nil

}
