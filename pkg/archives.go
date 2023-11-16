package armada

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
)

func ExtractTarGz(gzipStream io.Reader, dest string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)
	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(dest+"/"+header.Name, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(dest + "/" + header.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			outFile.Close()

		default:
			return err
		}
	}
	return nil
}

func (j *Job) DownloadExtractURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	// create temporary dir
	dir, err := os.MkdirTemp("", "armada")
	if err != nil {
		return "", err
	}

	if ExtractTarGz(resp.Body, dir) != nil {
		return "", err
	}
	return dir, nil
}
