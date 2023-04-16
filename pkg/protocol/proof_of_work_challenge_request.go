package protocol

type ProofOfWorkChallengeRequest struct {
	Difficulty uint8
	Token      []byte
}

func NewProofOfWorkChallengeRequest(difficulty uint8, token []byte) *ProofOfWorkChallengeRequest {
	return &ProofOfWorkChallengeRequest{
		Difficulty: difficulty,
		Token:      token,
	}
}

func (m *ProofOfWorkChallengeRequest) Size() uint64 {
	return uint64(SizeOfUint8InBytes + len(m.Token))
}

func (m *ProofOfWorkChallengeRequest) MessageType() MessageType {
	return MessagaTypeProofOfWorkChallengeRequest
}

func (m *ProofOfWorkChallengeRequest) Serialize(buf []byte) {
	buf[0] = m.Difficulty
	copy(buf[1:], m.Token)
}

func (m *ProofOfWorkChallengeRequest) Deserialize(buf []byte) {
	m.Difficulty = buf[0]
	m.Token = buf[1:]
}
