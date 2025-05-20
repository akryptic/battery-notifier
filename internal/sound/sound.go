package sound

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"

	"github.com/akryptic/battery-notifier/internal/config"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	//go:embed assets/low.mp3
	lowBatteryDefaultSound []byte

	//go:embed assets/overcharge.mp3
	overchargeDefaultSound []byte
)

func Play(soundType string, conf *config.Config) error {
	if !conf.EnableSound {
		return nil
	}

	var soundData []byte
	var err error

	switch soundType {
	case "low":
		if conf.LowSoundFile != "" {
			soundData, err = os.ReadFile(conf.LowSoundFile)
			if err != nil {
				return fmt.Errorf("failed to read low sound file: %w", err)
			}
		} else {
			soundData = lowBatteryDefaultSound
		}
	case "overcharge":
		if conf.OverchargeSoundFile != "" {
			soundData, err = os.ReadFile(conf.OverchargeSoundFile)
			if err != nil {
				return fmt.Errorf("failed to read overcharge sound file: %w", err)
			}
		} else {
			soundData = overchargeDefaultSound
		}
	default:
		return fmt.Errorf("invalid sound type: %s", soundType)
	}

	return playBytes(soundData, conf.SoundVolume)

}

func playBytes(data []byte, vol int) error {
	soundReader := bytes.NewReader(data)
	closer := io.NopCloser(soundReader)

	// Decode the MP3 data
	streamer, format, err := mp3.Decode(closer)
	if err != nil {
		return fmt.Errorf("failed to decode sound file, ensure it's an MP3: %w", err)
	}

	defer streamer.Close()

	initSpeaker(format)

	volumeControlledStreamer := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   userVolumeToDB(vol),
		Silent:   false,
	}

	done := make(chan bool)
	speaker.Clear()
	speaker.Play(beep.Seq(volumeControlledStreamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil

}

var speakerInitOnce sync.Once

func initSpeaker(format beep.Format) {
	speakerInitOnce.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	})
}

func userVolumeToDB(volume int) float64 {
	if volume <= 0 {
		return -math.MaxFloat64
	}
	return 20 * math.Log10(float64(volume)/100.0)
}
