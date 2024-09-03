package safelock

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	slErrs "github.com/mrf345/safelock-cli/slErrs"
	"golang.org/x/crypto/pbkdf2"
)

type asyncGcmItem struct {
	gcm   *cipher.AEAD
	nonce *[]byte
}

type asyncGcm struct {
	items  chan *asyncGcmItem
	pwd    string
	config EncryptionConfig
	errs   chan<- error
	done   chan bool
}

func newAsyncGcm(pwd string, config EncryptionConfig, errs chan<- error) *asyncGcm {
	ag := &asyncGcm{
		pwd:    pwd,
		config: config,
		errs:   errs,
		done:   make(chan bool, 1),
		items:  make(chan *asyncGcmItem, config.GcmBufferSize),
	}
	go ag.load()
	return ag
}

func (ag *asyncGcm) load() {
	var err error

	for {
		var gcm cipher.AEAD
		var nonce = make([]byte, ag.config.NonceLength)

		select {
		case <-ag.done:
			close(ag.items)
			close(ag.done)
			return
		default:
			if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
				ag.errs <- fmt.Errorf("failed to create random nonce > %w", err)
				return
			}

			if gcm, err = getGCM(ag.pwd, nonce, ag.config); err != nil {
				ag.errs <- fmt.Errorf("failed to create new GCM > %w", err)
				return
			}

			ag.items <- &asyncGcmItem{
				gcm:   &gcm,
				nonce: &nonce,
			}
		}
	}
}

func getGCM(pwd string, nonce []byte, config EncryptionConfig) (gcm cipher.AEAD, err error) {
	var block cipher.Block

	key := pbkdf2.Key([]byte(pwd), nonce, config.KeyLength, config.IterationCount, config.Hash)

	if block, err = aes.NewCipher(key); err != nil {
		err = fmt.Errorf("failed to create new cipher > %w", err)
		return
	}

	if gcm, err = cipher.NewGCM(block); err != nil {
		err = fmt.Errorf("failed to create new GCM > %w", err)
		return
	}

	return
}

func (ag *asyncGcm) encryptChunk(chunk []byte) []byte {
	item := <-ag.items
	return append(
		(*item.gcm).Seal(nil, *item.nonce, chunk[:], nil),
		*item.nonce...,
	)
}

func decryptChunk(chunk []byte, pwd string, limit int, config EncryptionConfig) (output []byte, err error) {
	var gcm cipher.AEAD

	encrypted := chunk[:limit-(config.NonceLength)]
	nonce := chunk[limit-(config.NonceLength) : limit]

	if gcm, err = getGCM(pwd, nonce, config); err != nil {
		err = fmt.Errorf("failed to create new GCM > %w", err)
		return
	}

	if output, err = gcm.Open(nil, nonce, encrypted, nil); err != nil {
		err = &slErrs.ErrFailedToAuthenticate{Msg: err.Error()}
		return
	}

	return
}
