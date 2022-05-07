package private_host_key

import (
	gossh "golang.org/x/crypto/ssh"
	"io/ioutil"
)

// ReadHostKeyFile 读取ssh私钥
func ReadHostKeyFile(filepath string) (gossh.Signer, error) {
	keyBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	key, err := gossh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}
