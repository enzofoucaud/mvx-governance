package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

const EXPLORER_API = "https://api.elrond.com"

type Governance struct {
	erd      string
	proposal string
	vote     string
	power    string
}

func main() {
	var (
		smartContract = "erd1qqqqqqqqqqqqqpgqdt9aady5jp7av97m7rqxh6e5ywyqsplz2jps5mw02n"
		page          int
		governance    []Governance
	)

	// Get total transactions count and init progress bar
	count, err := GetTransactionsCount(smartContract)
	if err != nil {
		log.Err(err).Msg("GetTransactionsCount err")
		panic(err)
	}
	bar := progressbar.Default(int64(count))

	for {
		txs, err := GetTransactionsAccounts(smartContract, strconv.Itoa(page))
		if err != nil {
			log.Err(err).Msg("GetStackingRewards() err")
			panic(err)
		}
		for _, tx := range txs {
			base64 := decodeBase64String(tx.Data)
			if len(strings.Split(base64, "@")) > 3 {
				governance = append(governance, Governance{
					erd:      tx.Sender,
					vote:     getVote(base64),
					proposal: getProposal(base64),
					power:    getPower(base64),
				})
			}
		}
		page = page + len(txs)
		// increment progress bar
		_ = bar.Add(len(txs))
		// check if all transactions are processed
		if len(txs) != 50 {
			break
		}
	}

	writeCSV(governance)
}

func decodeBase64String(base64String string) string {
	decodedString, _ := base64.StdEncoding.DecodeString(base64String)
	return string(decodedString)
}

func getProposal(decoded string) string {
	split := strings.Split(decoded, "@")
	return split[1]
}

func getVote(decoded string) string {
	split := strings.Split(decoded, "@")
	switch split[2] {
	case "":
		return "yes"
	case "01":
		return "no"
	case "03":
		return "blank"
	default:
		return ""
	}
}

func getPower(decoded string) string {
	split := strings.Split(decoded, "@")
	power, _ := strconv.ParseInt(split[3], 16, 64)
	return strconv.FormatInt(power, 10)
}

func GetTransactionsCount(erd string) (int, error) {
	var (
		url = EXPLORER_API + "/accounts/" + erd + "/transactions/count"
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Err(err).Msg("Error when client create GET request to " + url)
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("Error when client do request to " + url)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorElrond ErrorElrond
		body, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &errorElrond)
		return 0, errors.New(errorElrond.Message)
	}

	var count int
	err = json.NewDecoder(resp.Body).Decode(&count)
	if err != nil {
		log.Err(err).Msg("Error when Decode JSON")
		return 0, err
	}

	return count, nil
}

func GetTransactionsAccounts(erd, from string) ([]Transactions, error) {
	var (
		url = EXPLORER_API + "/accounts/" + erd + "/transactions?from=" + from + "&size=50&withScResults=true"
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Err(err).Msg("Error when client create GET request to " + url)
		return []Transactions{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("Error when client do request to " + url)
		return []Transactions{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorElrond ErrorElrond
		body, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &errorElrond)
		return []Transactions{}, errors.New(errorElrond.Message)
	}

	var transactions []Transactions
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	if err != nil {
		log.Err(err).Msg("Error when Decode JSON")
		return []Transactions{}, err
	}

	return transactions, nil
}

func writeCSV(governance []Governance) {
	file, err := os.Create("governance.csv")
	if err != nil {
		log.Err(err).Msg("Create file err")
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Comma = ';'

	for i, value := range governance {
		if i == 0 {
			err := writer.Write([]string{"erd", "proposal", "vote", "power"})
			if err != nil {
				log.Err(err).Msg("Write header err")
				panic(err)
			}
		}
		err := writer.Write([]string{value.erd, value.proposal, value.vote, value.power})
		if err != nil {
			log.Err(err).Msg("Write file err")
			panic(err)
		}
	}
}
