// Code generated by fastssz. DO NOT EDIT.
// Hash: 5a5523a7872ecb2adb360564e296c2d09bdd2d5cb471180ad271ca8e59ff7a78
// Version: 0.1.2
package types

import (
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the PartialSignatureMessages object
func (p *PartialSignatureMessages) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(p)
}

// MarshalSSZTo ssz marshals the PartialSignatureMessages object to a target array
func (p *PartialSignatureMessages) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(12)

	// Field (0) 'Type'
	dst = ssz.MarshalUint64(dst, uint64(p.Type))

	// Offset (1) 'Messages'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(p.Messages) * 136

	// Field (1) 'Messages'
	if size := len(p.Messages); size > 13 {
		err = ssz.ErrListTooBigFn("PartialSignatureMessages.Messages", size, 13)
		return
	}
	for ii := 0; ii < len(p.Messages); ii++ {
		if dst, err = p.Messages[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	return
}

// UnmarshalSSZ ssz unmarshals the PartialSignatureMessages object
func (p *PartialSignatureMessages) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 12 {
		return ssz.ErrSize
	}

	tail := buf
	var o1 uint64

	// Field (0) 'Type'
	p.Type = PartialSigMsgType(ssz.UnmarshallUint64(buf[0:8]))

	// Offset (1) 'Messages'
	if o1 = ssz.ReadOffset(buf[8:12]); o1 > size {
		return ssz.ErrOffset
	}

	if o1 < 12 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (1) 'Messages'
	{
		buf = tail[o1:]
		num, err := ssz.DivideInt2(len(buf), 136, 13)
		if err != nil {
			return err
		}
		p.Messages = make([]*PartialSignatureMessage, num)
		for ii := 0; ii < num; ii++ {
			if p.Messages[ii] == nil {
				p.Messages[ii] = new(PartialSignatureMessage)
			}
			if err = p.Messages[ii].UnmarshalSSZ(buf[ii*136 : (ii+1)*136]); err != nil {
				return err
			}
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the PartialSignatureMessages object
func (p *PartialSignatureMessages) SizeSSZ() (size int) {
	size = 12

	// Field (1) 'Messages'
	size += len(p.Messages) * 136

	return
}

// HashTreeRoot ssz hashes the PartialSignatureMessages object
func (p *PartialSignatureMessages) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(p)
}

// HashTreeRootWith ssz hashes the PartialSignatureMessages object with a hasher
func (p *PartialSignatureMessages) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Type'
	hh.PutUint64(uint64(p.Type))

	// Field (1) 'Messages'
	{
		subIndx := hh.Index()
		num := uint64(len(p.Messages))
		if num > 13 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range p.Messages {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 13)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the PartialSignatureMessages object
func (p *PartialSignatureMessages) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(p)
}

// MarshalSSZ ssz marshals the PartialSignatureMessage object
func (p *PartialSignatureMessage) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(p)
}

// MarshalSSZTo ssz marshals the PartialSignatureMessage object to a target array
func (p *PartialSignatureMessage) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'PartialSignature'
	if size := len(p.PartialSignature); size != 96 {
		err = ssz.ErrBytesLengthFn("PartialSignatureMessage.PartialSignature", size, 96)
		return
	}
	dst = append(dst, p.PartialSignature...)

	// Field (1) 'SigningRoot'
	dst = append(dst, p.SigningRoot[:]...)

	// Field (2) 'Signer'
	dst = ssz.MarshalUint64(dst, uint64(p.Signer))

	return
}

// UnmarshalSSZ ssz unmarshals the PartialSignatureMessage object
func (p *PartialSignatureMessage) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 136 {
		return ssz.ErrSize
	}

	// Field (0) 'PartialSignature'
	if cap(p.PartialSignature) == 0 {
		p.PartialSignature = make([]byte, 0, len(buf[0:96]))
	}
	p.PartialSignature = append(p.PartialSignature, buf[0:96]...)

	// Field (1) 'SigningRoot'
	copy(p.SigningRoot[:], buf[96:128])

	// Field (2) 'Signer'
	p.Signer = OperatorID(ssz.UnmarshallUint64(buf[128:136]))

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the PartialSignatureMessage object
func (p *PartialSignatureMessage) SizeSSZ() (size int) {
	size = 136
	return
}

// HashTreeRoot ssz hashes the PartialSignatureMessage object
func (p *PartialSignatureMessage) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(p)
}

// HashTreeRootWith ssz hashes the PartialSignatureMessage object with a hasher
func (p *PartialSignatureMessage) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'PartialSignature'
	if size := len(p.PartialSignature); size != 96 {
		err = ssz.ErrBytesLengthFn("PartialSignatureMessage.PartialSignature", size, 96)
		return
	}
	hh.PutBytes(p.PartialSignature)

	// Field (1) 'SigningRoot'
	hh.PutBytes(p.SigningRoot[:])

	// Field (2) 'Signer'
	hh.PutUint64(uint64(p.Signer))

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the PartialSignatureMessage object
func (p *PartialSignatureMessage) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(p)
}

// MarshalSSZ ssz marshals the SignedPartialSignatureMessage object
func (s *SignedPartialSignatureMessage) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(s)
}

// MarshalSSZTo ssz marshals the SignedPartialSignatureMessage object to a target array
func (s *SignedPartialSignatureMessage) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(108)

	// Offset (0) 'Message'
	dst = ssz.WriteOffset(dst, offset)
	offset += s.Message.SizeSSZ()

	// Field (1) 'Signature'
	if size := len(s.Signature); size != 96 {
		err = ssz.ErrBytesLengthFn("SignedPartialSignatureMessage.Signature", size, 96)
		return
	}
	dst = append(dst, s.Signature...)

	// Field (2) 'Signer'
	dst = ssz.MarshalUint64(dst, uint64(s.Signer))

	// Field (0) 'Message'
	if dst, err = s.Message.MarshalSSZTo(dst); err != nil {
		return
	}

	return
}

// UnmarshalSSZ ssz unmarshals the SignedPartialSignatureMessage object
func (s *SignedPartialSignatureMessage) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 108 {
		return ssz.ErrSize
	}

	tail := buf
	var o0 uint64

	// Offset (0) 'Message'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return ssz.ErrOffset
	}

	if o0 < 108 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (1) 'Signature'
	if cap(s.Signature) == 0 {
		s.Signature = make([]byte, 0, len(buf[4:100]))
	}
	s.Signature = append(s.Signature, buf[4:100]...)

	// Field (2) 'Signer'
	s.Signer = OperatorID(ssz.UnmarshallUint64(buf[100:108]))

	// Field (0) 'Message'
	{
		buf = tail[o0:]
		if err = s.Message.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the SignedPartialSignatureMessage object
func (s *SignedPartialSignatureMessage) SizeSSZ() (size int) {
	size = 108

	// Field (0) 'Message'
	size += s.Message.SizeSSZ()

	return
}

// HashTreeRoot ssz hashes the SignedPartialSignatureMessage object
func (s *SignedPartialSignatureMessage) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(s)
}

// HashTreeRootWith ssz hashes the SignedPartialSignatureMessage object with a hasher
func (s *SignedPartialSignatureMessage) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Message'
	if err = s.Message.HashTreeRootWith(hh); err != nil {
		return
	}

	// Field (1) 'Signature'
	if size := len(s.Signature); size != 96 {
		err = ssz.ErrBytesLengthFn("SignedPartialSignatureMessage.Signature", size, 96)
		return
	}
	hh.PutBytes(s.Signature)

	// Field (2) 'Signer'
	hh.PutUint64(uint64(s.Signer))

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the SignedPartialSignatureMessage object
func (s *SignedPartialSignatureMessage) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(s)
}
