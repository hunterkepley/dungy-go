package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
)

type musicType int

const (
	musicSampleRate = 22050 // Music sample rate

	typeOgg musicType = iota
	typeMP3
)

func (t musicType) String() string {
	switch t {
	case typeOgg:
		return "Ogg"
	case typeMP3:
		return "MP3"
	default:
		panic("not reached")
	}
}

// MusicContext is a context that holds the music player
type MusicContext struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
	musicType    musicType
}

func createMusicContext(audioContext *audio.Context, musicType musicType) (*MusicContext, error) {
	type audioStream interface {
		audio.ReadSeekCloser
		Length() int64
	}

	const bytesPerSample = 4 // TODO: This should be defined in audio package

	var s audioStream

	switch musicType {
	case typeOgg:
		var err error
		m, err := loadOggByte("./Assets/Music/underground_worm_song.ogg")

		s, err = vorbis.Decode(audioContext, audio.BytesReadSeekCloser(m))
		if err != nil {
			return nil, err
		}
	case typeMP3:
		// TODO: finish this if you add MP3!
		/*
			var err error
			s, err = mp3.Decode(audioContext, audio.BytesReadSeekCloser(raudio.Classic_mp3))
			if err != nil {
				return nil, err
			}*/
	default:
		panic("not reached")
	}
	p, err := audio.NewPlayer(audioContext, s)
	if err != nil {
		return nil, err
	}
	musicContext := &MusicContext{
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / musicSampleRate,
		volume128:    128,
		seCh:         make(chan []byte),
		musicType:    musicType,
	}
	if musicContext.total == 0 {
		musicContext.total = 1
	}
	musicContext.audioPlayer.Play()
	/*go func() {
		s, err := wav.Decode(audioContext, audio.BytesReadSeekCloser(raudio.Jab_wav))
		if err != nil {
			log.Fatal(err)
			return
		}
		b, err := ioutil.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		musicContext.seCh <- b
	}()*/
	return musicContext, nil
}

func (m *MusicContext) updateMusic() error {
	select {
	case m.seBytes = <-m.seCh:
		close(m.seCh)
		m.seCh = nil
	default:
	}

	if m.audioPlayer.IsPlaying() {
		m.current = m.audioPlayer.Current()
	}

	return nil
}

func (m *MusicContext) update(musicChan chan *MusicContext, errChan chan error) error {
	select {
	case p := <-musicChan:
		m = p
	case err := <-errChan:
		return err
	default:
	}

	/*if m != nil && inpututil.IsKeyJustPressed(ebiten.KeyA) {
		var t musicType
		switch m.musicType {
		case typeOgg:
			t = typeMP3
		case typeMP3:
			t = typeOgg
		default:
			panic("not reached")
		}

		//m.Close() ???
		m = nil

		go func() {
			p, err := createMusicContext(audio.CurrentContext(), t)
			if err != nil {
				errChan <- err
				return
			}
			musicChan <- p
		}()
	}*/

	if m != nil {
		if err := m.updateMusic(); err != nil {
			return err
		}
	}
	return nil
}
