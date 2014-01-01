package blkparser

import (
    "bytes"
    "encoding/binary"
    "encoding/hex"
)

type Block struct {
    Raw []byte
    Hash string
    Version uint32
    MerkleRoot string
    BlockTime uint32
    Bits uint32
    Nonce uint32
    Size uint32
    Parent string
    Txs []*Tx
}

func NewBlock(rawblock []byte) (block *Block, err error) {
    block = new(Block)
    block.Raw = rawblock

    block.Hash = GetShaString(rawblock[:80])
    block.Version = binary.LittleEndian.Uint32(rawblock[0:4])
    if !bytes.Equal(rawblock[4:36], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) {
        block.Parent = HashString(rawblock[4:36])
    }
    block.MerkleRoot = hex.EncodeToString(rawblock[36:68])
    block.BlockTime = binary.LittleEndian.Uint32(rawblock[68:72])
    block.Bits = binary.LittleEndian.Uint32(rawblock[72:76])
    block.Nonce = binary.LittleEndian.Uint32(rawblock[76:80])
    block.Size = uint32(len(rawblock))

    txs, _ := ParseTxs(rawblock[80:])

    block.Txs = txs

    return
}
