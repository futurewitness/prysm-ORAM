package lightclient

/*
	#include<stdlib.h>
*/
import "C"
import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"
	"unsafe"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/rpc/core"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v5/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v5/testing/util"
)

/*
type Blocker interface {
	Block(ctx context.Context, id []byte) (interfaces.ReadOnlySignedBeaconBlock, error)
	Blobs(ctx context.Context, id string, indices []uint64) ([]*blocks.VerifiedROBlob, *core.RpcError)
}

*/

// BeaconDbBlocker is an implementation of Blocker. It retrieves blocks from the beacon chain database.
type OramDB struct {
	Database OMapBindingSingleton
}

func bytesToUnit64(arr []byte) (uint64, error) {
	var num uint64
	err := binary.Read(bytes.NewBuffer(arr[:]), binary.LittleEndian, &num)
	return num, err
}

func NewOramDB() *OramDB {
	return &OramDB{
		NewOMapBindingSingleton(),
	}
}

func (oram_db *OramDB) Init(db_size uint) {
	oram_db.Database = NewOMapBindingSingleton()
	oram_db.Database.InitEmpty(db_size)

	for i := 0; i < int(db_size); i++ {
        oram_db.Database.Insert(uint64(i), 2)
    }	
}

func (oram_db *OramDB) Insert(key, val uint64) {
	oram_db.Database.Insert(key, val)
}

func makeDummyBlock() (interfaces.SignedBeaconBlock, error) {
	slot := primitives.Slot(params.BeaconConfig().AltairForkEpoch * primitives.Epoch(params.BeaconConfig().SlotsPerEpoch)).Add(1)

	fmt.Printf("Making Dummy block!!!!\n")
	b := util.NewBeaconBlockCapella()
	b.Block.StateRoot = bytesutil.PadTo([]byte("foo"), 32)
	b.Block.Slot = slot

	signedBlock, err := blocks.NewSignedBeaconBlock(b)

	return signedBlock, err
}

func (oram_db *OramDB) Block(ctx context.Context, id []byte) (interfaces.ReadOnlySignedBeaconBlock, error) {
	var findRes *uint64
	var x uint64
	findRes = (*uint64)(C.malloc(C.size_t(unsafe.Sizeof(x))))
	defer C.free(unsafe.Pointer(findRes))

	fmt.Printf("Looking for Block with key: %v\n", id)

	// Convert bytes to uint64 key
	converted, err := bytesToUnit64(id)
	if err != nil {
		fmt.Printf("Conversion failed for: %v\n", id)
		return nil, fmt.Errorf("Error in converting: %v\n", error.Error(err))
	}

	fmt.Printf("Converted key: %v\n", converted)
	fmt.Printf("id string: %s\n", string(id))

	dummyBlock, err := makeDummyBlock()

	foundFlag := true
	var totalDuration time.Duration
	for i := 0; i < int(10); i++ {
        startTime := time.Now()
		foundFlag = oram_db.Database.Find(uint64(i), findRes)
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		totalDuration += duration
		fmt.Println("\t\t[OMAP] time elapsed - ", duration)
    }
	fmt.Println("\t\t[OMAP] avg time elapsed - ", totalDuration / 10)


	if !foundFlag {
		fmt.Printf("\n\nERROR: Couldn't find block\n\n")
		return dummyBlock, fmt.Errorf("Error in ORAM lookup: %v\n", error.Error(err))
	} else {
		fmt.Printf("\n\nSuccess: %v\n\n", *findRes)
	}

	return dummyBlock, err
}

// Not actually used
func (db *OramDB) Blobs(ctx context.Context, id string, indices []uint64) ([]*blocks.VerifiedROBlob, *core.RpcError) {
	return nil, nil
}
