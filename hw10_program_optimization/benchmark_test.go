package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetDomain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, _ := zip.OpenReader("testdata/users.dat.zip")
		data, _ := r.File[0].Open()
		_, _ = GetDomainStat(data, "biz")
	}
}
