package DS18B20

import (
	"bufio"
	"errors"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type DS18B20 struct {
	address  string
	filePath string
}

func Init(address string) *DS18B20 {
	ds := new(DS18B20)
	ds.address = address
	ds.setSensorFilePath()

	return ds
}

func (ds *DS18B20) setSensorFilePath() {
	ds.filePath = "/sys/bus/w1/devices/" + ds.address + "/w1_slave"
}

func (ds *DS18B20) GetTemperature() (float64, error) {
	temperature, err := ds.getTemperatureFromFile()

	if err != nil {
		return 0, err
	}

	var comaIndex int = len(temperature) - 3
	var temperatureFixed string = temperature[:comaIndex] + "." + temperature[comaIndex:]

	parsed, err := strconv.ParseFloat(temperatureFixed, 64)

	if err != nil {
		return 0, errors.New("Unable to parse temperature.")
	}

	return math.Round(parsed*100) / 100, nil
}

func (ds *DS18B20) getTemperatureFromFile() (string, error) {
	var measuredTemperature string

	file, err := os.Open(ds.filePath)

	if err != nil {
		return "", errors.New("Failed to read the temperature. Cannot open the file.")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var line string = scanner.Text()

		re := regexp.MustCompile("t=.*")
		var temperature string = re.FindString(line)

		if len(temperature) > 0 {
			measuredTemperature = strings.Split(temperature, "=")[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return "", errors.New("Failed to read the temperature. Cannot read the file.")
	}

	if len(measuredTemperature) < 1 {
		return "", errors.New("Failed to read the temperature. Cannot parse the file.")
	}

	return measuredTemperature, nil
}
