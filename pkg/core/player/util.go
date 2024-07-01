package player

func GetTags(collections ...interface{ GetTags() []string }) []string {
	var tags []string
	for _, col := range collections {
		tags = append(tags, col.GetTags()...)
	}
	return tags
}
