package gutil

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Zip 压缩成Zip
func Zip(src, dest string) (err error) {
	// Create a file to be written.
	fw, err := os.Create(dest)
	defer fw.Close()
	if err != nil {
		return err
	}

	// Create zip.Write through fw.
	zw := zip.NewWriter(fw)
	defer func() {
		// Check if it is successfully closed.
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 打开文件夹
	dir, err := os.Open(src)
	defer dir.Close()
	if err != nil {
		return err
	}
	// 读取文件列表
	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	// 遍历文件列表
	for _, fi := range fis {
		// 逃过文件夹, 我这里就不递归了
		if fi.IsDir() {
			continue
		}
		// 打印文件名称
		fmt.Println(fi.Name())
		// 打开文件
		fr, err := os.Open(dir.Name() + "/" + fi.Name())
		if err != nil {
			return err
		}
		defer fr.Close()
		// Create zip file information through file information.
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// Write file information and return a Write structure.
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// If it is not a standard file, only the header information is written, and the file data is not written to w.
		// Such as a directory, there is no data to write.
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 写文件
		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}
	fmt.Println("Tar Success")
	return nil
}

// ZipRecursive Zip Compressed file.
func ZipRecursive(src, dest string) (err error) {
	// Create a file to be written.
	fw, err := os.Create(dest)
	defer fw.Close()
	if err != nil {
		return err
	}

	// Create zip.Write through fw.
	zw := zip.NewWriter(fw)
	defer func() {
		// Check if it is successfully closed.
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Write the file to zw, because there may be many directories and files, so recursive processing.
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		// Create zip file information through file information.
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		// Replace the file name in the file information
		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))

		// Judge it is not a directory.
		if fi.IsDir() {
			fh.Name += "/"
		}

		// Write file information and return a Write structure.
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		// If it is not a standard file, only the header information is written, and the file data is not written to w.
		// Such as a directory, there is no data to write.
		if !fh.Mode().IsRegular() {
			return nil
		}

		// Open the file to be compressed.
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return
		}

		// Copy the opened file to w .
		n, err := io.Copy(w, fr)
		if err != nil {
			return
		}
		// Output compressed content.
		fmt.Printf("Zip success %s, A total of %d characters of data are written\n", path, n)

		return nil
	})
}

// UnZip 解压缩文件.
func UnZip(dst, src string) (err error) {
	// Open compressed file.
	zr, err := zip.OpenReader(src)
	defer zr.Close()
	if err != nil {
		return
	}

	// If it is not placed in the current directory after decompression, create a directory according to the save directory.
	if dst != "" {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
	}

	// Traverse zr and write the file to disk
	for _, file := range zr.File {
		path := filepath.Join(dst, file.Name)

		// If it is a directory, create a directory.
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			// Because it is a directory, skip the current loop, because the following is the processing of files.
			continue
		}

		// Get the Reader.
		fr, err := file.Open()
		if err != nil {
			//Also close when abnormal
			close(fr, nil)
			return err
		}

		// Create the Write corresponding to the file to be written.
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			//Also close when abnormal
			close(fr, fw)
			return err
		}

		n, err := io.Copy(fw, fr)
		if err != nil {
			//Also close when abnormal
			close(fr, fw)
			return err
		}

		// Output the decompressed result.
		fmt.Printf("UnZip success %s A total of %d characters have been written\n", path, n)

		//Finally remember to close
		close(fr, fw)
	}
	return nil
}

// close close关闭文件流
func close(fr io.ReadCloser, fw *os.File) {
	if fr != nil {
		fr.Close()
	}
	if fw != nil {
		fw.Close()
	}

}

func Tar(src, dest string) (err error) {
	// file write
	fw, err := os.Create(dest)
	defer fw.Close()
	if err != nil {
		return err
	}
	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()
	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()
	// 打开文件夹
	dir, err := os.Open(src)
	defer dir.Close()
	if err != nil {
		return err
	}
	// 读取文件列表
	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	// 遍历文件列表
	for _, fi := range fis {
		// 逃过文件夹, 我这里就不递归了
		if fi.IsDir() {
			continue
		}
		// 打印文件名称
		fmt.Println(fi.Name())
		// 打开文件
		fr, err := os.Open(dir.Name() + "/" + fi.Name())
		if err != nil {
			return err
		}
		defer fr.Close()
		// 信息头
		h := new(tar.Header)
		h.Name = fi.Name()
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()
		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			return err
		}
		// 写文件
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
	}
	fmt.Println("Tar Success")
	return nil
}

func TarRecursive(src, dest string) (err error) {
	// Create a file to be written.
	fw, err := os.Create(dest)
	defer fw.Close()
	if err != nil {
		return err
	}

	// Create zip.Write through fw.
	gw := gzip.NewWriter(fw)
	defer func() {
		// Check if it is successfully closed.
		if err := gw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Write the file to zw, because there may be many directories and files, so recursive processing.
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		// Create zip file information through file information.
		fh, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return
		}

		// Replace the file name in the file information
		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))

		// Judge it is not a directory.
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 信息头
		h := new(tar.Header)
		h.Name = fi.Name()
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()
		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			return err
		}

		// Open the file to be compressed.
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return
		}

		// Copy the opened file to w .
		_, err = io.Copy(tw, fr)
		if err != nil {
			return
		}
		// Output compressed content.
		fmt.Printf("TarRecursive Success \n")

		return nil
	})
}

// UnTar 解压文件.
// src is the abbreviation of source, which is a directory of files to be decompressed.
// dst is the abbreviation of destination,  which is the decompressed file directory.
func UnTar(src, dest string) (err error) {
	// file read
	fr, err := os.Open(src)
	defer fr.Close()
	if err != nil {
		panic(err)
	}
	// tar read
	tr := tar.NewReader(fr)
	// 读取文件
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if h.FileInfo().IsDir() {
			err = os.MkdirAll(dest+h.Name, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}
		// 打开文件
		fw, err := os.OpenFile(dest+h.Name, os.O_CREATE|os.O_WRONLY, 0666 /*os.FileMode(h.Mode)*/)
		if err != nil {
			return err
		}
		defer fw.Close()
		// 写文件
		_, err = io.Copy(fw, tr)
		if err != nil {
			return err
		}
	}
	fmt.Println("UnTar Success!")
	return nil
}
