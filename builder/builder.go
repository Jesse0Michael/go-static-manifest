package builder

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafov/m3u8"
)

// Build will download the manifest segments to a relative path in the directory
// and rewrite the manifests to use the relative paths
func Build(manifest *url.URL, directory string) error {
	os.Mkdir(directory, os.ModePerm)

	resp, err := http.Get(manifest.String())
	if err != nil {
		return err
	}

	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return err
	}

	switch listType {
	case m3u8.MASTER:
		master := playlist.(*m3u8.MasterPlaylist)

		for i, v := range master.Variants {
			relative, err := url.Parse(v.URI)
			if err != nil {
				return err
			}

			err = Build(manifest.ResolveReference(relative), fmt.Sprintf("%s/variant%d", directory, i))
			if err != nil {
				return err
			}
			v.URI = fmt.Sprintf("variant%d/media.m3u8", i)

			for _, a := range v.Alternatives {
				alt, err := url.Parse(a.URI)
				if err != nil {
					return err
				}
				err = Build(manifest.ResolveReference(alt), fmt.Sprintf("%s/%s-%s-%s-%s", directory, a.Type, a.GroupId, a.Name, a.Language))
				if err != nil {
					return err
				}
				a.URI = fmt.Sprintf("%s-%s-%s-%s/media.m3u8", a.Type, a.GroupId, a.Name, a.Language)
			}
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/master.m3u8", directory), []byte(master.Encode().String()), os.ModePerm)
		if err != nil {
			return err
		}
	case m3u8.MEDIA:
		media := playlist.(*m3u8.MediaPlaylist)
		for i, s := range media.Segments[:media.Count()] {

			if s.Key != nil {
				keyURL, err := url.Parse(s.Key.URI)
				if err != nil {
					return err
				}
				resp, err := http.Get(manifest.ResolveReference(keyURL).String())
				if err != nil {
					return err
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				err = ioutil.WriteFile(fmt.Sprintf("%s/segment%d.key", directory, i), body, os.ModePerm)
				if err != nil {
					return err
				}
				s.Key.URI = fmt.Sprintf("segment%d.key", i)
			}

			relative, err := url.Parse(s.URI)
			if err != nil {
				return err
			}
			resp, err := http.Get(manifest.ResolveReference(relative).String())
			if err != nil {
				return err
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(fmt.Sprintf("%s/segment%d%s", directory, i, filepath.Ext(s.URI)), body, os.ModePerm)
			if err != nil {
				return err
			}
			s.URI = fmt.Sprintf("segment%d%s", i, filepath.Ext(s.URI))
		}

		// keys should already be downloaded a the m3u8 rewritten
		media.Key = nil

		err = ioutil.WriteFile(fmt.Sprintf("%s/media.m3u8", directory), []byte(media.Encode().String()), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// DecryptFile decrypts an encrypted file.
// Equivalent command: openssl aes-128-cbc -K <hex string> -iv <hex string> -d -in segment.ts -out segment-decrypted.ts
func DecryptFile(iv, key, inputFilepath, outputFilepath string) (err error) {
	// Strip off the prefix that indicates hex (if present)
	if strings.HasPrefix(iv, "0x") {
		iv = strings.TrimLeft(iv, "0x")
	}
	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return err
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return err
	}

	ciphertext, err := ioutil.ReadFile(inputFilepath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return err
	}

	stream := cipher.NewCBCDecrypter(block, ivBytes)
	stream.CryptBlocks(ciphertext, ciphertext)

	return ioutil.WriteFile(outputFilepath, ciphertext, 0444)
}

// EncryptFile encrypts a file.
func EncryptFile(iv, key, inputFilepath, outputFilepath string) error {
	// Strip off the prefix that indicates
	if strings.HasPrefix(iv, "0x") {
		iv = strings.TrimLeft(iv, "0x")
	}
	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return err
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return err
	}

	ciphertext, err := ioutil.ReadFile(inputFilepath)
	if err != nil {
		return err
	}

	stream := cipher.NewCBCEncrypter(block, ivBytes)
	stream.CryptBlocks(ciphertext, ciphertext)

	return ioutil.WriteFile(outputFilepath, ciphertext, 0444)
}
