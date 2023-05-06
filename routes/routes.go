package routes

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PraveenPin/GroupService/groupModels"
	"github.com/PraveenPin/TrackingService/controllers"
	"github.com/PraveenPin/TrackingService/services"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

const (
	PORT           = ":8082"
	AUTH0_AUDIENCE = "https://pin-swipe-group/"
	AUTH0_DOMAIN   = "dev-1fjzeag2i8t5jnxs.us.auth0.com"
)

type Dispatcher struct {
}

func HomeEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world :)")
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", AUTH0_DOMAIN))

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = groupModels.Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func (r *Dispatcher) Init(ctx context.Context, client *pubsub.Client, db *dynamodb.DynamoDB) {
	log.Println("Initialize the router")
	router := mux.NewRouter()

	log.Println("Securing all endpoints with jwt middleware")
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := AUTH0_AUDIENCE
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := fmt.Sprintf("https://%s/", AUTH0_DOMAIN)
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	trackingService := services.NewTrackingService(client, ctx, db)
	trackingController := controllers.NewTrackingController(trackingService)

	router.StrictSlash(true)
	router.HandleFunc("/", HomeEndpoint).Methods("GET")
	router.Handle("/publishMessage", jwtMiddleware.Handler(http.HandlerFunc(trackingController.PublishMessage))).Methods("POST")

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	log.Println("Add the listener to port", PORT)

	http.ListenAndServe(PORT, corsWrapper.Handler(router))
}
