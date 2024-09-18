package safelock

import (
	"crypto/cipher"
	"fmt"
	"io"

	slErrs "github.com/mrf345/safelock-cli/slErrs"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

type aeadWrapper struct {
	config    EncryptionConfig
	salt      []byte
	pwd       []byte
	errs      chan error
	counter   int
	aead      cipher.AEAD
	aeadReady bool
	aeadDone  chan bool
}

func newAeadWriter(pwd string, w io.Writer, config EncryptionConfig, errs chan error) *aeadWrapper {
	aw := &aeadWrapper{
		pwd:      []byte(pwd),
		config:   config,
		errs:     errs,
		aeadDone: make(chan bool, 2),
	}
	go aw.writeSaltAndLoad(w)
	return aw
}

func newAeadReader(pwd string, r InputReader, config EncryptionConfig, errs chan error) *aeadWrapper {
	aw := &aeadWrapper{
		pwd:      []byte(pwd),
		config:   config,
		errs:     errs,
		aeadDone: make(chan bool, 2),
	}
	aw.readSalt(r)
	go aw.loadAead()
	return aw
}

func (aw *aeadWrapper) getAead() cipher.AEAD {
	if !aw.aeadReady {
		aw.aeadReady = <-aw.aeadDone
	}

	return aw.aead
}

func (aw *aeadWrapper) writeSaltAndLoad(w io.Writer) {
	aw.salt = (<-aw.config.random)[:aw.config.SaltLength]

	if _, err := w.Write(aw.salt); err != nil {
		aw.errs <- fmt.Errorf("failed to write salt > %w", err)
		return
	}

	aw.loadAead()
}

func (aw *aeadWrapper) readSalt(r InputReader) {
	var err error
	var sought int

	if _, err = r.Seek(0, io.SeekStart); err != nil {
		aw.errs <- fmt.Errorf("failed to read input > %w", err)
		return
	}

	aw.salt = make([]byte, aw.config.SaltLength)

	if sought, err = r.Read(aw.salt); err != nil {
		aw.errs <- fmt.Errorf("failed to read salt from input > %w", err)
		return
	} else if sought != aw.config.SaltLength {
		aw.errs <- fmt.Errorf("invalid file or corrupted encryption (missing salt)")
		return
	}
}

func (aw *aeadWrapper) loadAead() {
	var err error

	key := argon2.IDKey(
		aw.pwd,
		aw.salt,
		aw.config.IterationCount,
		aw.config.MemSize,
		aw.config.Threads,
		aw.config.KeyLength,
	)

	if aw.aead, err = chacha20poly1305.NewX(key); err != nil {
		aw.errs <- fmt.Errorf("failed to create AEAD > %w", err)
		return
	}

	aw.aeadDone <- true
}

func (aw *aeadWrapper) encrypt(chunk []byte) []byte {
	idx := []byte(fmt.Sprintf("%d", aw.counter))
	aead := aw.getAead()
	nonce := (<-aw.config.random)[:aead.NonceSize()]
	aw.counter += 1
	return append(nonce, aead.Seal(nil, nonce, chunk, idx)...)
}

func (aw *aeadWrapper) decrypt(chunk []byte) (output []byte, err error) {
	aead := aw.getAead()

	if aead.NonceSize() > len(chunk) {
		err = &slErrs.ErrFailedToAuthenticate{Msg: "invalid chunk size"}
		aw.errs <- err
		return
	}

	idx := []byte(fmt.Sprintf("%d", aw.counter))
	nonce := chunk[:aead.NonceSize()]
	encrypted := chunk[aead.NonceSize():]

	if output, err = aead.Open(nil, nonce, encrypted, idx); err != nil {
		err = &slErrs.ErrFailedToAuthenticate{Msg: err.Error()}
		aw.errs <- err
		return
	}

	aw.counter += 1
	return
}
