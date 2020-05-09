package clients

import (
	"log"

	umAPI "github.com/influenzanet/user-management-service/pkg/api"
	"google.golang.org/grpc"
)

func connectToGRPCServer(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", addr, err)
	}
	return conn
}

func ConnectToUserManagement(addr string) (client umAPI.UserManagementApiClient, close func() error) {
	// Connect to user management service
	userManagementServerConn := connectToGRPCServer(addr)
	return umAPI.NewUserManagementApiClient(userManagementServerConn), userManagementServerConn.Close
}
