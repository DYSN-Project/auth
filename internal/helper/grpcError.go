package helper

import (
	"dysn/auth/internal/model/consts"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetGrpcUnauthenticatedError() error {
	st := status.New(codes.Unauthenticated, consts.ErrUnauthorized)

	return st.Err()
}

func IsGrpcInvalidData(err error) bool {
	return status.Code(err) == codes.InvalidArgument
}

func IsUnauthorizedData(err error) bool {
	return status.Code(err) == codes.Unauthenticated
}

func MakeGrpcValidationError(err error) error {
	return badRequestError(err, "validation error")
}

func MakeGrpcBadRequestError(err error) error {
	return badRequestError(err, "bad dto error")
}

func badRequestError(err error, message string) error {
	st := status.New(codes.InvalidArgument, message)
	v := &errdetails.BadRequest_FieldViolation{
		Field:       "toast",
		Description: err.Error(),
	}
	br := &errdetails.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, v)
	st, err = st.WithDetails(br)
	if err != nil {
		return err
	}

	return st.Err()
}
