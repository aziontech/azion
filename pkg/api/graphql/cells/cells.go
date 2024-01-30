package cells

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/machinebox/graphql"
)

type HTTPEvent struct {
	Status int `json:"status"`
	Count  int `json:"count"`
}

type Response struct {
	HTTPEvents []HTTPEvent `json:"httpEvents"`
}

func Option1(tokenvalue string) (Response, error) {
	graphqlClient := graphql.NewClient("https://api.azionapi.net/events/graphql")

	graphqlRequest := graphql.NewRequest(`
	query Top10StatusCodes {
		httpEvents(
		  limit: 5
		  filter: {
			tsRange: { begin:"2024-01-20T10:10:10", end:"2024-01-26T10:10:10" }
		  }
		  aggregate: {count: status}
		  groupBy: [status]
		  orderBy: [count_DESC]
		  )
		{
		  status
		  count
		}
	  }
`)
	token := "Token " + tokenvalue

	graphqlRequest.Header.Set("Authorization", token)

	// var graphqlResponse interface{}
	var response Response
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &response); err != nil {
		panic(err)
	}

	PrettyPrint(response)
	// spew.Dump(response)
	// fmt.Println(response)
	return response, nil
}

// print the contents of the obj
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
