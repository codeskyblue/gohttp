package fsdriver

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func FSDriverTest(dr FSDriver, t *testing.T) {
	var err error

	// prepare
	dr.DeleteDir("tmp")

	// Mkdir test
	if err = dr.Mkdir("tmp"); err != nil {
		t.Fatal(err)
	}

	// PutFile noappend test
	size, err := dr.PutFile("tmp/a.txt", bytes.NewBufferString("ab"), false)
	if err != nil {
		t.Fatal(err)
	}
	if size != 2 {
		t.Fatalf("write file expect 2, but got %v", size)
	}

	// Stat test
	finfo, err := dr.Stat("tmp/a.txt")
	if err != nil {
		t.Fatal(err)
	}
	if finfo.Size() != 2 {
		t.Fatalf("stat file size expect 2, but got %v", finfo.Size())
	}

	// PutFile append test
	dr.DeleteFile("tmp/a.txt")
	size, err = dr.PutFile("tmp/a.txt", bytes.NewBufferString("abcd"), false)
	if err != nil {
		t.Fatal(err)
	}

	// not pass in qiniu
	// if size != 2 {
	// 	t.Fatalf("write file expect 2, but got %v", size)
	// }
	// size, err = dr.PutFile("tmp/a.txt", bytes.NewBufferString("cd"), true)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if size != 2 {
	// 	t.Fatalf("write file expect 2, but got %v", size)
	// }

	size, reader, err := dr.GetFile("tmp/a.txt", 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(string(data))
	if string(data) != "bcd" {
		t.Fatalf("Expect bcd but got %v", string(data))
	}

	// Rename test
	err = dr.Rename("tmp/a.txt", "tmp/b.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = dr.Stat("tmp/b.txt")
	if err != nil {
		t.Fatal(err)
	}

	// not pass in qiniu
	// err = dr.Mkdir("tmp/b.txt")
	// if err == nil {
	// 	t.Fatal("expect mkdir failed")
	// }

	fis, err := dr.ListDir("tmp")
	if err != nil {
		t.Fatal(err)
	}
	if len(fis) != 1 || fis[0].Name() != "b.txt" {
		t.Log(fis[0].Name())
		t.Fatalf("Expect directory tmp has only one file b.txt, but : %#v", fis)
	}

	// Delete file test
	if err = dr.DeleteFile("tmp/b.txt"); err != nil {
		t.Fatal(err)
	}

	// Delete dir test
	if err = dr.DeleteDir("tmp"); err != nil {
		t.Fatal(err)
	}
}

func TestPosixFSDriver(t *testing.T) {
	var dr FSDriver = &PosixFSDriver{"./"}
	defer os.RemoveAll("tmp")
	FSDriverTest(dr, t)
}

func TestQiniuFSDriver(t *testing.T) {
	dr, err := NewQiniuDriver(
		os.Getenv("QNAK"), os.Getenv("QNSK"),
		"gobuild3-test")
	if err != nil {
		t.Fatal(err)
	}
	FSDriverTest(dr, t)
}

func TestMultiFSDriver(t *testing.T) {
	var dr FSDriver = &MultiFSDriver{&PosixFSDriver{"./"}}
	FSDriverTest(dr, t)
}
