//go:build !windows
// +build !windows

package common

// General compile function
func Compile(script string, silent bool) (ret []byte, err error) {
	if len(script) > 2 {
		byteCode, err := serpent.Compile(script)
		if err != nil {
			return nil, err
		}

		return byteCode, nil
	}

	return nil, nil
}
