// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command listresources lists the Google Cloud Monitoring v3 Environment against an authenticated user.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/monitoring/v3"
)

const metric = "compute.googleapis.com/instance/cpu/usage_time"

// listMonitoredResourceDescriptor lists all the resources available to be monitored in the API.
func listMonitoredResourceDescriptors(s *monitoring.Service, project string) error {
	resp, err := s.Projects.MonitoredResourceDescriptors.List(project).Do()
	if err != nil {
		return fmt.Errorf("Could not list time series: %v", err)
	}

	log.Printf("listMonitoredResourceDescriptors: %s\n", formatResource(resp))
	return nil
}

// listMetricDescriptors lists the metrics specified by the metric constant.
func listMetricDescriptors(s *monitoring.Service, projectResource string) error {
	resp, err := s.Projects.MetricDescriptors.List(projectResource).
		Filter(fmt.Sprintf("metric.type=%q", metric)).
		Do()
	if err != nil {
		return fmt.Errorf("Could not list metric descriptors: %v", err)
	}

	log.Printf("listMetricDescriptors %s\n", formatResource(resp))
	return nil
}

// listTimesSeries lists all the timeseries created for metric created in a 5
// minute interval an hour ago
func listTimeSeries(s *monitoring.Service, projectResource string) error {
	startTime := time.Now().UTC().Add(-time.Hour)
	endTime := startTime.Add(5 * time.Minute)

	resp, err := s.Projects.TimeSeries.List(projectResource).
		PageSize(3).
		Filter(fmt.Sprintf("metric.type=\"%s\"", metric)).
		IntervalStartTime(startTime.Format(time.RFC3339)).
		IntervalEndTime(endTime.Format(time.RFC3339)).
		Do()
	if err != nil {
		return fmt.Errorf("Could not list time series: %v", err)
	}

	log.Printf("listTimeseries %s\n", formatResource(resp))
	return nil
}

func main() {
	// walk through the basic calls of the Monitoring API
	client, err := google.DefaultClient(
		oauth2.NoContext,
		monitoring.CloudPlatformScope,
		monitoring.MonitoringScope,
		monitoring.MonitoringReadScope,
		monitoring.MonitoringWriteScope,
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) < 2 {
		fmt.Println("Usage: auth.go <project_id>")
		return
	}
	projectResource := "projects/" + os.Args[1]
	s, err := monitoring.New(client)
	if err != nil {
		log.Fatal(err)
	}

	if err := listMonitoredResourceDescriptors(s, projectResource); err != nil {
		log.Fatal(err)
	}
	if err := listMetricDescriptors(s, projectResource); err != nil {
		log.Fatal(err)
	}
	if err := listTimeSeries(s, projectResource); err != nil {
		log.Fatal(err)
	}
}

// printResource prints out our API response objects as JSON.
func formatResource(resource interface{}) []byte {
	b, err := json.MarshalIndent(resource, "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}
