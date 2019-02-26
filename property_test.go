package proprdr

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

const (
	fileName         = "./config.properties"
	unitTestFileName = "unit_test.properties"
)

func init() {
	copyPropertyFile()
}

func copyPropertyFile() {
	source, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Can not open file %s to copy... Error %s", fileName, err.Error())
	}
	defer source.Close()

	destination, err := os.Create(unitTestFileName)
	if err != nil {
		log.Fatalf("Can not create a new file %s... Error %s", unitTestFileName, err.Error())
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		log.Fatalf("Can not possible copy file %s to %s... Error %s", fileName, unitTestFileName, err.Error())
	}
}

func getNewFile() PropertyFile {
	pFile, err := New(unitTestFileName)
	if err != nil {
		log.Fatalf("PropertyFile creation failed... %s", err.Error())
	}

	return pFile
}

func changeFile(key, newValue string) {
	info, err := os.Lstat(unitTestFileName)
	if err != nil {
		log.Fatalf("Can not get FileMode for file %s... %s", unitTestFileName, err.Error())
	}

	file, err := os.OpenFile(unitTestFileName, os.O_APPEND|os.O_WRONLY, info.Mode())
	if err != nil {
		log.Fatalf("Can not open file %s for update... %s", unitTestFileName, err.Error())
	}
	defer file.Close()

	newLine := fmt.Sprintf("%s=%s\n", key, newValue)
	if _, err := file.Write([]byte(newLine)); err != nil {
		log.Fatal(err)
	}
}

func TestNew(t *testing.T) {
	pFile, err := New(unitTestFileName)
	if err != nil {
		t.Errorf("A new PropertyFile expected. Obtained error... %s", err.Error())
	}

	const expected = 6
	if pFile.Size() != expected {
		t.Errorf("Expected a size of %d elements. Obtained %d", expected, pFile.Size())
	}
}

func TestGet(t *testing.T) {
	pFile := getNewFile()

	expectedValue := "Property File Reader"
	if property, _ := pFile.Get("app.name"); property != expectedValue {
		t.Errorf("Expected %s. Obtained %s", expectedValue, property)
	}
}

func TestGetAsInt(t *testing.T) {
	pFile := getNewFile()

	intValue, err := pFile.GetAsInt("app.values.maxConnections")
	if err != nil {
		t.Errorf("Expected a int value. Obtained error... %s", err.Error())
	}

	expectedValue := 20
	if intValue != expectedValue {
		t.Errorf("Expected a value %d. Obtained %d", expectedValue, intValue)
	}
}

func TestGetAsFloat(t *testing.T) {
	pFile := getNewFile()

	floatValue, err := pFile.GetAsFloat("app.amount", 64)
	if err != nil {
		t.Errorf("An error has occurred... %s", err.Error())
	}

	const expected = 30.456
	if floatValue != expected {
		t.Errorf("Expected %f but obtained %f", expected, floatValue)
	}
}

func TestContains(t *testing.T) {
	pFile := getNewFile()

	expected := false
	if exist := pFile.Contains("undefinedKey"); exist != expected {
		t.Errorf("Expected %t but obtained %t", expected, exist)
	}

	expected = true
	if exist := pFile.Contains("app.values.maxConnections"); exist != expected {
		t.Errorf("Expected %t but obtained %t", expected, exist)
	}
}

func TestGetAsBool(t *testing.T) {
	pFile := getNewFile()

	expected := true
	if boolValue := pFile.GetAsBool("app.security.active"); expected != boolValue {
		t.Errorf("Expected %t but obtained %t", expected, boolValue)
	}

	expected = false
	if boolValue := pFile.GetAsBool("app.connections.retries"); expected != boolValue {
		t.Errorf("Expected %t but obtained %t", expected, boolValue)
	}
}

func TestGetAll(t *testing.T) {
	pFile := getNewFile()

	const expectedSize = 2
	subMap := pFile.GetAll("app.connections")
	if len(subMap) != expectedSize {
		t.Errorf("Expected subMap size of %d. Obtained %d", expectedSize, len(subMap))
	}

	const (
		retriesKey = "app.connections.retries"
		urlKey     = "app.connections.url"
	)

	retriesConnProperty, _ := pFile.Get(retriesKey)
	value, ok := subMap[retriesKey]
	if !ok {
		t.Errorf("Expected a property value %s but element not found in subMap", retriesConnProperty)
	}
	if value != retriesConnProperty {
		t.Errorf("Expected a property value %s. Obtained %s", retriesConnProperty, value)
	}

	urlConnProperty, _ := pFile.Get(urlKey)
	value, ok = subMap[urlKey]
	if !ok {
		t.Errorf("Expected a property value %s but element not found in subMap", urlConnProperty)
	}
	if value != urlConnProperty {
		t.Errorf("Expected a property value %s. Obtained %s", urlConnProperty, value)
	}
}

func TestHasChanged(t *testing.T) {
	pFile := getNewFile()

	if changed, _ := pFile.HasChanged(); changed != false {
		t.Errorf("Expected a %t. Obtained %t", false, changed)
	}

	// Wait a little time previous to change a file
	time.Sleep(100 * time.Millisecond)
	changeFile("app.newline", "This is a new line")

	if changed, _ := pFile.HasChanged(); changed != true {
		t.Errorf("Expected a %t. Obtained %t", true, changed)
	}
}

func TestRefresh(t *testing.T) {
	pFile := getNewFile()
	const propertyName = "app.changes"
	const propertyValueExpected = "true"

	// if property does not exist, return a error
	if foundValue, err := pFile.Get(propertyName); err == nil {
		t.Errorf("Expected 'Property not found' error and value ''. Obtained %s value", foundValue)
	}

	// add a new property
	changeFile(propertyName, propertyValueExpected)

	// refresh the file
	if err := pFile.Refresh(); err != nil {
		t.Errorf("An error has occurred... %s", err)
	}

	// if the file is refreshed correctly, now you should find the property
	foundValue, err := pFile.Get(propertyName)
	if err != nil {
		t.Errorf("Expected %s value. Obtained error... %s", propertyValueExpected, err.Error())
	}

	if foundValue != propertyValueExpected {
		t.Errorf("Expected %s value. Obtained %s", propertyValueExpected, foundValue)
	}
}

func TestUGet(t *testing.T) {

}
