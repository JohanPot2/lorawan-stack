// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package crypto

import (
	"crypto/aes"
	"encoding/binary"
	"errors"
)

var iv = [8]byte{0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6, 0xA6}

func concat(a, b [8]byte) []byte {
	c := make([]byte, aes.BlockSize)
	copy(c[:8], a[:])
	copy(c[8:], b[:])
	return c
}
func xor(b [8]byte, t uint64) (c [8]byte) {
	val := binary.BigEndian.Uint64(b[:]) ^ t
	binary.BigEndian.PutUint64(c[:], val)
	return
}
func msb(b [16]byte) (c [8]byte) { copy(c[:], b[:8]); return }
func lsb(b [16]byte) (c [8]byte) { copy(c[:], b[8:]); return }

// WrapKey implements the RFC 3394 Wrap algorithm
func WrapKey(plaintext, kek []byte) (ciphertext []byte, err error) {
	if len(plaintext)%8 != 0 {
		return nil, errors.New("pkg/keywrap: invalid plaintext length")
	}

	var n = len(plaintext) / 8
	if n < 2 {
		return nil, errors.New("pkg/keywrap: no key present")
	}

	cipher, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}

	// Set A to initial value
	var a = iv

	// Fill R blocks
	var r = make([][8]byte, n)
	for i := 0; i < n; i++ {
		copy(r[i][:], plaintext[i*8:(i+1)*8])
	}

	// Run the algorithm
	for j := 0; j <= 5; j++ {
		for i := 1; i <= n; i++ {
			var b [aes.BlockSize]byte
			cipher.Encrypt(b[:], concat(a, r[i-1]))
			a = xor(msb(b), uint64((n*j)+i))
			r[i-1] = lsb(b)
		}
	}

	// Build the result
	ciphertext = make([]byte, 0, 8*(n+1))
	ciphertext = append(ciphertext, a[:]...)
	for i := 0; i < n; i++ {
		ciphertext = append(ciphertext, r[i][:]...)
	}

	return ciphertext, nil
}

// UnwrapKey implements the RFC 3394 Unwrap algorithm
func UnwrapKey(ciphertext, kek []byte) (plaintext []byte, err error) {
	if len(ciphertext)%8 != 0 {
		return nil, errors.New("pkg/keywrap: invalid ciphertext length")
	}

	var n = (len(ciphertext) / 8) - 1
	if n < 2 {
		return nil, errors.New("pkg/keywrap: no key present")
	}

	cipher, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}

	// Set A to C[0]
	var a [8]byte
	copy(a[:], ciphertext[:8])

	// Fill R blocks
	var r = make([][8]byte, n)
	for i := 0; i < n; i++ {
		copy(r[i][:], ciphertext[(i+1)*8:(i+2)*8])
	}

	// Run the algorithm
	for j := 5; j >= 0; j-- {
		for i := n; i >= 1; i-- {
			var b [aes.BlockSize]byte
			cipher.Decrypt(b[:], concat(xor(a, uint64(n*j+i)), r[i-1]))
			a = msb(b)
			r[i-1] = lsb(b)
		}
	}

	// Check for corruption
	if a != iv {
		return nil, errors.New("pkg/keywrap: corrupt key data")
	}

	// Build the result
	plaintext = make([]byte, 0, 8*n)
	for i := 0; i < n; i++ {
		plaintext = append(plaintext, r[i][:]...)
	}

	return plaintext, nil
}
