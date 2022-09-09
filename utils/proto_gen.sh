protoc -I ../internal/transport/grpc/proto/ \
 auth.proto \
 notify.proto \
 --go-grpc_out=../internal/transport --go_out=../internal/transport