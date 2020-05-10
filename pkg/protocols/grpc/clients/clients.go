package clients

import (
	"log"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
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
	serverConn := connectToGRPCServer(addr)
	return umAPI.NewUserManagementApiClient(serverConn), serverConn.Close
}

func ConnectToStudyService(addr string) (client studyAPI.StudyServiceApiClient, close func() error) {
	// Connect to user management service
	serverConn := connectToGRPCServer(addr)
	return studyAPI.NewStudyServiceApiClient(serverConn), serverConn.Close
}
