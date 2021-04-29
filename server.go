package main

import (
	"github.com/nocensurasuritaly/reports/authentication"
	"github.com/nocensurasuritaly/reports/config"
	"github.com/nocensurasuritaly/reports/report"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type RootResolver struct{}

func (r *RootResolver) Reports() ([]report.Report, error) {
	reports := []report.Report{
		{
			Title:       "foo",
			Description: "bar",
		},
	}
	return reports, nil
}

func main() {
	schemaBytes, err := ioutil.ReadFile("./schema.graphql")
	if err != nil {
		log.Fatalln("unable to read GraphQL schema:", err)
	}
	schema := graphql.MustParseSchema(string(schemaBytes), &RootResolver{}, graphql.UseFieldResolvers())
	graphQLHandler := relay.Handler{Schema: schema}
	http.Handle("/graphql/protected", authentication.WithAuthentication(config.DecryptionKey, func(responseWriter http.ResponseWriter, request *http.Request) {
		graphQLHandler.ServeHTTP(responseWriter, request)
	}))
	err = http.ListenAndServe(":8080", nil)
	log.Fatalln("server shut down:", err)
}
