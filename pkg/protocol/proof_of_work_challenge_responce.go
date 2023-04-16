package protocol

import "encoding/binary"

type ProofOfWorkChallengeResponce struct {
	Nonce uint64
}

func NewProofOfWorkChallengeResponce(nonce uint64) *ProofOfWorkChallengeResponce {
	return &ProofOfWorkChallengeResponce{
		Nonce: nonce,
	}
}

func (m *ProofOfWorkChallengeResponce) Size() uint64 {
	return uint64(SizeOfUint64InBytes)
}

func (m *ProofOfWorkChallengeResponce) MessageType() MessageType {
	return MessagaTypeProofOfWorkChallengeResponce
}

func (m *ProofOfWorkChallengeResponce) Serialize(buf []byte) {
	binary.BigEndian.PutUint64(buf, m.Nonce)
}

func (m *ProofOfWorkChallengeResponce) Deserialize(buf []byte) {
	m.Nonce = binary.BigEndian.Uint64(buf)
}
