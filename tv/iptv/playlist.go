package iptv

import (
	"bufio"
)

type playlist struct {
	first      string
	streamList []stream
}

func newPlaylist(first string, sList []stream) playlist {
	return playlist{first: first, streamList: sList}
}

func (pl *playlist) fmt(w *bufio.Writer) {
	w.WriteString(pl.first)
	w.WriteByte('\n')
	for _, s := range pl.streamList {
		s.fmt(w)
	}
	w.Flush()
}

type stream struct {
	comment []string
	url     string
}

func newStream() *stream {
	return &stream{}
}

func (s *stream) fmt(w *bufio.Writer) {
	for _, v := range s.comment {
		w.WriteString(v)
		w.WriteByte('\n')
	}
	w.WriteString(s.url)
	w.WriteByte('\n')
}
