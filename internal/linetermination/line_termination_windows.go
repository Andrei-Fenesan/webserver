//go:build windows

package linetermination

func GetLineTermination() []byte {
	return []byte{'\r', '\n'}
}
