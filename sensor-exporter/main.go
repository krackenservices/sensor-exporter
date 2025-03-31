package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"

	"sensor-exporter/gosensors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fanspeedDesc = prometheus.NewDesc(
		"sensor_lm_fan_speed_rpm",
		"fan speed (rotations per minute).",
		[]string{"fantype", "chip", "adaptor"},
		nil)

	voltageDesc = prometheus.NewDesc(
		"sensor_lm_voltage_volts",
		"voltage in volts",
		[]string{"intype", "chip", "adaptor"},
		nil)

	powerDesc = prometheus.NewDesc(
		"sensor_lm_power_watts",
		"power in watts",
		[]string{"powertype", "chip", "adaptor"},
		nil)

	temperatureDesc = prometheus.NewDesc(
		"sensor_lm_temperature_celsius",
		"temperature in celsius",
		[]string{"temptype", "chip", "adaptor"},
		nil)

	hddTempDesc = prometheus.NewDesc(
		"sensor_hddsmart_temperature_celsius",
		"temperature in celsius",
		[]string{"device", "id"},
		nil)
)

func main() {
	var (
		listenAddress  = flag.String("web.listen-address", ":9255", "Address on which to expose metrics and web interface.")
		metricsPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		hddtempAddress = flag.String("hddtemp-address", "localhost:7634", "Address to fetch hdd metrics from.")
	)
	flag.Parse()

	hddcollector := NewHddCollector(*hddtempAddress)
	if err := hddcollector.Init(); err != nil {
		log.Printf("error reading hddtemps: %v", err)
	}
	prometheus.MustRegister(hddcollector)

	lmscollector := NewLmSensorsCollector()
	lmscollector.Init()
	prometheus.MustRegister(lmscollector)

	http.Handle(*metricsPath, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Sensor Exporter</title></head>
			<body>
			<h1>Sensor Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

type LmSensorsCollector struct{}

func NewLmSensorsCollector() *LmSensorsCollector {
	return &LmSensorsCollector{}
}

func (l *LmSensorsCollector) Init() {
	gosensors.Init()
}

func (l *LmSensorsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- fanspeedDesc
	ch <- powerDesc
	ch <- temperatureDesc
	ch <- voltageDesc
}

func (l *LmSensorsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, chip := range gosensors.GetDetectedChips() {
		chipName := chip.String()
		adaptorName := chip.AdapterName()
		for _, feature := range chip.GetFeatures() {
			value := feature.GetValue()
			label := feature.GetLabel()
			switch {
			case strings.HasPrefix(feature.Name, "fan"):
				ch <- prometheus.MustNewConstMetric(fanspeedDesc, prometheus.GaugeValue, value, label, chipName, adaptorName)
			case strings.HasPrefix(feature.Name, "temp"):
				ch <- prometheus.MustNewConstMetric(temperatureDesc, prometheus.GaugeValue, value, label, chipName, adaptorName)
			case strings.HasPrefix(feature.Name, "in"):
				ch <- prometheus.MustNewConstMetric(voltageDesc, prometheus.GaugeValue, value, label, chipName, adaptorName)
			case strings.HasPrefix(feature.Name, "power"):
				ch <- prometheus.MustNewConstMetric(powerDesc, prometheus.GaugeValue, value, label, chipName, adaptorName)
			}
		}
	}
}

type HddCollector struct {
	address string
	conn    net.Conn
	buf     bytes.Buffer
}

type HddTemperature struct {
	Device             string
	Id                 string
	TemperatureCelsius float64
}

func NewHddCollector(address string) *HddCollector {
	return &HddCollector{address: address}
}

func (h *HddCollector) Init() error {
	conn, err := net.Dial("tcp", h.address)
	if err != nil {
		return fmt.Errorf("error connecting to hddtemp address '%s': %v", h.address, err)
	}
	h.conn = conn
	return nil
}

func (h *HddCollector) readTempsFromConn() (string, error) {
	if h.conn == nil {
		if err := h.Init(); err != nil {
			return "", err
		}
	}
	h.buf.Reset()
	_, err := io.Copy(&h.buf, h.conn)
	if err != nil {
		return "", fmt.Errorf("error reading from hddtemp socket: %v", err)
	}
	return h.buf.String(), nil
}

func (h *HddCollector) Close() error {
	if h.conn != nil {
		return h.conn.Close()
	}
	return nil
}

func parseHddTemps(s string) ([]HddTemperature, error) {
	var hddtemps []HddTemperature
	if len(s) < 1 || s[0] != '|' {
		return nil, fmt.Errorf("error parsing output from hddtemp: %s", s)
	}
	for _, item := range strings.Split(s[1:len(s)-1], "||") {
		hddtemp, err := parseHddTemp(item)
		if err != nil {
			return nil, fmt.Errorf("error parsing hddtemp: %v", err)
		}
		hddtemps = append(hddtemps, hddtemp)
	}
	return hddtemps, nil
}

func parseHddTemp(s string) (HddTemperature, error) {
	pieces := strings.Split(s, "|")
	if len(pieces) != 4 {
		return HddTemperature{}, fmt.Errorf("invalid hddtemp format: %s", s)
	}
	dev, id, temp, unit := pieces[0], pieces[1], pieces[2], pieces[3]

	if unit == "*" {
		return HddTemperature{Device: dev, Id: id, TemperatureCelsius: -1}, nil
	}
	if unit != "C" {
		return HddTemperature{}, fmt.Errorf("unsupported unit from hddtemp: %s", s)
	}

	ftemp, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		return HddTemperature{}, fmt.Errorf("error parsing temperature as float: %s", temp)
	}

	return HddTemperature{Device: dev, Id: id, TemperatureCelsius: ftemp}, nil
}

func (h *HddCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- hddTempDesc
}

func (h *HddCollector) Collect(ch chan<- prometheus.Metric) {
	tempsString, err := h.readTempsFromConn()
	if err != nil {
		log.Printf("error reading temps from hddtemp daemon: %v", err)
		return
	}
	hddtemps, err := parseHddTemps(tempsString)
	if err != nil {
		log.Printf("error parsing temps from hddtemp daemon: %v", err)
		return
	}
	for _, ht := range hddtemps {
		ch <- prometheus.MustNewConstMetric(hddTempDesc, prometheus.GaugeValue, ht.TemperatureCelsius, ht.Device, ht.Id)
	}
}

