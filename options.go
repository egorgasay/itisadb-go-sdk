package itisadb

type GetOptions struct {
	Server int32
}

type SetOptions struct {
	Server   int32
	ReadOnly bool
	Unique   bool
	Level    Level
}

type DeleteOptions struct {
	Server int32
}

type Level byte

func (l Level) String() string {
	switch l {
	case DefaultLevel:
		return "Default"
	case RestrictedLevel:
		return "Restricted"
	case SecretLevel:
		return "Secret"
	}

	return "unknown"
}

type ObjectOptions struct {
	Server int32
	Level  Level
}

type ObjectToJSONOptions struct {
	Server int32
}

type DeleteObjectOptions struct {
	Server int32
}

type IsObjectOptions struct {
	Server int32
}

type SizeOptions struct {
	Server int32
}

type AttachToObjectOptions struct {
	Server int32
}

type SetToObjectOptions struct {
	Server   int32
	ReadOnly bool
}

type GetFromObjectOptions struct {
	Server int32
}

type ConnectOptions struct {
	Server int32
}

type DeleteKeyOptions struct {
	Server int32
}

type CreateUserOptions struct {
	Level Level
}
