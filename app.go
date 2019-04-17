package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"net/http"

	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/gorilla/mux"
)

// Block contains data that will be written to Blockchain
type Block struct {
	Position  int
	Data      Receipt
	Timestamp string
	Hash      string
	PreHash   string
}

// Receipt contains data for checked out week
type Receipt struct {
	ReceiptID string `json:"receipt_id"`
	Renter    string `json:"renter"`
	PayDate   string `json:"pay_date"`
	IsGenesis bool   `json:"is_genesis"`
}

// Renter contains data for a sample renter
type Renter struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	JoinDate string `json:"join_date"`
}

func (b *Block) generateHash() {
	// get string val of the Data
	bytes, _ := json.Marshal(b.Data)
	// concatenate the dataset
	data := string(b.Position) + b.Timestamp + string(bytes) + b.PreHash
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

// CreateBlock create a new Block. It requires 2 parameters: previous block and check out renter to add
func CreateBlock(preBlock *Block, receiptItem Receipt) *Block {
	block := &Block{}
	block.Position = preBlock.Position + 1
	block.Timestamp = time.Now().String()
	block.Data = receiptItem
	block.PreHash = preBlock.Hash
	block.generateHash()

	return block
}

// Blockchain is an ordered list of blocks
type Blockchain struct {
	blocks []*Block
}

// BlockChain is a global variable that will return the mutated Blockchain struct
var BlockChain *Blockchain

func (bc *Blockchain) AddBlock(data Receipt) {
	// get previous block
	preBlock := bc.blocks[len(bc.blocks)-1]
	// create new block
	block := CreateBlock(preBlock, data)
	if validBlock(block, preBlock) {
		bc.blocks = append(bc.blocks, block)
	}

}

// GenesisBlock checks for an existing block, if not, create one.
func GenesisBlock() *Block {
	return CreateBlock(&Block{}, Receipt{IsGenesis: true})
}

// NewBlockchain returns the Blockchain struct with a Genesis Block.
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func validBlock(block, preBlock *Block) bool {
	// Confirm the hashes
	if preBlock.Hash != block.PreHash {
		return false
	}

	// Confirm the block's hash is valid
	if !block.validateHash(block.Hash) {
		return false
	}

	// check the position to confirm its been incremented
	if preBlock.Position+1 != block.Position {
		return false
	}

	return true
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(BlockChain.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	// write JSOn string
	io.WriteString(w, string(jbytes))
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var receiptItem Receipt
	if err := json.NewDecoder(r.Body).Decode(&receiptItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not write Block: %v", err)
		w.Write([]byte("Could not write block"))
		return
	}

	// create block
	h := md5.New()
	receiptItem.ReceiptID = fmt.Sprintf("%x", h.Sum(nil))
	BlockChain.AddBlock(receiptItem)
	resp, err := json.MarshalIndent(receiptItem, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not marshal payload: %v", err)
		w.Write([]byte("Could not write block"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func newRenter(w http.ResponseWriter, r *http.Request) {
	var renter Renter
	if err := json.NewDecoder(r.Body).Decode(&renter); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not create: %v", err)
		w.Write([]byte("Could not create new Renter"))
		return
	}

	// We will create an ID
	h := md5.New()
	io.WriteString(h, renter.JoinDate)
	renter.ID = fmt.Sprintf("%x", h.Sum(nil))

	// send back payload
	resp, err := json.MarshalIndent(renter, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not marshal payload: %v", err)
		w.Write([]byte("Could not save renter data"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	// initialize the blockchain and store in var
	BlockChain = NewBlockchain()

	// register router
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newRenter).Methods("POST")

	// dump the state of the Blockchain to the console
	go func() {
		for _, block := range BlockChain.blocks {
			fmt.Printf("Pre.hash: %x\n", block.PreHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data: %v\n", string(bytes))
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Println()
		}
	}()

	log.Println("Listening on port 1911")

	log.Fatal(http.ListenAndServe(":1911", r))
}
