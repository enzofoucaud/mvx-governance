package main

type ErrorMultiversX struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type Transactions struct {
	TxHash        string      `json:"txHash"`
	GasLimit      int         `json:"gasLimit"`
	GasPrice      int         `json:"gasPrice"`
	GasUsed       int         `json:"gasUsed"`
	MiniBlockHash string      `json:"miniblockHash"`
	Nonce         int         `json:"nonce"`
	Receiver      string      `json:"receiver"`
	ReceiverShard int         `json:"receiverShard"`
	Round         int         `json:"round"`
	Sender        string      `json:"sender"`
	SenderShard   int         `json:"senderShard"`
	Signature     string      `json:"signature"`
	Status        string      `json:"status"`
	Value         string      `json:"value"`
	Fee           string      `json:"fee"`
	Timestamp     int         `json:"timestamp"`
	Data          string      `json:"data"`
	Function      string      `json:"function"`
	Action        Action      `json:"action"`
	Results       []Result    `json:"results"`
	Operations    []Operation `json:"operations"`
}
type Action struct {
	Category    string   `json:"category"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Arguments   Argument `json:"arguments"`
}

type Argument struct {
	Transfers      []Transfer `json:"transfers"`
	ProviderName   string     `json:"providerName"`
	ProviderAvatar string     `json:"providerAvatar"`
	Receiver       string     `json:"receiver"`
}

type Transfer struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	Ticker     string `json:"ticker"`
	Collection string `json:"collection"`
	Identifier string `json:"identifier"`
	Value      string `json:"value"`
	Token      string `json:"token"`
	Decimals   int    `json:"decimals"`
}

type Result struct {
	Hash          string `json:"hash"`
	Timestamp     int    `json:"timestamp"`
	Nonce         int    `json:"nonce"`
	GasLimit      int    `json:"gasLimit"`
	GasPrice      int    `json:"gasPrice"`
	Value         string `json:"value"`
	Sender        string `json:"sender"`
	Receiver      string `json:"receiver"`
	Data          string `json:"data"`
	PreTxHash     string `json:"preTxHash"`
	OriginTxHash  string `json:"originTxHash"`
	CallType      string `json:"callType"`
	MiniBlockHash string `json:"miniblockHash"`
}

type Operation struct {
	ID       string `json:"id"`
	Action   string `json:"action"`
	Type     string `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Value    string `json:"value"`
}
