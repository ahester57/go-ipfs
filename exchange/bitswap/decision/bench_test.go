package decision

import (
	"fmt"
	"math"
	"testing"

	"github.com/ahester57/go-ipfs/exchange/bitswap/wantlist"
	u "gx/ipfs/QmPdKqUcHGFdeSpvjVoaTRPPstGif9GBZb5Q56RVw9o69A/go-ipfs-util"
	cid "gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
	"gx/ipfs/QmcW4FGAt24fdK1jBgWQn3yP4R9ZLyWQqjozv9QK7epRhL/go-testutil"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// FWIW: At the time of this commit, including a timestamp in task increases
// time cost of Push by 3%.
func BenchmarkTaskQueuePush(b *testing.B) {
	q := newPRQ()
	peers := []peer.ID{
		testutil.RandPeerIDFatal(b),
		testutil.RandPeerIDFatal(b),
		testutil.RandPeerIDFatal(b),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := cid.NewCidV0(u.Hash([]byte(fmt.Sprint(i))))

		q.Push(&wantlist.Entry{Cid: c, Priority: math.MaxInt32}, peers[i%len(peers)])
	}
}
