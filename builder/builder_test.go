package builder

import (
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jesse0michael/go-rest-assured/assured"
	"github.com/stretchr/testify/require"
)

func TestBuild_SuccessEncrypted(t *testing.T) {
	ts := assured.NewDefaultClient()
	defer ts.Close()
	master, _ := ioutil.ReadFile("testdata/aes128/master.m3u8")
	media, _ := ioutil.ReadFile("testdata/aes128/media.m3u8")
	segment0, _ := ioutil.ReadFile("testdata/aes128/segment0.ts")
	segment1, _ := ioutil.ReadFile("testdata/aes128/segment1.ts")
	segmentKey, _ := ioutil.ReadFile("testdata/aes128/segment0.key")
	ts.Given(
		assured.Call{
			Path:     "master.m3u8",
			Response: master,
		},
		assured.Call{
			Path:     "media.m3u8",
			Response: media,
		},
		assured.Call{
			Path:     "segment0.ts",
			Response: segment0,
		},
		assured.Call{
			Path:     "segment1.ts",
			Response: segment1,
		},
		assured.Call{
			Path:     "segment0.key",
			Response: segmentKey,
		},
	)

	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	url, _ := url.Parse(ts.URL() + "/master.m3u8")

	err := Build(url, dir)

	require.NoError(t, err)
	require.FileExists(t, dir+"/master.m3u8")
	require.FileExists(t, dir+"/variant0/media.m3u8")
	require.FileExists(t, dir+"/variant0/segment0.ts")
	require.FileExists(t, dir+"/variant0/segment1.ts")
	require.FileExists(t, dir+"/variant0/segment0.key")

	calls, _ := ts.Verify(http.MethodGet, "master.m3u8")
	require.Len(t, calls, 1)
	calls, _ = ts.Verify(http.MethodGet, "media.m3u8")
	require.Len(t, calls, 1)
	calls, _ = ts.Verify(http.MethodGet, "segment0.ts")
	require.Len(t, calls, 1)
	calls, _ = ts.Verify(http.MethodGet, "segment1.ts")
	require.Len(t, calls, 1)
	calls, _ = ts.Verify(http.MethodGet, "segment0.key")
	require.Len(t, calls, 1)
}

func TestBuild_SuccessClear(t *testing.T) {
	ts := assured.NewDefaultClient()
	defer ts.Close()
	master, _ := ioutil.ReadFile("testdata/clear/master.m3u8")
	media, _ := ioutil.ReadFile("testdata/clear/media.m3u8")
	segment0, _ := ioutil.ReadFile("testdata/clear/segment0")
	segment1, _ := ioutil.ReadFile("testdata/clear/segment1")
	ts.Given(
		assured.Call{
			Path:     "master.m3u8",
			Response: master,
		},
		assured.Call{
			Path:     "media.m3u8",
			Response: media,
		},
		assured.Call{
			Path:     "segment0",
			Response: segment0,
		},
		assured.Call{
			Path:     "segment1",
			Response: segment1,
		},
	)

	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	url, _ := url.Parse(ts.URL() + "/master.m3u8")

	err := Build(url, dir)

	require.NoError(t, err)
	require.FileExists(t, dir+"/master.m3u8")
	require.FileExists(t, dir+"/variant0/media.m3u8")
	require.FileExists(t, dir+"/variant0/segment0")
	require.FileExists(t, dir+"/variant0/segment1")
	require.FileExists(t, dir+"/AUDIO-audio-0-en (Main)-en/media.m3u8")
	require.FileExists(t, dir+"/AUDIO-audio-0-en (Main)-en/segment0")
	require.FileExists(t, dir+"/AUDIO-audio-0-en (Main)-en/segment1")

	calls, _ := ts.Verify(http.MethodGet, "master.m3u8")
	require.Len(t, calls, 1)
	calls, _ = ts.Verify(http.MethodGet, "media.m3u8")
	require.Len(t, calls, 2)
	calls, _ = ts.Verify(http.MethodGet, "segment0")
	require.Len(t, calls, 2)
	calls, _ = ts.Verify(http.MethodGet, "segment1")
	require.Len(t, calls, 2)
}

func TestBuild_ErrorURL(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	url, _ := url.Parse("httz://localhost")

	err := Build(url, "")

	require.EqualError(t, err, `Get httz://localhost: unsupported protocol scheme "httz"`)
}

func TestBuild_ErrorNonM3U8(t *testing.T) {
	ts := assured.NewDefaultClient()
	defer ts.Close()
	nonMaster, _ := ioutil.ReadFile("testdata/aes128/segment0.ts")
	ts.Given(
		assured.Call{
			Path:     "master.m3u8",
			Response: nonMaster,
		},
	)

	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	url, _ := url.Parse(ts.URL() + "/master.m3u8")

	err := Build(url, dir)

	require.EqualError(t, err, "#EXTM3U absent")
}

func TestBuild_ErrorInvalidVariantURL(t *testing.T) {
	ts := assured.NewDefaultClient()
	defer ts.Close()
	badVariant, _ := ioutil.ReadFile("testdata/invalid/variant_url.m3u8")
	ts.Given(
		assured.Call{
			Path:     "master.m3u8",
			Response: badVariant,
		},
	)

	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	url, _ := url.Parse(ts.URL() + "/master.m3u8")

	err := Build(url, dir)

	require.EqualError(t, err, "parse ://media.m3u8: missing protocol scheme")
}

func TestBuild_ErrorInvalidKeyURL(t *testing.T) {
	ts := assured.NewDefaultClient()
	defer ts.Close()
	badKey, _ := ioutil.ReadFile("testdata/invalid/key_url.m3u8")
	ts.Given(
		assured.Call{
			Path:     "master.m3u8",
			Response: badKey,
		},
	)

	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	url, _ := url.Parse(ts.URL() + "/master.m3u8")

	err := Build(url, dir)

	require.EqualError(t, err, "parse ://segment0.key: missing protocol scheme")
}

func TestBuild_ErrorInvalidSegmentURL(t *testing.T) {
	ts := assured.NewDefaultClient()
	defer ts.Close()
	badSegment, _ := ioutil.ReadFile("testdata/invalid/segment_url.m3u8")
	ts.Given(
		assured.Call{
			Path:     "master.m3u8",
			Response: badSegment,
		},
	)

	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	url, _ := url.Parse(ts.URL() + "/master.m3u8")

	err := Build(url, dir)

	require.EqualError(t, err, "parse ://segment0: missing protocol scheme")
}

func TestBuild_ErrorInvalidAudioURL(t *testing.T) {
	ts := assured.NewDefaultClient()
	defer ts.Close()
	badAudio, _ := ioutil.ReadFile("testdata/invalid/audio_url.m3u8")
	media, _ := ioutil.ReadFile("testdata/clear/media.m3u8")
	ts.Given(
		assured.Call{
			Path:     "master.m3u8",
			Response: badAudio,
		},
		assured.Call{
			Path:     "media.m3u8",
			Response: media,
		},
	)

	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	url, _ := url.Parse(ts.URL() + "/master.m3u8")

	err := Build(url, dir)

	require.EqualError(t, err, "parse ://media.m3u8: missing protocol scheme")
}

func TestDecryptFile_Success(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	expected0, _ := ioutil.ReadFile("testdata/clear/segment0")
	expected1, _ := ioutil.ReadFile("testdata/clear/segment1")
	key, _ := ioutil.ReadFile("testdata/aes128/segment0.key")

	err := DecryptFile("0xeee0af7d41390eb6bb1405c66a029d3a", hex.EncodeToString(key), "testdata/aes128/segment0.ts", dir+"/segment0")

	require.NoError(t, err)
	actual, _ := ioutil.ReadFile(dir + "/segment0")
	require.Equal(t, expected0, actual)

	err = DecryptFile("eee0af7d41390eb6bb1405c66a029d3a", hex.EncodeToString(key), "testdata/aes128/segment1.ts", dir+"/segment1")

	require.NoError(t, err)
	actual, _ = ioutil.ReadFile(dir + "/segment1")
	require.Equal(t, expected1, actual)
}

func TestDecryptFile_ErrorInvalidIV(t *testing.T) {
	err := DecryptFile("`", "", "", "segment")

	require.EqualError(t, err, "encoding/hex: invalid byte: U+0060 '`'")
}

func TestDecryptFile_ErrorInvalidKey(t *testing.T) {
	err := DecryptFile("", "`", "", "segment")

	require.EqualError(t, err, "encoding/hex: invalid byte: U+0060 '`'")
}

func TestDecryptFile_ErrorInvalidFile(t *testing.T) {
	err := DecryptFile("", "", "", "segment")

	require.EqualError(t, err, "open : no such file or directory")
}

func TestDecryptFile_ErrorInvalidCipher(t *testing.T) {
	err := DecryptFile("", "", "testdata/aes128/segment1.ts", "segment")

	require.EqualError(t, err, "crypto/aes: invalid key size 0")
}

func TestEncryptFile_Success(t *testing.T) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	expected0, _ := ioutil.ReadFile("testdata/aes128/segment0.ts")
	expected1, _ := ioutil.ReadFile("testdata/aes128/segment1.ts")
	key, _ := ioutil.ReadFile("testdata/aes128/segment0.key")

	err := EncryptFile("0xeee0af7d41390eb6bb1405c66a029d3a", hex.EncodeToString(key), "testdata/clear/segment0", dir+"/segment0")

	require.NoError(t, err)
	actual, _ := ioutil.ReadFile(dir + "/segment0")
	require.Equal(t, expected0, actual)

	err = EncryptFile("eee0af7d41390eb6bb1405c66a029d3a", hex.EncodeToString(key), "testdata/clear/segment1", dir+"/segment1")

	require.NoError(t, err)
	actual, _ = ioutil.ReadFile(dir + "/segment1")
	require.Equal(t, expected1, actual)
}

func TestEncryptFile_ErrorInvalidIV(t *testing.T) {
	err := EncryptFile("`", "", "", "segment")

	require.EqualError(t, err, "encoding/hex: invalid byte: U+0060 '`'")
}

func TestEncryptFile_ErrorInvalidKey(t *testing.T) {
	err := EncryptFile("", "`", "", "segment")

	require.EqualError(t, err, "encoding/hex: invalid byte: U+0060 '`'")
}

func TestEncryptFile_ErrorInvalidFile(t *testing.T) {
	key, _ := ioutil.ReadFile("testdata/aes128/segment0.key")
	err := EncryptFile("0xeee0af7d41390eb6bb1405c66a029d3a", hex.EncodeToString(key), "", "segment")

	require.EqualError(t, err, "open : no such file or directory")
}

func TestEncryptFile_ErrorInvalidCipher(t *testing.T) {
	err := EncryptFile("", "", "testdata/clear/segment0", "segment")

	require.EqualError(t, err, "crypto/aes: invalid key size 0")
}
