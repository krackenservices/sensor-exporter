package gosensors

// #cgo LDFLAGS: -lsensors
// #include <stdlib.h>
// #include <stdio.h>
// #include <sensors/sensors.h>
import "C"

import (
	"fmt"
	"log"
	"unsafe"
)

// SubFeature represents a sensor sub-feature
type SubFeature struct {
	Name    string
	Number  int32
	Type    SubFeatureType
	Mapping int32
	Flags   uint32
	chip    *C.struct_sensors_chip_name
}

// GetValue returns the value of the sub-feature
func (s SubFeature) GetValue() float64 {
	var value C.double
	ret := C.sensors_get_value(s.chip, C.int(s.Number), &value)
	if ret != 0 {
		// In Go 1.21, it's better to handle potential errors
		return 0.0
	}
	return float64(value)
}

// SubFeatureType represents the type of a sensor sub-feature
type SubFeatureType int32

// Constants for sub-feature types
const (
	SubFeatureTypeInInput      SubFeatureType = C.SENSORS_SUBFEATURE_IN_INPUT
	SubFeatureTypeInMin        SubFeatureType = C.SENSORS_SUBFEATURE_IN_MIN
	SubFeatureTypeInMax        SubFeatureType = C.SENSORS_SUBFEATURE_IN_MAX
	SubFeatureTypeInLCrit      SubFeatureType = C.SENSORS_SUBFEATURE_IN_LCRIT
	SubFeatureTypeInCrit       SubFeatureType = C.SENSORS_SUBFEATURE_IN_CRIT
	SubFeatureTypeInAverage    SubFeatureType = C.SENSORS_SUBFEATURE_IN_AVERAGE
	SubFeatureTypeInLowest     SubFeatureType = C.SENSORS_SUBFEATURE_IN_LOWEST
	SubFeatureTypeInHighest    SubFeatureType = C.SENSORS_SUBFEATURE_IN_HIGHEST
	SubFeatureTypeInAlarm      SubFeatureType = C.SENSORS_SUBFEATURE_IN_ALARM
	SubFeatureTypeInMinAlarm   SubFeatureType = C.SENSORS_SUBFEATURE_IN_MIN_ALARM
	SubFeatureTypeInMaxAlarm   SubFeatureType = C.SENSORS_SUBFEATURE_IN_MAX_ALARM
	SubFeatureTypeInBeep       SubFeatureType = C.SENSORS_SUBFEATURE_IN_BEEP
	SubFeatureTypeInLCritAlarm SubFeatureType = C.SENSORS_SUBFEATURE_IN_LCRIT_ALARM
	SubFeatureTypeInCritAlarm  SubFeatureType = C.SENSORS_SUBFEATURE_IN_CRIT_ALARM

	SubFeatureTypeFanInput    SubFeatureType = C.SENSORS_SUBFEATURE_FAN_INPUT
	SubFeatureTypeFanMin      SubFeatureType = C.SENSORS_SUBFEATURE_FAN_MIN
	SubFeatureTypeFanMax      SubFeatureType = C.SENSORS_SUBFEATURE_FAN_MAX
	SubFeatureTypeFanAlarm    SubFeatureType = C.SENSORS_SUBFEATURE_FAN_ALARM
	SubFeatureTypeFanFault    SubFeatureType = C.SENSORS_SUBFEATURE_FAN_FAULT
	SubFeatureTypeFanDiv      SubFeatureType = C.SENSORS_SUBFEATURE_FAN_DIV
	SubFeatureTypeFanBeep     SubFeatureType = C.SENSORS_SUBFEATURE_FAN_BEEP
	SubFeatureTypeFanPulses   SubFeatureType = C.SENSORS_SUBFEATURE_FAN_PULSES
	SubFeatureTypeFanMinAlarm SubFeatureType = C.SENSORS_SUBFEATURE_FAN_MIN_ALARM
	SubFeatureTypeFanMaxAlarm SubFeatureType = C.SENSORS_SUBFEATURE_FAN_MAX_ALARM

	SubFeatureTypeTempInput          SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_INPUT
	SubFeatureTypeTempMax            SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_MAX
	SubFeatureTypeTempMaxHyst        SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_MAX_HYST
	SubFeatureTypeTempMin            SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_MIN
	SubFeatureTypeTempCrit           SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_CRIT
	SubFeatureTypeTempCritHyst       SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_CRIT_HYST
	SubFeatureTypeTempLCrit          SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_LCRIT
	SubFeatureTypeTempEmergency      SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_EMERGENCY
	SubFeatureTypeTempEmergencyHyst  SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_EMERGENCY_HYST
	SubFeatureTypeTempLowest         SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_LOWEST
	SubFeatureTypeTempHighest        SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_HIGHEST
	SubFeatureTypeTempMinHyst        SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_MIN_HYST
	SubFeatureTypeTempLCritHyst      SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_LCRIT_HYST
	SubFeatureTypeTempAlarm          SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_ALARM
	SubFeatureTypeTempMaxAlarm       SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_MAX_ALARM
	SubFeatureTypeTempMinAlarm       SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_MIN_ALARM
	SubFeatureTypeTempCritAlarm      SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_CRIT_ALARM
	SubFeatureTypeTempFault          SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_FAULT
	SubFeatureTypeTempType           SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_TYPE
	SubFeatureTypeTempOffset         SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_OFFSET
	SubFeatureTypeTempBeep           SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_BEEP
	SubFeatureTypeTempEmergencyAlarm SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_EMERGENCY_ALARM
	SubFeatureTypeTempLCritAlarm     SubFeatureType = C.SENSORS_SUBFEATURE_TEMP_LCRIT_ALARM

	SubFeatureTypePowerAverage         SubFeatureType = C.SENSORS_SUBFEATURE_POWER_AVERAGE
	SubFeatureTypePowerHighest         SubFeatureType = C.SENSORS_SUBFEATURE_POWER_AVERAGE_HIGHEST
	SubFeatureTypePowerLowest          SubFeatureType = C.SENSORS_SUBFEATURE_POWER_AVERAGE_LOWEST
	SubFeatureTypePowerInput           SubFeatureType = C.SENSORS_SUBFEATURE_POWER_INPUT
	SubFeatureTypePowerInputHighest    SubFeatureType = C.SENSORS_SUBFEATURE_POWER_INPUT_HIGHEST
	SubFeatureTypePowerInputLowest     SubFeatureType = C.SENSORS_SUBFEATURE_POWER_INPUT_LOWEST
	SubFeatureTypePowerCap             SubFeatureType = C.SENSORS_SUBFEATURE_POWER_CAP
	SubFeatureTypePowerCapHyst         SubFeatureType = C.SENSORS_SUBFEATURE_POWER_CAP_HYST
	SubFeatureTypePowerMax             SubFeatureType = C.SENSORS_SUBFEATURE_POWER_MAX
	SubFeatureTypePowerCrit            SubFeatureType = C.SENSORS_SUBFEATURE_POWER_CRIT
	SubFeatureTypePowerAverageInterval SubFeatureType = C.SENSORS_SUBFEATURE_POWER_AVERAGE_INTERVAL
	SubFeatureTypePowerAlarm           SubFeatureType = C.SENSORS_SUBFEATURE_POWER_ALARM
	SubFeatureTypePowerCapAlarm        SubFeatureType = C.SENSORS_SUBFEATURE_POWER_CAP_ALARM
	SubFeatureTypePowerMaxAlarm        SubFeatureType = C.SENSORS_SUBFEATURE_POWER_MAX_ALARM
	SubFeatureTypePowerCritAlarm       SubFeatureType = C.SENSORS_SUBFEATURE_POWER_CRIT_ALARM

	SubFeatureTypeEnergyInput SubFeatureType = C.SENSORS_SUBFEATURE_ENERGY_INPUT

	SubFeatureTypeCurrInput      SubFeatureType = C.SENSORS_SUBFEATURE_CURR_INPUT
	SubFeatureTypeCurrMin        SubFeatureType = C.SENSORS_SUBFEATURE_CURR_MIN
	SubFeatureTypeCurrMax        SubFeatureType = C.SENSORS_SUBFEATURE_CURR_MAX
	SubFeatureTypeCurrLCrit      SubFeatureType = C.SENSORS_SUBFEATURE_CURR_LCRIT
	SubFeatureTypeCurrCrit       SubFeatureType = C.SENSORS_SUBFEATURE_CURR_CRIT
	SubFeatureTypeCurrAverage    SubFeatureType = C.SENSORS_SUBFEATURE_CURR_AVERAGE
	SubFeatureTypeCurrLowest     SubFeatureType = C.SENSORS_SUBFEATURE_CURR_LOWEST
	SubFeatureTypeCurrHighest    SubFeatureType = C.SENSORS_SUBFEATURE_CURR_HIGHEST
	SubFeatureTypeCurrAlarm      SubFeatureType = C.SENSORS_SUBFEATURE_CURR_ALARM
	SubFeatureTypeCurrMinAlarm   SubFeatureType = C.SENSORS_SUBFEATURE_CURR_MIN_ALARM
	SubFeatureTypeCurrMaxAlarm   SubFeatureType = C.SENSORS_SUBFEATURE_CURR_MAX_ALARM
	SubFeatureTypeCurrBeep       SubFeatureType = C.SENSORS_SUBFEATURE_CURR_BEEP
	SubFeatureTypeCurrLCritAlarm SubFeatureType = C.SENSORS_SUBFEATURE_CURR_LCRIT_ALARM
	SubFeatureTypeCurrCritAlarm  SubFeatureType = C.SENSORS_SUBFEATURE_CURR_CRIT_ALARM

	SubFeatureTypeHumidityInput SubFeatureType = C.SENSORS_SUBFEATURE_HUMIDITY_INPUT

	SubFeatureTypeVid SubFeatureType = C.SENSORS_SUBFEATURE_VID

	SubFeatureTypeIntrusionAlarm SubFeatureType = C.SENSORS_SUBFEATURE_INTRUSION_ALARM
	SubFeatureTypeIntrusionBeep  SubFeatureType = C.SENSORS_SUBFEATURE_INTRUSION_BEEP

	SubFeatureTypeBeepEnable SubFeatureType = C.SENSORS_SUBFEATURE_BEEP_ENABLE

	SubFeatureTypeUnknown SubFeatureType = C.SENSORS_SUBFEATURE_UNKNOWN
)

// Feature represents a sensor feature
type Feature struct {
	Name    string
	Number  int32
	Type    FeatureType
	chip    *C.struct_sensors_chip_name
	feature *C.struct_sensors_feature
}

// GetSubFeatures returns all sub-features of a feature
func (f Feature) GetSubFeatures() []SubFeature {
	var subfeatures []SubFeature
	var count C.int = 0

	for {
		resp := C.sensors_get_all_subfeatures(f.chip, f.feature, &count)
		if resp == nil {
			break
		}

		subfeature := SubFeature{
			Name:    C.GoString(resp.name),
			Number:  int32(resp.number),
			Type:    SubFeatureType(resp._type),
			Mapping: int32(resp.mapping),
			Flags:   uint32(resp.flags),
			chip:    f.chip,
		}

		subfeatures = append(subfeatures, subfeature)
	}

	return subfeatures
}

// FeatureType represents the type of a sensor feature
type FeatureType int32

// Constants for feature types
const (
	FeatureTypeIn         FeatureType = C.SENSORS_FEATURE_IN
	FeatureTypeFan        FeatureType = C.SENSORS_FEATURE_FAN
	FeatureTypeTemp       FeatureType = C.SENSORS_FEATURE_TEMP
	FeatureTypePower      FeatureType = C.SENSORS_FEATURE_POWER
	FeatureTypeEnergy     FeatureType = C.SENSORS_FEATURE_ENERGY
	FeatureTypeCurr       FeatureType = C.SENSORS_FEATURE_CURR
	FeatureTypeHumidity   FeatureType = C.SENSORS_FEATURE_HUMIDITY
	FeatureTypeMaxMain    FeatureType = C.SENSORS_FEATURE_MAX_MAIN
	FeatureTypeVid        FeatureType = C.SENSORS_FEATURE_VID
	FeatureTypeIntrusion  FeatureType = C.SENSORS_FEATURE_INTRUSION
	FeatureTypeMaxOther   FeatureType = C.SENSORS_FEATURE_MAX_OTHER
	FeatureTypeBeepEnable FeatureType = C.SENSORS_FEATURE_BEEP_ENABLE
	FeatureTypeMax        FeatureType = C.SENSORS_FEATURE_MAX
	FeatureTypeUnknown    FeatureType = C.SENSORS_FEATURE_UNKNOWN
)

// GetLabel returns the label of a feature
func (f Feature) GetLabel() string {
	clabel := C.sensors_get_label(f.chip, f.feature)
	if clabel == nil {
		return "" // Handle potential nil pointer
	}
	golabel := C.GoString(clabel)
	C.free(unsafe.Pointer(clabel))
	return golabel
}

// GetValue returns the value of a feature
func (f Feature) GetValue() float64 {
	subfeatures := f.GetSubFeatures()
	if len(subfeatures) == 0 {
		return 0.0 // Return default value if no subfeatures found
	}
	return subfeatures[0].GetValue()
}

// Bus represents a sensor bus
type Bus struct {
	Type int16
	Nr   int16
	bus  *C.struct_sensors_bus_id
}

// String returns the string representation of a bus
func (b Bus) String() string {
	if b.Type == -1 {
		return "*"
	}
	// In Go 1.21, we should handle potential nil pointers
	if b.bus == nil {
		return fmt.Sprintf("bus-%d-%d", b.Type, b.Nr)
	}
	adapterName := C.sensors_get_adapter_name(b.bus)
	if adapterName == nil {
		return fmt.Sprintf("bus-%d-%d", b.Type, b.Nr)
	}
	return C.GoString(adapterName)
}

// Chip represents a sensor chip
type Chip struct {
	Prefix string
	Bus    Bus
	Addr   int32
	Path   string
	chip   *C.struct_sensors_chip_name
}

// String returns the string representation of a chip
func (c Chip) String() string {
	var buffer [200]C.char
	
	// In Go 1.21, we should handle potential errors more explicitly
	if c.chip == nil {
		return c.Prefix
	}

	len := C.sensors_snprintf_chip_name(&buffer[0], C.size_t(len(buffer)), c.chip)
	if len <= 0 {
		return c.Prefix
	}

	return C.GoStringN(&buffer[0], len)
}

// AdapterName returns the adapter name of a chip
func (c Chip) AdapterName() string {
	return c.Bus.String()
}

// GetFeatures returns all features of a chip
func (c Chip) GetFeatures() []Feature {
	var features []Feature
	var count C.int = 0

	for {
		resp := C.sensors_get_features(c.chip, &count)
		if resp == nil {
			break
		}

		feature := Feature{
			Name:    C.GoString(resp.name),
			Number:  int32(resp.number),
			Type:    FeatureType(resp._type),
			chip:    c.chip,
			feature: resp,
		}

		features = append(features, feature)
	}

	return features
}

// Init initializes the sensors library
func Init() {
	// In Go 1.21, let's use a more robust approach to handle config file
	filename := C.CString("/etc/sensors3.conf")
	defer C.free(unsafe.Pointer(filename))

	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))

	fp, err := C.fopen(filename, mode)
	
	// If we can't open the config file, try initializing without it
	if fp == nil {
		// Log the error but don't fatal - library may still work with defaults
		log.Printf("Warning: Could not open sensors config file: %v", err)
		C.sensors_init(nil)
		return
	}
	
	defer C.fclose(fp)
	ret := C.sensors_init(fp)
	
	if ret != 0 {
		log.Printf("Warning: sensors_init returned error code: %d", ret)
	}
}

// Cleanup cleans up the sensors library
func Cleanup() {
	C.sensors_cleanup()
}

// GetDetectedChips returns all detected sensor chips
func GetDetectedChips() []Chip {
	var chips []Chip
	var count C.int = 0

	for {
		resp := C.sensors_get_detected_chips(nil, &count)
		if resp == nil {
			break
		}

		bus := Bus{
			Type: int16(resp.bus._type),
			Nr:   int16(resp.bus.nr),
			bus:  &resp.bus,
		}

		prefix := ""
		if resp.prefix != nil {
			prefix = C.GoString(resp.prefix)
		}
		
		path := ""
		if resp.path != nil {
			path = C.GoString(resp.path)
		}

		chip := Chip{
			Prefix: prefix,
			Bus:    bus,
			Addr:   int32(resp.addr),
			Path:   path,
			chip:   resp,
		}

		chips = append(chips, chip)
	}

	return chips
}
