package log_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/PermissionData/log"
)

func ExampleNew_customFilters() {
	logger := log.New(log.Config{
		Threshold: log.TraceLevel,
		Encoder:   json.NewEncoder(os.Stdout),
		Filters: []log.Filter{
			log.DefaultFilter,
			func(lvl, threshold log.Level, data log.Data) log.Data {
				data["hey"] = &struct{ Ho bool }{true}
				return data
			},
			func(lvl, threshold log.Level, data log.Data) log.Data {
				if data == nil {
					return nil
				}
				data["@timestamp"] = nil
				return data
			},
		},
	})

	logger.Log(log.InfoLevel, log.Data{
		"foo": "bar",
		"pi":  3.14,
	})
	// Output:
	// {"@timestamp":null,"@version":"1","foo":"bar","hey":{"Ho":true},"log_level":"Info","pi":3.14}
}

func ExampleNew_toFile() {
	logFile, err := ioutil.TempFile("", "service.log")
	if err != nil {
		panic(err)
	}
	defer os.Remove(logFile.Name()) // cleanup for example

	defer logFile.Close()

	logger := log.New(log.Config{
		Threshold: log.TraceLevel,
		Encoder:   json.NewEncoder(logFile),
		Filters: []log.Filter{
			log.DefaultFilter,
			func(lvl, threshold log.Level, data log.Data) log.Data {
				// just for the example output
				if data == nil {
					return nil
				}
				data["@timestamp"] = nil
				return data
			},
		},
	})

	logger.Log(log.InfoLevel, log.Data{
		"foo": "bar",
		"pi":  3.14,
	})

	// output for example
	outFile, err := os.Open(logFile.Name())
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	if _, err := io.Copy(os.Stdout, outFile); err != nil {
		panic(err)
	}
	// Output:
	// {"@timestamp":null,"@version":"1","foo":"bar","log_level":"Info","pi":3.14}
}

func ExampleWithLevels() {
	// given that...
	log.DefaultEncoder = json.NewEncoder(os.Stdout)

	// Setting up a new LevelLogger
	logger := log.WithLevels(log.New(log.Config{
		Filters: []log.Filter{
			log.DefaultFilter,
			func(lvl, threshold log.Level, data log.Data) log.Data {
				data["hey"] = &struct{ Ho bool }{true}
				return data
			},
			func(lvl, threshold log.Level, data log.Data) log.Data {
				if data == nil {
					return nil
				}
				data["@timestamp"] = nil
				return data
			},
		},
	}))

	// Calling,
	logger.Fatal(log.Data{
		"foo": "bar",
		"pi":  3.14,
	})
	// should be the same as calling,
	logger.Log(log.FatalLevel, log.Data{
		"foo": "bar",
		"pi":  3.14,
	})
	// and by default, nothing should be output for,
	logger.Error(log.Data{
		"foo": "bar",
		"pi":  3.14,
	})
	// or,
	logger.Info(log.Data{
		"foo": "bar",
		"pi":  3.14,
	})
	// or,
	logger.Trace(log.Data{
		"foo": "bar",
		"pi":  3.14,
	})
	// Output:
	// {"@timestamp":null,"@version":"1","foo":"bar","hey":{"Ho":true},"log_level":"Fatal","pi":3.14}
	// {"@timestamp":null,"@version":"1","foo":"bar","hey":{"Ho":true},"log_level":"Fatal","pi":3.14}
}
