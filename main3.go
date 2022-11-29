package main

import (
	"fmt"
	_ "github.com/sjwhitworth/golearn"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/ensemble"
	"github.com/sjwhitworth/golearn/evaluation"
	"os"
)

func main() {
	rawData, err := base.ParseCSVToInstances("db/out.csv", true)
	if err != nil {
		panic(err)
	}

	//fmt.Println(rawData)

	//cls := knn.NewKnnClassifier("euclidean", "linear", 2)
	cls := ensemble.NewRandomForest(10, 5)
	trainData, valData := base.InstancesTrainTestSplit(rawData, 0.50)
	err = cls.Fit(trainData)
	if err != nil {
		println("FIT ERROR", err.Error())
	}
	validations, err := cls.Predict(valData)
	if err != nil {
		println("TRAIN PREDICT ERROR", err.Error())
	}

	confusionMat, err := evaluation.GetConfusionMatrix(valData, validations)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))
	err = cls.Save("model.h")
	if err != nil {
		println("MODEL SAVE ERROR", err.Error())
	}
	testData, err := base.ParseCSVToInstances("db/out2.csv", true)
	cls2 := ensemble.RandomForest{}
	err = cls2.Load("model.h")
	if err != nil {
		println("COULD NOT LOAD MODEL", err.Error())
	}
	predictions, err := cls2.Predict(testData)
	if err != nil {
		println("TEST PREDICT ERROR", err.Error())
	}
	println(predictions.Size())
	writer, _ := os.OpenFile("db/out3.csv", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	err = base.SerializeInstancesToCSVStream(predictions, writer)
	if err != nil {
		println("CSV SAVE ERROR", err.Error())
	}
	writer.Close()

}
