package blkparser

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

type Blockchain struct {
	Path        string
	Magic       [4]byte
	CurrentFile *os.File
	CurrentId   uint32
}

func NewBlockchain(path string, magic [4]byte) (blockchain *Blockchain, err error) {
	blockchain = new(Blockchain)
	blockchain.Path = path
	blockchain.Magic = magic
	blockchain.CurrentId = 0

	f, err := os.Open(blkfilename(path, 0))
	if err != nil {
		return
	}

	blockchain.CurrentFile = f
	return
}

func (blockchain *Blockchain) NextBlock() (block *Block, err error) {
	rawblock, err := blockchain.FetchNextBlock()
	if err != nil {
		newblkfile, err2 := os.Open(blkfilename(blockchain.Path, blockchain.CurrentId+1))
		if err2 != nil {
			return nil, err2
		}
		blockchain.CurrentId++
		blockchain.CurrentFile.Close()
		blockchain.CurrentFile = newblkfile
		rawblock, err = blockchain.FetchNextBlock()
		if err != nil {
			return nil, err
		}
	}

	block, err = NewBlock(rawblock)
	if err != nil {
		return
	}

	return
}

func (blockchain *Blockchain) SkipBlock() (err error) {
	_, err = blockchain.FetchNextBlock()
	if err != nil {
		newblkfile, err2 := os.Open(blkfilename(blockchain.Path, blockchain.CurrentId+1))
		if err2 != nil {
			return err2
		}
		blockchain.CurrentId++
		blockchain.CurrentFile.Close()
		blockchain.CurrentFile = newblkfile
		_, err = blockchain.FetchNextBlock()

	}

	return
}

func (blockchain *Blockchain) FetchNextBlock() (rawblock []byte, err error) {
	buf := [4]byte{}
	_, err = blockchain.CurrentFile.Read(buf[:])
	if err != nil {
		return
	}

	if !bytes.Equal(buf[:], blockchain.Magic[:]) {
		err = errors.New("Bad magic")
		return
	}

	_, err = blockchain.CurrentFile.Read(buf[:])
	if err != nil {
		return
	}

	blocksize := uint32(blksize(buf[:]))

	rawblock = make([]byte, blocksize)

	_, err = blockchain.CurrentFile.Read(rawblock[:])
	if err != nil {
		return
	}
	return
}

// Convenience method to skip directly to the given blkfile / offset,
// you must take care of the height
func (blockchain *Blockchain) SkipTo(blkId uint32, offset int64) (err error) {
	blockchain.CurrentId = blkId
	f, err := os.Open(blkfilename(blockchain.Path, blkId))
	if err != nil {
		return
	}
	blockchain.CurrentFile = f
	_, err = blockchain.CurrentFile.Seek(offset, 0)
	return
}

func blkfilename(path string, id uint32) string {
	return fmt.Sprintf("%s/blk%05d.dat", path, id)
}

func blksize(buf []byte) (size uint64) {
	for i := 0; i < len(buf); i++ {
		size |= (uint64(buf[i]) << uint(i*8))
	}
	return
}
