package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"bufio"
	"encoding/base64"
	"time"
	"regexp"
)

type blockSignatires struct {
	Data struct {
		AggregateBlockSignatures []string `json:"aggregateBlockSignatures"`
	} `json:"data"`
}

type blocks struct {
	Data struct {
		AggregateBlocks []string `json:"aggregateBlocks"`
	} `json:"data"`
}
type p34 struct {
	Data struct {
		Blocks []struct {
			Master struct {
				Config struct {
					P34 struct {
						List []struct {
							AdnlAddr  string `json:"adnl_addr"`
							PublicKey string `json:"public_key"`
						} `json:"list"`
						UtimeUntil int `json:"utime_until"`
					} `json:"p34"`
				} `json:"config"`
			} `json:"master"`
		} `json:"blocks"`
	} `json:"data"`
}

type seqno struct {
	Data struct {
		SeqNo []struct {
			SeqNo             int    `json:"seq_no"`
			PrevKeyBlockSeqno int    `json:"prev_key_block_seqno"`
			Typename          string `json:"__typename"`
		} `json:"seq_no"`
	} `json:"data"`
}

type p15 struct {
	Data struct {
		Blocks []struct {
			Master struct {
				Config struct {
					P15 struct {
						ElectionsEndBefore   int `json:"elections_end_before"`
						ElectionsStartBefore int `json:"elections_start_before"`
					} `json:"p15"`
				} `json:"config"`
			} `json:"master"`
		} `json:"blocks"`
	} `json:"data"`
}

func graph(q string, t int) string {
	url := os.Args[2]
	var jsonStr = []byte(q)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var response string

	switch t {
	case 0:
		var signs blocks
		json.Unmarshal([]byte(body), &signs)
		str := signs.Data.AggregateBlocks[0]
		response = str
	case 1:
		var uptime p34
		json.Unmarshal([]byte(body), &uptime)
		u := uptime.Data.Blocks[0].Master.Config.P34.UtimeUntil
		response = strconv.Itoa(u)
	case 2:
		var eTime p15
		json.Unmarshal([]byte(body), &eTime)
		e := eTime.Data.Blocks[0].Master.Config.P15.ElectionsStartBefore
		response = strconv.Itoa(e)
	case 3:
		var eTime p15
		json.Unmarshal([]byte(body), &eTime)
		e := eTime.Data.Blocks[0].Master.Config.P15.ElectionsEndBefore
		response = strconv.Itoa(e)
	}
	return response

}

func signs(q string) string {
	url := os.Args[2]
	var jsonStr = []byte(q)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var signs blockSignatires
	json.Unmarshal([]byte(body), &signs)
	str := signs.Data.AggregateBlockSignatures[0]
	return str
}

func getSignedBlocks(nid string) string {

	signedBlocks := signs(`{"query": "query{ aggregateBlockSignatures( filter: { signatures: { any: { node_id: { eq: \"` + nid + `\" } } } } fields: [ { fn: COUNT } ] ) }"}`)

	return signedBlocks
}

func getSeqno() string {
	url := os.Args[2]
	var jsonStr = []byte(`{"query": "query { seq_no: blocks( filter: { workchain_id: { eq: -1 } } orderBy: { path: \"seq_no\" direction: DESC } limit: 1 ) { seq_no prev_key_block_seqno __typename } }"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var seq seqno
	json.Unmarshal([]byte(body), &seq)
	str := seq.Data.SeqNo[0].PrevKeyBlockSeqno
	return strconv.Itoa(str)
}

func checkValidator(s string) [3]string {
	url := os.Args[2]
	adnl := getAdnl()
	oldAdnl := getOldAdnl()
	var jsonStr = []byte(`{"query": "query { blocks(filter: { seq_no: { eq: ` + s + ` } workchain_id: { eq: -1 } }) { master { config { p34 { utime_until, list { public_key, adnl_addr } } } } } }"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var res [3]string
	res[0] = "0"
	var validators p34
	json.Unmarshal([]byte(body), &validators)
	arr := validators.Data.Blocks[0].Master.Config.P34.List
	res[2] = strconv.Itoa(len(arr))
	for _, value := range arr {
		if value.AdnlAddr == oldAdnl {
			res[0] = "1"
			res[1] = value.PublicKey
		} else if value.AdnlAddr == adnl {
			res[0] = "1"
			res[1] = value.PublicKey
			w := writeToFile("/var/tmp/old_adnl", adnl)
			if w != nil {
				log.Fatal(w)
			}
		}
	}

	return res
}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func nodeIDCount(n string) string {

	publicKey := n
	salt := []byte{198, 180, 19, 72}

	decoded, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	res := append(salt, decoded...)

	h := sha256.New()
	h.Write(res)
	return hex.EncodeToString(h.Sum(nil))
}

func getOldAdnl() string {
	data, err := ioutil.ReadFile("/var/tmp/old_adnl")
	if err != nil {
		log.Fatal(err)
	}
	adnl := string(data)

	return adnl
}

func getAdnl() string {
	file, err := os.Open(os.Args[3])
	scanner := bufio.NewScanner(file)
	i := 0
	s := ""
	for scanner.Scan() {
	 if strings.Contains(scanner.Text(), "validator_adnl_key_id") {
	  if i != 0 {
	   s = scanner.Text()
	  }
	  i++
	  // return line, nil
	 }
	}
	re := regexp.MustCompile(`"(.*?)"`)
	matched := re.FindAllString(s, 2)
	key := matched[1]
	p, err := base64.StdEncoding.DecodeString(key[1 : len(key)-1])
	if err != nil {
	 // handle error
	}
	h := hex.EncodeToString(p)
   
	fmt.Println(h)
	return h
   }

func getElectionStatus(s string) string {
	u, _ := strconv.Atoi(graph(`{"query": "query { blocks(filter: { seq_no: { eq: `+s+` } workchain_id: { eq: -1 } }) { master { config { p34 { utime_until, list { public_key, adnl_addr } } } } } }"}`, 1))
	eS, _ := strconv.Atoi(graph(`{"query": "query { blocks(filter: { seq_no: { eq: `+s+` } workchain_id: { eq: -1 } }) { master { config { p15 { elections_end_before, elections_start_before} } } } }"}`, 2))
	eE, _ := strconv.Atoi(graph(`{"query": "query { blocks(filter: { seq_no: { eq: `+s+` } workchain_id: { eq: -1 } }) { master { config { p15 { elections_end_before, elections_start_before} } } } }"}`, 3))

	electionsStart := u - eS
	electoinsEnd := u - eE

	currTime := int(time.Now().Unix())

	var resp string
	if electionsStart < currTime && currTime < electoinsEnd {
		resp = "1"
	} else {
		resp = "0"
	}

	return resp

}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("No arguments were specified!\nPlease add 1 - target file 2 - network (net/main) 3 - election key file")
		return
	}
	seqno := getSeqno()
	v := checkValidator(seqno)
	eStatus := getElectionStatus(seqno)
	nodeID := nodeIDCount(v[1])
	agg := graph(`{"query": "query { aggregateBlocks( filter: { workchain_id: {eq: -1} gen_utime: {ge: 1598478058} } ) }"}`, 0)
	sign := getSignedBlocks(nodeID)

	w := writeToFile(os.Args[1], "ton_aggregated_blocks "+agg+"\nton_signed_blocks "+sign+"\nton_validator "+v[0]+"\nton_total_validators "+v[2]+"\nton_elections "+eStatus+"\n")
	if w != nil {
		log.Fatal(w)
	}
}