package proprdr

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

const (
	fileName = "./config.properties"
)

func changeFile(key, newValue string) {
	info, err := os.Lstat(fileName)
	if err != nil {
		log.Fatalf("Can not get FileMode for file %s... %s", fileName, err.Error())
	}

	file, err := os.OpenFile(fileName, os.O_APPEND, info.Mode())
	if err != nil {
		log.Fatalf("Can not open file %s for update... %s", fileName, err.Error())
	}
	defer file.Close()

	newLine := fmt.Sprintf("\n%s=%s\n", key, newValue)
	if _, err := file.Write([]byte(newLine)); err != nil {
		log.Fatal(err)
	}
}

func TestNew(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("A new PropertyFile expected. Obtained error... %s", err.Error())
	}

	const expected = 6
	if propertyFile.Size() != expected {
		t.Errorf("Expected a size of %d elements. Obtained %d", expected, propertyFile.Size())
	}
}

func TestGet(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("PropertyFile creation failed... %s", err.Error())
	}

	expectedValue := "Property File Reader"
	if property, _ := propertyFile.Get("app.name"); property != expectedValue {
		t.Errorf("Expected %s. Obtained %s", expectedValue, property)
	}
}

func TestGetAsInt(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("PropertyFile creation failed... %s", err.Error())
	}

	intValue, err := propertyFile.GetAsInt("app.values.maxConnections")
	if err != nil {
		t.Errorf("Expected a int value. Obtained error... %s", err.Error())
	}

	expectedValue := 20
	if intValue != expectedValue {
		t.Errorf("Expected a value %d. Obtained %d", expectedValue, intValue)
	}
}

func TestGetAsFloat(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("PropertyFile creation failed... %s", err.Error())
	}

	floatValue, err := propertyFile.GetAsFloat("app.amount", 64)
	if err != nil {
		t.Errorf("An error has occurred... %s", err.Error())
	}

	const expected = 30.456
	if floatValue != expected {
		t.Errorf("Expected %f but obtained %f", expected, floatValue)
	}
}

func TestContains(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("PropertyFile creation failed... %s", err.Error())
	}

	expected := false
	if exist := propertyFile.Contains("undefinedKey"); exist != expected {
		t.Errorf("Expected %t but obtained %t", expected, exist)
	}

	expected = true
	if exist := propertyFile.Contains("app.values.maxConnections"); exist != expected {
		t.Errorf("Expected %t but obtained %t", expected, exist)
	}
}

func TestGetAsBool(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("PropertyFile creation failed... %s", err.Error())
	}

	expected := true
	if boolValue := propertyFile.GetAsBool("app.security.active"); expected != boolValue {
		t.Errorf("Expected %t but obtained %t", expected, boolValue)
	}

	expected = false
	if boolValue := propertyFile.GetAsBool("app.connections.retries"); expected != boolValue {
		t.Errorf("Expected %t but obtained %t", expected, boolValue)
	}
}

func TestGetAll(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("PropertyFile creation failed... %s", err.Error())
	}

	const expectedSize = 2
	subMap := propertyFile.GetAll("app.connections")
	if len(subMap) != expectedSize {
		t.Errorf("Expected subMap size of %d. Obtained %d", expectedSize, len(subMap))
	}

	const (
		retriesKey = "app.connections.retries"
		urlKey     = "app.connections.url"
	)

	retriesConnProperty, _ := propertyFile.Get(retriesKey)
	value, ok := subMap[retriesKey]
	if !ok {
		t.Errorf("Expected a property value %s but element not found in subMap", retriesConnProperty)
	}
	if value != retriesConnProperty {
		t.Errorf("Expected a property value %s. Obtained %s", retriesConnProperty, value)
	}

	urlConnProperty, _ := propertyFile.Get(urlKey)
	value, ok = subMap[urlKey]
	if !ok {
		t.Errorf("Expected a property value %s but element not found in subMap", urlConnProperty)
	}
	if value != urlConnProperty {
		t.Errorf("Expected a property value %s. Obtained %s", urlConnProperty, value)
	}
}

func TestHasChanged(t *testing.T) {
	propertyFile, err := New(fileName)
	if err != nil {
		t.Errorf("PropertyFile creation failed... %s", err.Error())
	}

	if changed, _ := propertyFile.HasChanged(); changed != false {
		t.Errorf("Expected a %t. Obtained %t", false, changed)
	}

	// Wait a little time previous to change a file
	time.Sleep(100 * time.Millisecond)
	changeFile("app.newline", "This is a new line")

	if changed, _ := propertyFile.HasChanged(); changed != true {
		t.Errorf("Expected a %t. Obtained %t", true, changed)
	}
}
