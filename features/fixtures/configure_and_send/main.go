package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func main() {
	testcase := flag.String("case", "", "test case to run")
	send := flag.String("send", "", "whether to send a session/error or both")
	flag.Parse()

	// Increase publish rate for testing
	bugsnag.DefaultSessionPublishInterval = time.Millisecond * 20

	sentError := false

	switch *testcase {
	case "default":
		caseDefault()
	case "app version":
		caseAppVersion()
	case "app type":
		caseAppType()
	case "legacy endpoint":
		caseLegacyEndpoint()
	case "hostname":
		caseHostname()
	case "release stage":
		caseNotifyReleaseStage()
	case "on before notify":
		caseOnBeforeNotify()
		sentError = true
	case "params filters":
		caseParamsFilters()
	case "project packages":
		caseProjectPackages()
	default:
		panic("No valid test case: " + *testcase)
	}

	if *send == "error" {
		if !sentError {
			sendError()
		}
	} else if *send == "session" {
		bugsnag.StartSession(context.Background())
		time.Sleep(100 * time.Millisecond)
	} else {
		panic("No valid send case: " + *send)
	}
}

func newDefaultConfig() bugsnag.Configuration {
	return bugsnag.Configuration{
		APIKey: os.Getenv("API_KEY"),
		Endpoints: bugsnag.Endpoints{
			Notify:   os.Getenv("NOTIFY_ENDPOINT"),
			Sessions: os.Getenv("SESSIONS_ENDPOINT"),
		},
	}
}

func sendError() {
	notifier := bugsnag.New()
	notifier.NotifySync(fmt.Errorf("oops"), true, bugsnag.MetaData{
		"Account": {
			"Name":           "Company XYZ",
			"Price(dollars)": "1 Million",
		},
	})
}

func caseDefault() {
	config := newDefaultConfig()
	bugsnag.Configure(config)
}

func caseAppVersion() {
	config := newDefaultConfig()
	config.AppVersion = os.Getenv("APP_VERSION")
	bugsnag.Configure(config)
}

func caseAppType() {
	config := newDefaultConfig()
	config.AppType = os.Getenv("APP_TYPE")
	bugsnag.Configure(config)
}

func caseLegacyEndpoint() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:   os.Getenv("API_KEY"),
		Endpoint: os.Getenv("NOTIFY_ENDPOINT"),
	})
}

func caseHostname() {
	config := newDefaultConfig()
	config.Hostname = os.Getenv("HOSTNAME")
	bugsnag.Configure(config)
}

func caseNotifyReleaseStage() {
	config := newDefaultConfig()
	notifyReleaseStages := os.Getenv("NOTIFY_RELEASE_STAGES")
	if notifyReleaseStages != "" {
		config.NotifyReleaseStages = strings.Split(notifyReleaseStages, ",")
	}
	releaseStage := os.Getenv("RELEASE_STAGE")
	if releaseStage != "" {
		config.ReleaseStage = releaseStage
	}
	bugsnag.Configure(config)
}

func caseOnBeforeNotify() {
	config := newDefaultConfig()
	bugsnag.Configure(config)
	bugsnag.OnBeforeNotify(
		func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			if event.Message == "Ignore this error" {
				return fmt.Errorf("not sending errors to ignore")
			}
			// continue notifying as normal
			if event.Message == "Change error message" {
				event.Message = "Error message was changed"
			}
			return nil
		})

	notifier := bugsnag.New()
	notifier.NotifySync(fmt.Errorf("Don't ignore this error"), true)
	notifier.NotifySync(fmt.Errorf("Ignore this error"), true)
	notifier.NotifySync(fmt.Errorf("Change error message"), true)
}

func caseParamsFilters() {
	config := newDefaultConfig()
	paramsFilters := os.Getenv("PARAMS_FILTERS")
	if paramsFilters != "" {
		config.ParamsFilters = strings.Split(paramsFilters, ",")
	}
	bugsnag.Configure(config)
}

func caseProjectPackages() {
	config := newDefaultConfig()
	projectPackages := os.Getenv("PROJECT_PACKAGES")
	if projectPackages != "" {
		config.ProjectPackages = strings.Split(projectPackages, ",")
	}
	bugsnag.Configure(config)
}
