// Code generated by fastssz. DO NOT EDIT.
// Hash: 6753f14e9ee5a735aa5d2d414c2bdf106d9137c9f93c28bae7d5df850e44abdd
// Version: 0.1.2
package types

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the ExitMessage object
func (e *ExitMessage) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(e)
}

// MarshalSSZTo ssz marshals the ExitMessage object to a target array
func (e *ExitMessage) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(52)

	// Field (0) 'ValidatorPubKey'
	if size := len(e.ValidatorPubKey); size != 48 {
		err = ssz.ErrBytesLengthFn("ExitMessage.ValidatorPubKey", size, 48)
		return
	}
	dst = append(dst, e.ValidatorPubKey...)

	// Offset (1) 'Message'
	dst = ssz.WriteOffset(dst, offset)
	if e.Message == nil {
		e.Message = new(phase0.VoluntaryExit)
	}
	offset += e.Message.SizeSSZ()

	// Field (1) 'Message'
	if dst, err = e.Message.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the ExitMessage object
func (e *ExitMessage) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 52 {
		return ssz.ErrSize
	}

	tail := buf
	var o1 uint64

	// Field (0) 'ValidatorPubKey'
	if cap(e.ValidatorPubKey) == 0 {
		e.ValidatorPubKey = make([]byte, 0, len(buf[0:48]))
	}
	e.ValidatorPubKey = append(e.ValidatorPubKey, buf[0:48]...)

	// Offset (1) 'Message'
	if o1 = ssz.ReadOffset(buf[48:52]); o1 > size {
		return ssz.ErrOffset
	}

	if o1 < 52 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (1) 'Message'
	{
		buf = tail[o1:]
		if e.Message == nil {
			e.Message = new(phase0.VoluntaryExit)
		}
		if err = e.Message.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the ExitMessage object
func (e *ExitMessage) SizeSSZ() (size int) {
	size = 52

	// Field (1) 'Message'
	if e.Message == nil {
		e.Message = new(phase0.VoluntaryExit)
	}
	size += e.Message.SizeSSZ()

	return
}

// HashTreeRoot ssz hashes the ExitMessage object
func (e *ExitMessage) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(e)
}

// HashTreeRootWith ssz hashes the ExitMessage object with a hasher
func (e *ExitMessage) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'ValidatorPubKey'
	if size := len(e.ValidatorPubKey); size != 48 {
		err = ssz.ErrBytesLengthFn("ExitMessage.ValidatorPubKey", size, 48)
		return
	}
	hh.PutBytes(e.ValidatorPubKey)

	// Field (1) 'Message'
	if err = e.Message.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the ExitMessage object
func (e *ExitMessage) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(e)
}

// MarshalSSZ ssz marshals the SignedExitMessage object
func (s *SignedExitMessage) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(s)
}

// MarshalSSZTo ssz marshals the SignedExitMessage object to a target array
func (s *SignedExitMessage) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(69)

	// Field (0) 'Signature'
	dst = append(dst, s.Signature[:]...)

	// Offset (1) 'Message'
	dst = ssz.WriteOffset(dst, offset)
	offset += s.Message.SizeSSZ()

	// Field (1) 'Message'
	if dst, err = s.Message.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the SignedExitMessage object
func (s *SignedExitMessage) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 69 {
		return ssz.ErrSize
	}

	tail := buf
	var o1 uint64

	// Field (0) 'Signature'
	copy(s.Signature[:], buf[0:65])

	// Offset (1) 'Message'
	if o1 = ssz.ReadOffset(buf[65:69]); o1 > size {
		return ssz.ErrOffset
	}

	if o1 < 69 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (1) 'Message'
	{
		buf = tail[o1:]
		if err = s.Message.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the SignedExitMessage object
func (s *SignedExitMessage) SizeSSZ() (size int) {
	size = 69

	// Field (1) 'Message'
	size += s.Message.SizeSSZ()

	return
}

// HashTreeRoot ssz hashes the SignedExitMessage object
func (s *SignedExitMessage) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(s)
}

// HashTreeRootWith ssz hashes the SignedExitMessage object with a hasher
func (s *SignedExitMessage) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Signature'
	hh.PutBytes(s.Signature[:])

	// Field (1) 'Message'
	if err = s.Message.HashTreeRootWith(hh); err != nil {
		return
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the SignedExitMessage object
func (s *SignedExitMessage) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(s)
}
