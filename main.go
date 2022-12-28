package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Sioro-Neoku/go-peerflix/search"
	"github.com/olekukonko/tablewriter"
)

// Exit statuses.
const (
	_ = iota
	exitNoTorrentProvided
	exitErrorInClient
)

func isTorrent(query string) bool {
	return strings.HasPrefix(query, "magnet:") || strings.HasPrefix(query, "http") || strings.HasSuffix(query, ".torrent")
}

func main() {
	// Parse flags.
	player := flag.String("player", "vlc", "Open the stream with a video player ("+joinPlayerNames()+")")
	cfg := NewClientConfig()
	flag.IntVar(&cfg.Port, "port", cfg.Port, "Port to stream the video on")
	flag.IntVar(&cfg.TorrentPort, "torrent-port", cfg.TorrentPort, "Port to listen for incoming torrent connections")
	flag.BoolVar(&cfg.Seed, "seed", cfg.Seed, "Seed after finished downloading")
	flag.IntVar(&cfg.MaxConnections, "conn", cfg.MaxConnections, "Maximum number of connections")
	flag.BoolVar(&cfg.TCP, "tcp", cfg.TCP, "Allow connections via TCP")
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(exitNoTorrentProvided)
	}

	query := strings.Join(flag.Args(), " ")
	if isTorrent(query) {
		cfg.TorrentPath = flag.Arg(0)
	} else {
		torrents := search.Search(query)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"#", "Name", "Size", "Seeders", "Leechers"})
		for i, torrent := range torrents {
			table.Append([]string{
				strconv.Itoa(i + 1),
				torrent.Name,
				torrent.Size,
				strconv.Itoa(torrent.Seeders),
				strconv.Itoa(torrent.Leechers),
			})
		}

		table.Render()

		var index int
		for {
			fmt.Print("Select a torrent: ")
			_, err := fmt.Scanf("%d", &index)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if index < 1 || index > len(torrents) {
				fmt.Println("Invalid index")
				continue
			}

			break
		}

		cfg.TorrentPath = torrents[index-1].MagnetURI()
		fmt.Println(cfg.TorrentPath)
	}

	// Start up the torrent client.
	client, err := NewClient(cfg)
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(exitErrorInClient)
	}

	// Http handler.
	go func() {
		http.HandleFunc("/", client.GetFile)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil))
	}()

	// Open selected video player
	if *player != "" {
		go func() {
			for !client.ReadyForPlayback() {
				time.Sleep(time.Second)
			}
			openPlayer(*player, cfg.Port)
		}()
	}

	// Handle exit signals.
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func(interruptChannel chan os.Signal) {
		for range interruptChannel {
			log.Println("Exiting...")
			client.Close()
			os.Exit(0)
		}
	}(interruptChannel)

	// Cli render loop.
	for {
		client.Render()
		time.Sleep(time.Second)
	}
}
