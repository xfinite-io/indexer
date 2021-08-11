package importer

import (
	"bytes"
	"fmt"
	"encoding/json"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"

	"github.com/algorand/indexer/idb"
	"github.com/algorand/indexer/types"
	"github.com/google/uuid"
)

// Importer is used to import blocks into an idb.IndexerDb object.
type Importer interface {
	ImportBlock(blockbytes []byte) (txCount int, err error)
	ImportDecodedBlock(block *types.EncodedBlockCert) (txCount int, err error)
}

type dbImporter struct {
	db idb.IndexerDb
}

type noteField struct {
	org string
}

// TypeEnumMap is used to convert type strings into idb types.
/*
var TypeEnumMap = map[string]int{
	"pay":    idb.TypeEnumPay,
	"keyreg": idb.TypeEnumKeyreg,
	"acfg":   idb.TypeEnumAssetConfig,
	"axfer":  idb.TypeEnumAssetTransfer,
	"afrz":   idb.TypeEnumAssetFreeze,
	"appl":   idb.TypeEnumApplication,
}

// TypeEnumString is used for an error message somewhere.
var TypeEnumString string

func init() {
	TypeEnumString = util.KeysStringInt(TypeEnumMap)
}
*/
var zeroAddr = [32]byte{}

func participate(participants [][]byte, addr []byte) [][]byte {
	if len(addr) == 0 || bytes.Equal(addr, zeroAddr[:]) {
		return participants
	}
	for _, pp := range participants {
		if bytes.Equal(pp, addr) {
			return participants
		}
	}
	return append(participants, addr)
}

// ImportBlock processes a block and adds it to the IndexerDb
func (imp *dbImporter) ImportBlock(blockbytes []byte) (txCount int, err error) {
	var blockContainer types.EncodedBlockCert
	err = msgpack.Decode(blockbytes, &blockContainer)
	if err != nil {
		return txCount, fmt.Errorf("error decoding blockbytes, %v", err)
	}
	return imp.ImportDecodedBlock(&blockContainer)
}

// ImportBlock processes a block and adds it to the IndexerDb
func (imp *dbImporter) ImportDecodedBlock(blockContainer *types.EncodedBlockCert) (txCount int, err error) {
	txCount = 0
	proto, err := types.Protocol(string(blockContainer.Block.CurrentProtocol))
	if err != nil {
		return txCount, fmt.Errorf("block %d, %v", blockContainer.Block.Round, err)
	}
	err = imp.db.StartBlock()
	if err != nil {
		return txCount, fmt.Errorf("error starting block, %v", err)
	}
	block := blockContainer.Block
	round := uint64(block.Round)
	for intra := range block.Payset {
		stxn := &block.Payset[intra]
		txtypeenum, ok := idb.GetTypeEnum(stxn.Txn.Type)
		if !ok {
			return txCount,
				fmt.Errorf("%d:%d unknown txn type %v", round, intra, stxn.Txn.Type)
		}
		assetid := uint64(0)
		switch txtypeenum {
		case 3:
			assetid = uint64(stxn.Txn.ConfigAsset)
			if assetid == 0 {
				assetid = block.TxnCounter - uint64(len(block.Payset)) + uint64(intra) + 1
			}
		case 4:
			assetid = uint64(stxn.Txn.XferAsset)
		case 5:
			assetid = uint64(stxn.Txn.FreezeAsset)
		case 6:
			assetid = uint64(stxn.Txn.ApplicationID)
			if assetid == 0 {
				assetid = block.TxnCounter - uint64(len(block.Payset)) + uint64(intra) + 1
			}
		}
		if stxn.HasGenesisID {
			stxn.Txn.GenesisID = block.GenesisID
		}
		if stxn.HasGenesisHash || proto.RequireGenesisHash {
			stxn.Txn.GenesisHash = block.GenesisHash
		}
		stxnad := stxn.SignedTxnWithAD
		var note map[string]interface{}
		if stxn.Txn.Note != nil {
			//fmt.Println(string(stxn.Txn.Note))
			json.Unmarshal([]byte(stxn.Txn.Note), &note)
			fmt.Println(note["org"])
			fmt.Println(note["type"])
			if note["org"].(string) == "mzaalo" {
				participants := make([][]byte, 0, 10)
				participants = participate(participants, stxn.Txn.Sender[:])
				participants = participate(participants, stxn.Txn.Receiver[:])
				participants = participate(participants, stxn.Txn.CloseRemainderTo[:])
				participants = participate(participants, stxn.Txn.AssetSender[:])
				participants = participate(participants, stxn.Txn.AssetReceiver[:])
				participants = participate(participants, stxn.Txn.AssetCloseTo[:])
				participants = participate(participants, stxn.Txn.FreezeAccount[:])
				err = imp.db.AddTransaction(round, intra, int(txtypeenum), assetid, stxnad, participants, note["type"].(string), uuid.MustParse(note["id"].(string)), string(stxn.Txn.Note))
				if err != nil {
					return txCount, fmt.Errorf("error importing txn r=%d i=%d, %v", round, intra, err)
				}
				txCount++
			}
		}
	}
	blockheaderBytes := msgpack.Encode(block.BlockHeader)
	err = imp.db.CommitBlock(round, block.TimeStamp, block.RewardsLevel, blockheaderBytes)
	if err != nil {
		return txCount, fmt.Errorf("error committing block, %v", err)
	}
	return
}

// NewDBImporter creates a new importer object.
func NewDBImporter(db idb.IndexerDb) Importer {
	return &dbImporter{db: db}
}
