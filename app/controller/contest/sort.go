package contest

import (
	"sports_service/server/models/mcontest"
)

type SortContestLive []*mcontest.ContestLiveInfo

func (cl SortContestLive) Len() int {
	return len(cl)
}

func (cl SortContestLive) Less(i, j int) bool {
	if cl[i].Index >= cl[j].Index {
		return true
	}

	return false
}

func (cl SortContestLive) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
	return
}




