package sound

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
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
				log.Printf("ERROR: Failed to read custom low sound file '%s': %v", conf.LowSoundFile, err)
				return fmt.Errorf("failed to read low sound file: %w", err)
			}
		} else {
			soundData = lowBatteryDefaultSound
		}
	case "overcharge":
		if conf.OverchargeSoundFile != "" {
			soundData, err = os.ReadFile(conf.OverchargeSoundFile)
			if err != nil {
				log.Printf("ERROR: Failed to read custom overcharge sound file '%s': %v", conf.OverchargeSoundFile, err)
				return fmt.Errorf("failed to read overcharge sound file: %w", err)
			}
		} else {
			soundData = overchargeDefaultSound
		}
	default:
		log.Printf("ERROR: Invalid sound type requested: %s", soundType)
		return fmt.Errorf("invalid sound type: %s", soundType)
	}

	err = playBytes(soundData, conf.SoundVolume)
	if err != nil {
		log.Printf("ERROR: Sound playback failed: %v", err)
		return err
	}

	return nil
}

func playBytes(data []byte, vol int) error {
	soundReader := bytes.NewReader(data)
	closer := io.NopCloser(soundReader)

	streamer, format, err := mp3.Decode(closer)
	if err != nil {
		log.Printf("ERROR: Failed to decode MP3 data: %v", err)
		return fmt.Errorf("failed to decode sound file, ensure it's an MP3: %w", err)
	}

	defer streamer.Close()

	err = initSpeakerWithTimeout(format)
	if err != nil {
		log.Printf("ERROR: Failed to initialize speaker: %v", err)
		return fmt.Errorf("audio system unavailable: %w", err)
	}

	volumeDB := userVolumeToDB(vol)

	volumeControlledStreamer := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   volumeDB,
		Silent:   false,
	}

	done := make(chan bool)
	timeout := make(chan bool, 1)

	go func() {
		time.Sleep(10 * time.Second)
		select {
		case timeout <- true:
			log.Println("WARNING: Audio playback timed out after 10 seconds")
		default:
		}
	}()

	speaker.Clear()
	speaker.Play(beep.Seq(volumeControlledStreamer, beep.Callback(func() {
		select {
		case done <- true:
		default:
		}
	})))

	select {
	case <-done:
		return nil
	case <-timeout:
		log.Println("WARNING: Sound playback timed out - likely running without audio session")
		speaker.Clear()
		return fmt.Errorf("audio playback timed out - no audio session available")
	}
}

var speakerInitOnce sync.Once
var speakerInitError error

func initSpeakerWithTimeout(format beep.Format) error {
	speakerInitOnce.Do(func() {
		done := make(chan bool, 1)
		timeout := make(chan bool, 1)

		go func() {
			time.Sleep(5 * time.Second)
			select {
			case timeout <- true:
			default:
			}
		}()

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("ERROR: Speaker initialization failed: %v", r)
					speakerInitError = fmt.Errorf("speaker initialization failed: %v", r)
				}
				select {
				case done <- true:
				default:
				}
			}()

			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		}()

		select {
		case <-done:
			// initialization completed (successfully or with panic)
		case <-timeout:
			log.Println("ERROR: Speaker initialization timed out")
			speakerInitError = fmt.Errorf("speaker initialization timed out")
		}
	})

	return speakerInitError
}

func userVolumeToDB(volume int) float64 {
	if volume <= 0 {
		return -math.MaxFloat64
	}
	return 20 * math.Log10(float64(volume)/100.0)
}
