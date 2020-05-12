package foods

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"sort"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

var (
	graphModel   *tf.Graph
	sessionModel *tf.Session
	labels       []string
)

func init() {
	if err := loadModel(); err != nil {
		log.Fatal(err)
		return
	}
}

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

// TODO Recognize...
func Recognize(image *bytes.Buffer, imageFormat string) (string, float32) {
	tensor, err := makeTensorFromImage(image, imageFormat)
	if err != nil {
		log.Println(err)
		return "", 0
	}

	// Run inference
	output, err := sessionModel.Run(
		map[tf.Output]*tf.Tensor{
			graphModel.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graphModel.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		log.Println(err)
		return "", 0
	}

	if len(output) == 0 {
		return "", 0
	}

	result := findBestLabel(output[0].Value().([][]float32)[0])
	return result.Label, result.Probability
}

// TODO loadModel...
func loadModel() error {
	// Load inception model
	model, err := ioutil.ReadFile("./foods/model/tensorflow_inception_graph.pb")
	if err != nil {
		return err
	}
	graphModel = tf.NewGraph()
	if err := graphModel.Import(model, ""); err != nil {
		return err
	}

	sessionModel, err = tf.NewSession(graphModel, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Load labels
	labelsFile, err := os.Open("./foods/model/imagenet_comp_graph_label_strings.txt")
	if err != nil {
		return err
	}
	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)
	// Labels are separated by newlines
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func findBestLabel(probabilities []float32) LabelResult {
	var result []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		result = append(result, LabelResult{Label: labels[i], Probability: p})
	}

	sort.Sort(byProbability(result))

	return result[0]
}

type byProbability []LabelResult

func (a byProbability) Len() int           { return len(a) }
func (a byProbability) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byProbability) Less(i, j int) bool { return a[i].Probability > a[j].Probability }
