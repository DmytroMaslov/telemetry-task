package crypto

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Encode_Decode(t *testing.T) {
	os.Setenv("SECRET_KEY", "1234567890123456") // Set a 16-byte key for AES-128
	defer os.Unsetenv("SECRET_KEY")

	en, err := NewEncryptor()
	require.NoError(t, err)
	de, err := NewDecryptor()
	require.NoError(t, err)

	message := "4505412866205444854,demo,1750926183300"

	encoded, err := en.EncryptMessage(message)
	require.NoError(t, err)

	decoded, err := de.DecryptMessage(encoded)
	require.NoError(t, err)
	require.Equal(t, message, decoded, "Decoded message should match the original message")
}
