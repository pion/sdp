package sdp

type HashFunc int

const (
	HashFuncSHA1 HashFunc = iota + 1
	HashFuncSHA224
	HashFuncSHA256
	HashFuncSHA384
	HashFuncSHA512
	HashFuncMD5
	HashFuncMD2
)

const (
	hashFuncSHA1Str   = "sha-1"
	hashFuncSHA224Str = "sha-224"
	hashFuncSHA256Str = "sha-256"
	hashFuncSHA384Str = "sha-384"
	hashFuncSHA512Str = "sha-512"
	hashFuncMD5Str    = "md5"
	hashFuncMD2Str    = "md2"
)

// NewHashFunc defines a procedure for creating a new hash-func enum from a raw
// string.
func NewHashFunc(raw string) HashFunc {
	switch raw {
	case hashFuncSHA1Str:
		return HashFuncSHA1
	case hashFuncSHA224Str:
		return HashFuncSHA224
	case hashFuncSHA256Str:
		return HashFuncSHA256
	case hashFuncSHA384Str:
		return HashFuncSHA384
	case hashFuncSHA512Str:
		return HashFuncSHA512
	case hashFuncMD5Str:
		return HashFuncMD5
	case hashFuncMD2Str:
		return HashFuncMD2
	default:
		return HashFunc(unknown)
	}
}

func (t HashFunc) String() string {
	switch t {
	case HashFuncSHA1:
		return hashFuncSHA1Str
	case HashFuncSHA224:
		return hashFuncSHA224Str
	case HashFuncSHA256:
		return hashFuncSHA256Str
	case HashFuncSHA384:
		return hashFuncSHA384Str
	case HashFuncSHA512:
		return hashFuncSHA512Str
	case HashFuncMD5:
		return hashFuncMD5Str
	case HashFuncMD2:
		return hashFuncMD2Str
	default:
		return unknownStr
	}
}
