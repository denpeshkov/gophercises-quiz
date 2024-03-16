package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	path := flag.String("csv", "problems.csv", "a `path` to the .csv file with quiz questions/answers pairs")
	limit := flag.Duration("limit", 30*time.Second, "`time` limit for a quiz")
	flag.Parse()

	log.SetFlags(0)

	f, err := os.Open(*path)
	if err != nil {
		log.Fatalf("Couldn't open the file: %q\n", *path)
	}
	defer f.Close()

	in, out := make(chan struct{}), make(chan string)

	go func() {
		defer close(out)

		sc := bufio.NewScanner(os.Stdin)
		for range in {
			sc.Scan()
			if sc.Err() != nil {
				log.Fatalf("Error reading user input: %v", sc.Err())
			}
			out <- sc.Text()
		}
	}()

	fmt.Println("Press ENTER to start the quiz...")
	fmt.Scanln()
	tCh := time.After(*limit)

	correct, total := 0, 0
	r := csv.NewReader(f)
readLoop:
	for {
		line, err := r.Read()
		if err == io.EOF {
			break readLoop
		}
		if err != nil {
			log.Fatal("Failed to parse the CSV record")
		}
		problem := problem{q: line[0], a: line[1]}

		total++
		fmt.Printf("Problem #%d: %v = ", total, problem.q)

		in <- struct{}{}
		select {
		case v := <-out:
			if v == problem.a {
				correct++
			}
		case <-tCh:
			break readLoop
		}
	}
	close(in)
	fmt.Printf("\nYou scored %d out of %d\n", correct, total)
}

type problem struct {
	q, a string
}
