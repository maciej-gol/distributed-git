package main

import (
	"fmt"
	g2g "gopkg.in/libgit2/git2go.v24"
	"time"
)

func main() {
	repo, err := g2g.InitRepository("/tmp/sample-repo", true)
	fmt.Println(err)
	oid, err := repo.CreateBlobFromBuffer([]byte("XOXOXO"))
	fmt.Println(err)
	tb, err := repo.TreeBuilder()
	fmt.Println(err)
	err = tb.Insert("somefile.txt", oid, g2g.FilemodeBlob)
	fmt.Println(err)
	tbOid, err := tb.Write()
	fmt.Println(err)
	tree, err := repo.LookupTree(tbOid)
	fmt.Println(err)

	t := time.Date(2019, time.November, 2, 23, 0, 0, 0, time.UTC)
	author := &g2g.Signature{"Me", "me@example.com", t}
	cOid, err := repo.CreateCommit("HEAD", author, author, "Some msg", tree)
	fmt.Println(err)
	fmt.Println(cOid)
}
