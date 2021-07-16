package comment

import (
	"sports_service/server/models/mcomment"
)

type SortComment []*mcomment.CommentList

func (cm SortComment) Len() int {
	return len(cm)
}

func (cm SortComment) Less(i, j int) bool {
	if cm[i].LikeNum >= cm[j].LikeNum {
		return true
	}

	return false
}

func (cm SortComment) Swap(i, j int) {
	cm[i], cm[j] = cm[j], cm[i]
	return
}



