package tests

import (
	_struct "github.com/golang/protobuf/ptypes/struct"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/protoc-gen-ext/tests/api"
)

var _ = Describe("hash", func() {
	testRuns := 100
	Context("simple (non-nested)", func() {
		It("can handle different simple objects", func() {
			testCases := []*api.Simple{
				{
					Str:          "",
					Byt:          []byte("hsajkldhaks"),
					TestUint32:   10,
					TestUint64:   5,
					TestBool:     true,
					DoubleTest:   7,
					FloatTest:    10,
					Int32Test:    12,
					Int64Test:    13,
					Sint32Test:   14,
					Sint64Test:   15,
					Fixed32Test:  15,
					Fixed64Test:  16,
					Sfixed32Test: 18,
					Sfixed64Test: 18,
				},
				{
					Str:          "sadfas",
					Byt:          []byte("sadhjk"),
					TestUint32:   22,
					TestUint64:   30,
					TestBool:     false,
					DoubleTest:   21,
					FloatTest:    22,
					Int32Test:    23,
					Int64Test:    24,
					Sint32Test:   25,
					Sint64Test:   26,
					Fixed32Test:  2,
					Fixed64Test:  27,
					Sfixed32Test: 28,
					Sfixed64Test: 29,
				},
			}
			hashes := make([]uint64, len(testCases))
			for i, v := range testCases {
				hash, err := v.Hash(nil)
				Expect(err).NotTo(HaveOccurred())
				hashes[i] = hash
			}
			for idx := 0; idx < testRuns; idx++ {
				newHashes := make([]uint64, len(testCases))
				for i, v := range testCases {
					hash, err := v.Hash(nil)
					Expect(err).NotTo(HaveOccurred())
					Expect(hash).To(Equal(hashes[i]))
					newHashes[i] = hash
				}
				Expect(allDifferent(newHashes)).To(BeTrue())
			}
		})
	})

	FContext("complex (nested)", func() {
		It("can produce the same object across many runs", func() {
			testObject := &api.Nested{
				Simple:  &api.Simple{},
				Details: &_struct.Struct{},
				Hello:   []string{"hello", "world"},
				X:       []*api.Simple{&api.Simple{}, nil},
				Initial: map[string]*api.Simple{
					"hello": &api.Simple{},
					"world": nil,
				},
			}
			initialHash, err := testObject.Hash(nil)
			Expect(err).NotTo(HaveOccurred())
			for idx := 1; idx < testRuns; idx++ {
				hash, err := testObject.Hash(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(hash).To(Equal(initialHash))
			}
		})

		FIt("can produce the same object across many runs", func() {
			testCases := []*api.Nested{
				{
					Simple:  &api.Simple{},
					Details: &_struct.Struct{},
					Hello:   []string{"hello", "world"},
					X:       []*api.Simple{&api.Simple{}, nil},
					Initial: map[string]*api.Simple{
						"hello": &api.Simple{},
						"world": nil,
					},
				},
				{
					Simple:               nil,
					Details:              nil,
					Hello:                nil,
					X:                    nil,
					Initial:              nil,
				},
				{
					Simple:               &api.Simple{
						Str:                  "gello",
						Byt:                  []byte("world"),
						TestUint32:           123,
						TestUint64:           32,
						TestBool:             true,
					},
					Details:              nil,
					Hello:   			[]string{"", "test"},
					X:                    []*api.Simple{nil, nil, &api.Simple{}},
					Initial:              nil,
				},
			}
			hashes := make([]uint64, len(testCases))
			for i, v := range testCases {
				hash, err := v.Hash(nil)
				Expect(err).NotTo(HaveOccurred())
				hashes[i] = hash
			}
			for idx := 0; idx < testRuns; idx++ {
				newHashes := make([]uint64, len(testCases))
				for i, v := range testCases {
					hash, err := v.Hash(nil)
					Expect(err).NotTo(HaveOccurred())
					Expect(hash).To(Equal(hashes[i]))
					newHashes[i] = hash
				}
				Expect(allDifferent(newHashes)).To(BeTrue())
			}
		})
	})
})

func allDifferent(hashes []uint64) bool {
	present := make(map[uint64]struct{})
	for _, v := range hashes {
		_, ok := present[v]
		if ok {
			return false
		}
		present[v] = struct{}{}
	}
	return true
}