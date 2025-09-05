package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	pb "attendance/proto" // generated from attendance.proto

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// MongoDB document structure
type AttendanceRecord struct {
	ID           int64      `bson:"_id"`
	UserID       string     `bson:"user_id"`
	Username     string     `bson:"username"`
	CheckinTime  time.Time  `bson:"checkin_time"`
	CheckoutTime *time.Time `bson:"checkout_time,omitempty"`
}

type server struct {
	pb.UnimplementedAttendanceServiceServer
	collection *mongo.Collection
	ist        *time.Location
}

// --- Helper: Get Next ID ---
func (s *server) getNextID(ctx context.Context) (int64, error) {
	var latest AttendanceRecord
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})

	err := s.collection.FindOne(ctx, bson.D{}, opts).Decode(&latest)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		}
		return 0, status.Errorf(codes.Internal, "Could not fetch latest record ID: %v", err)
	}
	nextID := latest.ID + 1
	if nextID >= 1000 {
		return 0, status.Errorf(codes.ResourceExhausted, "Cannot create new record, maximum ID of 999 reached")
	}
	return nextID, nil
}

// --- gRPC Methods ---

// CheckIn
func (s *server) CheckIn(ctx context.Context, req *pb.CheckInRequest) (*pb.AttendanceRecordResponse, error) {
	log.Printf("Received CheckIn request for user_id: %v", req.UserId)
	newID, err := s.getNextID(ctx)
	if err != nil {
		return nil, err
	}

	record := AttendanceRecord{
		ID:          newID,
		UserID:      req.UserId,
		Username:    req.Username,
		CheckinTime: time.Now().UTC(),
	}
	_, err = s.collection.InsertOne(ctx, record)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not insert record: %v", err)
	}

	checkinIST := record.CheckinTime.In(s.ist).Format("2006-01-02 15:04:05 IST")

	return &pb.AttendanceRecordResponse{
		Id:            strconv.FormatInt(record.ID, 10),
		UserId:        record.UserID,
		Username:      record.Username,
		CheckinTime:   checkinIST,
		StatusMessage: "User checked in successfully.",
	}, nil
}

// CheckOut
func (s *server) CheckOut(ctx context.Context, req *pb.CheckOutRequest) (*pb.AttendanceRecordResponse, error) {
	log.Printf("Received CheckOut request for record_id: %v", req.RecordId)

	recordID, err := strconv.ParseInt(req.RecordId, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid record ID format")
	}

	update := bson.M{"$set": bson.M{"checkout_time": time.Now().UTC()}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated AttendanceRecord
	if err = s.collection.FindOneAndUpdate(ctx, bson.M{"_id": recordID}, update, opts).Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "Record not found")
		}
		return nil, status.Errorf(codes.Internal, "Could not update record: %v", err)
	}

	checkinIST := updated.CheckinTime.In(s.ist).Format("2006-01-02 15:04:05 IST")
	checkoutIST := ""
	if updated.CheckoutTime != nil {
		checkoutIST = updated.CheckoutTime.In(s.ist).Format("2006-01-02 15:04:05 IST")
	}

	return &pb.AttendanceRecordResponse{
		Id:            strconv.FormatInt(updated.ID, 10),
		UserId:        updated.UserID,
		Username:      updated.Username,
		CheckinTime:   checkinIST,
		CheckoutTime:  checkoutIST,
		StatusMessage: "User checked out successfully.",
	}, nil
}

// GetAttendance
func (s *server) GetAttendance(ctx context.Context, req *pb.GetAttendanceRequest) (*pb.AttendanceRecordResponse, error) {
	log.Printf("Received GetAttendance request for user_id: %v", req.UserId)

	filter := bson.M{"user_id": req.UserId}
	opts := options.FindOne().SetSort(bson.D{{Key: "checkin_time", Value: -1}})

	var latest AttendanceRecord
	if err := s.collection.FindOne(ctx, filter, opts).Decode(&latest); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "No records found for this user")
		}
		return nil, status.Errorf(codes.Internal, "Could not fetch record: %v", err)
	}

	checkinIST := latest.CheckinTime.In(s.ist).Format("2006-01-02 15:04:05 IST")
	checkoutIST := ""
	if latest.CheckoutTime != nil {
		checkoutIST = latest.CheckoutTime.In(s.ist).Format("2006-01-02 15:04:05 IST")
	}

	return &pb.AttendanceRecordResponse{
		Id:            strconv.FormatInt(latest.ID, 10),
		UserId:        latest.UserID,
		Username:      latest.Username,
		CheckinTime:   checkinIST,
		CheckoutTime:  checkoutIST,
		StatusMessage: "Record found.",
	}, nil
}

// GetAllAttendance
func (s *server) GetAllAttendance(ctx context.Context, req *pb.GetAllAttendanceRequest) (*pb.GetAllAttendanceResponse, error) {
	log.Println("Received GetAllAttendance request")

	cursor, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not fetch records: %v", err)
	}
	defer cursor.Close(ctx)

	var allRecords []*pb.AttendanceRecordResponse
	for cursor.Next(ctx) {
		var record AttendanceRecord
		if err := cursor.Decode(&record); err != nil {
			log.Printf("Error decoding record: %v", err)
			continue
		}

		checkinIST := record.CheckinTime.In(s.ist).Format("2006-01-02 15:04:05 IST")
		checkoutIST := ""
		if record.CheckoutTime != nil {
			checkoutIST = record.CheckoutTime.In(s.ist).Format("2006-01-02 15:04:05 IST")
		}

		allRecords = append(allRecords, &pb.AttendanceRecordResponse{
			Id:            strconv.FormatInt(record.ID, 10),
			UserId:        record.UserID,
			Username:      record.Username,
			CheckinTime:   checkinIST,
			CheckoutTime:  checkoutIST,
			StatusMessage: "Record retrieved.",
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Cursor error: %v", err)
	}

	return &pb.GetAllAttendanceResponse{Records: allRecords}, nil
}

// --- REST Gateway ---
func runHTTPGateway(grpcAddr, httpPort string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterAttendanceServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	log.Printf("✅ HTTP Gateway is listening on port %s", httpPort)
	if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
		log.Fatalf("failed to serve HTTP gateway: %v", err)
	}
}

// --- Main ---
func main() {
	// Check for go.mod
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		log.Fatalf("go.mod not found. Please run 'go mod init attendance-service' in your project root.")
	}

	// MongoDB connection
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	collection := client.Database("attendance_db").Collection("records")

	// Load IST timezone
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalf("Failed to load IST timezone: %v", err)
	}

	// Start gRPC
	grpcPort := "50051"
	httpPort := "8080"
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterAttendanceServiceServer(s, &server{collection: collection, ist: ist})
	reflection.Register(s)

	log.Printf("✅ gRPC Service is listening on %s", grpcPort)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	runHTTPGateway("localhost:"+grpcPort, httpPort)
}

// package main

// import (
// 	"context"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"time"

// 	pb "attendance/proto" // generated from attendance.proto

// 	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/credentials/insecure"
// 	"google.golang.org/grpc/reflection"
// 	"google.golang.org/grpc/status"
// )

// // MongoDB document structure
// type AttendanceRecord struct {
// 	ID           int64      `bson:"_id"`
// 	UserID       string     `bson:"user_id"`
// 	Username     string     `bson:"username"`
// 	CheckinTime  time.Time  `bson:"checkin_time"`
// 	CheckoutTime *time.Time `bson:"checkout_time,omitempty"`
// }

// type server struct {
// 	pb.UnimplementedAttendanceServiceServer
// 	collection *mongo.Collection
// 	ist        *time.Location
// }

// // getNextID finds the highest ID and returns the next one.
// func (s *server) getNextID(ctx context.Context) (int64, error) {
// 	var latest AttendanceRecord
// 	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
// 	err := s.collection.FindOne(ctx, bson.D{}, opts).Decode(&latest)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return 1, nil
// 		}
// 		return 0, status.Errorf(codes.Internal, "Could not fetch latest record ID: %v", err)
// 	}
// 	return latest.ID + 1, nil
// }

// // --- gRPC Methods ---

// func (s *server) CheckIn(ctx context.Context, req *pb.CheckInRequest) (*pb.AttendanceRecordResponse, error) {
// 	log.Printf("CheckIn request: user_id=%v username=%v", req.UserId, req.Username)

// 	newID, err := s.getNextID(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	record := AttendanceRecord{
// 		ID:          newID,
// 		UserID:      req.UserId,
// 		Username:    req.Username,
// 		CheckinTime: time.Now().UTC(),
// 	}

// 	_, err = s.collection.InsertOne(ctx, record)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Could not insert record: %v", err)
// 	}

// 	return &pb.AttendanceRecordResponse{
// 		Id:            strconv.FormatInt(record.ID, 10),
// 		UserId:        record.UserID,
// 		Username:      record.Username,
// 		CheckinTime:   record.CheckinTime.In(s.ist).Format("2006-01-02 15:04:05 IST"),
// 		StatusMessage: "User checked in successfully.",
// 	}, nil
// }

// func (s *server) CheckOut(ctx context.Context, req *pb.CheckOutRequest) (*pb.AttendanceRecordResponse, error) {
// 	log.Printf("CheckOut request: user_id=%v", req.UserId)

// 	filter := bson.M{"user_id": req.UserId, "checkout_time": bson.M{"$exists": false}}
// 	update := bson.M{"$set": bson.M{"checkout_time": time.Now().UTC()}}
// 	opts := options.FindOneAndUpdate().SetSort(bson.D{{Key: "checkin_time", Value: -1}}).SetReturnDocument(options.After)

// 	var updated AttendanceRecord
// 	if err := s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated); err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, status.Errorf(codes.NotFound, "No active check-in found for user %s", req.UserId)
// 		}
// 		return nil, status.Errorf(codes.Internal, "Could not update record: %v", err)
// 	}

// 	return &pb.AttendanceRecordResponse{
// 		Id:            strconv.FormatInt(updated.ID, 10),
// 		UserId:        updated.UserID,
// 		Username:      updated.Username,
// 		CheckinTime:   updated.CheckinTime.In(s.ist).Format("2006-01-02 15:04:05 IST"),
// 		CheckoutTime:  updated.CheckoutTime.In(s.ist).Format("2006-01-02 15:04:05 IST"),
// 		StatusMessage: "User checked out successfully.",
// 	}, nil
// }

// func (s *server) GetAttendance(ctx context.Context, req *pb.GetAttendanceRequest) (*pb.AttendanceRecordResponse, error) {
// 	log.Printf("GetAttendance request: user_id=%v", req.UserId)

// 	filter := bson.M{"user_id": req.UserId}
// 	opts := options.FindOne().SetSort(bson.D{{Key: "checkin_time", Value: -1}})

// 	var latest AttendanceRecord
// 	if err := s.collection.FindOne(ctx, filter, opts).Decode(&latest); err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, status.Errorf(codes.NotFound, "No records found for this user")
// 		}
// 		return nil, status.Errorf(codes.Internal, "Could not fetch record: %v", err)
// 	}

// 	resp := &pb.AttendanceRecordResponse{
// 		Id:            strconv.FormatInt(latest.ID, 10),
// 		UserId:        latest.UserID,
// 		Username:      latest.Username,
// 		CheckinTime:   latest.CheckinTime.In(s.ist).Format("2006-01-02 15:04:05 IST"),
// 		StatusMessage: "Record found.",
// 	}
// 	if latest.CheckoutTime != nil {
// 		resp.CheckoutTime = latest.CheckoutTime.In(s.ist).Format("2006-01-02 15:04:05 IST")
// 	}
// 	return resp, nil
// }

// func (s *server) GetAllAttendance(ctx context.Context, req *pb.GetAllAttendanceRequest) (*pb.GetAllAttendanceResponse, error) {
// 	log.Println("GetAllAttendance request")

// 	cursor, err := s.collection.Find(ctx, bson.D{})
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "Could not fetch records: %v", err)
// 	}
// 	defer cursor.Close(ctx)

// 	var records []*pb.AttendanceRecordResponse
// 	for cursor.Next(ctx) {
// 		var r AttendanceRecord
// 		if err := cursor.Decode(&r); err != nil {
// 			continue
// 		}
// 		rec := &pb.AttendanceRecordResponse{
// 			Id:            strconv.FormatInt(r.ID, 10),
// 			UserId:        r.UserID,
// 			Username:      r.Username,
// 			CheckinTime:   r.CheckinTime.In(s.ist).Format("2006-01-02 15:04:05 IST"),
// 			StatusMessage: "Record retrieved.",
// 		}
// 		if r.CheckoutTime != nil {
// 			rec.CheckoutTime = r.CheckoutTime.In(s.ist).Format("2006-01-02 15:04:05 IST")
// 		}
// 		records = append(records, rec)
// 	}
// 	return &pb.GetAllAttendanceResponse{Records: records}, nil
// }

// // --- REST Gateway ---
// func runHTTPGateway(grpcAddr, httpPort string) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	mux := runtime.NewServeMux()
// 	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
// 	if err := pb.RegisterAttendanceServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
// 		log.Fatalf("failed to register gateway: %v", err)
// 	}
// 	log.Printf("✅ HTTP Gateway listening on %s", httpPort)
// 	if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
// 		log.Fatalf("failed to serve HTTP gateway: %v", err)
// 	}
// }

// func main() {
// 	mongoURI := os.Getenv("MONGO_URI")
// 	if mongoURI == "" {
// 		mongoURI = "mongodb://localhost:27017"
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
// 	if err != nil {
// 		log.Fatalf("Failed to connect to MongoDB: %v", err)
// 	}
// 	collection := client.Database("attendance_db").Collection("records")

// 	ist, _ := time.LoadLocation("Asia/Kolkata")

// 	grpcPort := "50051"
// 	httpPort := "8080"

// 	lis, err := net.Listen("tcp", ":"+grpcPort)
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	s := grpc.NewServer()
// 	pb.RegisterAttendanceServiceServer(s, &server{collection: collection, ist: ist})
// 	reflection.Register(s)

// 	log.Printf("✅ gRPC Service listening on %s", grpcPort)
// 	go func() {
// 		if err := s.Serve(lis); err != nil {
// 			log.Fatalf("failed to serve gRPC: %v", err)
// 		}
// 	}()

// 	runHTTPGateway("localhost:"+grpcPort, httpPort)
// }
