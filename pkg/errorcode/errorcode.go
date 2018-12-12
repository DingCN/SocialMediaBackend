package errorcode

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// type error interface {
// 	Error() string
// }

// type ErrInvalidUsername struct {
// 	message string
// }
// type ErrUsernameTaken struct {
// 	message string
// }
// type ErrInvalidPassword struct {
// 	message string
// }

// func (e *ErrInvalidUsername) Error() string {
// 	return fmt.Sprintf("Username Invalid")
// }
// func (e *ErrUsernameTaken) Error() string {
// 	return fmt.Sprintf("Username already taken")
// }
// func (e *ErrInvalidPassword) Error() string {
// 	return fmt.Sprintf("Password Invalid")
// }

var (
	ErrInvalidUsername   = status.New(codes.InvalidArgument, "Storage: Username Invalid").Err()
	ErrInvalidPassword   = status.New(codes.InvalidArgument, "Storage: Password Invalid").Err()
	ErrUsernameTaken     = status.New(codes.FailedPrecondition, "Storage: ErrUsernameTaken").Err()
	ErrUserNotExist      = status.New(codes.FailedPrecondition, "Storage: User does not exist").Err()
	ErrIncorrectPassword = status.New(codes.FailedPrecondition, "Module: ErrIncorrectPassword").Err()
	ErrRPCConnectionLost = status.New(codes.FailedPrecondition, "Module: ErrRPCConnectionLost").Err()
	errStringToError     = map[string]error{
		ErrorDesc(ErrInvalidUsername):   ErrInvalidUsername,
		ErrorDesc(ErrInvalidPassword):   ErrInvalidPassword,
		ErrorDesc(ErrUsernameTaken):     ErrUsernameTaken,
		ErrorDesc(ErrUserNotExist):      ErrUserNotExist,
		ErrorDesc(ErrIncorrectPassword): ErrIncorrectPassword,
		ErrorDesc(ErrRPCConnectionLost): ErrRPCConnectionLost,
	}
)

type Errorcode struct {
	code codes.Code
	desc string
}

// Code returns grpc/codes.Code.
// TODO: define clientv3/codes.Code.
func (e Errorcode) Code() codes.Code {
	return e.code
}

func (e Errorcode) Error() string {
	return e.desc
}

func Error(err error) error {
	if err == nil {
		return nil
	}
	verr, ok := errStringToError[ErrorDesc(err)]
	if !ok { // not gRPC error
		return err
	}
	ev, ok := status.FromError(verr)
	var desc string
	if ok {
		desc = ev.Message()
	} else {
		desc = verr.Error()
	}
	return Errorcode{code: ev.Code(), desc: desc}
}

func ErrorDesc(err error) string {
	if s, ok := status.FromError(err); ok {
		return s.Message()
	}
	return err.Error()
}
