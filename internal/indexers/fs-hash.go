package indexers

import (
	"crypto/md5"
	"df/internal/nodes"
	"df/internal/sources"
	"encoding/hex"
	"log"
	"math"
	"os"
)

type HashFsIndexer struct {
	percentage     float64
	suppressErrors bool
}

func NewHashFsIndexer(readPercentage float64, suppressErrors bool) *HashFsIndexer {
	if readPercentage > 0.5 || readPercentage < 0 {
		panic("readPercentage of NewHashFsIndexer must be 0 > readPercentage > 0.5")
	}

	return &HashFsIndexer{
		percentage:     readPercentage,
		suppressErrors: suppressErrors,
	}
}

func (ix *HashFsIndexer) Index(node *nodes.Node[sources.FsData]) (interface{}, error) {
	type segment struct {
		start int64
		end   int64
	}

	f, err := os.Open(node.Payload.Path)
	if err != nil {
		if ix.suppressErrors {
			log.Println(err)
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		if ix.suppressErrors {
			log.Println(err)
			return nil, nil
		}
		return nil, err
	}

	size := info.Size()
	segments := []segment{
		{0, int64(math.Max(math.Round(float64(size)*ix.percentage), float64(size)))},
		{
			int64(math.Max(math.Round(float64(size)*0.5), float64(size))),
			int64(math.Max(math.Round(float64(size)*(0.5+ix.percentage)), float64(size))),
		},
		{
			int64(math.Max(math.Round(float64(size)*(1.0-ix.percentage)), float64(size))),
			size,
		},
	}

	hasher := md5.New()

	for _, sgmt := range segments {
		if _, err := f.Seek(sgmt.start, 0); err != nil {
			if ix.suppressErrors {
				log.Println(err)
				return nil, nil
			}
			return nil, err
		}

		buf := make([]byte, sgmt.end-sgmt.start)
		if _, err := f.Read(buf); err != nil {
			if ix.suppressErrors {
				log.Println(err)
				return nil, nil
			}
			return nil, err
		}

		if _, err := hasher.Write(buf); err != nil {
			if ix.suppressErrors {
				log.Println(err)
				return nil, nil
			}
			return nil, err
		}
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
