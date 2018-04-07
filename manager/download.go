package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/brandoneprice31/freedrive/service"
)

func (m *Manager) Download(prefix, downloadToPath string) {
	k, err := loadKey(m.config.KeyPath)
	if err != nil {
		panic(err)
	}

	err = createDownloadFolder(downloadToPath)
	if err != nil {
		panic(err)
	}

	for _, s := range m.services {
		b, err := service.NewDownloadBuffer(s)
		if err != nil {
			panic(err)
		}

		m.addBuffer(b)
	}

	ff, dd := []file{}, []directory{}

	for _, ssd := range k.storedServiceData {
		err := m.bufs[ssd.ServiceType].Download(service.ServiceData{
			Data: ssd.ServiceData,
		})
		if err != nil {
			panic(err)
		}
	}

	raw, err := m.flushDownload()
	if err != nil {
		panic(err)
	}

	for _, r := range raw {
		var f file
		err := json.Unmarshal(r, &f)
		if err == nil && f.Path != "" {
			ff = append(ff, f)
			continue
		}

		var d directory
		err = json.Unmarshal(r, &d)
		if err == nil {
			dd = append(dd, d)
			continue
		}

		fmt.Println(err.Error())
	}

	for _, d := range dd {
		err := createDirectory(prefix, downloadToPath, d)
		if err != nil {
			panic(err)
		}
	}

	for _, f := range ff {
		err := createFile(prefix, downloadToPath, f)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("finished download")
}

func createDirectory(prefix, downloadToPath string, d directory) error {
	suffix := []rune(d.Path)[len(prefix)+1:]
	path := fmt.Sprintf("%s/%s", downloadToPath, string(suffix))
	return os.MkdirAll(path, 0777)
}

func createFile(prefix, downloadToPath string, f file) error {
	suffix := []rune(f.Path)[len(prefix)+1:]
	path := fmt.Sprintf("%s/%s", downloadToPath, string(suffix))
	return ioutil.WriteFile(path, f.Data, 0777)
}

func createDownloadFolder(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return os.Mkdir(path, 0777)
}

func loadKey(path string) (*key, error) {
	data, err := ioutil.ReadFile(path)
	if pErr, ok := err.(*os.PathError); ok && pErr.Err.Error() == "no such file or directory" {
		return nil, ErrNoKey
	} else if err != nil {
		return nil, err
	}

	var ssd []storedServiceData
	err = json.Unmarshal(data, &ssd)
	if err != nil {
		return nil, err
	}

	return &key{
		path:              path,
		storedServiceData: ssd,
	}, nil
}

func (m *Manager) flushDownload() ([][]byte, error) {
	rr := [][]byte{}
	for _, b := range m.bufs {
		raw, err := b.FlushDownload()
		if err != nil {
			return nil, err
		}
		rr = append(rr, raw...)
	}

	return rr, nil
}