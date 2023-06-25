package pow

import (
	"crypto/sha1"
	"fmt"
)

const zeroByte = 48

// Hashcash - struct with fields of Hashcash
type Hashcash struct {
	Version    int
	ZerosCount int
	Date       int64
	Resource   string
	Rand       string
	Counter    int
}

func (h *Hashcash) Stringify() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter)
}

// CalcSha1 - calculates sha1 hash from given string
func CalcSha1(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func IsHashCorrect(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}
	prefix := hash[:zerosCount]
	for _, ch := range prefix {
		if ch != zeroByte {
			return false
		}
	}
	return true
}

// CalculateHashcash - calculates correct hashcash by bruteforce
func (h Hashcash) CalculateHashcash(maxIterations int) (Hashcash, error) {
	for {
		if maxIterations > 0 && h.Counter > maxIterations {
			break
		}
		header := h.Stringify()
		hash := CalcSha1(header)
		if IsHashCorrect(hash, h.ZerosCount) {
			return h, nil
		}

		h.Counter++
	}
	return h, fmt.Errorf("max iterations exceeded")
}
