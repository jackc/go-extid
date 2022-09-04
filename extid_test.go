package extid_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/jackc/go-extid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeKnownValues(t *testing.T) {
	prefix := "user"
	key := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	et, err := extid.NewType(prefix, key)
	require.NoError(t, err)

	for i, tt := range []struct {
		id  int64
		xid string
	}{
		{id: math.MinInt64, xid: "user_4399572cd6ea5341b8d35876a7098af7"},
		{id: -1, xid: "user_25d4e948bd5e1296afc0bf87095a7248"},
		{id: 0, xid: "user_c6a13b37878f5b826f4f8162a1c8d879"},
		{id: 1, xid: "user_13189a6ae4ab07ae70a3aabd30be99de"},
		{id: math.MaxInt64, xid: "user_edc17bee21fb24e211e6419412e1c32e"},
	} {
		xid := et.Encode(tt.id)
		assert.Equal(t, tt.xid, xid, "%d", i)
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	prefix := "user"
	key := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	et, err := extid.NewType(prefix, key)
	require.NoError(t, err)

	for i, tt := range []struct {
		id int64
	}{
		{id: math.MinInt64},
		{id: -1},
		{id: 0},
		{id: 1},
		{id: math.MaxInt64},
	} {
		xid := et.Encode(tt.id)
		assert.Truef(t, strings.HasPrefix(xid, prefix+"_"), "%d. xid: %s", i, xid)
		n, err := et.Decode(xid)
		assert.NoErrorf(t, err, "%d", i)
		assert.Equalf(t, tt.id, n, "%d", i)
	}
}

func FuzzEncodeDecode(f *testing.F) {

	testcases := []struct {
		id int64
	}{
		{id: math.MinInt64},
		{id: -1},
		{id: 0},
		{id: 1},
		{id: math.MaxInt64},
	}
	for _, tc := range testcases {
		f.Add(tc.id)
	}

	f.Fuzz(func(t *testing.T, id int64) {
		prefix := "user"
		key := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

		et, err := extid.NewType(prefix, key)
		require.NoError(t, err)

		xid := et.Encode(id)
		require.True(t, strings.HasPrefix(xid, prefix+"_"))
		roundTripID, err := et.Decode(xid)
		require.NoError(t, err)
		require.Equal(t, id, roundTripID)
	})
}

var benchXID string

func BenchmarkEncode(b *testing.B) {
	prefix := "user"
	key := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	et, err := extid.NewType(prefix, key)
	require.NoError(b, err)

	b.ResetTimer()
	var xid string
	for i := 0; i < b.N; i++ {
		xid = et.Encode(int64(i))
	}

	benchXID = xid
}

var benchID int64

func BenchmarkDecode(b *testing.B) {
	prefix := "user"
	key := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	et, err := extid.NewType(prefix, key)
	require.NoError(b, err)

	xids := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		xids[i] = et.Encode(int64(i))
	}

	b.ResetTimer()
	var id int64
	for i := 0; i < b.N; i++ {
		id, err = et.Decode(xids[i])
		if err != nil {
			b.Fatalf("et.Decode failed: i: %d; err: %v", i, err)
		}
		if id != int64(i) {
			b.Fatalf("et.Decode should have been %d but got %d", i, id)
		}
	}

	benchID = id
}

func ExampleType_Encode() {
	prefix := "user"
	key := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	et, err := extid.NewType(prefix, key)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 5; i++ {
		fmt.Println(i, et.Encode(int64(i)))
	}

	// Output:
	// 0 user_c6a13b37878f5b826f4f8162a1c8d879
	// 1 user_13189a6ae4ab07ae70a3aabd30be99de
	// 2 user_c76e8fcf7ad0fe9b39e083739cbe26c2
	// 3 user_90cb45611c3105c84624b2ac12cb5b74
	// 4 user_79d0782401799c8fd25121e869b5b532
}
