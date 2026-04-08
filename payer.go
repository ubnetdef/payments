package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ubnetdef/ScoreTrak/pkg/auth"
	"github.com/ubnetdef/ScoreTrak/pkg/storage"
	"github.com/ubnetdef/payments/report"
	authv1 "github.com/ubnetdef/scoretrakapis/gen/proto/go/scoretrak/auth/v1"
	reportv1 "github.com/ubnetdef/scoretrakapis/gen/proto/go/scoretrak/report/v1"
	"google.golang.org/grpc"
)

var conf *Config

func main() {
	var err error
	confPath := flag.String("config", "configs/config.yml", "Please enter a path to config file")
	flag.Parse()
	conf, err = NewConfig(*confPath)
	if err != nil {
		panic(err)
	}
	db, err := storage.NewDB(conf.DB)
	if err != nil {
		panic(err)
	}
	err = initTables(db)
	if err != nil {
		panic(err)
	}
	cc, err := grpc.Dial(conf.ScoreTrakAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	authClient := authv1.NewAuthServiceClient(cc)
	log.Println("Requesting authentication token")
	t, err := getAuth(authClient)
	if err != nil {
		panic(err)
	}
	token, _, err := new(jwt.Parser).ParseUnverified(t, &auth.UserClaims{})
	if err != nil {
		panic(err)
	}
	claims, ok := token.Claims.(*auth.UserClaims)
	if !ok {
		log.Fatalf("invalid token claims")
	}
	log.Println("Authentication token is valid until " + claims.ExpiresAt.Time.String())
	log.Println("Requesting report")
	reportClient := reportv1.NewReportServiceClient(cc)
	latestReport, err := getReport(reportClient, t)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Got report for round %d\n", latestReport.Round)
	current := report.PointsPerTeam(latestReport)
	previous := getMaxPointsPerTeam(db)
	diff := computePointsDiff(current, previous)
	client := http.Client{Timeout: time.Duration(conf.ClientTimeout) * time.Second}
	makeDeposit(client, diff)
	updateMaxPointsPerTeam(db, current)
}
