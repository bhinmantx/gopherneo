package gopherneo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Connection struct {
	httpClient       *http.Client
	Uri              string
	AuthTokenEncoded string
	Version          string `json:"neo4j_version"`
	NodeURI          string `json:"node"`
	NodeLabelsURI    string `json:"node_labels"`
	CypherURI        string `json:"cypher"`
	TransactionURI   string `json:"transaction"`
	//  Extensions     interface{} `json:"extensions"`
	// RefNodeURI     string      `json:"reference_node"`
	// NodeIndexURI   string      `json:"node_index"`
	// RelIndexURI    string      `json:"relationship_index"`
	// ExtInfoURI     string      `json:"extensions_info"`
	// RelTypesURI    string      `json:"relationship_types"`
	// BatchURI       string      `json:"batch"`
}

type ErrorResponse struct {
	Errors         []ErrorMessage `json:"errors"`
	Authentication string         `json:"authentication"`
}

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TransactionStatement struct {
	Cypher      string                 `json:"statement"`
	Params      map[string]interface{} `json:"parameters"`
	ResultTypes []string               `json:"resultDataContents"`
}

type TransactionResponse struct {
	Results []TransactionResult `json:"results"`
	Errors  []TransactionError  `json:"errors"`
}

type TransactionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// represents a transactional response result
type TransactionResult struct {
	Columns []string                        `json:"columns"`
	Data    []map[string][]*json.RawMessage `json:"data"`
}

type CypherResult struct {
	ColumnNames []string
	Rows        [][]*json.RawMessage
}

func NewConnection(baseUri string) (c *Connection, err error) {
	return NewConnectionWithToken(baseUri, "")
}

// get the Neo4j "service root"
// http://docs.neo4j.org/chunked/stable/rest-api-service-root.html
func NewConnectionWithToken(baseUri string, token string) (c *Connection, err error) {

	rootUri := fmt.Sprintf("%v/db/data/", baseUri) // WARNING: stupid, but trailing '/' is req

	c = &Connection{httpClient: &http.Client{}, Uri: rootUri}
	if len(token) > 0 {
		s := fmt.Sprintf(":%s", token)
		c.AuthTokenEncoded = base64.StdEncoding.EncodeToString([]byte(s))
	}

	// prepare request
	req, err := http.NewRequest("GET", rootUri, nil)
	if err != nil {
		return
	}
	c.addDefaultHeaders(req)

	// perform request
	data, err := c.performRequest(req) // gets []byte

	// check for errors
	e := &ErrorResponse{}
	err = json.Unmarshal(data, &e)
	if err != nil {
		return
	}
	if len(e.Errors) > 0 {
		err = fmt.Errorf("%s: '%s'", e.Errors[0].Code, e.Errors[0].Message)
		return
	}

	// unmarshal to Connection obj
	err = json.Unmarshal(data, &c)
	if err != nil {
		return
	}
	return
}

// ExecuteCypher will return a slice of "rows", each "row" is a []*json.RawMessage representing
// a slice of node properties that the user can unmarshal themselves
func (c *Connection) ExecuteCypher(cypher string, params *map[string]interface{}) (cr CypherResult, err error) {

	statement := &TransactionStatement{
		Cypher:      cypher,
		Params:      *params,
		ResultTypes: []string{"ROW"},
	}

	// create a new transaction for one single statement
	// http://neo4j.com/docs/stable/rest-api-transactional.html#rest-api-begin-and-commit-a-transaction-in-one-request
	transaction := struct {
		Statements []*TransactionStatement `json:"statements"`
	}{}
	transaction.Statements = []*TransactionStatement{statement}

	// prepare request
	reqData, err := json.Marshal(transaction)
	if err != nil {
		return
	}
	reqBuf := bytes.NewBuffer(reqData)
	uri := joinPath([]string{c.TransactionURI, "commit"})
	req, err := http.NewRequest("POST", uri, reqBuf)
	if err != nil {
		return
	}
	c.addDefaultHeaders(req)

	// make request
	data, err := c.performRequest(req)
	if err != nil {
		return
	}
	resp := &TransactionResponse{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return
	}
	if len(resp.Errors) > 0 {
		err = fmt.Errorf("%v: %v", resp.Errors[0].Code, resp.Errors[0].Message)
		return
	}
	if len(resp.Results) == 0 {
		return
	}

	// expecting only one result, since it's a single statement transaction
	tr := resp.Results[0] // TransactionResult

	// copy cols and rows into a CypherResult
	cr.ColumnNames = tr.Columns
	cr.Rows = make([][]*json.RawMessage, len(tr.Data))
	for i, rType := range tr.Data {
		if val, ok := rType["row"]; ok {
			cr.Rows[i] = val
		}
	}
	return
}

func (c *Connection) performRequest(req *http.Request) (data []byte, err error) {
	// perform request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	// get bytes from body
	data, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return
	}
	return
}

func (c *Connection) addDefaultHeaders(req *http.Request) {
	// add headers used in all Neo4j requests
	req.Header.Add("Accept", "application/json; charset=UTF-8")
	req.Header.Add("Content-Type", "application/json")
	if len(c.AuthTokenEncoded) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Basic realm=\"Neo4j\" %s", c.AuthTokenEncoded))
	}
}
