//go:build windows
// +build windows

package common

// General compile function
func Compile(script string, silent bool) (ret []byte, err error) {
	if len(script) > 2 {
		return nil, nil
	}

	return nil, nil
}
