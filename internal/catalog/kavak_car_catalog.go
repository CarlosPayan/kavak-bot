package catalog

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Car struct {
	StockID   string
	KM        int
	Price     float64
	Make      string
	Model     string
	Year      int
	Version   string
	Bluetooth string
	Lenght    float64
	Widht     float64
	Height    float64
	CarPlay   string
	Embedding []float32
}

type Catalog struct {
	client *openai.Client
	cars   []Car
}

func NewCatalog(apiKey, path string) (*Catalog, error) {
	cli := openai.NewClient(apiKey)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	if _, err := r.Read(); err != nil {
		return nil, fmt.Errorf("error reading headers of CSV: %w", err)
	}

	var cars []Car
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading row of CSV: %w", err)
		}
		km, _ := strconv.Atoi(record[1])
		price, _ := strconv.ParseFloat(record[2], 64)
		year, _ := strconv.Atoi(record[5])
		lenght, _ := strconv.ParseFloat(record[8], 64)
		width, _ := strconv.ParseFloat(record[9], 64)
		height, _ := strconv.ParseFloat(record[10], 64)

		car := Car{
			StockID:   record[0],
			KM:        km,
			Price:     price,
			Make:      record[3],
			Model:     record[4],
			Year:      year,
			Version:   record[6],
			Bluetooth: record[7],
			Lenght:    lenght,
			Widht:     width,
			Height:    height,
			CarPlay:   record[11],
		}

		embedText := strings.Join([]string{
			car.Make,
			car.Model,
			car.Version,
			strconv.Itoa(car.Year),
			fmt.Sprintf("$%.0f", car.Price),
			fmt.Sprintf("%d km", car.KM),
		}, " ")

		resp, err := cli.CreateEmbeddings(
			context.Background(),
			openai.EmbeddingRequest{
				Model: openai.AdaEmbeddingV2,
				Input: []string{embedText},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error calculating embedding for %s %s: %w", car.Make, car.Model, err)
		}
		car.Embedding = resp.Data[0].Embedding
		cars = append(cars, car)
	}

	return &Catalog{
		client: cli,
		cars:   cars,
	}, nil
}

func cosine(a, b []float32) float32 {
	var dot, normaA, normaB float32
	for i := range a {
		dot += a[i] * b[i]
		normaA += a[i] * a[i]
		normaB += b[i] * b[i]
	}
	return dot / (float32(math.Sqrt(float64(normaA))) * float32(math.Sqrt(float64(normaB))))
}

func (c *Catalog) Search(ctx context.Context, query string, topN int) ([]Car, error) {
	resp, err := c.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: []string{query},
	})
	if err != nil {
		return nil, fmt.Errorf("error calculating embedding for query: %w", err)
	}
	qEmb := resp.Data[0].Embedding

	type scoredCar struct {
		Car
		Score float32
	}
	var scoredList []scoredCar
	scoredList = make([]scoredCar, 0, len(c.cars))

	for _, car := range c.cars {
		sim := cosine(qEmb, car.Embedding)
		scoredList = append(scoredList, scoredCar{
			Car:   car,
			Score: sim,
		})
	}

	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].Score > scoredList[j].Score
	})

	limit := topN
	if limit > len(scoredList) {
		limit = len(scoredList)
	}
	result := make([]Car, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, scoredList[i].Car)
	}
	return result, nil
}
