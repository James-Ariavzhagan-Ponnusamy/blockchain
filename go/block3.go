package main
import "fmt"
import "bytes"
import "crypto/sha256"
import "math/big"
import "log"
import "encoding/binary"
import "encoding/gob"
import "math"
//import "strconv"

//import "github.com/dgraph-io/badger"

const Difficulty = 12
const dbPath = "C:/tmp/blocks"


/******************************************************************************/
type Block struct{
	Data 	     [] byte
	Hash 	     [] byte
	PrevHash   [] byte
	Nonce         int
}

func (b *Block) Serialize()[]byte{
	var res bytes.Buffer
	encoder := gob.NewEncoder (&res)
	err :=encoder.Encode(b)
	if ( err != nil){
		log.Panic(err)
	}
	return res.Bytes()
}

func DeSerialize(data []byte) *Block{
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if ( err != nil){
		log.Panic(err)
	}
	return &block
}

func (b *Block) DeriveHash(){
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
  b.Nonce = 0;
}

/******************************************************************************/

type BlockChain struct {
  lastHash []byte
	blocks []*Block
}

var chain1 BlockChain

func (chain *BlockChain) CreateBlock(data string, phash []byte) *Block{
	block := &Block{[]byte(data), []byte{}, phash, 0}
  fmt.Println(data)
  proof.SetBlock(block)

  block.DeriveHash()

  nonce, hash := proof.CalculateHash()
	block.Hash = hash[:]
	block.Nonce = nonce
  chain.blocks = append ( chain.blocks, block)
  chain.lastHash = block.Hash
	return block
}

func (chain *BlockChain) Genesis() {
	chain.CreateBlock("Genesis", []byte{})
}

func (chain *BlockChain) InitBlockChain(){
  chain.Genesis()
}

func (chain *BlockChain) AddBlock(data string){
    chain.CreateBlock(data, chain.lastHash)
}

/*****************************************************************************/

type ProofOfWork struct{
	block *Block
	Target *big.Int // used to calculate proof-of-work hash
}

var proof ProofOfWork

func (pw *ProofOfWork) SetBlock(block *Block){
  pw.block = block
}

func (pw *ProofOfWork) Init () {
	pw.Target = big.NewInt ( 1 )
	pw.Target.Lsh(pw.Target, (256-Difficulty)) // max value to those many bits zero
}

func (pw *ProofOfWork)ToHex ( num int64) []byte {
	buff:= new (bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if ( err != nil){
		log.Panic (err)
	}
	return buff.Bytes()
}

func (pw *ProofOfWork) InitData (nonce int) []byte{
	data := bytes.Join([][]byte{
		pw.block.PrevHash,
		pw.block.Data,
		pw.ToHex(int64(nonce)),
		pw.ToHex(int64(Difficulty)),
		}, []byte{})
	return data
}

func (pw *ProofOfWork) CalculateHash() (int, []byte){
	var intHash big.Int
	var hash[32] byte
	nonce := 0
	i := 0
	for nonce < math.MaxInt64 {
		data := pw.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("nonce hash is %x\r", hash)

		intHash.SetBytes(hash[:])
		if (intHash.Cmp(pw.Target) == -1){
			break
		}else{
			nonce ++
		}
    i++;
	}
	fmt.Printf("\n\n\r\rNumber of iterations for the hash is %d\n", i)
	return nonce, hash[:]
}

func (pw *ProofOfWork) ValidateHash(block *Block) bool{
  var intHash big.Int
	var hash[32] byte

  pw.block = block
  data := pw.InitData(block.Nonce)
  hash = sha256.Sum256(data)
  intHash.SetBytes(hash[:])
  if (intHash.Cmp(pw.Target) == -1){
    return true
  }
  return false;
}

/*****************************************************************************/

func HandleErr ( err error) {
	if ( err != nil){
  	log.Panic(err)
  }
}
/*****************************************************************************/
func main(){
  proof.Init()

	chain1.InitBlockChain()

	chain1.AddBlock("First Block after Genesis")
	chain1.AddBlock("Second Block after Genesis")
	chain1.AddBlock("Third Block after Genesis")
	fmt.Println()
	for _, block := range chain1.blocks {
    test:=block.Serialize()
    block1 := DeSerialize(test)
		fmt.Printf("Data in the Block           : %s\n", block1.Data, )
		fmt.Printf("Hash in the Block           : %x\n", block1.Hash)
		fmt.Printf("Prev Hash in the Block      : %x\n", block1.PrevHash)
    fmt.Printf("Nonce     in the Block      : %d\n", block1.Nonce)

    if (proof.ValidateHash(block1)  == true){
      fmt.Printf("Validate                    : %s\n",  "true" )
    }

    if (proof.ValidateHash(block1)  == false){
      fmt.Printf("Validate                    : %s\n",  "false" )
    }

		fmt.Println()
	}

	fmt.Println("block3.go")

}
