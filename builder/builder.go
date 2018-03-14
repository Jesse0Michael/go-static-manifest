package builder

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/grafov/m3u8"
)

// Build will download the manifest segments to a relative path in the directory
// and rewrite the manifests to use the relative paths
func Build(manifest *url.URL, directory string) error {
	os.Mkdir(directory, os.ModePerm)
	fmt.Println("building manifest: " + manifest.String())

	resp, err := http.Get(manifest.String())
	if err != nil {
		return err
	}

	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		fmt.Println("decode error: " + err.Error())
		return err
	}

	switch listType {
	case m3u8.MASTER:
		fmt.Println("start master")
		master := playlist.(*m3u8.MasterPlaylist)

		for i, v := range master.Variants {
			fmt.Println("start variant")
			relative, err := url.Parse(v.URI)
			if err != nil {
				return err
			}

			err = Build(manifest.ResolveReference(relative), fmt.Sprintf("%s/variant%d", directory, i))
			if err != nil {
				return err
			}
			v.URI = fmt.Sprintf("variant%d/variant.m3u8", i)

			for _, a := range v.Alternatives {
				fmt.Println("start alternate")
				alt, err := url.Parse(a.URI)
				if err != nil {
					return err
				}
				err = Build(manifest.ResolveReference(alt), fmt.Sprintf("%s/%s-%s-%s-%s", directory, a.Type, a.GroupId, a.Name, a.Language))
				if err != nil {
					return err
				}
				a.URI = fmt.Sprintf("%s-%s-%s-%s/variant.m3u8", a.Type, a.GroupId, a.Name, a.Language)
				fmt.Println("end alternate")
			}
			fmt.Println("end variant")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/master.m3u8", directory), []byte(master.Encode().String()), os.ModePerm)
		if err != nil {
			return err
		}
		fmt.Println("end master")
	case m3u8.MEDIA:
		fmt.Println("start media")
		media := playlist.(*m3u8.MediaPlaylist)
		for i, s := range media.Segments {
			if s != nil {
				fmt.Println("start segment")
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
				fmt.Println("end segment")
			}
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/media.m3u8", directory), []byte(media.Encode().String()), os.ModePerm)
		if err != nil {
			return err
		}
		fmt.Println("end media")
	}
	return nil
}