package sdp

const (
	exampleAttrExtmap1     = "extmap:1 http://example.com/082005/ext.htm#ttime"
	exampleAttrExtmap1Line = exampleAttrExtmap1
	exampleAttrExtmap2     = "extmap:2/sendrecv http://example.com/082005/ext.htm#xmeta short"
	exampleAttrExtmap2Line = exampleAttrExtmap2
	failingAttrExtmap1     = "extmap:257/sendrecv http://example.com/082005/ext.htm#xmeta short"
	failingAttrExtmap1Line = attributeKey + failingAttrExtmap1
	failingAttrExtmap2     = "extmap:2/blorg http://example.com/082005/ext.htm#xmeta short"
	failingAttrExtmap2Line = attributeKey + failingAttrExtmap2
)
